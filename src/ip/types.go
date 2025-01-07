package ip

import (
	"context"
	"io"
)

type Config struct {
	Ctx      context.Context
	IpChan   <-chan string
	Filename string // file that holds the most recent IP value
	Token    string // cloudflare api token to update the DNS record
}

// not passing the context to monitorAndUpdate() because once the check and
// update IP logic starts, we do not want to shutdown until the fn has returned
type config struct {
	ip       string
	filename string
	token    string
}

type rwSeekTrunc interface {
	io.ReadWriteSeeker
	Truncate(size int64) error
}
