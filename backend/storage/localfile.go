package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/rs/zerolog/log"
)

// LocalFilePersister provides local persist service.
type LocalFilePersister struct {
	root string
}

// NewLocalFilePersister returns local persist service instance.
// For unix-like user, your provided path must meet unix naming rules for path.
func NewLocalFilePersister(path string) *LocalFilePersister {
	return &LocalFilePersister{
		root: path,
	}
}

// Init inits local persist service.
func (f *LocalFilePersister) Init() error {
	return os.MkdirAll(f.root, 0755)
}

// Destroy destroys local persist service.
// Not implemented for LocalFilePersister.
func (f *LocalFilePersister) Destroy() error {
	return nil
}

func (f *LocalFilePersister) localPath(file IndexFile) string {
	if file.Temporary {
		return filepath.Join(f.root, fmt.Sprintf("MR-Map-%d-Reduce-%d.tmp", file.MapIdx, file.ReduceIdx))
	}
	return filepath.Join(f.root, fmt.Sprintf("MR-Reduce-%d.final", file.ReduceIdx))
}

// RetrieveMRIdx retrieves MapIdx and ReduceIdx from given path.
// Path must meet naming rule "/path/to/MR-Map-%d-Reduce-%d.tmp".
func (f *LocalFilePersister) RetrieveMRIdx(path string) (int32, int32) {
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
func (f *LocalFilePersister) Writable(file IndexFile) (string, error) {
	return f.localPath(file), nil
}

// Commit writes the index file right now.
func (f *LocalFilePersister) Commit(file IndexFile) (string, error) {
	path := f.localPath(file)

	fd, err := os.Create(path)
	if err != nil {
		log.Error().Err(err).Msgf("cannot write local index file, file=%s", path)
		return "", err
	}
	defer fd.Close()

	enc := json.NewEncoder(fd)
	for _, kv := range *(file.Body) {
		if err := enc.Encode(&kv); err != nil {
			log.Error().Err(err).Msgf("cannot write local index file, file=%s", path)
			return "", err
		}
	}

	log.Debug().Msgf("write local index file, file=%s", path)

	return path, nil
}

// Readable checks whether the index file is readable right now or not.
func (f *LocalFilePersister) Readable(file IndexFile) (string, error) {
	path := f.localPath(file)

	if _, err := os.Stat(path); err != nil {
		log.Error().Err(err).Msgf("cannot stat local index file, file=%s", path)
		return "", err
	}

	return path, nil
}

// Request reads the index file right now.
func (f *LocalFilePersister) Request(file IndexFile) (string, error) {
	path := f.localPath(file)

	fd, err := os.Open(path)
	if err != nil {
		log.Error().Err(err).Msgf("cannot read local index file, file=%s", path)
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
				log.Error().Err(err).Msgf("cannot read local index file, file=%s", path)
			}
		}
		*(file.Body) = append(*(file.Body), kv)
	}

	return path, nil
}

// Abort aborts the index file right now.
// Not implemented for LocalFilePersister.
func (f *LocalFilePersister) Abort(file IndexFile) error {
	return nil
}

// Delete deletes the index file right now.
func (f *LocalFilePersister) Delete(file IndexFile) error {
	path := f.localPath(file)

	if err := os.Remove(path); err != nil {
		log.Error().Err(err).Msgf("cannot delete local index file, file=%s", path)
		return err
	}

	return nil
}
