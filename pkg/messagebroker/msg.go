package messagebroker

// Msg
type Msg interface {
	Marshal() ([]byte, error)
}
