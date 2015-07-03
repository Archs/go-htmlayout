package main

import (
	"fmt"
	"runtime"
	"syscall"

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
	// 2. 创建窗口
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
	// ui(wnd)

	// 3. 主消息循环
	var msg win.MSG
	msg.Message = win.WM_QUIT + 1 // 让它不等于 win.WM_QUIT

	for win.GetMessage(&msg, 0, 0, 0) > 0 {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}

	return int32(msg.WParam)
}
