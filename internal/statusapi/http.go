package statusapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pfarrer/foghorn/scheduler"
	"github.com/pfarrer/foghorn/state"
)

const (
	StatusPath         = "/v1/status"
	HistoryPath        = "/v1/history"
	DefaultListenAddr  = "127.0.0.1:7676"
	DefaultBaseURL     = "http://127.0.0.1:7676"
	defaultReadTimeout = 2 * time.Second
)

type HistoryQuerier interface {
	Query(checkName string, start, end time.Time) ([]state.Record, error)
}

func NewHandler(snapshotFn func() scheduler.Snapshot, historyQuerier HistoryQuerier) http.Handler {
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
	mux.HandleFunc(HistoryPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if historyQuerier == nil {
			http.Error(w, "history not available", http.StatusServiceUnavailable)
			return
		}

		checkName := r.URL.Query().Get("check")
		var start, end time.Time
		if s := r.URL.Query().Get("start"); s != "" {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				http.Error(w, "invalid start time format (use RFC3339)", http.StatusBadRequest)
				return
			}
			start = t
		}
		if e := r.URL.Query().Get("end"); e != "" {
			t, err := time.Parse(time.RFC3339, e)
			if err != nil {
				http.Error(w, "invalid end time format (use RFC3339)", http.StatusBadRequest)
				return
			}
			end = t
		}

		records, err := historyQuerier.Query(checkName, start, end)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(records); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}

func StartServer(addr string, snapshotFn func() scheduler.Snapshot, historyQuerier HistoryQuerier) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           NewHandler(snapshotFn, historyQuerier),
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
