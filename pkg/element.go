package htmlayout
/*
#cgo CFLAGS: -I../htmlayout/include
#cgo LDFLAGS: ../htmlayout/lib/HTMLayout.lib
#include <stdlib.h>
#include <htmlayout.h>
*/
import "C"

import (
	"fmt"
	"strconv"
	"unsafe"
	"utf16"
)

const (
	HLDOM_OK = C.HLDOM_OK
	HLDOM_INVALID_HWND = C.HLDOM_INVALID_HWND
	HLDOM_INVALID_HANDLE = C.HLDOM_INVALID_HANDLE
	HLDOM_PASSIVE_HANDLE = C.HLDOM_PASSIVE_HANDLE
	HLDOM_INVALID_PARAMETER = C.HLDOM_INVALID_PARAMETER
	HLDOM_OPERATION_FAILED = C.HLDOM_OPERATION_FAILED
	HLDOM_OK_NOT_HANDLED = C.int(-1)
)

var errorToString = map[C.HLDOM_RESULT]string {
	C.HLDOM_OK: "HLDOM_OK",
	C.HLDOM_INVALID_HWND: "HLDOM_INVALID_HWND",
	C.HLDOM_INVALID_HANDLE: "HLDOM_INVALID_HANDLE",
	C.HLDOM_PASSIVE_HANDLE: "HLDOM_PASSIVE_HANDLE",
	C.HLDOM_INVALID_PARAMETER: "HLDOM_INVALID_PARAMETER",
	C.HLDOM_OPERATION_FAILED: "HLDOM_OPERATION_FAILED",
	C.HLDOM_OK_NOT_HANDLED: "HLDOM_OK_NOT_HANDLED",
}

type DomError struct {
	Result C.HLDOM_RESULT
	Message string
}

func (self DomError) String() string {
	return fmt.Sprintf( "%s: %s", errorToString[self.Result], self.Message )
}

func domPanic(result C.HLDOM_RESULT, message string) {
	panic(DomError{result, message})
}

// StringToUTF16 returns the UTF-16 encoding of the UTF-8 string s,
// with a terminating NUL added.
func StringToUTF16(s string) []uint16 { return utf16.Encode([]int(s + "\x00")) }

// UTF16ToString returns the UTF-8 encoding of the UTF-16 sequence s,
// with a terminating NUL removed.
func UTF16ToString(s *uint16) string {
	if s == nil {
		panic("null cstring")
	}
	us := make([]uint16, 0, 256) 
	for p := uintptr(unsafe.Pointer(s)); ; p += 2 { 
		u := *(*uint16)(unsafe.Pointer(p)) 
		if u == 0 { 
			return string(utf16.Decode(us)) 
		} 
		us = append(us, u) 
	} 
	return ""
}

// StringToUTF16Ptr returns pointer to the UTF-16 encoding of
// the UTF-8 string s, with a terminating NUL added.
func StringToUTF16Ptr(s string) *uint16 { return &StringToUTF16(s)[0] }




type Handle C.HELEMENT

func use(handle Handle) {
	if dr := C.HTMLayout_UseElement(handle); dr != HLDOM_OK {
		domPanic(dr, "UseElement");
	}
}

func unuse(handle Handle) {
	if handle != nil {
		if dr := C.HTMLayout_UnuseElement(handle); dr != HLDOM_OK {
			domPanic(dr, "UnuseElement");
		}
	}
}


/*
Element

Represents a single DOM element, owns and manages a Handle
*/
type Element struct {
	handle Handle
}

func (e *Element) set(h Handle) {
	use(h)
	unuse(e.handle)
	e.handle = h
}

func (e *Element) Release() {
	unuse(e.handle)
	e.handle = nil
}

func (e *Element) GetHandle() Handle {
	return e.handle
}


// HTML attribute accessors/modifiers:

func (e *Element) GetAttr(key string) *string {
	szValue := (*C.WCHAR)(nil)
	szKey := C.CString(key)
	ret := C.HTMLayoutGetAttributeByName(e.handle, (*C.CHAR)(szKey), (*C.LPCWSTR)(&szValue))
	C.free(unsafe.Pointer(szKey))
	if ret != HLDOM_OK {
		domPanic(ret, "failed to get attribute: "+key)
	}
	if szValue != nil {
		s := UTF16ToString((*uint16)(szValue))
		return &s
	}
	return nil;
}

func (e *Element) GetAttrAsFloat(key string) *float32 {
	if s := e.GetAttr(key); s != nil {
		if f, err := strconv.Atof32(*s); err != nil {
			panic(err)
		} else {
			return &f
		}
	}
	return nil
}

func (e *Element) GetAttrAsInt(key string) *int {
	if s := e.GetAttr(key); s != nil {
		if i, err := strconv.Atoi(*s); err != nil {
			panic(err)
		} else {
			return &i
		}
	}
	return nil
}

