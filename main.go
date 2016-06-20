// procsnitch daemon - UNIX domain socket service providing process information for local network connections

package main

import (
	"flag"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/op/go-logging"
	"github.com/subgraph/go-procsnitch"
	"github.com/subgraph/procsnitchd/protocol"
	"github.com/subgraph/procsnitchd/service"
)

var log = logging.MustGetLogger("procsnitchd")

var logFormat = logging.MustStringFormatter(
	"%{level:.4s} %{id:03x} %{message}",
)
var ttyFormat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

const ioctlReadTermios = 0x5401

func isTerminal(fd int) bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	return err == 0
}

func setupLoggerBackend() logging.LeveledBackend {
	format := logFormat
	if isTerminal(int(os.Stderr.Fd())) {
		format = ttyFormat
	}
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	formatter := logging.NewBackendFormatter(backend, format)
	leveler := logging.AddModuleLevel(formatter)
	leveler.SetLevel(logging.NOTICE, "procsnitchd")
	return leveler
}

func main() {
	socketFile := flag.String("socket", "", "UNIX domain socket file")
	group := flag.String("group", "", "Group ownership of the socket file")

	logBackend := setupLoggerBackend()
	log.SetBackend(logBackend)
	procsnitch.SetLogger(log)
	protocol.SetLogger(log)
	service.SetLogger(log)

	if os.Geteuid() != 0 {
		log.Error("Must be run as root")
		os.Exit(1)
	}

	flag.Parse()
	if *socketFile == "" {
		log.Critical("UNIX domain socket file must be specified!")
		os.Exit(1)
	}
	if *group == "" {
		log.Critical("group ownership of our UNIX domain socket file must be specified!")
		os.Exit(1)
	}

	procInfo := procsnitch.SystemProcInfo{}
	service := service.NewMortalService("unix", *socketFile, protocol.ConnectionHandlerFactory(procInfo))
	service.Start()
	log.Notice("procsnitchd starting")

	// change the group ownership / permissions of the UNIX domain socket
	cmd := exec.Command("/bin/chgrp", *group, *socketFile)
	err := cmd.Run()
	if err != nil {
		log.Criticalf("failed to chmod socket: %s", err)
		panic("wtf")
	}
	mode := 0775
	err = os.Chmod(*socketFile, os.FileMode(mode))
	if err != nil {
		log.Critical("cannot chmod socket file")
		panic("wtf")
	}

	// wait for a control-c or kill signal
	sigKillChan := make(chan os.Signal, 1)
	signal.Notify(sigKillChan, os.Interrupt, os.Kill)
	for {
		select {
		case <-sigKillChan:
			log.Notice("procsnitchd stopping")
			service.Stop()
			return
		}
	}
}
