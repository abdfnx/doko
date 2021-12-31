package stream

import (
	"context"

	"github.com/docker/docker/api/types"
)

type ResizeContainer func(ctx context.Context, id string, options types.ResizeOptions) error