func (e *Element) SetAttr(key string, value interface{}) {
	szKey := C.CString(key)
	var ret C.HLDOM_RESULT = HLDOM_OK
	if v, ok := value.(string); ok {
		ret = C.HTMLayoutSetAttributeByName(e.handle, (*C.CHAR)(szKey), (*C.WCHAR)(StringToUTF16Ptr(v)))
	} else if v, ok := value.(float32); ok {
		ret = C.HTMLayoutSetAttributeByName(e.handle, (*C.CHAR)(szKey), (*C.WCHAR)(StringToUTF16Ptr(strconv.Ftoa32(v, 'e', 6))))
	} else if v, ok := value.(int); ok {
		ret = C.HTMLayoutSetAttributeByName(e.handle, (*C.CHAR)(szKey), (*C.WCHAR)(StringToUTF16Ptr(strconv.Itoa(v))))
	} else if value == nil {
		ret = C.HTMLayoutSetAttributeByName(e.handle, (*C.CHAR)(szKey), nil)
	} else {
		panic("don't know how to format this argument type")
	}
	C.free(unsafe.Pointer(szKey))
	if ret != HLDOM_OK {
		domPanic(ret, "failed to set attribute: "+key)
	}
}

func (e *Element) RemoveAttr(key string) {
	e.SetAttr(key, nil)
}

func (e *Element) GetAttrValueByIndex(index int) string {
	szValue := (*C.WCHAR)(nil)
	if ret := C.HTMLayoutGetNthAttribute(e.handle, (C.UINT)(index), nil, (*C.LPCWSTR)(&szValue)); ret != HLDOM_OK {
		domPanic(ret, fmt.Sprintf("failed to get attribute name by index: %d", index))
	}
	return UTF16ToString((*uint16)(szValue))
}

func (e *Element) GetAttrNameByIndex(index int) string {
	szName := (*C.CHAR)(nil)
	if ret := C.HTMLayoutGetNthAttribute(e.handle, (C.UINT)(index), (*C.LPCSTR)(&szName), nil); ret != HLDOM_OK {
		domPanic(ret, fmt.Sprintf("failed to get attribute name by index: %d", index))	
	}
	return C.GoString((*C.char)(szName))
}

func (e *Element) GetAttrCount(index int) int {
	var count C.UINT = 0
	if ret := C.HTMLayoutGetAttributeCount(e.handle, &count); ret != HLDOM_OK {
		domPanic(ret, "failed to get attribute count")
	}
	return int(count)
}



// CSS style attribute accessors/mutators

func (e *Element) GetStyle(key string) *string {
	szValue := (*C.WCHAR)(nil)
	szKey := C.CString(key)
	ret := C.HTMLayoutGetStyleAttribute(e.handle, (*C.CHAR)(szKey), (*C.LPCWSTR)(&szValue))
	C.free(unsafe.Pointer(szKey))
	if ret != HLDOM_OK {
		domPanic(ret, "failed to get style: "+key)
	}
	if szValue != nil {
		s := UTF16ToString((*uint16)(szValue))
		return &s
	}
	return nil;
}

func (e *Element) GetStyleAsFloat(key string) *float32 {
	if s := e.GetStyle(key); s != nil {
		if f, err := strconv.Atof32(*s); err != nil {
			panic(err)
		} else {
			return &f
		}
	}
	return nil
}

func (e *Element) GetStyleAsInt(key string) *int {
	if s := e.GetStyle(key); s != nil {
		if i, err := strconv.Atoi(*s); err != nil {
			panic(err)
		} else {
			return &i
		}
	}
	return nil
}

func (e *Element) SetStyle(key string, value interface{}) {
	szKey := C.CString(key)
	var ret C.HLDOM_RESULT = HLDOM_OK
	if v, ok := value.(string); ok {
		ret = C.HTMLayoutSetStyleAttribute(e.handle, (*C.CHAR)(szKey), (*C.WCHAR)(StringToUTF16Ptr(v)))
	} else if v, ok := value.(float32); ok {
		ret = C.HTMLayoutSetStyleAttribute(e.handle, (*C.CHAR)(szKey), (*C.WCHAR)(StringToUTF16Ptr(strconv.Ftoa32(v, 'e', 6))))
	} else if v, ok := value.(int); ok {
		ret = C.HTMLayoutSetStyleAttribute(e.handle, (*C.CHAR)(szKey), (*C.WCHAR)(StringToUTF16Ptr(strconv.Itoa(v))))
	} else if value == nil {
		ret = C.HTMLayoutSetStyleAttribute(e.handle, (*C.CHAR)(szKey), nil)
	} else {
		panic("don't know how to format this argument type")
	}
	C.free(unsafe.Pointer(szKey))
	if ret != HLDOM_OK {
		domPanic(ret, "failed to set style: "+key)
	}
}

func (e *Element) RemoveStyle(key string) {
	e.SetStyle(key, nil)
}

func (e *Element) ClearStyles(key string) {
	if ret := C.HTMLayoutSetStyleAttribute(e.handle, nil, nil); ret != HLDOM_OK {
		domPanic(ret, "failed to clear all styles")
	}
}



// Constructor

func NewElement(h Handle) *Element {
	e := &Element{nil}
	e.set(h)
	return e
}






