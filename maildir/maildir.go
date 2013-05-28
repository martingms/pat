package maildir

import (
	"path"
	"os"
	"io/ioutil"
	"strings"
	"sync"
	"errors"
	"net/mail"
)

var (
	ErrInvalidMsgKey = errors.New("maildir: msg not in maildir")
)

type Maildir struct {
	path string
	paths map[string]string
	msgs map[string]*MaildirMessage
}

// Create a new Maildir-instance from existing maildir path.
func NewMaildir(path string) (*Maildir, error) {
	m := new(Maildir)

	// TODO(mg): Support creating new maildirs?
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	m.path = path
	m paths = map[string]string {
		'cur' : path.Join(m.path, 'cur'),
		'new' : path.Join(m.path, 'new')
	}

	return m, nil
}

// List all sub-maildirs inside a Maildir.
func (m *Maildir) ListMaildirs() ([]*Maildir, error) {
	fis, err := ioutil.ReadDir(m.path)
	if err != nil {
		return nil, err
	}

	dirs := []Maildir{}
	for i, fi := range fis {
		if len(fi.Name()) > 1 && fi.Name()[0] == '.' && fi.IsDir() {
			dir, err := NewMaildir(path.Join(m.path, fi.Name()))
			if err != nil {
				return dirs, err
			}
			dirs = append(dirs, dir)
		}
	}
	return dirs, nil
}

func (m *Maildir) GetMessage(key string) (*mail.Message, error) {
	maildirMsg, err := m.getMaildirMessage(key)
	if err != nil {
		return nil, err
	}

	return getMailMessage(maildirMsg, key)
}

func (m *Maildir) getMailMessage(maildirMsg *MaildirMessage, key string) (*mail.Message, error) {
	path := path.Join(m.path + maildirMsg.subdir + key + maildirMsg.flags)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	// TODO(mg): defer file.Close()? The reader is still used by the returning
	// message... Not sure how to handle this tbh.

	msg, err := mail.ReadMessage(file)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *Maildir) getMaildirMessage(key string) (*MaildirMessage, error) {
	msg, ok := m.msgs[key]
	if !ok {
		return nil, ErrInvalidMsgKey
	} else {
		return msg, nil
	}
}

// TODO(mg): Should provide keys aswell, so one can update, delete mails etc.
func (m *Maildir) GetAllMessages() (map[string]*mail.Message, error) {
	m.refreshMsgs()
	msgMap := map[string]*mail.Message{}
	for key, msg := range m.msgs {
		mailMsg, err := getMailMessage(msg, key)
		if err != nil {
			return nil, err
		}
		msgMap[key] = mailMsg
	}

	return msgMap, nil
}

// Refreshes messages by rebuilding the entire msgs map.
func (m *Maildir) refreshMsgs() {
	// TODO(mg): Find some heuristic to not always update everything.
	m.msgs = map[string]*MaildirMessage{}

	// Check both "cur" and "new".
	// TODO(mg): If we ever add "tmp" to write messages as well, rewrite this.
	for _, path := range m.paths {
		fis, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err) // TODO(mg): Return error instead.
		}

		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}

			name = strings.Split(fi.Name(), ":")
			m.msgs[name[0]], err = NewMailDirMessage(name[1], path)
			if err != nil {
				// TODO(mg): Do something, probably log the offending files name.
			}
		}
	}
}

// Refreshes messages in all sub-maildirs of m.
// TODO(mg): Find a way to return errors from .refreshMsgs()
func (m *Maildir) refreshMsgsSubdirs() {
	dirs := m.ListMailDirs()

	wg := sync.WaitGroup()
	wg.Add(len(dirs))
	for _, dir := range dirs {
		go func(dir *Maildir) {
			dir.refreshMsgs()
			wg.Done()
		}(dir)	
	}
	wg.Wait()
}

// TODO(mg): Set flags according to spec.
// TODO(mg): Should this be exported?
type MaildirMessage struct {
	flags string
	subdir string
}

func NewMaildirMessage(flags, subdir string) (*MaildirMessage, error) {
	msg := new(MaildirMessage)
	msg.flags = flags
	msg.subdir = subdir

	return msg, nil
}
