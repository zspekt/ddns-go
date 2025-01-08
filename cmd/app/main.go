package main

import (
	"github.com/zspekt/ddns-go/src/run"
	"github.com/zspekt/ddns-go/src/setup"
)

func main() {
	run.Start(setup.Config())
}
