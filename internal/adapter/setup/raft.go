package setup

// Raft configuration
type Raft struct {
	Bootstrap bool   `hcl:"bootstrap"`
	LocalID   string `hcl:"local_id"`
	BaseDir   string `hcl:"base_dir"`
}
