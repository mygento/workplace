package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/mygento/workplace/internal/nginx"
	"github.com/mygento/workplace/internal/php"
)

type Options struct {
	Host string
	Port int
}

func Start() error {
	projectDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf(
			"get project directory: %w",
			err,
		)
	}

	runtimeDir := filepath.Join(
		projectDir,
		".workplace",
	)

	if err := os.MkdirAll(
		filepath.Join(
			runtimeDir,
			"php",
		),
		0755,
	); err != nil {
		return fmt.Errorf(
			"create runtime directory: %w",
			err,
		)
	}

	// ---------------------------------------------------------
	// PHP-FPM
	// ---------------------------------------------------------

	fmt.Println("Starting PHP-FPM...")

	fpm, err := php.Start(filepath.Join(
		runtimeDir,
		"php",
	))
	if err != nil {
		return err
	}

	defer func() {
		if err := fpm.Stop(); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"failed to stop PHP-FPM: %v\n",
				err,
			)
		}
	}()

	// ---------------------------------------------------------
	// Nginx
	// ---------------------------------------------------------

	fmt.Println("Starting Nginx...")
	if err := checkPort(
		"127.0.0.1",
		8000,
	); err != nil {
		return err
	}

	nginxServer, err := nginx.Start(
		nginx.Config{
			RuntimeDir: filepath.Join(
				runtimeDir,
				"nginx",
			),
			Root:      projectDir,
			Host:      "127.0.0.1",
			Port:      8000,
			FPMSocket: fpm.Socket(),
		},
	)

	if err != nil {
		return err
	}

	defer func() {
		if err := nginxServer.Stop(); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"failed to stop Nginx: %v\n",
				err,
			)
		}
	}()

	// ---------------------------------------------------------
	// Ready
	// ---------------------------------------------------------

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println(" Development server started")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop.")
	fmt.Println()

	// ---------------------------------------------------------
	// Wait for Ctrl+C
	// ---------------------------------------------------------

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	<-ctx.Done()

	fmt.Println()
	fmt.Println("Stopping development environment...")

	return nil
}

func checkPort(host string, port int) error {
	address := fmt.Sprintf(
		"%s:%d",
		host,
		port,
	)

	listener, err := net.Listen(
		"tcp",
		address,
	)

	if err != nil {
		return fmt.Errorf(
			"cannot listen on %s: %w",
			address,
			err,
		)
	}

	return listener.Close()
}
