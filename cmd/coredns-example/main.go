package main

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"

	// We silently import the plugin here so that it is included in the binary
	_ "github.com/coredns/coredns/plugin/forward"

	// We import the plugin here so that it is included in the binary
	_ "github.com/tcrisseapp/coredns-example/adblocker"
)

var (
	// We need to tell CoreDNS which plugins are available and in which order to call them.
	directives = []string{
		"adblocker",
		"forward",
	}
)

func main() {
	// We pass the directives to CoreDNS.
	dnsserver.Directives = directives
	coremain.Run()
}
