package worker

import "github.com/amazingchow/mapreduce/backend/storage"

// ServiceConfig defines worker config.
type ServiceConfig struct {
	ID                 string            `json:"id"`
	MasterGRPCEndpoint string            `json:"master_grpc_endpoint"`
	S3                 *storage.S3Config `json:"s3"`
	DumpRootPath       string            `json:"dump_root_path"`
}
