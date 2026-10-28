package runtime

import (
	"context"

	"github.com/mygento/workplace/internal/config"
)

type Runtime interface {
	Start(ctx context.Context, cfg config.Config) error
	Stop(ctx context.Context, cfg config.Config) error
	Delete(ctx context.Context, cfg config.Config) error
}
