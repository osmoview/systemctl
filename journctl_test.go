package systemctl

import "testing"

func TestJournalctl(t *testing.T) {
	j := NewDefaultJournal()
	msgs, err := j.Get(JournalGetOpt{Unit: "systemd-logind"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(msgs) == 0 {
		t.Error("msgs is empty")
	}

	t.Log(msgs)
}
