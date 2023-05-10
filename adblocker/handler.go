package adblocker

import (
	"context"
	"strings"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
)

// The Name function should return the name of the plugin
func (p *Plug) Name() string {
	return "adblocker"
}

// The ServeDNS function is where the magic happens. This is where we implement our plugin logic.
func (p *Plug) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	if len(r.Question) > 0 {
		q := r.Question[0]
		if _, ok := p.BlockEntries[strings.TrimSuffix(q.Name, ".")]; ok {
			clog.Infof("Blocking %s", q.Name)
			return dns.RcodeBadKey, nil
		}
		clog.Infof("Not blocking %s", q.Name)
	}

	return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
}
