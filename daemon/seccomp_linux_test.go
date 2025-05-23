package daemon // import "github.com/DevanshMathur19/docker-v23/daemon"

import (
	"testing"

	coci "github.com/containerd/containerd/oci"
	containertypes "github.com/DevanshMathur19/docker-v23/api/types/container"
	"github.com/DevanshMathur19/docker-v23/container"
	dconfig "github.com/DevanshMathur19/docker-v23/daemon/config"
	"github.com/DevanshMathur19/docker-v23/oci"
	"github.com/DevanshMathur19/docker-v23/pkg/sysinfo"
	"github.com/DevanshMathur19/docker-v23/profiles/seccomp"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"gotest.tools/v3/assert"
)

func TestWithSeccomp(t *testing.T) {
	type expected struct {
		daemon  *Daemon
		c       *container.Container
		inSpec  coci.Spec
		outSpec coci.Spec
		err     string
		comment string
	}

	for _, x := range []expected{
		{
			comment: "unconfined seccompProfile runs unconfined",
			daemon: &Daemon{
				sysInfo: &sysinfo.SysInfo{Seccomp: true},
			},
			c: &container.Container{
				SeccompProfile: dconfig.SeccompProfileUnconfined,
				HostConfig: &containertypes.HostConfig{
					Privileged: false,
				},
			},
			inSpec:  oci.DefaultLinuxSpec(),
			outSpec: oci.DefaultLinuxSpec(),
		},
		{
			comment: "privileged container w/ custom profile runs unconfined",
			daemon: &Daemon{
				sysInfo: &sysinfo.SysInfo{Seccomp: true},
			},
			c: &container.Container{
				SeccompProfile: "{ \"defaultAction\": \"SCMP_ACT_LOG\" }",
				HostConfig: &containertypes.HostConfig{
					Privileged: true,
				},
			},
			inSpec:  oci.DefaultLinuxSpec(),
			outSpec: oci.DefaultLinuxSpec(),
		},
		{
			comment: "privileged container w/ default runs unconfined",
			daemon: &Daemon{
				sysInfo: &sysinfo.SysInfo{Seccomp: true},
			},
			c: &container.Container{
				SeccompProfile: "",
				HostConfig: &containertypes.HostConfig{
					Privileged: true,
				},
			},
			inSpec:  oci.DefaultLinuxSpec(),
			outSpec: oci.DefaultLinuxSpec(),
		},
		{
			comment: "privileged container w/ daemon profile runs unconfined",
			daemon: &Daemon{
				sysInfo:        &sysinfo.SysInfo{Seccomp: true},
				seccompProfile: []byte("{ \"defaultAction\": \"SCMP_ACT_ERRNO\" }"),
			},
			c: &container.Container{
				SeccompProfile: "",
				HostConfig: &containertypes.HostConfig{
					Privileged: true,
				},
			},
			inSpec:  oci.DefaultLinuxSpec(),
			outSpec: oci.DefaultLinuxSpec(),
		},
		{
			comment: "custom profile when seccomp is disabled returns error",
			daemon: &Daemon{
				sysInfo: &sysinfo.SysInfo{Seccomp: false},
			},
			c: &container.Container{
				SeccompProfile: "{ \"defaultAction\": \"SCMP_ACT_ERRNO\" }",
				HostConfig: &containertypes.HostConfig{
					Privileged: false,
				},
			},
			inSpec:  oci.DefaultLinuxSpec(),
			outSpec: oci.DefaultLinuxSpec(),
			err:     "seccomp is not enabled in your kernel, cannot run a custom seccomp profile",
		},
		{
			comment: "empty profile name loads default profile",
			daemon: &Daemon{
				sysInfo: &sysinfo.SysInfo{Seccomp: true},
			},
			c: &container.Container{
				SeccompProfile: "",
				HostConfig: &containertypes.HostConfig{
					Privileged: false,
				},
			},
			inSpec: oci.DefaultLinuxSpec(),
			outSpec: func() coci.Spec {
				s := oci.DefaultLinuxSpec()
				profile, _ := seccomp.GetDefaultProfile(&s)
				s.Linux.Seccomp = profile
				return s
			}(),
		},
		{
			comment: "load container's profile",
			daemon: &Daemon{
				sysInfo: &sysinfo.SysInfo{Seccomp: true},
			},
			c: &container.Container{
				SeccompProfile: "{ \"defaultAction\": \"SCMP_ACT_ERRNO\" }",
				HostConfig: &containertypes.HostConfig{
					Privileged: false,
				},
			},
			inSpec: oci.DefaultLinuxSpec(),
			outSpec: func() coci.Spec {
				s := oci.DefaultLinuxSpec()
				profile := &specs.LinuxSeccomp{
					DefaultAction: specs.LinuxSeccompAction("SCMP_ACT_ERRNO"),
				}
				s.Linux.Seccomp = profile
				return s
			}(),
		},
		{
			comment: "load daemon's profile",
			daemon: &Daemon{
				sysInfo:        &sysinfo.SysInfo{Seccomp: true},
				seccompProfile: []byte("{ \"defaultAction\": \"SCMP_ACT_ERRNO\" }"),
			},
			c: &container.Container{
				SeccompProfile: "",
				HostConfig: &containertypes.HostConfig{
					Privileged: false,
				},
			},
			inSpec: oci.DefaultLinuxSpec(),
			outSpec: func() coci.Spec {
				s := oci.DefaultLinuxSpec()
				profile := &specs.LinuxSeccomp{
					DefaultAction: specs.LinuxSeccompAction("SCMP_ACT_ERRNO"),
				}
				s.Linux.Seccomp = profile
				return s
			}(),
		},
		{
			comment: "load prioritise container profile over daemon's",
			daemon: &Daemon{
				sysInfo:        &sysinfo.SysInfo{Seccomp: true},
				seccompProfile: []byte("{ \"defaultAction\": \"SCMP_ACT_ERRNO\" }"),
			},
			c: &container.Container{
				SeccompProfile: "{ \"defaultAction\": \"SCMP_ACT_LOG\" }",
				HostConfig: &containertypes.HostConfig{
					Privileged: false,
				},
			},
			inSpec: oci.DefaultLinuxSpec(),
			outSpec: func() coci.Spec {
				s := oci.DefaultLinuxSpec()
				profile := &specs.LinuxSeccomp{
					DefaultAction: specs.LinuxSeccompAction("SCMP_ACT_LOG"),
				}
				s.Linux.Seccomp = profile
				return s
			}(),
		},
	} {
		x := x
		t.Run(x.comment, func(t *testing.T) {
			opts := WithSeccomp(x.daemon, x.c)
			err := opts(nil, nil, nil, &x.inSpec)

			assert.DeepEqual(t, x.inSpec, x.outSpec)
			if x.err != "" {
				assert.Error(t, err, x.err)
			} else {
				assert.NilError(t, err)
			}
		})
	}
}
