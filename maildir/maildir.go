package maildir

import (
	"errors"
	"io/ioutil"
	"net/mail"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
)

var (
	ErrInvalidMsgKey = errors.New("maildir: msg not in maildir")
)

type Maildir struct {
	Name  string
	path  string
	paths []string
	msgs  map[string]*maildirMessage
	sync.RWMutex
}

// Create a new Maildir-instance from existing maildir path.
func NewMaildir(pathstr string) (*Maildir, error) {
	m := new(Maildir)

	// TODO(mg): Support creating new maildirs?
	fi, err := os.Stat(pathstr)
	if err != nil {
		return nil, err
	}
	m.Name = fi.Name()

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
	m.RLock()
	_, ok := m.msgs[key]
	m.RUnlock()
	if !ok {
		return nil, ErrInvalidMsgKey
	}

	return m.getMailMessage(key)
}

// TODO(mg): Check if parallelism has any gains on large amounts of mail.
// Overhead to large for small amounts at least, so this needs to be checked.
// Might have to revert the entire thing.
//
// Get every message in the maildir.
func (m *Maildir) GetAllMessages() (map[string]*mail.Message, error) {
	m.refreshMsgs()
	mutex := new(sync.RWMutex) // For locking the msgMap
	msgMap := map[string]*mail.Message{}

	m.RLock()
	wg := new(sync.WaitGroup)
	wg.Add(len(m.msgs))
	for key := range m.msgs {
		go func(key string) {
			mailMsg, err := m.getMailMessage(key)
			if err != nil {
				// TODO(mg): Handle errors better.
				panic(err)
			}
			mutex.Lock()
			msgMap[key] = mailMsg
			mutex.Unlock()
			wg.Done()
		}(key)
	}
	m.RUnlock()

	wg.Wait()

	return msgMap, nil
}

func (m *Maildir) getMailMessage(key string) (*mail.Message, error) {
	m.RLock()
	maildirMsg := m.msgs[key]
	m.RUnlock()
	file, err := os.Open(maildirMsg.curName)
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
// TODO(mg): Control if parallelism actually works any better.
func (m *Maildir) refreshMsgs() {
	// TODO(mg): Find some heuristic to not always update everything.
	m.Lock()
	m.msgs = map[string]*maildirMessage{}
	m.Unlock()

	// Check both "cur" and "new".
	// TODO(mg): If we ever add "tmp" to write messages as well, rewrite this.
	wg := new(sync.WaitGroup)
	for _, dir := range m.paths {
		fis, err := ioutil.ReadDir(path.Join(m.path, dir))
		if err != nil {
			panic(err) // TODO(mg): Return error instead.
		}

		wg.Add(len(fis))
		for _, fi := range fis {
			go func(dir string, fi os.FileInfo) {
				if fi.IsDir() {
					return
				}

				name := strings.Split(fi.Name(), ":")
				flags := ""
				// If not the message has no flags.
				if len(name) > 1 {
					flags = name[1]
				}
				m.Lock()
				m.msgs[name[0]], err = newMaildirMessage(m.path, dir, name[0], flags)
				m.Unlock()
				if err != nil {
					// TODO(mg): Do something, probably log the offending files name.
					panic(err)
				}

				wg.Done()
			}(dir, fi)
		}
	}
	wg.Wait()
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
	flags      map[rune]bool
	path       string
	flagPrefix string
	curName    string
	subdir     string
	key        string
}

func newMaildirMessage(fpath, subdir, key, flags string) (*maildirMessage, error) {
	mmsg := new(maildirMessage)
	mmsg.flags = map[rune]bool{}
	if strings.Contains(flags, "2,") {
		mmsg.flagPrefix = ":2,"
		flagStr := strings.Split(flags, "2,")[1]
		mmsg.setFlagsFromStr(flagStr)
	}

	mmsg.path = fpath
	mmsg.subdir = subdir
	mmsg.key = key
	mmsg.curName = path.Join(fpath, subdir, key+mmsg.getFlagStr())

	return mmsg, nil
}

func (mmsg *maildirMessage) getFlagStr() string {
	flagList := []string{}

	for char := range mmsg.flags {
		flagList = append(flagList, string(char))
	}

	// Spec says flags must be sorted.
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

func (mmsg *maildirMessage) flush() error {
	newName := path.Join(mmsg.path, mmsg.subdir, mmsg.key+mmsg.getFlagStr())
	err := os.Rename(mmsg.curName, newName)
	if err != nil {
		return err
	}
	mmsg.curName = newName

	return nil
}
