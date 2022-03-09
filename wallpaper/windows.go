//go:build windows
// +build windows

package wallpaper

import (
	"syscall"
	"unsafe"
)

var (
	moduser32                = syscall.NewLazyDLL("user32.dll")
	procSystemParametersInfo = moduser32.NewProc("SystemParametersInfoW")
)

func setWallpaper(path string) error {
	pathUTF16Ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	procSystemParametersInfo.Call(
		uintptr(0x0014),
		uintptr(0),
		uintptr(unsafe.Pointer(pathUTF16Ptr)),
		uintptr(0x01|0x2),
	)

	return nil
}
