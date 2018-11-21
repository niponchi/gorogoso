package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/radovskyb/watcher"
)

func cmdLogHandler(pid chan int, cmd *exec.Cmd) {

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	var errStdout, errStderr error
	var stdoutBuf, stderrBuf bytes.Buffer
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err := cmd.Start()
	fmt.Printf("Reload... pid %d\n\n", cmd.Process.Pid)
	pid <- cmd.Process.Pid
	cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with '%s'\n", err)
	}

}

func watchGlob(evn chan bool, glob string) {
	w := watcher.New()
	go func() {
		for {
			time.Sleep(1000 * time.Millisecond)
			select {
			case event := <-w.Event:
				fmt.Println(event)
				evn <- true
			case <-w.Closed:
				return
			}
		}
	}()

	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)
	paths, _ := filepath.Glob(glob)
	for _, path := range paths {
		if err := w.Add(path); err != nil {
			panic(err)
		}
	}
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func runCMDAndWatch(reload chan bool, name string, cmdArgs []string, watchGlobPattern string) <-chan int {
	pid := make(chan int)

	cmd := exec.Command(name, cmdArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	go cmdLogHandler(pid, cmd)
	go watchGlob(reload, watchGlobPattern)

	go func() {
		for {
			<-reload
			if pgid, err := syscall.Getpgid(cmd.Process.Pid); err == nil {
				fmt.Printf("kill server %d\n", cmd.Process.Pid)
				syscall.Kill(-pgid, syscall.SIGKILL)
			} else {
				fmt.Println(err)
			}
			cmd = exec.Command(name, cmdArgs...)
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
			go cmdLogHandler(pid, cmd)
		}
	}()
	return pid
}

// GoroGoso watch file by glob and
// reload entry file
// everytime on change
func GoroGoso(glob string, entry string) <-chan int {
	reload := make(chan bool)
	fmt.Printf("Watch files: %s\n", glob)
	fmt.Printf("Run entrypoint at: %s\n\n", entry)
	return runCMDAndWatch(reload, "go", []string{"run", entry}, glob)
}

func main() {

	watchPattern := flag.String("watch", "*.go", "Glob pattern you want to watch")
	entryfile := flag.String("entry", "main.go", "entryfile path")
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	watchPathPattern := path.Join(dir, *watchPattern)
	entryfilePath := path.Join(dir, *entryfile)

	pid := GoroGoso(watchPathPattern, entryfilePath)
	for {
		<-pid
	}

}
