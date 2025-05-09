package daemon

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/DevanshMathur19/docker-v23/api/types"
	"github.com/DevanshMathur19/docker-v23/api/types/swarm"
	"gotest.tools/v3/assert"
)

// NodeConstructor defines a swarm node constructor
type NodeConstructor func(*swarm.Node)

// GetNode returns a swarm node identified by the specified id
func (d *Daemon) GetNode(t testing.TB, id string, errCheck ...func(error) bool) *swarm.Node {
	t.Helper()
	cli := d.NewClientT(t)
	defer cli.Close()

	node, _, err := cli.NodeInspectWithRaw(context.Background(), id)
	if err != nil {
		for _, f := range errCheck {
			if f(err) {
				return nil
			}
		}
	}
	assert.NilError(t, err, "[%s] (*Daemon).GetNode: NodeInspectWithRaw(%q) failed", d.id, id)
	assert.Check(t, node.ID == id)
	return &node
}

// RemoveNode removes the specified node
func (d *Daemon) RemoveNode(t testing.TB, id string, force bool) {
	t.Helper()
	cli := d.NewClientT(t)
	defer cli.Close()

	options := types.NodeRemoveOptions{
		Force: force,
	}
	err := cli.NodeRemove(context.Background(), id, options)
	assert.NilError(t, err)
}

// UpdateNode updates a swarm node with the specified node constructor
func (d *Daemon) UpdateNode(t testing.TB, id string, f ...NodeConstructor) {
	t.Helper()
	cli := d.NewClientT(t)
	defer cli.Close()

	for i := 0; ; i++ {
		node := d.GetNode(t, id)
		for _, fn := range f {
			fn(node)
		}

		err := cli.NodeUpdate(context.Background(), node.ID, node.Version, node.Spec)
		if i < 10 && err != nil && strings.Contains(err.Error(), "update out of sequence") {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		assert.NilError(t, err)
		return
	}
}

// ListNodes returns the list of the current swarm nodes
func (d *Daemon) ListNodes(t testing.TB) []swarm.Node {
	t.Helper()
	cli := d.NewClientT(t)
	defer cli.Close()

	nodes, err := cli.NodeList(context.Background(), types.NodeListOptions{})
	assert.NilError(t, err)

	return nodes
}
