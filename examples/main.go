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
	fmt.Println("WinMain函数返回", r)
}

func WinMain(Inst win.HINSTANCE) int32 {
	// 1. 注册窗口类
	atom := MyRegisterClass(Inst)
	if atom == 0 {
		fmt.Println("注册窗口类失败:", win.GetLastError())
		return 0
	}
	fmt.Println("注册窗口类成功", atom)

	// 2. 创建窗口
	wnd := win.CreateWindowEx(win.WS_EX_APPWINDOW,
		syscall.StringToUTF16Ptr("主窗口类"),
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
		fmt.Println("创建窗口失败", win.GetLastError())
		return 0
	}
	fmt.Println("创建窗口成功", wnd)
	win.ShowWindow(wnd, win.SW_SHOW)
	win.UpdateWindow(wnd)
	// load file
	gohl.EnableDebug()
	if err := gohl.LoadFile(wnd, "a.html"); err != nil {
		println("LoadFile failed", err.Error())
		return 0
	}
	ui(wnd)

	// 3. 主消息循环
	var msg win.MSG
	msg.Message = win.WM_QUIT + 1 // 让它不等于 win.WM_QUIT

	for win.GetMessage(&msg, 0, 0, 0) > 0 {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}

	return int32(msg.WParam)
}

var (
	handler = &gohl.EventHandler{
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
	gohl.AttachWindowEventHandler(hwnd, handler)
	el := gohl.GetRootElement(hwnd)
	rs := el.Select("#button")
	el = rs[0]
	// el.AttachHandler(handler)
}

// func DefWindowProc(hWnd HWND, Msg uint32, wParam, lParam uintptr) uintptr
func WndProc(hWnd win.HWND, message uint32, wParam uintptr, lParam uintptr) uintptr {
	ret, handled := gohl.ProcNoDefault(hWnd, message, wParam, lParam)
	if handled {
		return uintptr(ret)
	}
	switch message {
	// case win.WM_CREATE:
	// 	ui(hWnd)
	default:
		return win.DefWindowProc(hWnd, message, wParam, lParam)
	}
	return 0
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
	wc.LpszClassName = syscall.StringToUTF16Ptr("主窗口类")
	wc.HIconSm = win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_APPLICATION))
	wc.HIcon = win.LoadIcon(0, win.MAKEINTRESOURCE(win.IDI_APPLICATION))
	wc.HCursor = win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW))

	return win.RegisterClassEx(&wc)
}
