package runner

import (
	"io/ioutil"
	"path"
	"runtime"
	"testing"
)

func TestWatcher(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	t.Log("Current test filename: " + filename)

	watchPathPattern := path.Join(path.Dir(filename), "../test/*.go")
	entryfilePath := path.Join(path.Dir(filename), "../test/main.go")

	if b, err := ioutil.ReadFile(entryfilePath); err != nil {
		panic(err)
	} else {

		next := []byte("// hi")
		newfile := append(b, next...)

		pid := Monit(watchPathPattern, entryfilePath)
		PID := <-pid
		t.Logf("first reload %d", PID)
		go func() {
			t.Log("Update file...")
			if err := ioutil.WriteFile(entryfilePath, newfile, 0644); err != nil {
				panic(err)
			}
		}()
		// check if instance reload
		// reset file
		if PID2 := <-pid; PID2 == PID {
			t.Errorf("Process not update")
		} else {
			t.Logf("second reload %d", PID2)
		}
		t.Log("Rollback file....")
		if err := ioutil.WriteFile(entryfilePath, b, 0644); err != nil {
			panic(err)
		}

	}

}
