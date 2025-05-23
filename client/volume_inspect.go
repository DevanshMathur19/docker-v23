package client // import "github.com/DevanshMathur19/docker-v23/client"

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/DevanshMathur19/docker-v23/api/types/volume"
)

// VolumeInspect returns the information about a specific volume in the docker host.
func (cli *Client) VolumeInspect(ctx context.Context, volumeID string) (volume.Volume, error) {
	vol, _, err := cli.VolumeInspectWithRaw(ctx, volumeID)
	return vol, err
}

// VolumeInspectWithRaw returns the information about a specific volume in the docker host and its raw representation
func (cli *Client) VolumeInspectWithRaw(ctx context.Context, volumeID string) (volume.Volume, []byte, error) {
	if volumeID == "" {
		return volume.Volume{}, nil, objectNotFoundError{object: "volume", id: volumeID}
	}

	var vol volume.Volume
	resp, err := cli.get(ctx, "/volumes/"+volumeID, nil, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return vol, nil, err
	}

	body, err := io.ReadAll(resp.body)
	if err != nil {
		return vol, nil, err
	}
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&vol)
	return vol, body, err
}
