package systemctl

import (
	"bufio"
	"testing"
	"time"
)

func TestJournalctl(t *testing.T) {
	j := NewDefaultJournal()
	msgs, close, err := j.Stream(JournalGetOpt{Unit: "systemd-logind"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	buf := bufio.NewScanner(msgs)

	go func() {
		time.Sleep(2 * time.Second)
		close()
	}()
	
	for buf.Scan() {
		m, err := j.DecodeMsgString(buf.Bytes())
		if err != nil {
			t.Error(err)
		}
		t.Log(m)
	}
}
