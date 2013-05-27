package maildir

import (
	"path"
	"os"
	"io/ioutil"
	"strings"
	"sync"
)

type Maildir struct {
	path string
	paths map[string]string
	msgs map[string]*MaildirMessage
}

// Create a new Maildir-instance from existing maildir path.
func NewMaildir(path string) *Maildir {
	m := new(Maildir)

	_, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	m.path = path
	m paths = map[string]string {
		'cur' : path.Join(m.path, 'cur'),
		'new' : path.Join(m.path, 'new')
	}

	return m
}

// List all sub-maildirs inside a Maildir.
func (m *Maildir) ListMaildirs() []*Maildir {
	fis, err := ioutil.ReadDir(m.path)
	if err != nil {
		panic(err)
	}

	folders := []Maildir{}
	for i, fi := range fis {
		if len(fi.Name()) > 1 && fi.Name()[0] == '.' && fi.IsDir() {
			folders = append(folders, NewMaildir(path.Join(m.path, fi.Name())))
		}
	}
	return folders
}

func (m *Maildir) updateMsgs() {
	// TODO(mg): Find some heuristic to not always update everything.
	m.msgs = map[string]*MaildirMessage{}

	// Check both "cur" and "new".
	// TODO(mg): If we ever add "tmp" to write messages as well, rewrite this.
	for _, path := range m.paths {
		fis, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}

		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}

			name = strings.Split(fi.Name(), ":")
			m.msgs[name[0]] = NewMailDirMessage(name[1])
		}
	}
}

func (m *Maildir) updateMsgsSubdirs() {
	dirs := m.ListMailDirs()

	wg := sync.WaitGroup()
	wg.Add(len(dirs))
	for _, dir := range dirs {
		go func(dir *Maildir) {
			dir.updateMsgs()
			wg.Done()
		}(dir)	
	}
	wg.Wait()
}

type MaildirMessage struct {
	flags string
}

func NewMaildirMessage(flags string) *MaildirMessage {
	msg := new(MaildirMessage)
	msg.flags = flags

	return msg
}
