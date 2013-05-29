package maildir

// TODO(mg): Rewrite all these tests with a generated Maildir.

import (
	"errors"
	"fmt"
	"testing"
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

//func ExamplePrint() {
func TestPrint(t *testing.T) {
	return // Shouldn't you be able to .Skip()?
	maildir, err := NewMaildir(TEST_MAILDIR)
	if err != nil {
		panic(err)
	}

	submaildirs, err := maildir.ListMaildirs()
	if err != nil {
		panic(err)
	}

	fmt.Println("Maildir       HasNewMail")
	fmt.Println("========================")
	for _, dir := range submaildirs {
		fmt.Println(dir.name, " ", dir.HasNewMail())
	}

	fmt.Println("\n.personlig")
	fmt.Println("==========")

	mails, err := submaildirs[1].GetAllMessages()
	if err != nil {
		panic(err)
	}

	for _, mail := range mails {
		from := mail.Header.Get("From")
		subj := mail.Header.Get("Subject")

		fmt.Println(from, " ", subj)
	}
}
