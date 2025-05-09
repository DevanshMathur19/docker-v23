package images // import "github.com/DevanshMathur19/docker-v23/daemon/images"

import (
	"context"
	"fmt"
	"os"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/leases"
	"github.com/DevanshMathur19/docker-v23/api/types"
	"github.com/DevanshMathur19/docker-v23/api/types/filters"
	"github.com/DevanshMathur19/docker-v23/container"
	daemonevents "github.com/DevanshMathur19/docker-v23/daemon/events"
	"github.com/DevanshMathur19/docker-v23/distribution/metadata"
	"github.com/DevanshMathur19/docker-v23/distribution/xfer"
	"github.com/DevanshMathur19/docker-v23/image"
	"github.com/DevanshMathur19/docker-v23/layer"
	dockerreference "github.com/DevanshMathur19/docker-v23/reference"
	"github.com/DevanshMathur19/docker-v23/registry"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

type containerStore interface {
	// used by image delete
	First(container.StoreFilter) *container.Container
	// used by image prune, and image list
	List() []*container.Container
	// TODO: remove, only used for CommitBuildStep
	Get(string) *container.Container
}

// ImageServiceConfig is the configuration used to create a new ImageService
type ImageServiceConfig struct {
	ContainerStore            containerStore
	DistributionMetadataStore metadata.Store
	EventsService             *daemonevents.Events
	ImageStore                image.Store
	LayerStore                layer.Store
	MaxConcurrentDownloads    int
	MaxConcurrentUploads      int
	MaxDownloadAttempts       int
	ReferenceStore            dockerreference.Store
	RegistryService           registry.Service
	ContentStore              content.Store
	Leases                    leases.Manager
	ContentNamespace          string
}

// NewImageService returns a new ImageService from a configuration
func NewImageService(config ImageServiceConfig) *ImageService {
	return &ImageService{
		containers:                config.ContainerStore,
		distributionMetadataStore: config.DistributionMetadataStore,
		downloadManager:           xfer.NewLayerDownloadManager(config.LayerStore, config.MaxConcurrentDownloads, xfer.WithMaxDownloadAttempts(config.MaxDownloadAttempts)),
		eventsService:             config.EventsService,
		imageStore:                &imageStoreWithLease{Store: config.ImageStore, leases: config.Leases, ns: config.ContentNamespace},
		layerStore:                config.LayerStore,
		referenceStore:            config.ReferenceStore,
		registryService:           config.RegistryService,
		uploadManager:             xfer.NewLayerUploadManager(config.MaxConcurrentUploads),
		leases:                    config.Leases,
		content:                   config.ContentStore,
		contentNamespace:          config.ContentNamespace,
	}
}

// ImageService provides a backend for image management
type ImageService struct {
	containers                containerStore
	distributionMetadataStore metadata.Store
	downloadManager           *xfer.LayerDownloadManager
	eventsService             *daemonevents.Events
	imageStore                image.Store
	layerStore                layer.Store
	pruneRunning              int32
	referenceStore            dockerreference.Store
	registryService           registry.Service
	uploadManager             *xfer.LayerUploadManager
	leases                    leases.Manager
	content                   content.Store
	contentNamespace          string
	usage                     singleflight.Group
}

// DistributionServices provides daemon image storage services
type DistributionServices struct {
	DownloadManager   *xfer.LayerDownloadManager
	V2MetadataService metadata.V2MetadataService
	LayerStore        layer.Store
	ImageStore        image.Store
	ReferenceStore    dockerreference.Store
}

// DistributionServices return services controlling daemon image storage
func (i *ImageService) DistributionServices() DistributionServices {
	return DistributionServices{
		DownloadManager:   i.downloadManager,
		V2MetadataService: metadata.NewV2MetadataService(i.distributionMetadataStore),
		LayerStore:        i.layerStore,
		ImageStore:        i.imageStore,
		ReferenceStore:    i.referenceStore,
	}
}

// CountImages returns the number of images stored by ImageService
// called from info.go
func (i *ImageService) CountImages() int {
	return i.imageStore.Len()
}

// Children returns the children image.IDs for a parent image.
// called from list.go to filter containers
// TODO: refactor to expose an ancestry for image.ID?
func (i *ImageService) Children(id image.ID) []image.ID {
	return i.imageStore.Children(id)
}

// CreateLayer creates a filesystem layer for a container.
// called from create.go
// TODO: accept an opt struct instead of container?
func (i *ImageService) CreateLayer(container *container.Container, initFunc layer.MountInit) (layer.RWLayer, error) {
	var layerID layer.ChainID
	if container.ImageID != "" {
		img, err := i.imageStore.Get(container.ImageID)
		if err != nil {
			return nil, err
		}
		layerID = img.RootFS.ChainID()
	}

	rwLayerOpts := &layer.CreateRWLayerOpts{
		MountLabel: container.MountLabel,
		InitFunc:   initFunc,
		StorageOpt: container.HostConfig.StorageOpt,
	}

	return i.layerStore.CreateRWLayer(container.ID, layerID, rwLayerOpts)
}

// GetLayerByID returns a layer by ID
// called from daemon.go Daemon.restore(), and Daemon.containerExport()
func (i *ImageService) GetLayerByID(cid string) (layer.RWLayer, error) {
	return i.layerStore.GetRWLayer(cid)
}

// LayerStoreStatus returns the status for each layer store
// called from info.go
func (i *ImageService) LayerStoreStatus() [][2]string {
	return i.layerStore.DriverStatus()
}

// GetLayerMountID returns the mount ID for a layer
// called from daemon.go Daemon.Shutdown(), and Daemon.Cleanup() (cleanup is actually continerCleanup)
// TODO: needs to be refactored to Unmount (see callers), or removed and replaced with GetLayerByID
func (i *ImageService) GetLayerMountID(cid string) (string, error) {
	return i.layerStore.GetMountID(cid)
}

// Cleanup resources before the process is shutdown.
// called from daemon.go Daemon.Shutdown()
func (i *ImageService) Cleanup() error {
	if err := i.layerStore.Cleanup(); err != nil {
		return errors.Wrap(err, "error during layerStore.Cleanup()")
	}
	return nil
}

// GraphDriverName returns the name of the graph drvier
// moved from Daemon.GraphDriverName, used by:
// - newContainer
// - to report an error in Daemon.Mount(container)
func (i *ImageService) GraphDriverName() string {
	return i.layerStore.DriverName()
}

// ReleaseLayer releases a layer allowing it to be removed
// called from delete.go Daemon.cleanupContainer(), and Daemon.containerExport()
func (i *ImageService) ReleaseLayer(rwlayer layer.RWLayer) error {
	metadata, err := i.layerStore.ReleaseRWLayer(rwlayer)
	layer.LogReleaseMetadata(metadata)
	if err != nil && !errors.Is(err, layer.ErrMountDoesNotExist) && !errors.Is(err, os.ErrNotExist) {
		return errors.Wrapf(err, "driver %q failed to remove root filesystem",
			i.layerStore.DriverName())
	}
	return nil
}

// LayerDiskUsage returns the number of bytes used by layer stores
// called from disk_usage.go
func (i *ImageService) LayerDiskUsage(ctx context.Context) (int64, error) {
	ch := i.usage.DoChan("LayerDiskUsage", func() (interface{}, error) {
		var allLayersSize int64
		layerRefs := i.getLayerRefs()
		allLayers := i.layerStore.Map()
		for _, l := range allLayers {
			select {
			case <-ctx.Done():
				return allLayersSize, ctx.Err()
			default:
				size := l.DiffSize()
				if _, ok := layerRefs[l.ChainID()]; ok {
					allLayersSize += size
				}
			}
		}
		return allLayersSize, nil
	})
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case res := <-ch:
		if res.Err != nil {
			return 0, res.Err
		}
		return res.Val.(int64), nil
	}
}

