package systemctl

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
)

const journalctlExec = "journalctl"

type Journalctl struct {
	AsUser bool
}

// NewDefaultJournal returns journal istance to reads log from the system
func NewDefaultJournal() Journalctl {
	return Journalctl{}
}

// NewUserJournal returns journal prepared to reads logs from the current user
func NewUserJournal() Journalctl {
	return Journalctl{AsUser: true}
}

type JournalGetOpt struct {
	// Show logs from the specified unit
	Unit string

	// Number of journal entries to show
	Lines string

	// Show entries not older than the specified date
	Since string

	// Show entries after the specified cursor
	AfterCursor string

	// TODO: add additional options
}

func (opt JournalGetOpt) toArgs() (args []string) {
	if opt.Unit != "" {
		args = append(args, "-u", opt.Unit)
	}

	if opt.Lines != "" {
		args = append(args, "-n", opt.Lines)
	} else {
		args = append(args, "-n", "20")
	}

	if opt.Since != "" {
		args = append(args, "--since", opt.Since)
	}

	if opt.AfterCursor != "" {
		args = append(args, opt.AfterCursor)
	}

	return args
}

type JournalMsg struct {
	Message    string `json:"MESSAGE"`
	Timestamp  string `json:"__REALTIME_TIMESTAMP"`
	JobType    string `json:"JOB_TYPE"`
	Transport  string `json:"_TRANSPORT"`
	Cursor     string `json:"__CURSOR"`
	ExitStatus string `json:"EXIT_STATUS"`
	ExitCode   string `json:"EXIT_CODE"`
}

// Get journal messages by options
func (j Journalctl) Get(opt JournalGetOpt) (msgs []JournalMsg, err error) {
	cmd, stdout, err := j.execJournalctl(opt.toArgs())
	if err != nil {
		return nil, err
	}

	// all errors occured when read stdout
	var errs []error

	// read stdout and parse journal messages
	go func() {
		s := bufio.NewScanner(stdout)

		for s.Scan() {
			var msg JournalMsg
			if err = json.Unmarshal(s.Bytes(), &msg); err != nil {
				errs = append(errs, err)
				continue
			}

			msgs = append(msgs, msg)
		}
	}()

	if err := cmd.Wait(); err != nil {
		errs = append(errs, err)
	}

	return msgs, errors.Join(errs...)
}

func (j Journalctl) execJournalctl(args []string) (cmd *exec.Cmd, stdout io.ReadCloser, err error) {
	args = append(args, "--output", "json")
	if j.AsUser {
		args = append(args, "--user")
	}

	// TODO: add context

	cmd = exec.Command(journalctlExec, args...)
	stdout, _ = cmd.StdoutPipe()

	return cmd, stdout, cmd.Start()
}
