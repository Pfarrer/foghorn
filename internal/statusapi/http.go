package statusapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pfarrer/foghorn/scheduler"
)

const (
	StatusPath         = "/v1/status"
	DefaultListenAddr  = "127.0.0.1:7676"
	DefaultBaseURL     = "http://127.0.0.1:7676"
	defaultReadTimeout = 2 * time.Second
)

func NewHandler(snapshotFn func() scheduler.Snapshot) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(StatusPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(snapshotFn()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}

func StartServer(addr string, snapshotFn func() scheduler.Snapshot) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           NewHandler(snapshotFn),
		ReadHeaderTimeout: defaultReadTimeout,
	}
}

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (c *Client) GetStatus(ctx context.Context) (scheduler.Snapshot, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+StatusPath, nil)
	if err != nil {
		return scheduler.Snapshot{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return scheduler.Snapshot{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return scheduler.Snapshot{}, fmt.Errorf("status endpoint returned %s", resp.Status)
	}

	var snapshot scheduler.Snapshot
	if err := json.NewDecoder(resp.Body).Decode(&snapshot); err != nil {
		return scheduler.Snapshot{}, err
	}
	return snapshot, nil
}
