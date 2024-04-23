package setup

// Server configuration
type Server struct {
	Addr     string `hcl:"bind_addr"`
	HTTPPort string `hcl:"SERVER_HTTP_PORT"`
	GRPCPort int    `hcl:"SERVER_GRPC_PORT"`
}
