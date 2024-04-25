package setup

// Ports server nested block configuration
type Ports struct {
	HTTP int `hcl:"http"`
	GRPC int `hcl:"grpc"`
}

// Server configuration
type Server struct {
	BindAddr string `hcl:"bind_addr"`
	Ports    Ports  `hcl:"ports,block"`
}
