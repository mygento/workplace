package nginx

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type Server struct {
	cmd *exec.Cmd
}

func Start(cfg Config) (*Server, error) {
	if _, err := os.Stat(cfg.Root); err != nil {
		return nil, fmt.Errorf(
			"document root does not exist: %s",
			cfg.Root,
		)
	}

	binary, err := exec.LookPath("nginx")
	if err != nil {
		return nil, fmt.Errorf(
			"nginx not found; install Nginx and make sure it is available in PATH",
		)
	}

	configPath, err := WriteConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Validate configuration before starting Nginx.
	testCmd := exec.Command(
		binary,
		"-t",
		"-c",
		configPath,
		"-p",
		cfg.RuntimeDir,
	)

	output, err := testCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(
			"nginx configuration is invalid:\n%s",
			string(output),
		)
	}

	cmd := exec.Command(
		binary,
		"-c",
		configPath,
		"-p",
		cfg.RuntimeDir,
		"-g",
		"daemon off;",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf(
			"start Nginx: %w",
			err,
		)
	}

	server := &Server{
		cmd: cmd,
	}

	// Give Nginx a moment to fail on things such as an occupied port.
	time.Sleep(100 * time.Millisecond)

	if cmd.ProcessState != nil &&
		cmd.ProcessState.Exited() {
		return nil, fmt.Errorf(
			"nginx exited during startup",
		)
	}

	fmt.Printf(
		"Nginx started (pid %d)\n",
		cmd.Process.Pid,
	)

	return server, nil
}

func (s *Server) Stop() error {
	if s == nil ||
		s.cmd == nil ||
		s.cmd.Process == nil {
		return nil
	}

	if s.cmd.ProcessState != nil &&
		s.cmd.ProcessState.Exited() {
		return nil
	}

	fmt.Println("Stopping Nginx...")

	// SIGQUIT = graceful Nginx shutdown.
	if err := s.cmd.Process.Signal(
		syscall.SIGQUIT,
	); err != nil {
		return err
	}

	done := make(chan error, 1)

	go func() {
		done <- s.cmd.Wait()
	}()

	select {
	case err := <-done:
		return err

	case <-time.After(10 * time.Second):
		fmt.Println(
			"Nginx did not stop gracefully; killing it",
		)

		_ = s.cmd.Process.Kill()

		<-done

		return nil
	}
}
