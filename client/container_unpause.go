package client // import "github.com/DevanshMathur19/docker-v23/client"

import "context"

// ContainerUnpause resumes the process execution within a container
func (cli *Client) ContainerUnpause(ctx context.Context, containerID string) error {
	resp, err := cli.post(ctx, "/containers/"+containerID+"/unpause", nil, nil, nil)
	ensureReaderClosed(resp)
	return err
}
