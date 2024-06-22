package systemctl

import (
	"errors"
	"testing"
)

var s = NewAsUser()
var serviceName = "testdate"

func TestSaveService(t *testing.T) {
	serv := Service{
		ExecStart:   "/usr/bin/date",
		Description: "Just a test",
	}

	if err := s.SaveService(serviceName, serv); err != nil {
		t.Error(err)
	}
}

func TestDaemonReload(t *testing.T) {
	if _, err := s.DaemonReload(); err != nil {
		t.Error(err)
	}
}

func TestUnits(t *testing.T) {
	list, err := s.Units()
	if err != nil {
		t.Error(err)
	}

	if len(list) == 0 {
		t.Error("units not found")
	}
}

func TestStartService(t *testing.T) {
	if _, err := s.Start(serviceName); err != nil {
		t.Error(err)
	}
}

func TestUnitsPattern(t *testing.T) {
	list, err := s.Units("testdate*")
	if err != nil {
		t.Error(err)
	}

	if len(list) == 0 {
		t.Error("units not found")
	}
}

func TestStatusService(t *testing.T) {
	s, err := s.Status(serviceName)
	if err != nil {
		if !errors.Is(err, ErrUnitIsNotActive) {
			t.Error(err)
		}
	}
	t.Log(s)
}

func TestShowService(t *testing.T) {
	props, err := s.Show(serviceName)
	if err != nil {
		t.Error(err)
	}

	if len(props) == 0 {
		t.Error("props not loaded")
	}
}

func TestStopService(t *testing.T) {
	if _, err := s.Stop(serviceName); err != nil {
		t.Error(err)
	}
}

func TestEnableService(t *testing.T) {
	if _, err := s.Enable(serviceName); err != nil {
		t.Error(err)
	}
}

func TestDisableService(t *testing.T) {
	if _, err := s.Disable(serviceName); err != nil {
		t.Error(err)
	}
}

func TestRemoveService(t *testing.T) {
	if _, err := s.Remove(serviceName); err != nil {
		t.Error(err)
	}
}

func TestResetFailed(t *testing.T) {
	if _, err := s.ResetFailed(); err != nil {
		t.Error(err)
	}
}
