//go:build windows
// +build windows

package process

import (
	"syscall"
	"unsafe"
)

var (
	modkernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procCreateToolhelp32Snapshot = modkernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modkernel32.NewProc("Process32FirstW")
	procProcess32Next            = modkernel32.NewProc("Process32NextW")
	procCloseHandle              = modkernel32.NewProc("CloseHandle")
)

type PROCESSENTRY32 struct {
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	CntThreads        uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [260]uint16
}

type windowsProcess struct {
	Pid         int
	ParentPid   int
	ProcessName string
}

func (p *windowsProcess) pid() int            { return p.Pid }
func (p *windowsProcess) parentPid() int      { return p.ParentPid }
func (p *windowsProcess) processName() string { return p.ProcessName }

func newWindowsProcess(e *PROCESSENTRY32) *windowsProcess {
	// Find when the string ends for decoding
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return &windowsProcess{
		Pid:         int(e.ProcessID),
		ParentPid:   int(e.ParentProcessID),
		ProcessName: syscall.UTF16ToString(e.ExeFile[:end]),
	}
}

func listProcesses() ([]iProcess, error) {
	handle, _, _ := procCreateToolhelp32Snapshot.Call(uintptr(0x00000002), uintptr(0x00000000))
	defer procCloseHandle.Call(handle)

	var entry PROCESSENTRY32
	entry.Size = uint32(unsafe.Sizeof(entry))
	res, _, _ := procProcess32First.Call(handle, uintptr(unsafe.Pointer(&entry)))
	if res == 0 {
		return nil, syscall.GetLastError()
	}

	results := make([]iProcess, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		res, _, _ = procProcess32Next.Call(handle, uintptr(unsafe.Pointer(&entry)))
		if res == 0 {
			break
		}
	}

	return results, nil
}
