package master

// ServiceConfig defines master config
type ServiceConfig struct {
	GRPCEndpoint string `json:"grpc_endpoint"`
	HTTPEndpoint string `json:"http_endpoint"`
}
