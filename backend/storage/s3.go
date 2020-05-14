package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/cenkalti/backoff"
	minio "github.com/minio/minio-go"
	"github.com/rs/zerolog/log"
)

// S3Persister provides s3 persist service.
type S3Persister struct {
	conf   *S3Config
	root   string
	tmp    string
	client *minio.Client
}

// S3Config defines s3 persist service config.
type S3Config struct {
	Endpoint          string `json:"endpoint"`
	AccessKeyID       string `json:"access_key_id"`
	SecretAccessKeyID string `json:"secret_access_key_id"`
	UseSSL            bool   `json:"use_ssl"`
	Bucket            string `json:"bucket"`
}

// NewS3Persister returns s3 persist service instance.
func NewS3Persister(conf *S3Config) (*S3Persister, error) {
	client, err := minio.New(conf.Endpoint, conf.AccessKeyID, conf.SecretAccessKeyID, conf.UseSSL)
	if err != nil {
		log.Error().Err(err).Msg("cannot create s3 client")
		return nil, err
	}

	tmpDir, err := ioutil.TempDir("", "mapreduce")
	if err != nil {
		log.Error().Err(err).Msg("cannot create temporary dir '/tmp/mapreduce'")
		return nil, err
	}

	return &S3Persister{
		conf:   conf,
		root:   "mapreduce",
		tmp:    tmpDir,
		client: client,
	}, nil
}

// Init inits s3 persist service.
func (s *S3Persister) Init() error {
	ok, err := s.client.BucketExists(s.conf.Bucket)
	if err != nil || !ok {
		log.Error().Err(err).Msgf("bucket <%s> not exist", s.conf.Bucket)
		return fmt.Errorf("bucket <%s> not exist", s.conf.Bucket)
	}
	return nil
}

// Destroy destroys s3 persist service.
// Remove all local index files.
func (s *S3Persister) Destroy() error {
	return os.RemoveAll(s.tmp)
}

func (s *S3Persister) remotePath(file IndexFile) string {
	if file.Temporary {
		return filepath.Join(s.root, fmt.Sprintf("MR-Map-%d-Reduce-%d.tmp", file.MapIdx, file.ReduceIdx))
	}
	return filepath.Join(s.root, fmt.Sprintf("MR-Reduce-%d.final", file.ReduceIdx))
}

func (s *S3Persister) localPath(file IndexFile) string {
	if file.Temporary {
		return filepath.Join(s.tmp, fmt.Sprintf("MR-Map-%d-Reduce-%d.tmp", file.MapIdx, file.ReduceIdx))
	}
	return filepath.Join(s.tmp, fmt.Sprintf("MR-Reduce-%d.final", file.ReduceIdx))
}

// RetrieveMRIdx retrieves MapIdx and ReduceIdx from given path.
// Path must meet naming rule "/path/to/MR-Map-%d-Reduce-%d.tmp".
func (s *S3Persister) RetrieveMRIdx(path string) (int32, int32) {
	reg := regexp.MustCompile(`[0-9]`)
	ret := reg.FindAllString(path, -1)
	if len(ret) != 2 {
		return -1, -1
	}
	mIdx, _ := strconv.Atoi(ret[0])
	rIdx, _ := strconv.Atoi(ret[1])
	return int32(mIdx), int32(rIdx)
}

// Writable checks whether the index file is writable right now or not.
// Create local index file.
func (s *S3Persister) Writable(file IndexFile) (string, error) {
	lPath := s.localPath(file)

	fd, err := os.Create(lPath)
	if err != nil {
		log.Error().Err(err).Msgf("cannot create local index file, file=%s", lPath)
		return "", err
	}
	defer fd.Close()

	enc := json.NewEncoder(fd)
	for _, kv := range *(file.Body) {
		if err := enc.Encode(&kv); err != nil {
			log.Error().Err(err).Msgf("cannot create local index file, file=%s", lPath)
			return "", err
		}
	}

	return lPath, nil
}

// Commit writes the index file right now.
// Upload local index file to s3.
func (s *S3Persister) Commit(file IndexFile) (string, error) {
	rPath := s.remotePath(file)
	lPath := s.localPath(file)

	retry := 0
	operation := func() error {
		n, err := s.client.FPutObject(s.conf.Bucket, rPath, lPath, minio.PutObjectOptions{})
		if err != nil {
			log.Warn().Err(err).Msgf("cannot write index file to s3, retry=%d, object=%s, file=%s, file size=%d, uploaded=%d",
				retry, rPath, lPath, fileSize(lPath), n)
			retry++
			return err
		}
		log.Debug().Msgf("write index file to s3, object=%s", rPath)
		return nil
	}

	notify := func(err error, sec time.Duration) {
		if err != nil {
			log.Info().Msgf("will retry in %.1fs", sec.Seconds())
		}
	}

	if err := backoff.RetryNotify(operation, backoffConfig(), notify); err != nil {
		log.Error().Err(err).Msgf("cannot write index file to s3, object=%s", rPath)
		return "", err
	}

	return lPath, nil
}

// Readable checks whether the index file is readable right now or not.
// Download remote index file at s3 to local.
func (s *S3Persister) Readable(file IndexFile) (string, error) {
	rPath := s.remotePath(file)

	obj, err := s.client.GetObject(s.conf.Bucket, rPath, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}

	lPath := s.localPath(file)
	fd, err := os.Create(lPath)
	if err != nil {
		log.Error().Err(err).Msgf("cannot create local index file, file=%s", lPath)
		return "", err
	}
	defer fd.Close()

	if _, err = io.Copy(fd, obj); err != nil {
		log.Error().Err(err).Msgf("cannot copy remote index file to local index file, object=%s, file=%s", rPath, lPath)
		return "", err
	}

	return lPath, nil
}

// Request reads the local index file right now.
func (s *S3Persister) Request(file IndexFile) (string, error) {
	lPath := s.localPath(file)

	fd, err := os.Open(lPath)
	if err != nil {
		log.Error().Err(err).Msgf("cannot read local index file, file=%s", lPath)
		return "", err
	}
	defer fd.Close()

	dec := json.NewDecoder(fd)
	for {
		var kv KeyValue
		if err := dec.Decode(&kv); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Error().Err(err).Msgf("cannot read local index file, file=%s", lPath)
			}
		}
		*(file.Body) = append(*(file.Body), kv)
	}

	return lPath, nil
}

// Abort aborts the local index file right now.
func (s *S3Persister) Abort(file IndexFile) error {
	lPath := s.localPath(file)

	if err := os.Remove(lPath); err != nil {
		log.Warn().Err(err).Msgf("cannot delete local index file, file=%s", lPath)
		return err
	}

	return nil
}

// Delete deletes the remote index file right now.
func (s *S3Persister) Delete(file IndexFile) error {
	rPath := s.remotePath(file)

	if err := s.client.RemoveObject(s.conf.Bucket, rPath); err != nil {
		log.Warn().Err(err).Msgf("cannot delete remote index file, object=%s", rPath)
		return err
	}

	return nil
}

func fileSize(file string) int64 {
	fd, err := os.Stat(file)
	if err != nil {
		log.Error().Err(err).Msgf("cannot stat local index file, file=%s", file)
		return 0
	}
	return fd.Size()
}

func backoffConfig() backoff.BackOff {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 600 * time.Millisecond
	bo.Multiplier = 10.0
	bo.MaxInterval = 30 * time.Second
	bo.MaxElapsedTime = 90 * time.Second

	forwarder_bo := backoff.WithMaxRetries(bo, 5)
	forwarder_bo.Reset()
	return forwarder_bo
}
