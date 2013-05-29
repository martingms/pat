package maildir

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
