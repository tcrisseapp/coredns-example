package adblocker

import (
	"bufio"
	"os"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() {
	// The registration part where we tell CoreDNS that if adblocker is inside the Corefile, then it should call the setup function
	plugin.Register("adblocker", setup)
}

type Plug struct {
	Next           plugin.Handler // The next plugin in the chain
	BlockListParam string
	BlockEntries   map[string]struct{}
}

func setup(c *caddy.Controller) error {
	// The setup function where we parse the Corefile plugin configuration and set up the plugin
	plug := &Plug{
		BlockEntries: map[string]struct{}{},
	}

	// Here we parse the configuration and set up the plugin
	for c.Next() {
		args := c.RemainingArgs()
		if len(args) > 0 {
			plug.BlockListParam = args[0]
		}
	}

	// Load the blocklist entries
	if err := plug.loadBlockEntries(); err != nil {
		return err
	}

	// This is where we tell CoreDNS what the next plugin is in the chain
	dnsConfig := dnsserver.GetConfig(c)
	dnsConfig.AddPlugin(func(next plugin.Handler) plugin.Handler {
		plug.Next = next
		return plug
	})

	return nil
}

func (p *Plug) loadBlockEntries() error {
	file, err := os.Open(p.BlockListParam)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		p.BlockEntries[strings.TrimSpace(scanner.Text())] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
