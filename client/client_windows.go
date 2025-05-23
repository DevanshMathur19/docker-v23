package client // import "github.com/DevanshMathur19/docker-v23/client"

// DefaultDockerHost defines OS-specific default host if the DOCKER_HOST
// (EnvOverrideHost) environment variable is unset or empty.
const DefaultDockerHost = "npipe:////./pipe/docker_engine"

const defaultProto = "npipe"
const defaultAddr = "//./pipe/docker_engine"
