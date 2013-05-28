package maildir

import (
	"path"
	"os"
	"io/ioutil"
	"strings"
	"sync"
	"errors"
	"net/mail"
	"sort"
)

var (
	ErrInvalidMsgKey = errors.New("maildir: msg not in maildir")
)

type Maildir struct {
	path string
	paths []string
	msgs map[string]*maildirMessage
  name string
}

// Create a new Maildir-instance from existing maildir path.
func NewMaildir(pathstr string) (*Maildir, error) {
	m := new(Maildir)

	// TODO(mg): Support creating new maildirs?
	fi, err := os.Stat(pathstr)
	if err != nil {
		return nil, err
	}
  m.name = fi.Name()

	m.path = pathstr
	m.paths = []string{"cur", "new"}

	return m, nil
}

// List all sub-maildirs inside a Maildir.
func (m *Maildir) ListMaildirs() ([]*Maildir, error) {
	fis, err := ioutil.ReadDir(m.path)
	if err != nil {
		return nil, err
	}

	dirs := []*Maildir{}
	for _, fi := range fis {
		if len(fi.Name()) > 1 && fi.Name()[0] == '.' && fi.IsDir() {
			dir, err := NewMaildir(path.Join(m.path, fi.Name()))
			if err != nil {
				return dirs, err
			}
			dirs = append(dirs, dir)
		}
	}
	return dirs, nil
}

func (m *Maildir) HasNewMail() bool {
	fis, err := ioutil.ReadDir(path.Join(m.path, "new"))
	if err != nil {
		panic(err)
	}

	if len(fis) != 0 {
		return true
	}

	return false
}

// Get a single mail.Message given the maildir key.
func (m *Maildir) GetMessage(key string) (*mail.Message, error) {
	maildirMsg, ok := m.msgs[key]
	if !ok {
		return nil, ErrInvalidMsgKey
	}

	return m.getMailMessage(maildirMsg, key)
}

// Get every message in the maildir.
func (m *Maildir) GetAllMessages() (map[string]*mail.Message, error) {
	m.refreshMsgs()
	msgMap := map[string]*mail.Message{}
	for key, msg := range m.msgs {
		mailMsg, err := m.getMailMessage(msg, key)
		if err != nil {
			return nil, err
		}
		msgMap[key] = mailMsg
	}

	return msgMap, nil
}

func (m *Maildir) getMailMessage(maildirMsg *maildirMessage, key string) (*mail.Message, error) {
	path := path.Join(m.path, maildirMsg.subdir, key + maildirMsg.getFlagStr())
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

// Refreshes messages by rebuilding the entire msgs map.
// TODO(mg): Make it concurrent.
func (m *Maildir) refreshMsgs() {
	// TODO(mg): Find some heuristic to not always update everything.
	m.msgs = map[string]*maildirMessage{}

	// Check both "cur" and "new".
	// TODO(mg): If we ever add "tmp" to write messages as well, rewrite this.
	for _, dir := range m.paths {
		fis, err := ioutil.ReadDir(path.Join(m.path, dir))
		if err != nil {
			panic(err) // TODO(mg): Return error instead.
		}

		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}

			name := strings.Split(fi.Name(), ":")
      flags := ""
      // If not the message has no flags.
      if len(name) > 1 {
        flags = name[1]
      }
			m.msgs[name[0]], err = newMaildirMessage(flags, dir)
			if err != nil {
				// TODO(mg): Do something, probably log the offending files name.
			}
		}
	}
}

// Refreshes messages in all sub-maildirs of m.
// TODO(mg): Find a way to return errors from .refreshMsgs()
// NOT IN USE
func (m *Maildir) refreshMsgsSubdirs() {
	// TODO(mg): Handle error here.
	dirs, _ := m.ListMaildirs()

	wg := new(sync.WaitGroup)
	wg.Add(len(dirs))
	for _, dir := range dirs {
		go func(dir *Maildir) {
			dir.refreshMsgs()
			wg.Done()
		}(dir)	
	}
	wg.Wait()
}

type maildirMessage struct {
	flags map[rune]bool
	subdir string
	flagPrefix string
}

func newMaildirMessage(flags, subdir string) (*maildirMessage, error) {
	msg := new(maildirMessage)
	msg.flags = map[rune]bool{}
	if strings.Contains(flags, "2,") {
		msg.flagPrefix = ":2,"
		flagStr := strings.Split(flags, "2,")[1]
		msg.setFlagsFromStr(flagStr)
	}

	msg.subdir = subdir

	return msg, nil
}

func (mmsg *maildirMessage) getFlagStr() string {
	flagList := []string{}

	for char := range mmsg.flags {
		flagList = append(flagList, string(char))
	}

	sort.Strings(flagList)

	// TODO(mg): Find better way to string concatenate.
	flagStr := mmsg.flagPrefix
	for _, char := range flagList {
		flagStr = flagStr + char
	}

	return flagStr
}

func (mmsg *maildirMessage) setFlagsFromStr(flagStr string) {
	for _, f := range flagStr {
		mmsg.flags[f] = true
	}	
}

func (mmsg *maildirMessage) removeFlagsFromStr(flagStr string) {
	for _, f := range flagStr {
		delete(mmsg.flags, f)
	}
}
