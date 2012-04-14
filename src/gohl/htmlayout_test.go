package gohl

import (
	"log"
	"syscall"
	"testing"
	"unsafe"
)

const (
	WM_CREATE     = 1
	WM_DESTROY    = 2
	WM_CLOSE      = 16
	WM_QUIT       = 0x0012
	WM_ERASEBKGND = 0x0014
	WM_SHOWWINDOW = 0x0018
	ERROR_SUCCESS = 0
)

var (
	moduser32 = syscall.NewLazyDLL("user32.dll")

	procRegisterClassExW = moduser32.NewProc("RegisterClassExW")
	procCreateWindowExW  = moduser32.NewProc("CreateWindowExW")
	procDefWindowProcW   = moduser32.NewProc("DefWindowProcW")
	procDestroyWindow    = moduser32.NewProc("DestroyWindow")
	procPostQuitMessage  = moduser32.NewProc("PostQuitMessage")
	procGetMessageW      = moduser32.NewProc("GetMessageW")
	procTranslateMessage = moduser32.NewProc("TranslateMessage")
	procDispatchMessageW = moduser32.NewProc("DispatchMessageW")
	procSendMessageW     = moduser32.NewProc("SendMessageW")
	procPostMessageW     = moduser32.NewProc("PostMessageW")

	classRegistered = false
)

type Wndclassex struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uint32
	Icon       uint32
	Cursor     uint32
	Background uint32
	MenuName   *uint16
	ClassName  *uint16
	IconSm     uint32
}

type Msg struct {
	Hwnd    uint32
	Message uint32
	Wparam  int32
	Lparam  int32
	Time    uint32
	Pt      Point
}

