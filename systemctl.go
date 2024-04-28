package systemctl

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	systemctlExec      = "systemctl"
	defaultServicesDir = "/etc/systemd/system/"
	userServicesDir    = ".local/share/systemd/user/"
)

var ErrUnitUnused = errors.New("unit unused")
var ErrUnitIsNotActive = errors.New("unit is not active")
var ErrNoSuchUnit = errors.New("no such unit")

// Unit of the systemd
type Unit struct {
	Unit        string `json:"unit"`
	Load        string `json:"load"`
	Active      string `json:"active"`
	Sub         string `json:"sub"`
	Description string `json:"description"`
}

type Systemctl struct {
	Dir    string
	AsUser bool
}

func NewDefault() *Systemctl {
	return &Systemctl{
		Dir: defaultServicesDir,
	}
}

func NewAsUser() *Systemctl {
	dir, _ := os.UserHomeDir()
	if dir == "" {
		dir = "~"
	}

	return &Systemctl{
		Dir:    filepath.Join(dir, userServicesDir),
		AsUser: true,
	}
}

// Units returns list of units
func (s *Systemctl) Units() (list []Unit, err error) {
	err = s.execSystemctlJSON(&list)
	return
}

// Start service
func (s *Systemctl) Start(name string) (string, error) {
	return s.execSystemctl("start", name)
}

// Stop service
func (s *Systemctl) Stop(name string) (string, error) {
	return s.execSystemctl("stop", name)
}

// Restart service
func (s *Systemctl) Restart(name string) (string, error) {
	return s.execSystemctl("restart", name)
}

// Enable service to autorin with the OS
func (s *Systemctl) Enable(name string) (string, error) {
	return s.execSystemctl("enable", name)
}

// Disable service from autorun
func (s *Systemctl) Disable(name string) (string, error) {
	return s.execSystemctl("disable", name)
}

// Status returns status of the service
func (s *Systemctl) Status(name string) (string, error) {
	out, err := s.execSystemctl("status", name)
	if exiterr, ok := err.(*exec.ExitError); ok {
		switch exiterr.ExitCode() {
		case 2:
			return out, ErrUnitUnused
		case 3:
			return out, ErrUnitIsNotActive
		case 4:
			return out, ErrNoSuchUnit
		}
	}

	return out, err
}

func (s *Systemctl) DaemonReload() (string, error) {
	return s.execSystemctl("daemon-reload")
}

func (s *Systemctl) ResetFailed() (string, error) {
	return s.execSystemctl("reset-failed")
}

// Remove service file and execute daemon reload
func (s *Systemctl) Remove(name string) (string, error) {
	name = checkServiceExtension(name)

	if err := os.Remove(filepath.Join(s.Dir, name)); err != nil {
		return "", err
	}

	return s.DaemonReload()
}

func (s *Systemctl) SaveService(name string, serv Service) error {
	fp := filepath.Join(s.Dir, checkServiceExtension(name))

	f, err := os.OpenFile(fp, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	return serv.WriteServiceFile(f)
}

//
//
//

func (s *Systemctl) execSystemctlJSON(v interface{}, args ...string) error {
	args = append(args, "--output", "json")
	if s.AsUser {
		args = append(args, "--user")
	}

	out, err := exec.Command(systemctlExec, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s, %v", out, err)
	}

	return json.Unmarshal(out, &v)
}

func (s *Systemctl) execSystemctl(args ...string) (string, error) {
	if s.AsUser {
		args = append(args, "--user")
	}

	out, err := exec.Command(systemctlExec, args...).CombinedOutput()
	return fmt.Sprintf("%s", out), err
}
