package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

import (
	"github.com/DevanshMathur19/docker-v23/api/types"
	"github.com/DevanshMathur19/docker-v23/api/types/backend"
	"github.com/DevanshMathur19/docker-v23/api/types/versions/v1p19"
	"github.com/DevanshMathur19/docker-v23/container"
	"github.com/DevanshMathur19/docker-v23/daemon/exec"
)

// This sets platform-specific fields
func setPlatformSpecificContainerFields(container *container.Container, contJSONBase *types.ContainerJSONBase) *types.ContainerJSONBase {
	contJSONBase.AppArmorProfile = container.AppArmorProfile
	contJSONBase.ResolvConfPath = container.ResolvConfPath
	contJSONBase.HostnamePath = container.HostnamePath
	contJSONBase.HostsPath = container.HostsPath

	return contJSONBase
}

// containerInspectPre120 gets containers for pre 1.20 APIs.
func (daemon *Daemon) containerInspectPre120(name string) (*v1p19.ContainerJSON, error) {
	ctr, err := daemon.GetContainer(name)
	if err != nil {
		return nil, err
	}

	ctr.Lock()
	defer ctr.Unlock()

	base, err := daemon.getInspectData(ctr)
	if err != nil {
		return nil, err
	}

	volumes := make(map[string]string)
	volumesRW := make(map[string]bool)
	for _, m := range ctr.MountPoints {
		volumes[m.Destination] = m.Path()
		volumesRW[m.Destination] = m.RW
	}

	config := &v1p19.ContainerConfig{
		Config:          ctr.Config,
		MacAddress:      ctr.Config.MacAddress,
		NetworkDisabled: ctr.Config.NetworkDisabled,
		ExposedPorts:    ctr.Config.ExposedPorts,
		VolumeDriver:    ctr.HostConfig.VolumeDriver,
		Memory:          ctr.HostConfig.Memory,
		MemorySwap:      ctr.HostConfig.MemorySwap,
		CPUShares:       ctr.HostConfig.CPUShares,
		CPUSet:          ctr.HostConfig.CpusetCpus,
	}
	networkSettings := daemon.getBackwardsCompatibleNetworkSettings(ctr.NetworkSettings)

	return &v1p19.ContainerJSON{
		ContainerJSONBase: base,
		Volumes:           volumes,
		VolumesRW:         volumesRW,
		Config:            config,
		NetworkSettings:   networkSettings,
	}, nil
}

func inspectExecProcessConfig(e *exec.Config) *backend.ExecProcessConfig {
	return &backend.ExecProcessConfig{
		Tty:        e.Tty,
		Entrypoint: e.Entrypoint,
		Arguments:  e.Args,
		Privileged: &e.Privileged,
		User:       e.User,
	}
}
