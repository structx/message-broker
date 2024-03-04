package setup

// Server configuration
type Server struct {
	HTTPPort string `env:"SERVER_HTTP_PORT"`
	GRPCPort int    `env:"SERVER_GRPC_PORT"`
}
