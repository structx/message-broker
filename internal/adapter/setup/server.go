package setup

// Server
type Server struct {
	HttpPort string `env:"SERVER_HTTP_PORT"`
	GrpcPort int    `env:"SERVER_GRPC_PORT"`
}
