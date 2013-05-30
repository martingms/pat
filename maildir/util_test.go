package maildir

import (
	"os"
	"path/filepath"
	"testing"
)

func testWalkFuncBuilder(m map[string]error) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		m[path] = err
		return nil
	}
}

// We assume that filepath.Walk works as expected,
// and that we can check the output directly against that.
func TestConcurrentWalk(t *testing.T) {
	return // Currently deadlocks on MAXPROCS > 1, and is not used so we skip it.
	builtinWalkMap := map[string]error{}
	concurrentWalkMap := map[string]error{}

	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	// TODO(mg): Should use some standard test directories/files.
	err = filepath.Walk(wd, testWalkFuncBuilder(builtinWalkMap))
	if err != nil {
		t.Error(err)
	}

	err = ConcurrentWalk(wd, testWalkFuncBuilder(concurrentWalkMap))
	if err != nil {
		t.Error(err)
	}

	for key, err := range builtinWalkMap {
		e, ok := concurrentWalkMap[key]
		if !ok || e != err {
			t.FailNow()
		}
	}

	for key, err := range concurrentWalkMap {
		e, ok := builtinWalkMap[key]
		if !ok || e != err {
			t.FailNow()
		}
	}
}
