package run

import (
	"github.com/zspekt/ddns-go/pkg/utils"
	"github.com/zspekt/ddns-go/src/ip"
	"github.com/zspekt/ddns-go/src/listener"
	"github.com/zspekt/ddns-go/src/setup"
)

func Start(c *setup.Cfg) {
	go utils.Shutdown(c.Sigs, c.Cancel)
	go listener.Listen(c.Listener, c.Ch, c.Ctx)

	ip.MonitorAndUpdate(&ip.Config{
		Ctx:      c.Ctx,
		IpChan:   c.Ch,
		Filename: c.Filename,
		Api:      c.Api,
	})
}
