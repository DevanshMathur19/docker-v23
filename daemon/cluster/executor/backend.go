package executor // import "github.com/DevanshMathur19/docker-v23/daemon/cluster/executor"

import (
	"context"
	"io"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/DevanshMathur19/docker-v23/api/types"
	"github.com/DevanshMathur19/docker-v23/api/types/backend"
	"github.com/DevanshMathur19/docker-v23/api/types/container"
	"github.com/DevanshMathur19/docker-v23/api/types/events"
	"github.com/DevanshMathur19/docker-v23/api/types/filters"
	"github.com/DevanshMathur19/docker-v23/api/types/network"
	swarmtypes "github.com/DevanshMathur19/docker-v23/api/types/swarm"
	"github.com/DevanshMathur19/docker-v23/api/types/volume"
	containerpkg "github.com/DevanshMathur19/docker-v23/container"
	clustertypes "github.com/DevanshMathur19/docker-v23/daemon/cluster/provider"
	networkSettings "github.com/DevanshMathur19/docker-v23/daemon/network"
	"github.com/DevanshMathur19/docker-v23/image"
	"github.com/DevanshMathur19/docker-v23/libnetwork"
	"github.com/DevanshMathur19/docker-v23/libnetwork/cluster"
	networktypes "github.com/DevanshMathur19/docker-v23/libnetwork/types"
	"github.com/DevanshMathur19/docker-v23/plugin"
	volumeopts "github.com/DevanshMathur19/docker-v23/volume/service/opts"
	"github.com/moby/swarmkit/v2/agent/exec"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// Backend defines the executor component for a swarm agent.
type Backend interface {
	CreateManagedNetwork(clustertypes.NetworkCreateRequest) error
	DeleteManagedNetwork(networkID string) error
	FindNetwork(idName string) (libnetwork.Network, error)
	SetupIngress(clustertypes.NetworkCreateRequest, string) (<-chan struct{}, error)
	ReleaseIngress() (<-chan struct{}, error)
	CreateManagedContainer(config types.ContainerCreateConfig) (container.CreateResponse, error)
	ContainerStart(name string, hostConfig *container.HostConfig, checkpoint string, checkpointDir string) error
	ContainerStop(ctx context.Context, name string, config container.StopOptions) error
	ContainerLogs(ctx context.Context, name string, config *types.ContainerLogsOptions) (msgs <-chan *backend.LogMessage, tty bool, err error)
	ConnectContainerToNetwork(containerName, networkName string, endpointConfig *network.EndpointSettings) error
	ActivateContainerServiceBinding(containerName string) error
	DeactivateContainerServiceBinding(containerName string) error
	UpdateContainerServiceConfig(containerName string, serviceConfig *clustertypes.ServiceConfig) error
	ContainerInspectCurrent(name string, size bool) (*types.ContainerJSON, error)
	ContainerWait(ctx context.Context, name string, condition containerpkg.WaitCondition) (<-chan containerpkg.StateStatus, error)
	ContainerRm(name string, config *types.ContainerRmConfig) error
	ContainerKill(name string, sig string) error
	SetContainerDependencyStore(name string, store exec.DependencyGetter) error
	SetContainerSecretReferences(name string, refs []*swarmtypes.SecretReference) error
	SetContainerConfigReferences(name string, refs []*swarmtypes.ConfigReference) error
	SystemInfo() *types.Info
	Containers(config *types.ContainerListOptions) ([]*types.Container, error)
	SetNetworkBootstrapKeys([]*networktypes.EncryptionKey) error
	DaemonJoinsCluster(provider cluster.Provider)
	DaemonLeavesCluster()
	IsSwarmCompatible() error
	SubscribeToEvents(since, until time.Time, filter filters.Args) ([]events.Message, chan interface{})
	UnsubscribeFromEvents(listener chan interface{})
	UpdateAttachment(string, string, string, *network.NetworkingConfig) error
	WaitForDetachment(context.Context, string, string, string, string) error
	PluginManager() *plugin.Manager
	PluginGetter() *plugin.Store
	GetAttachmentStore() *networkSettings.AttachmentStore
	HasExperimental() bool
}

// VolumeBackend is used by an executor to perform volume operations
type VolumeBackend interface {
	Create(ctx context.Context, name, driverName string, opts ...volumeopts.CreateOption) (*volume.Volume, error)
}

// ImageBackend is used by an executor to perform image operations
type ImageBackend interface {
	PullImage(ctx context.Context, image, tag string, platform *specs.Platform, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error
	GetRepository(context.Context, reference.Named, *types.AuthConfig) (distribution.Repository, error)
	GetImage(refOrID string, platform *specs.Platform) (retImg *image.Image, retErr error)
}
