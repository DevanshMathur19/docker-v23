package client // import "github.com/DevanshMathur19/docker-v23/client"

import "context"

// ContainerPause pauses the main process of a given container without terminating it.
func (cli *Client) ContainerPause(ctx context.Context, containerID string) error {
	resp, err := cli.post(ctx, "/containers/"+containerID+"/pause", nil, nil, nil)
	ensureReaderClosed(resp)
	return err
}
