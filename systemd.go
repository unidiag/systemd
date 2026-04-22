package systemd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Service struct {
	Name        string // service name (without .service)
	Description string
	ExecStart   string
	WorkingDir  string
	UserMode    bool // true = systemctl --user
	Restart     string
	After       string
}

// Create writes the service file
func (s *Service) Create() (string, error) {
	if s.Name == "" || s.ExecStart == "" {
		return "", fmt.Errorf("Name and ExecStart are required")
	}

	if s.Restart == "" {
		s.Restart = "always"
	}

	if s.After == "" {
		s.After = "network.target"
	}

	var serviceDir string

	if s.UserMode {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		serviceDir = filepath.Join(home, ".config/systemd/user")
	} else {
		serviceDir = "/etc/systemd/system"
	}

	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		return "", err
	}

	path := filepath.Join(serviceDir, s.Name+".service")

	content := fmt.Sprintf(`[Unit]
Description=%s
After=%s

[Service]
ExecStart=%s
WorkingDirectory=%s
Restart=%s
RestartSec=3

[Install]
WantedBy=multi-user.target
`, s.Description, s.After, s.ExecStart, s.WorkingDir, s.Restart)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return "", err
	}

	return path, nil
}

// Run systemctl command
func (s *Service) systemctl(args ...string) error {
	cmdArgs := args

	if s.UserMode {
		cmdArgs = append([]string{"--user"}, args...)
	}

	cmd := exec.Command("systemctl", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Enable service
func (s *Service) Enable() error {
	return s.systemctl("enable", s.Name)
}

// Start service
func (s *Service) Start() error {
	return s.systemctl("start", s.Name)
}

// Restart service
func (s *Service) RestartService() error {
	return s.systemctl("restart", s.Name)
}

// Stop service
func (s *Service) Stop() error {
	return s.systemctl("stop", s.Name)
}

// Reload daemon
func (s *Service) Reload() error {
	return s.systemctl("daemon-reexec")
}

// Full install (create + reload + enable + restart)
func (s *Service) InstallAndStart() error {
	_, err := s.Create()
	if err != nil {
		return err
	}

	if err := s.Reload(); err != nil {
		return err
	}

	if err := s.Enable(); err != nil {
		return err
	}

	return s.RestartService()
}

// Detect if running as root
func IsRoot() bool {
	return os.Geteuid() == 0
}

// Helper: create service from current binary
func NewFromCurrentBinary(name string) (*Service, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	execPath, err = filepath.Abs(execPath)
	if err != nil {
		return nil, err
	}

	return &Service{
		Name:        name,
		Description: name,
		ExecStart:   execPath,
		WorkingDir:  filepath.Dir(execPath),
		UserMode:    !IsRoot(),
	}, nil
}
