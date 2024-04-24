package setup

// Logger configuration
type Logger struct {
	Level string `hcl:"log_level"`
	Path  string `hcl:"log_path"`
}
