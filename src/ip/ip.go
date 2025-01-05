package ip

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/netip"
	"os"
	"regexp"
	"strings"

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
			err := handleIPCheck(&config{
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

// checks the new IP against the stored value. if it differs, it will call
// dns.UpdateDnsRecord()
func handleIPCheck(c *config) error {
	f, err := os.OpenFile(c.filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
		}
		return err
	}
	defer f.Close()

	if !ipHasChanged(c.ip, f) {
		slog.Debug("IP hasn't changed", "IP", c.ip)
		return nil
	}

	err = dns.UpdateRecord(c.ip, c.token)
	if err != nil {
		return err
	}

	err = updateIP(f, c.ip)
	if err != nil {
		return err
	}
	return nil
}

func ipHasChanged(newIp string, r io.Reader) bool {
	slog.Debug("ipHasChanged() called...")
	buf := make([]byte, 15)
	_, err := r.Read(buf)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			slog.Error("ipHasChanged(): error reading file", "error", err)
			return false
		}
	}

	c := strings.TrimSuffix(string(buf), "\n\x00")
	oldIp, err := netip.ParseAddr(c)
	if err != nil {
		slog.Error("ipHasChanged(): error parsing address", "error", err, "ip", string(buf))
	}
	return !strings.EqualFold(newIp, oldIp.String())
}

func ipRegexp(b []byte) []byte {
	// netip.ParseAddr(s string)
	r := regexp.MustCompile(`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`)
	return r.Find(b)
}

func updateIP(f *os.File, ip string) error {
	err := f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(ip))
	if err != nil {
		return err
	}
	return nil
}