func RegisterClassEx(wndclass *Wndclassex) (atom uint16, err syscall.Errno) {
	r0, _, e1 := syscall.Syscall(procRegisterClassExW.Addr(), 1, uintptr(unsafe.Pointer(wndclass)), 0, 0)
	atom = uint16(r0)
	if atom == 0 {
		if e1 != 0 {
			err = syscall.Errno(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CreateWindowEx(exstyle uint32, classname *uint16, windowname *uint16, style uint32, x int32, y int32, width int32, height int32, wndparent uint32, menu uint32, instance uint32, param uintptr) (hwnd uint32, err syscall.Errno) {
	r0, _, e1 := syscall.Syscall12(procCreateWindowExW.Addr(), 12, uintptr(exstyle), uintptr(unsafe.Pointer(classname)), uintptr(unsafe.Pointer(windowname)), uintptr(style), uintptr(x), uintptr(y), uintptr(width), uintptr(height), uintptr(wndparent), uintptr(menu), uintptr(instance), uintptr(param))
	hwnd = uint32(r0)
	if hwnd == 0 {
		if e1 != 0 {
			err = syscall.Errno(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func DefWindowProc(hwnd uint32, msg uint32, wparam uintptr, lparam uintptr) (lresult int32) {
	r0, _, _ := syscall.Syscall6(procDefWindowProcW.Addr(), 4, uintptr(hwnd), uintptr(msg), uintptr(wparam), uintptr(lparam), 0, 0)
	lresult = int32(r0)
	return
}

func DestroyWindow(hwnd uint32) (err syscall.Errno) {
	r1, _, e1 := syscall.Syscall(procDestroyWindow.Addr(), 1, uintptr(hwnd), 0, 0)
	if int(r1) == 0 {
		if e1 != 0 {
			err = syscall.Errno(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func PostQuitMessage(exitcode int32) {
	syscall.Syscall(procPostQuitMessage.Addr(), 1, uintptr(exitcode), 0, 0)
	return
}

func GetMessage(msg *Msg, hwnd uint32, MsgFilterMin uint32, MsgFilterMax uint32) (ret int32, err syscall.Errno) {
	r0, _, e1 := syscall.Syscall6(procGetMessageW.Addr(), 4, uintptr(unsafe.Pointer(msg)), uintptr(hwnd), uintptr(MsgFilterMin), uintptr(MsgFilterMax), 0, 0)
	ret = int32(r0)
	if ret == -1 {
		if e1 != 0 {
			err = syscall.Errno(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func TranslateMessage(msg *Msg) (done bool) {
	r0, _, _ := syscall.Syscall(procTranslateMessage.Addr(), 1, uintptr(unsafe.Pointer(msg)), 0, 0)
	done = bool(r0 != 0)
	return
}

func DispatchMessage(msg *Msg) (ret int32) {
	r0, _, _ := syscall.Syscall(procDispatchMessageW.Addr(), 1, uintptr(unsafe.Pointer(msg)), 0, 0)
	ret = int32(r0)
	return
}

func SendMessage(hwnd uint32, msg uint32, wparam uintptr, lparam uintptr) (lresult int32) {
	r0, _, _ := syscall.Syscall6(procSendMessageW.Addr(), 4, uintptr(hwnd), uintptr(msg), uintptr(wparam), uintptr(lparam), 0, 0)
	lresult = int32(r0)
	return
}

func PostMessage(hwnd uint32, msg uint32, wparam uintptr, lparam uintptr) (err syscall.Errno) {
	r1, _, e1 := syscall.Syscall6(procPostMessageW.Addr(), 4, uintptr(hwnd), uintptr(msg), uintptr(wparam), uintptr(lparam), 0, 0)
	if int(r1) == 0 {
		if e1 != 0 {
			err = syscall.Errno(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}


// Utility functions for creating a testing window, etc

func extendHandlerMap(original, extras MsgHandlerMap) MsgHandlerMap {
	m := make(MsgHandlerMap, 32)
	for k, v := range original {
		m[k] = v
	}
	for k, v := range extras {
		m[k] = v
	}
	return m
}

func makeWindow(callbacks MsgHandlerMap) {

	wproc := syscall.NewCallback(func(hwnd, msg uint32, wparam uintptr, lparam uintptr) uintptr {
		if result, handled := ProcNoDefault(hwnd, msg, wparam, lparam); handled {
			return result
		}

		var rc interface{} = nil
		if cb := callbacks[msg]; cb != nil {
			rc = cb(hwnd)
		}

		// Handler provided a return code
		if rc != nil {
			if code, ok := rc.(int); !ok {
				panic("window msg response should be int")
			} else {
				return uintptr(code)
			}
		}

		// Handler did not provide a return code, use the default window procedure
		code := DefWindowProc(hwnd, msg, wparam, lparam)
		return uintptr(code)
	})

	// RegisterClassEx
	wcname := stringToUtf16Ptr("gohlTesting")

	if !classRegistered {
		var wc Wndclassex
		wc.Size = uint32(unsafe.Sizeof(wc))
		wc.WndProc = wproc
		wc.Instance = 0
		wc.Icon = 0
		wc.Cursor = 0
		wc.Background = 0
		wc.MenuName = nil
		wc.ClassName = wcname
		wc.IconSm = 0

		if _, errno := RegisterClassEx(&wc); errno != ERROR_SUCCESS {
			log.Panic(errno)
		}

		classRegistered = true
	}

	_, errno := CreateWindowEx(
		0,
		wcname,
		stringToUtf16Ptr("Gohl Test App"),
		0,
		0, 0, 20, 10,
		0, 0, 0, 0)
	if errno != ERROR_SUCCESS {
		log.Panic(errno)
	}
}

func pump() {
	var m Msg
	for {
		if r, errno := GetMessage(&m, 0, 0, 0); errno != ERROR_SUCCESS {
			panic(errno)
		} else if r == 0 {
			break
		}
		TranslateMessage(&m)
		DispatchMessage(&m)
	}
}

func testWithHtml(html string, test func(hwnd uint32)) {
	m := extendHandlerMap(defaultHandlerMap, MsgHandlerMap{
		WM_CREATE: func(hwnd uint32) interface{} {
			ret := defaultHandlerMap[WM_CREATE](hwnd)
			if err := LoadHtml(hwnd, []byte(html), ""); err != nil {
				log.Panic(err)
			}
			test(hwnd)
			PostMessage(hwnd, WM_CLOSE, 0, 0)
			return ret
		},
	})
	makeWindow(m)
	pump()
}




// Variables and types for testing

type MsgHandler func(uint32) interface{}
type MsgHandlerMap map[uint32]MsgHandler

var defaultHandlerMap = MsgHandlerMap{
	WM_CREATE: func(hwnd uint32) interface{} {
		//log.Print("WM_CREATE, hwnd = ", hwnd)
		AttachNotifyHandler(hwnd, notifyHandler)
		AttachWindowEventHandler(hwnd, windowEventHandler)
		return 0
	},
	WM_SHOWWINDOW: func(hwnd uint32) interface{} {
		return 0
	},
	WM_ERASEBKGND: func(hwnd uint32) interface{} {
		return 0
	},
	WM_CLOSE: func(hwnd uint32) interface{} {
		//log.Print("WM_CLOSE, hwnd = ", hwnd)
		DetachWindowEventHandler(hwnd)
		DetachNotifyHandler(hwnd)
		DestroyWindow(hwnd)
		return nil
	},
	WM_DESTROY: func(hwnd uint32) interface{} {
		//log.Print("WM_DESTROY, hwnd = ", hwnd)
		//DumpObjectCounts()
		PostQuitMessage(0)
		return 0
	},
	WM_QUIT: func(hwnd uint32) interface{} {
		log.Print("hai, quitting")
		return nil
	},
}

// Page templates used for various tests
var pages = map[string]string{
	"empty":    ``,
	"page":     `<html><body></body></html>`,
	"one-div":  `<div id="a">a</div>`,
	"two-divs": `<div id="a">a</div><div id="b">b</div>`,
	"nested-divs": `<div id="a">a<div id="b">b</div></div>`,
}

// Notify handler deals with WM_NOTIFY messages sent by htmlayout
var notifyHandler = &NotifyHandler{
	OnLoadData: func(params *NmhlLoadData) uintptr {
		relativePath := utf16ToString(params.Uri)
		log.Print("Load resource request: ", relativePath)
		return 0
	},
}

// Window event handler gets first and last chance to process events
var windowEventHandler = &EventHandler{}





// Tests:

func TestBasicWindow(t *testing.T) {
	m := extendHandlerMap(defaultHandlerMap, MsgHandlerMap{
		WM_CREATE: func(hwnd uint32) interface{} {
			ret := defaultHandlerMap[WM_CREATE](hwnd)
			PostMessage(hwnd, WM_CLOSE, 0, 0)
			return ret
		},
	})
	makeWindow(m)
	pump()
}

func TestLoadHtml(t *testing.T) {
	testWithHtml(pages["page"], func(hwnd uint32) {})
}

func TestLoadHtmlEmptyString(t *testing.T) {
	testWithHtml(pages["empty"], func(hwnd uint32) {})
}

func TestRootElement(t *testing.T) {
	testWithHtml(pages["one-div"], func(hwnd uint32) {
		if e := RootElement(hwnd); e == nil {
			t.Fatal("Could not get root elem")
		}
	})
}

func TestHandle(t *testing.T) {
	testWithHtml(pages["one-div"], func(hwnd uint32) {
		e := RootElement(hwnd)
		if h := e.Handle(); h == nil {
			t.Fatal("Handle was nil")
		}
	})
}

func TestRelease(t *testing.T) {
	testWithHtml(pages["one-div"], func(hwnd uint32) {
		e := RootElement(hwnd)
		e.Release()
		if h := e.Handle(); h != nil {
			t.Fatal("Released but handle is not nil, finalizer not called?")
		}
	})
}

func TestChildCount(t *testing.T) {
	testWithHtml(pages["two-divs"], func(hwnd uint32) {
		root := RootElement(hwnd)
		if count := root.ChildCount(); count != 2 {
			t.Fatal("Expected two divs as children")
		}
	})
}

func TestChildCount2(t *testing.T) {
	testWithHtml(pages["nested-divs"], func(hwnd uint32) {
		root := RootElement(hwnd)
		if count := root.ChildCount(); count != 1 {
			t.Fatal("Expected one divs as child")
		}
	})
}

func TestChild(t *testing.T) {
	testWithHtml(pages["two-divs"], func(hwnd uint32) {
		root := RootElement(hwnd)
		d1 := root.Child(0)
		d2 := root.Child(1)
		if d1 == nil || d2 == nil {
			t.Fatal("A child element could not be retrieved by index")
		}
	})
}

func TestIndex(t *testing.T) {
	testWithHtml(pages["two-divs"], func(hwnd uint32) {
		root := RootElement(hwnd)
		d1 := root.Child(0)
		d2 := root.Child(1)
		if d1.Index() != 0 {
			t.Fatal("Expected index 0")
		}
		if d2.Index() != 1 {
			t.Fatal("Expected index 1")
		}
	})
}

func TestEquals(t *testing.T) {
	testWithHtml(pages["one-div"], func(hwnd uint32) {
		root := RootElement(hwnd)
		d1 := root.Child(0)
		if root.Equals(d1) {
			t.Fatal("Distinct elems should not be equal")
		}
		if !d1.Equals(d1) {
			t.Fatal("Same elements should be equal")
		}
	})
}

func TestParent(t *testing.T) {
	testWithHtml(pages["nested-divs"], func(hwnd uint32) {
		root := RootElement(hwnd)
		d1 := root.Child(0)
		d2 := d1.Child(0)
		
		if !d2.Parent().Equals(d1) {
			t.Fatal("Parent was not the expected elem")
		}
		if !d1.Parent().Equals(root) {
			t.Fatal("Parent was not the expected elem")
		}
		if root.Parent() != nil {
			t.Fatal("Root's parent should be nil")
		}
	})
}



