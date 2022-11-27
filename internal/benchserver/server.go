package benchserver

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/troydai/http-over-uds/internal/payload"
)

var (
	_payload = payload.MustGenRandomBytes(512)
)

const (
	_socketPath = "socket/server.socks"
)

func New() (*http.Server, error) {
	if _, err := os.Stat(_socketPath); os.IsExist(err) {
		if err := os.Remove(_socketPath); err != nil {
			return nil, fmt.Errorf("fail to remove pre-exist socket path")
		}
	}

	listener, err := net.Listen("unix", _socketPath)
	if err != nil {
		return nil, fmt.Errorf("fail to create TCP listner: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(_payload)
	})

	server := http.Server{
		Handler: mux,
	}

	go func() {
		fmt.Fprintln(os.Stderr, "start serving traffic at", _socketPath)
		server.Serve(listener)
	}()

	return &server, nil
}
