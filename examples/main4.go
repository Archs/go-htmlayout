package main

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/Archs/go-htmlayout"
	"github.com/lxn/win"
)

var (
	wndClassName = "mytestingclass"
)

func main() {
	inst := win.GetModuleHandle(nil)
	WinMain(inst)
}

func WinMain(Inst win.HINSTANCE) int32 {
	// RegisterClass
	atom := MyRegisterClass(Inst)
	if atom == 0 {
		fmt.Println("RegisterClass failed:", win.GetLastError())
		return 0
	}
	fmt.Println("RegisterClass ok", atom)

	// CreateWindowEx
	wnd := win.CreateWindowEx(win.WS_EX_APPWINDOW,
		syscall.StringToUTF16Ptr(wndClassName),
		nil,
		win.WS_OVERLAPPEDWINDOW|win.WS_CLIPSIBLINGS,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		win.CW_USEDEFAULT,
		0,
		0,
		Inst,
		nil)
	if wnd == 0 {
		fmt.Println("CreateWindowEx failed:", win.GetLastError())
		return 0
	}
	fmt.Println("CreateWindowEx done", wnd)
	win.ShowWindow(wnd, win.SW_SHOW)
	win.UpdateWindow(wnd)
	// load file
	gohl.EnableDebug()
	if err := gohl.LoadFile(wnd, "a.html"); err != nil {
		println("LoadFile failed", err.Error())
		return 0
	}

	// main message loop
	var msg win.MSG

	for win.GetMessage(&msg, 0, 0, 0) > 0 {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}

	return int32(msg.WParam)
}

func WndProc(hWnd win.HWND, message uint32, wParam uintptr, lParam uintptr) uintptr {
	ret, handled := gohl.ProcNoDefault(hWnd, message, wParam, lParam)
	if handled {
		return uintptr(ret)
	}
	switch message {
	case win.WM_CREATE:
		println("win.WM_CREATE called", win.WM_CREATE)
	}
	return win.DefWindowProc(hWnd, message, wParam, lParam)
}

func MyRegisterClass(hInstance win.HINSTANCE) (atom win.ATOM) {
	var wc win.WNDCLASSEX
	wc.CbSize = uint32(unsafe.Sizeof(wc))
	wc.Style = win.CS_HREDRAW | win.CS_VREDRAW
	wc.LpfnWndProc = syscall.NewCallback(WndProc)
	wc.CbClsExtra = 0
	wc.CbWndExtra = 0
	wc.HInstance = hInstance
	wc.HbrBackground = win.GetSysColorBrush(win.COLOR_WINDOWFRAME)
	wc.LpszMenuName = syscall.StringToUTF16Ptr("")
	wc.LpszClassName = syscall.StringToUTF16Ptr(wndClassName)
	wc.HIconSm = win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_APPLICATION))
	wc.HIcon = win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_APPLICATION))
	wc.HCursor = win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW))

	return win.RegisterClassEx(&wc)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
