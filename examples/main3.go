package main

import (
	"fmt"
	"log"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/Archs/go-htmlayout"
	"github.com/lxn/win"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	inst := win.GetModuleHandle(nil)
	r := WinMain(inst)
	fmt.Println("WinMain:", r)
}

var (
	handler = gohl.EventHandler{
		OnBehaviorEvent: func(el *gohl.Element, params *gohl.BehaviorEventParams) bool {
			log.Println("OnBehaviorEvent:", el, params, "|", gohl.BUTTON_CLICK)
			if params.Cmd == gohl.BUTTON_CLICK {
				log.Println("button clicked")
			}
			return false
		},
		OnScriptCall: func(el *gohl.Element, params *gohl.XcallParams) bool {
			log.Printf("xcall: %s %v\n\t", params.MethodName, params)
			for i, v := range params.Argv {
				log.Print("\t", i, uintptr(unsafe.Pointer(v)), v, v.IsString())
			}
			log.Println()
			return true
		},
	}
)

func ui(hwnd win.HWND) {
	gohl.AttachWindowEventHandler(hwnd, &handler)
	el := gohl.GetRootElement(hwnd)
	rs := el.Select("#button")
	el = rs[0]
	// el.AttachHandler(handler)
}

func WinMain(Inst win.HINSTANCE) int32 {
	// CreateWindowEx
	wnd := win.CreateWindowEx(win.WS_EX_APPWINDOW,
		syscall.StringToUTF16Ptr(gohl.GetClassName()),
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
	fmt.Println("ok CreateWindowEx", wnd)
	win.ShowWindow(wnd, win.SW_SHOW)
	win.UpdateWindow(wnd)
	// load file
	gohl.EnableDebug()
	if err := gohl.LoadFile(wnd, "a.html"); err != nil {
		println("LoadFile failed", err.Error())
		return 0
	}
	ui(wnd)

	// main loop
	var msg win.MSG

	for win.GetMessage(&msg, 0, 0, 0) > 0 {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}

	return int32(msg.WParam)
}
