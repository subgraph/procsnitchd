package client

import (
	"fmt"
	"net"
	"testing"

	"github.com/subgraph/go-procsnitch"
	"github.com/subgraph/procsnitchd/protocol"
	"github.com/subgraph/procsnitchd/service"
)

type MockProcInfo struct {
	procInfo *procsnitch.Info
}

func NewMockProcInfo(procInfo *procsnitch.Info) MockProcInfo {
	p := MockProcInfo{
		procInfo: procInfo,
	}
	return p
}

func (r MockProcInfo) Set(procInfo *procsnitch.Info) {
	r.procInfo = procInfo
}

func (r MockProcInfo) LookupTCPSocketProcess(srcPort uint16, dstAddr net.IP, dstPort uint16) *procsnitch.Info {
	return r.procInfo
}

func (r MockProcInfo) LookupUNIXSocketProcess(socketFile string) *procsnitch.Info {
	return r.procInfo
}

func (r MockProcInfo) LookupUDPSocketProcess(srcPort uint16) *procsnitch.Info {
	return r.procInfo
}

func TestSnitchClientGetUnixSocket(t *testing.T) {
	var err error = nil
	ricochetProcInfo := procsnitch.Info{
		UID:       1,
		Pid:       1,
		ParentPid: 1,
		ExePath:   "/usr/local/bin/ricochet",
		CmdLine:   "testing_cmd_line",
	}

	mockProcInfo := NewMockProcInfo(&ricochetProcInfo)
	socketFile := "procsnitchunixtest.socket"

	service := service.NewMortalService("unix", socketFile, protocol.ConnectionHandlerFactory(mockProcInfo))
	err = service.Start()
	if err != nil {
		t.Errorf("service listener failed to start: %s", err)
		t.Fail()
	}
	defer service.Stop()

	client := NewSnitchClient(socketFile)
	err = client.Start()
	if err != nil {
		t.Errorf("client failed to connect: %s", err)
		t.Fail()
	}
	defer client.Stop()

	info := client.LookupUNIXSocketProcess(socketFile)
	if *info != ricochetProcInfo {
		t.Error("proc info mismatch")
		t.Fail()
	}

	fmt.Println("PROC INFO", info)
}
