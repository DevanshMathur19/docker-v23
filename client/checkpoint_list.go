package client // import "github.com/DevanshMathur19/docker-v23/client"

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/DevanshMathur19/docker-v23/api/types"
)

// CheckpointList returns the checkpoints of the given container in the docker host
func (cli *Client) CheckpointList(ctx context.Context, container string, options types.CheckpointListOptions) ([]types.Checkpoint, error) {
	var checkpoints []types.Checkpoint

	query := url.Values{}
	if options.CheckpointDir != "" {
		query.Set("dir", options.CheckpointDir)
	}

	resp, err := cli.get(ctx, "/containers/"+container+"/checkpoints", query, nil)
	defer ensureReaderClosed(resp)
	if err != nil {
		return checkpoints, err
	}

	err = json.NewDecoder(resp.body).Decode(&checkpoints)
	return checkpoints, err
}
