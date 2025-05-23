package client // import "github.com/DevanshMathur19/docker-v23/client"

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/DevanshMathur19/docker-v23/api/types/volume"
	"github.com/DevanshMathur19/docker-v23/errdefs"
)

func TestVolumeCreateError(t *testing.T) {
	client := &Client{
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	}

	_, err := client.VolumeCreate(context.Background(), volume.CreateOptions{})
	if !errdefs.IsSystem(err) {
		t.Fatalf("expected a Server Error, got %[1]T: %[1]v", err)
	}
}

func TestVolumeCreate(t *testing.T) {
	expectedURL := "/volumes/create"

	client := &Client{
		client: newMockClient(func(req *http.Request) (*http.Response, error) {
			if !strings.HasPrefix(req.URL.Path, expectedURL) {
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			}

			if req.Method != http.MethodPost {
				return nil, fmt.Errorf("expected POST method, got %s", req.Method)
			}

			content, err := json.Marshal(volume.Volume{
				Name:       "volume",
				Driver:     "local",
				Mountpoint: "mountpoint",
			})
			if err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(content)),
			}, nil
		}),
	}

	vol, err := client.VolumeCreate(context.Background(), volume.CreateOptions{
		Name:   "myvolume",
		Driver: "mydriver",
		DriverOpts: map[string]string{
			"opt-key": "opt-value",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if vol.Name != "volume" {
		t.Fatalf("expected volume.Name to be 'volume', got %s", vol.Name)
	}
	if vol.Driver != "local" {
		t.Fatalf("expected volume.Driver to be 'local', got %s", vol.Driver)
	}
	if vol.Mountpoint != "mountpoint" {
		t.Fatalf("expected volume.Mountpoint to be 'mountpoint', got %s", vol.Mountpoint)
	}
}
