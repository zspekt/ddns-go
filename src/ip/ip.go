package ip

import (
	"errors"
	"io"
	"log"
	"log/slog"
	"net/netip"
	"os"
	"strings"
)

// listens on channels c.ipChan for new IP value to hand off to monitorAndUpdate()
// and c.ctx.Done() to gracefully shut down
func MonitorAndUpdate(c *Config) {
	for {
		select {
		case ip := <-c.IpChan:
			err := handleIPCheck(&config{
				ip:       ip,
				filename: c.Filename,
				api:      c.Api,
			})
			if err != nil {
				// TODO: handle error
				log.Fatal(err)
			}
		case <-c.Ctx.Done():
			// TODO: shutdown logic
		}
	}
}

// checks the new IP against the stored value. if it differs, it will call
// dns.UpdateDnsRecord()
func handleIPCheck(c *config) error {
	f, code, err := openOrCreate(c.filename)
	if err != nil {
		return err
	}

	if code == FILE_EXISTS && !ipHasChanged(c.ip, f) {
		slog.Debug("IP hasn't changed", "IP", c.ip)
		return nil
	}

	err = c.api.UpdateRecord(c.ip)
	if err != nil {
		return err
	}

	err = updateIP(f, c.ip)
	if err != nil {
		return err
	}
	return nil
}

const (
	FILE_EXISTS  int = 1
	FILE_CREATED int = 2
)

func openOrCreate(filename string) (f *os.File, code int, err error) {
	slog.Debug("openOrCreate(): called...")
	_, err = os.Stat(filename)
	flags := os.O_RDWR

	switch err.(type) {
	case nil:
		slog.Debug("openOrCreate(): file already exists")
		code = FILE_EXISTS
	case *os.PathError:
		slog.Debug("openOrCreate(): file doesn't exist")
		flags |= os.O_CREATE
		code = FILE_CREATED
	default:
		slog.Error("openOrCreate(): unexpected error. returning...", "error", err)
		return nil, 0, err
	}
	f, err = os.OpenFile(filename, flags, 0666)
	if err != nil {
		return nil, 0, err
	}
	return f, code, nil
}

func ipHasChanged(newIp string, r io.ReadSeeker) bool {
	slog.Debug("ipHasChanged(): called...")
	buf := make([]byte, 64)

	_, err := r.Seek(0, 0)
	if err != nil {
		return true
	}

	n, err := r.Read(buf)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			slog.Error("ipHasChanged(): error reading from reader", "error", err)
			return true
		}
		slog.Info("ipHasChanged(): reader is empty")
		return true
	}
	buf = buf[:n]

	parsed := strings.ReplaceAll(strings.TrimSpace(string(buf)), "\n", "")

	oldIp, err := netip.ParseAddr(parsed)
	if err != nil {
		slog.Error("ipHasChanged(): error parsing address", "error", err, "ip", string(buf))
	}
	return !strings.EqualFold(newIp, oldIp.String())
}

func updateIP(f rwSeekTrunc, ip string) error {
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
