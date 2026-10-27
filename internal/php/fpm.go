package php

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

type FPM struct {
	cmd        *exec.Cmd
	socketPath string
}

func Start(runtimeDir string) (*FPM, error) {
	if err := os.MkdirAll(runtimeDir, 0755); err != nil {
		return nil, fmt.Errorf("create PHP runtime directory: %w", err)
	}

	socketPath := filepath.Join(runtimeDir, "php-fpm.sock")

	// Remove stale socket.
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("remove stale PHP-FPM socket: %w", err)
	}

	configPath, err := WriteConfig(Config{
		RuntimeDir: runtimeDir,
		SocketPath: socketPath,
	})
	if err != nil {
		return nil, err
	}

	binary, err := findBinary()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(
		binary,
		"--nodaemonize",
		"--fpm-config",
		configPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start PHP-FPM: %w", err)
	}

	fpm := &FPM{
		cmd:        cmd,
		socketPath: socketPath,
	}

	if err := fpm.waitReady(10 * time.Second); err != nil {
		if stopErr := fpm.Stop(); stopErr != nil {
			err = errors.Join(err, fmt.Errorf("stop PHP-FPM after startup failure: %w", stopErr))
		}

		return nil, fmt.Errorf("PHP-FPM startup failed: %w", err)
	}

	fmt.Printf("PHP-FPM started (pid %d)\n", cmd.Process.Pid)

	return fpm, nil
}

func findBinary() (string, error) {
	names := []string{"php-fpm"}

	for _, name := range names {
		path, err := exec.LookPath(name)
		if err == nil {
			return path, nil
		}
	}

	return "", errors.New(
		"PHP-FPM not found; install PHP-FPM and make sure it is available in PATH",
	)
}

func (f *FPM) Socket() string {
	return f.socketPath
}

func (f *FPM) waitReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		// Did FPM exit?
		if f.cmd.ProcessState != nil && f.cmd.ProcessState.Exited() {
			return errors.New("PHP-FPM exited before becoming ready")
		}

		// Does the socket exist?
		if _, err := os.Stat(f.socketPath); err == nil {
			// Can we connect?
			conn, err := net.DialTimeout(
				"unix",
				f.socketPath,
				250*time.Millisecond,
			)

			if err == nil {
				if closeErr := conn.Close(); closeErr != nil {
					return fmt.Errorf("close PHP-FPM readiness connection: %w", closeErr)
				}
				return nil
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("PHP-FPM socket did not become ready: %s", f.socketPath)
}

func (f *FPM) Stop() error {
	if f == nil || f.cmd == nil || f.cmd.Process == nil {
		return nil
	}

	if f.cmd.ProcessState != nil && f.cmd.ProcessState.Exited() {
		return removeSocket(f.socketPath)
	}

	fmt.Println("Stopping PHP-FPM...")

	// SIGQUIT = graceful PHP-FPM shutdown.
	if err := f.cmd.Process.Signal(syscall.SIGQUIT); err != nil {
		return err
	}

	done := make(chan error, 1)

	go func() {
		done <- f.cmd.Wait()
	}()

	select {
	case err := <-done:
		if removeErr := removeSocket(f.socketPath); removeErr != nil {
			return errors.Join(err, removeErr)
		}
		return err

	case <-time.After(10 * time.Second):
		fmt.Println("PHP-FPM did not stop gracefully; killing it")

		killErr := f.cmd.Process.Kill()
		waitErr := <-done

		removeErr := removeSocket(f.socketPath)

		// The process may have exited between the timeout and Kill.
		if killErr != nil {
			if waitErr != nil {
				killErr = fmt.Errorf("kill PHP-FPM: %w (wait: %v)", killErr, waitErr)
			} else {
				killErr = nil
			}
		} else {
			// Wait commonly returns an exit error after a successful kill.
			waitErr = nil
		}

		return errors.Join(killErr, removeErr, waitErr)
	}
}

func removeSocket(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove PHP-FPM socket: %w", err)
	}
	return nil
}
