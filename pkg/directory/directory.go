package directory // import "github.com/DevanshMathur19/docker-v23/pkg/directory"

import "context"

// Size walks a directory tree and returns its total size in bytes.
func Size(ctx context.Context, dir string) (int64, error) {
	return calcSize(ctx, dir)
}
