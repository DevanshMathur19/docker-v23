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

	"github.com/DevanshMathur19/docker-v23/api/types"
	"github.com/DevanshMathur19/docker-v23/errdefs"
)

func TestDiskUsageError(t *testing.T) {
	client := &Client{
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	}
	_, err := client.DiskUsage(context.Background(), types.DiskUsageOptions{})
	if !errdefs.IsSystem(err) {
		t.Fatalf("expected a Server Error, got %[1]T: %[1]v", err)
	}
}

func TestDiskUsage(t *testing.T) {
	expectedURL := "/system/df"
	client := &Client{
		client: newMockClient(func(req *http.Request) (*http.Response, error) {
			if !strings.HasPrefix(req.URL.Path, expectedURL) {
				return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
			}

			du := types.DiskUsage{
				LayersSize: int64(100),
				Images:     nil,
				Containers: nil,
				Volumes:    nil,
			}

			b, err := json.Marshal(du)
			if err != nil {
				return nil, err
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(b)),
			}, nil
		}),
	}
	if _, err := client.DiskUsage(context.Background(), types.DiskUsageOptions{}); err != nil {
		t.Fatal(err)
	}
}
