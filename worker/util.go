package worker

import "github.com/amazingchow/mapreduce/backend/storage"

type KeyValueList []storage.KeyValue

func (kvs KeyValueList) Len() int           { return len(kvs) }
func (kvs KeyValueList) Swap(i, j int)      { kvs[i], kvs[j] = kvs[j], kvs[i] }
func (kvs KeyValueList) Less(i, j int) bool { return kvs[i].Key < kvs[j].Key }
