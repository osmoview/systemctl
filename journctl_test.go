package systemctl

import (
	"bufio"
	"testing"
)

func TestJournalctl(t *testing.T) {
	j := NewDefaultJournal()
	msgs, _, err := j.Stream(JournalGetOpt{Unit: "systemd-logind"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	buf := bufio.NewScanner(msgs)

	for buf.Scan() {
		m, err := j.DecodeMsgString(buf.Bytes())
		if err != nil {
			t.Error(err)
		}
		t.Log(m)
	}
}
