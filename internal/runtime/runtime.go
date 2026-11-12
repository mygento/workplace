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

// LogsRuntime is an optional capability for runtimes
// that support streaming service logs.
type LogsRuntime interface {
	Logs(
		ctx context.Context,
		cfg config.Config,
		services []string,
		follow bool,
		tail string,
		timestamps bool,
	) error
}
