package systemctl

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os/exec"

	"github.com/acarl005/stripansi"
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

	Follow bool
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

	if opt.Follow {
		args = append(args, "-f")
	}

	return args
}

type JournalMsg struct {
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
	JobType    string `json:"job_type"`
	Transport  string `json:"transport"`
	Cursor     string `json:"cursor"`
	ExitStatus string `json:"exit_status"`
	ExitCode   string `json:"exit_code"`
}

type journalMsgFields struct {
	Timestamp  string `json:"__REALTIME_TIMESTAMP"`
	JobType    string `json:"JOB_TYPE"`
	Transport  string `json:"_TRANSPORT"`
	Cursor     string `json:"__CURSOR"`
	ExitStatus string `json:"EXIT_STATUS"`
	ExitCode   string `json:"EXIT_CODE"`
}

// Get journal messages by options;
// Return stream reader, close function and error.
// To close stream, use close function.
func (j Journalctl) Stream(opt JournalGetOpt) (io.ReadCloser, func(), error) {
	opt.Follow = true
	ctx, cancel := context.WithCancel(context.Background())

	cmd, stdout, err := j.execJournalctlWithContext(ctx, opt.toArgs())
	if err != nil {
		return nil, cancel, err
	}

	var close = func() {
		cancel()
		if cmd == nil || cmd.Process == nil {
			return
		}
		cmd.Process.Kill()
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			return
		}
		stdout.Close()
	}()

	return stdout, close, err
}

func (j Journalctl) DecodeMsgString(line []byte) (JournalMsg, error) {
	message, err := j.decodeMsgString(line)
	if err != nil {
		return JournalMsg{}, err
	}

	var rawmsg journalMsgFields
	if err = json.Unmarshal(line, &rawmsg); err != nil {
		return JournalMsg{}, err
	}

	msg := JournalMsg{
		Message:    stripansi.Strip(message),
		Timestamp:  rawmsg.Timestamp,
		JobType:    rawmsg.JobType,
		Transport:  rawmsg.Timestamp,
		Cursor:     rawmsg.Cursor,
		ExitStatus: rawmsg.Cursor,
		ExitCode:   rawmsg.ExitCode,
	}

	return msg, nil
}

type journalMsgString struct {
	Message string `json:"MESSAGE"`
}

type journalMsgBytes struct {
	Message []byte `json:"MESSAGE"`
}

func (j Journalctl) decodeMsgString(line []byte) (message string, err error) {
	// try to decode message as string
	var strmsg journalMsgString
	if err = json.Unmarshal(line, &strmsg); err == nil {
		return strmsg.Message, nil
	}

	// try to decode message as bytes
	var bytesmsg journalMsgBytes
	if err = json.Unmarshal(line, &bytesmsg); err == nil {
		return string(bytesmsg.Message), nil
	}

	return "", errors.New("failed to decode message text")
}

func (j Journalctl) execJournalctlWithContext(ctx context.Context, args []string) (cmd *exec.Cmd, stdout io.ReadCloser, err error) {
	args = append(args, "--output", "json")
	if j.AsUser {
		args = append(args, "--user")
	}

	// TODO: add context

	cmd = exec.CommandContext(ctx, journalctlExec, args...)
	stdout, _ = cmd.StdoutPipe()

	return cmd, stdout, cmd.Start()
}
