package ip

import (
	"context"
	"log"
	"os"

	"github.com/zspekt/ddns-go/src/dns"
)

type Config struct {
	ctx      context.Context
	ipChan   <-chan string
	filename string // file that holds the most recent IP value
	token    string // cloudflare api token to update the DNS record
}

// not passing the context to monitorAndUpdate() because once the check and
// update IP logic starts, we do not want to shutdown until the fn has returned
type config struct {
	ip       string
	filename string
	token    string
}

// listens on channels c.ipChan for new IP value to hand off to monitorAndUpdate()
// and c.ctx.Done() to gracefully shut down
func MonitorAndUpdate(c *Config) {
	for {
		select {
		case ip := <-c.ipChan:
			err := monitorAndUpdate(&config{
				ip:       ip,
				filename: c.filename,
				token:    c.token,
			})
			if err != nil {
				// TODO: handle error
			}
		case <-c.ctx.Done():
			// TODO: shutdown logic
		}
	}
}

// compares given IP to the stored value, updating it if necessary. returns true
// only if given IP matches the value. if no value was stored (file didn't exist),
// or it differs from the one passed in, it returns false.
func compareAndStoreIP(ip string) bool { return false } // TODO: implement

func monitorAndUpdate(c *config) error {
	f, err := os.OpenFile(c.filename, os.O_RDWR, 0666)
	if err != nil {
		return err // TODO: decide what to do with errors
	}

	if !ipHasChanged(c.ip, f) {
		// if it hasn't changed we do nothing and return
		return nil
	}

	dns.UpdateDnsRecord()
}

func ipHasChanged(ip string, f *os.File) bool {
	f.Read(nil)
	// TODO: implement
}
