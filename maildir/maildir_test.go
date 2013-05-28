package maildir
// TODO(mg): Rewrite all these tests with a generated Maildir.

import (
  "testing"
	"errors"
)

const (
	TEST_MAILDIR = "/home/mg/dev/pat/maildir/test_mails/Mail"
)

func setup() {}

func TestNewMaildir(t *testing.T) {
	// TODO(mg): Use setup() to create test-environment from scratch.
	maildir, err := NewMaildir(TEST_MAILDIR)
	if err != nil {
		t.Fatal(err)
	}

	if maildir.name != "Mail" {
		t.Fatal("incorrect maildir name. Should be Mail, was ", maildir.name)
	}
}

func TestListMaildirs(t *testing.T) {
	// TODO(mg): Use setup() to create test-environment from scratch.
	maildir, err := NewMaildir(TEST_MAILDIR)
	if err != nil {
		t.Fatal(err)
	}

	submaildirs, err := maildir.ListMaildirs()
	if err != nil {
		t.Fatal(err)
	}
	// TODO(mg): Use setup() to create test-environment from scratch.
	if len(submaildirs) != 3 {
		t.Fatal(errors.New("found wrong number of subdirs"))
	}
}

func TestGetAllMessages(t *testing.T) {
	maildir, err := NewMaildir(TEST_MAILDIR)
	if err != nil {
		t.Fatal(err)
	}

	submaildirs, err := maildir.ListMaildirs()
	if err != nil {
		t.Fatal(err)
	}

	for _, dir := range submaildirs {
		msgs, err := dir.GetAllMessages()
		if err != nil {
			t.Fatal(err)
		}

		if len(msgs) == 0 {
			t.Fatal("found zero messages")
		}
	}
}
