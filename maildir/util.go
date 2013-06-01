package maildir

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"net/mail"
)

// ConcurrentWalk walks the file tree rooted at root, calling walkFn for each
// file or directory in the tree, including root. All errors that arise visiting
// files and directories are filtered by walkFn. It is identical to
// filepath.Walk, other than that is is concurrent and therefore does not
// process directories in a lexical manner.
// Walk does not follow symbolic links.
func ConcurrentWalk(root string, walkFn filepath.WalkFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		return walkFn(root, nil, err)
	}
	return concurrentWalk(root, info, walkFn)
}

// concurrentWalk recursively descends path, calling walkFn.
func concurrentWalk(path string, info os.FileInfo, walkFn filepath.WalkFunc) error {
	err := walkFn(path, info, nil)
	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	list, err := ioutil.ReadDir(path)
	if err != nil {
		return walkFn(path, info, err)
	}

	active := len(list)
	c := make(chan error)

	for _, fileInfo := range list {
		go func(fileInfo os.FileInfo) {
			err = concurrentWalk(filepath.Join(path, fileInfo.Name()), fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir {
					c <- err
					return
				}
			}
			c <- nil
		}(fileInfo)
	}

	for ; active > 0; active-- {
		select {
		case err := <-c:
			if err != nil {
				return err
			}
		}
	}
	return nil
}


// Copied from pat utils to avoid circular dependencies.
var formats = []string{
  "Mon, 2 Jan 2006 15:04:05 -0700",
  "Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
  time.ANSIC,
  time.UnixDate,
  time.RubyDate,
  time.RFC822,
  time.RFC822Z,
  time.RFC850,
  time.RFC1123,
  time.RFC1123Z,
  time.RFC3339,
  time.RFC3339Nano, // lol
}

// Copied from pat utils to avoid circular dependencies.
func parseDate(h *mail.Header) (t time.Time, err error) {
  // We prefer Date header to Delivery-date.
  // TODO(mg): Should we?
  dateHeaders := []string{h.Get("Date"), h.Get("Delivery-date")}

header_loop:
  for _, val := range dateHeaders {
    for _, format := range formats {
      t, err = time.Parse(format, val)
      if err == nil {
        break header_loop
      }
    }
  }

  return t, err 
}