func (i *ImageService) getLayerRefs() map[layer.ChainID]int {
	tmpImages := i.imageStore.Map()
	layerRefs := map[layer.ChainID]int{}
	for id, img := range tmpImages {
		dgst := digest.Digest(id)
		if len(i.referenceStore.References(dgst)) == 0 && len(i.imageStore.Children(id)) != 0 {
			continue
		}

		rootFS := *img.RootFS
		rootFS.DiffIDs = nil
		for _, id := range img.RootFS.DiffIDs {
			rootFS.Append(id)
			chid := rootFS.ChainID()
			layerRefs[chid]++
		}
	}

	return layerRefs
}

// ImageDiskUsage returns information about image data disk usage.
func (i *ImageService) ImageDiskUsage(ctx context.Context) ([]*types.ImageSummary, error) {
	ch := i.usage.DoChan("ImageDiskUsage", func() (interface{}, error) {
		// Get all top images with extra attributes
		images, err := i.Images(ctx, types.ImageListOptions{
			Filters:        filters.NewArgs(),
			SharedSize:     true,
			ContainerCount: true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve image list: %v", err)
		}
		return images, nil
	})
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-ch:
		if res.Err != nil {
			return nil, res.Err
		}
		return res.Val.([]*types.ImageSummary), nil
	}
}

// UpdateConfig values
//
// called from reload.go
func (i *ImageService) UpdateConfig(maxDownloads, maxUploads int) {
	if i.downloadManager != nil && maxDownloads != 0 {
		i.downloadManager.SetConcurrency(maxDownloads)
	}
	if i.uploadManager != nil && maxUploads != 0 {
		i.uploadManager.SetConcurrency(maxUploads)
	}
}
