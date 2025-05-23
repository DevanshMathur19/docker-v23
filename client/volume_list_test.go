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

	"github.com/DevanshMathur19/docker-v23/api/types/filters"
	"github.com/DevanshMathur19/docker-v23/api/types/volume"
	"github.com/DevanshMathur19/docker-v23/errdefs"
)

func TestVolumeListError(t *testing.T) {
	client := &Client{
		client: newMockClient(errorMock(http.StatusInternalServerError, "Server error")),
	}

	_, err := client.VolumeList(context.Background(), filters.NewArgs())
	if !errdefs.IsSystem(err) {
		t.Fatalf("expected a Server Error, got %[1]T: %[1]v", err)
	}
}

func TestVolumeList(t *testing.T) {
	expectedURL := "/volumes"

	noDanglingFilters := filters.NewArgs()
	noDanglingFilters.Add("dangling", "false")

	danglingFilters := filters.NewArgs()
	danglingFilters.Add("dangling", "true")

	labelFilters := filters.NewArgs()
	labelFilters.Add("label", "label1")
	labelFilters.Add("label", "label2")

	listCases := []struct {
		filters         filters.Args
		expectedFilters string
	}{
		{
			filters:         filters.NewArgs(),
			expectedFilters: "",
		}, {
			filters:         noDanglingFilters,
			expectedFilters: `{"dangling":{"false":true}}`,
		}, {
			filters:         danglingFilters,
			expectedFilters: `{"dangling":{"true":true}}`,
		}, {
			filters:         labelFilters,
			expectedFilters: `{"label":{"label1":true,"label2":true}}`,
		},
	}

	for _, listCase := range listCases {
		client := &Client{
			client: newMockClient(func(req *http.Request) (*http.Response, error) {
				if !strings.HasPrefix(req.URL.Path, expectedURL) {
					return nil, fmt.Errorf("Expected URL '%s', got '%s'", expectedURL, req.URL)
				}
				query := req.URL.Query()
				actualFilters := query.Get("filters")
				if actualFilters != listCase.expectedFilters {
					return nil, fmt.Errorf("filters not set in URL query properly. Expected '%s', got %s", listCase.expectedFilters, actualFilters)
				}
				content, err := json.Marshal(volume.ListResponse{
					Volumes: []*volume.Volume{
						{
							Name:   "volume",
							Driver: "local",
						},
					},
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

		volumeResponse, err := client.VolumeList(context.Background(), listCase.filters)
		if err != nil {
			t.Fatal(err)
		}
		if len(volumeResponse.Volumes) != 1 {
			t.Fatalf("expected 1 volume, got %v", volumeResponse.Volumes)
		}
	}
}
