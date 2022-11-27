package benchclient

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"
)

type Agent struct {
	httpClient *http.Client
}

func New(address string, timeout time.Duration) *Agent {
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, "unix", address)
			},
		},
		Timeout: timeout,
	}

	return &Agent{httpClient: httpClient}
}

func (a *Agent) Start(ctx context.Context, wg *sync.WaitGroup, report func(*http.Response, error, time.Time)) {
	go func() {
		defer wg.Done()
		for {
			if ctx.Err() != nil {
				return
			}

			beginning := time.Now()
			resp, err := a.httpClient.Get("http://localhost")
			report(resp, err, beginning)
		}
	}()
}
