package ip

import (
	"context"
	"io"

	"github.com/zspekt/ddns-go/src/dns"
)

type Config struct {
	Ctx      context.Context
	IpChan   <-chan string
	Filename string // file that holds the most recent IP value
	Api      dns.Api
}

// not passing the context to monitorAndUpdate() because once the check and
// update IP logic starts, we do not want to shutdown until the fn has returned
type config struct {
	ip       string
	filename string
	api      dns.Api
}

type rwSeekTrunc interface {
	io.ReadWriteSeeker
	Truncate(size int64) error
}
