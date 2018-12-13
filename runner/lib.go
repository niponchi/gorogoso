package runner

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bmatcuk/doublestar"
	"github.com/radovskyb/watcher"
)

// CMDLogHandler handle result std output
// from os.Process
func CMDLogHandler(pid chan int, cmd *exec.Cmd) {

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
	fmt.Printf("[gorogoso] Reload... pid %d\n\n", cmd.Process.Pid)
	pid <- cmd.Process.Pid
	cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with '%s'\n", err)
	}

}

func watchGlob(reload chan bool, glob string) {
	w := watcher.New()
	go func() {
		for {
			time.Sleep(1000 * time.Millisecond)
			select {
			case event := <-w.Event:
				fmt.Printf("\n[gosogoso] watcher tigger...\n")
				fmt.Printf("[gorogoso] %s\n", event)
				reload <- true
			case <-w.Closed:
				return
			}
		}
	}()

	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)
	globList := strings.Split(glob, ",")
	fmt.Println("[gorogoso] watch file list")
	for _, g := range globList {
		paths, _ := doublestar.Glob(g)
		for _, path := range paths {
			fmt.Println(path)
			if err := w.Add(path); err != nil {
				panic(err)
			}
		}
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func killCMD(cmd exec.Cmd) {
	if pgid, err := syscall.Getpgid(cmd.Process.Pid); err == nil {
		fmt.Printf("[gorogoso] kill server %d\n", cmd.Process.Pid)
		syscall.Kill(-pgid, syscall.SIGKILL)
	} else {
		fmt.Println(err)
	}
}

// RunCMDAndWatch create new child process
// and monit it
func RunCMDAndWatch(name string, cmdArgs []string, watchGlobPattern string) <-chan int {
	reload := make(chan bool)
	pid := make(chan int)

	cmd := exec.Command(name, cmdArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	go CMDLogHandler(pid, cmd)
	go watchGlob(reload, watchGlobPattern)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill)
	signal.Notify(sig, os.Interrupt)
	go func() {
		for {
			select {
			case <-sig:
				killCMD(*cmd)
				os.Exit(0)
				continue
			case <-reload:
				fmt.Println("[gorogoso] Reload.....")
				killCMD(*cmd)
				cmd = exec.Command(name, cmdArgs...)
				cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
				go CMDLogHandler(pid, cmd)
			}
		}
	}()

	return pid
}

// Monit watch file by glob and
// reload entry file
// everytime on change
func Monit(glob string, entry string) <-chan int {
	fmt.Printf("Watch files: %s\n", glob)
	fmt.Printf("Run entrypoint at: %s\n\n", entry)

	return RunCMDAndWatch("go", []string{"run", entry}, glob)
}
