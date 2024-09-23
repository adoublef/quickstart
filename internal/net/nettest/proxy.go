package nettest

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/Shopify/toxiproxy/v2"
	"github.com/Shopify/toxiproxy/v2/toxics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type (
	// Defines how packets flow through a ToxicStub.
	Toxic = toxics.Toxic
	// The BandwidthToxic passes data through at a limited rate.
	BandwidthToxic = toxics.BandwidthToxic
	// The LatencyToxic passes data through with the a delay of latency +/- jitter added.
	LatencyToxic = toxics.LatencyToxic
	// LimitDataToxic has limit in bytes.
	LimitDataToxic = toxics.LimitDataToxic
	// The TimeoutToxic stops any data from flowing through, and will close the connection after a timeout. If the timeout is set to 0, then the connection will not be closed.
	TimeoutToxic = toxics.TimeoutToxic
	// The SlicerToxic slices data into multiple smaller packets to simulate real-world TCP behaviour.
	SlicerToxic = toxics.SlicerToxic
	// The SlowCloseToxic stops the TCP connection from closing until after a delay.
	SlowCloseToxic = toxics.SlowCloseToxic
)

// Proxy
type Proxy struct {
	p *toxiproxy.Proxy
	s *toxiproxy.ApiServer
}

// Listen
func (p *Proxy) Listen() string {
	if p.p == nil {
		return ""
	}
	return p.p.Listen
}

// AddToxic
func (p *Proxy) AddToxic(typ, stream string, toxic toxics.Toxic) error {
	b, err := json.Marshal(map[string]any{
		"name":       typ + "_" + stream,
		"type":       typ,
		"stream":     stream,
		"attributes": toxic,
	})
	if err != nil {
		return err
	}

	_, err = p.p.Toxics.AddToxicJson(bytes.NewReader(b))
	if err != nil {
		return err
	}
	return nil
}

// RemoveToxic
func (p *Proxy) RemoveToxic(name string) error {
	err := p.p.Toxics.RemoveToxic(context.Background(), name)
	if err != nil {
		return err
	}
	return nil
}

// Close
func (p *Proxy) Close() error {
	if p.p != nil {
		p.p.Stop()
	}
	if p.s != nil {
		return p.s.Shutdown()
	}
	return nil
}

// NewProxy returns a new [Proxy]
func NewProxy(name, upstream string) *Proxy {
	xm := toxiproxy.NewMetricsContainer(prometheus.NewRegistry())
	xl := zerolog.New(os.Stderr).Level(zerolog.ErrorLevel)
	s := toxiproxy.NewServer(xm, xl)
	proxy := toxiproxy.NewProxy(s, name, "localhost:0", upstream)
	proxy.Start()
	return &Proxy{p: proxy, s: s}
}
