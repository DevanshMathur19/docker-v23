package client // import "github.com/DevanshMathur19/docker-v23/client"

import (
	"context"
	"net/url"

	"github.com/DevanshMathur19/docker-v23/api/types"
)

// CheckpointDelete deletes the checkpoint with the given name from the given container
func (cli *Client) CheckpointDelete(ctx context.Context, containerID string, options types.CheckpointDeleteOptions) error {
	query := url.Values{}
	if options.CheckpointDir != "" {
		query.Set("dir", options.CheckpointDir)
	}

	resp, err := cli.delete(ctx, "/containers/"+containerID+"/checkpoints/"+options.CheckpointID, query, nil)
	ensureReaderClosed(resp)
	return err
}
