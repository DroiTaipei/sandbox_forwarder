package memento

import "time"

// Config defines the items which can be set in the config file
type Config struct {
	Addr    string
	MaxConn int
	// grpc timeout
	Timeout time.Duration
}
