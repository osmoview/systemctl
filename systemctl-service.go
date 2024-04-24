package systemctl

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"strings"
)

//go:embed service.tmpl
var serviceTmpl string
var ext = ".service"

type Service struct {
	// Commands that are executed when this service is started.
	ExecStart string

	// Sets the working directory for executed processes.
	WorkingDirectory string

	// Just a service description
	Description string

	// Name of units after which the service should start
	After string
}

// WriteServiceFile build service template and write it to the writer
func (s Service) WriteServiceFile(w io.Writer) error {
	if err := s.IsValid(); err != nil {
		return err
	}

	tmpl, err := template.New("service").Parse(serviceTmpl)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, s); err != nil {
		return err
	}

	return err
}

// IsValid checks if all required fields are has values
func (s Service) IsValid() error {
	if s.ExecStart == "" {
		return fmt.Errorf("service 'ExecStart' is required")
	}

	return nil
}

func checkServiceExtension(name string) string {
	if !strings.HasSuffix(name, ext) {
		name += ext
	}

	return name
}
