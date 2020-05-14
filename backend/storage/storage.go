package storage

type KeyValue struct {
	Key   string
	Value string
}

// IndexFile defines mapreduce persist file.
type IndexFile struct {
	Temporary bool
	MapIdx    int32
	ReduceIdx int32
	Body      *[]KeyValue
}

// Persister defines persist interface methods.
type Persister interface {
	Init() (err error)
	Destroy() (err error)
	RetrieveMRIdx(path string) (mIdx int32, rIdx int32)
	Writable(file IndexFile) (path string, err error)
	Commit(file IndexFile) (path string, err error)
	Readable(file IndexFile) (path string, err error)
	Request(file IndexFile) (path string, err error)
	Abort(file IndexFile) (err error)
	Delete(file IndexFile) (err error)
}
