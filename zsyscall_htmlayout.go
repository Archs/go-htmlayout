// mksyscall_dll htmlayout.go element.go
// MACHINE GENERATED BY THE COMMAND ABOVE; DO NOT EDIT

package gohl

import (
	"syscall"
	"unsafe"
)

var (
	modhtmlayout = syscall.NewLazyDLL("htmlayout.dll")

	procHTMLayoutProcND = modhtmlayout.NewProc("HTMLayoutProcND")
	procHTMLayoutLoadHtmlEx = modhtmlayout.NewProc("HTMLayoutLoadHtmlEx")
	procHTMLayoutLoadFile = modhtmlayout.NewProc("HTMLayoutLoadFile")
	procHTMLayoutDataReady = modhtmlayout.NewProc("HTMLayoutDataReady")
	procHTMLayoutAttachEventHandler = modhtmlayout.NewProc("HTMLayoutAttachEventHandler")
	procHTMLayoutDetachEventHandler = modhtmlayout.NewProc("HTMLayoutDetachEventHandler")
	procHTMLayoutAttachEventHandlerEx = modhtmlayout.NewProc("HTMLayoutAttachEventHandlerEx")
	procHTMLayoutWindowAttachEventHandler = modhtmlayout.NewProc("HTMLayoutWindowAttachEventHandler")
	procHTMLayoutWindowDetachEventHandler = modhtmlayout.NewProc("HTMLayoutWindowDetachEventHandler")
	procHTMLayoutSetCallback = modhtmlayout.NewProc("HTMLayoutSetCallback")
	procHTMLayout_UseElement = modhtmlayout.NewProc("HTMLayout_UseElement")
	procHTMLayout_UnuseElement = modhtmlayout.NewProc("HTMLayout_UnuseElement")
	procHTMLayoutCreateElement = modhtmlayout.NewProc("HTMLayoutCreateElement")
	procHTMLayoutGetRootElement = modhtmlayout.NewProc("HTMLayoutGetRootElement")
	procHTMLayoutGetFocusElement = modhtmlayout.NewProc("HTMLayoutGetFocusElement")
	procHTMLayoutUpdateElementEx = modhtmlayout.NewProc("HTMLayoutUpdateElementEx")
	procHTMLayoutSetCapture = modhtmlayout.NewProc("HTMLayoutSetCapture")
	procHTMLayoutSelectElementsW = modhtmlayout.NewProc("HTMLayoutSelectElementsW")
	procHTMLayoutSelectParentW = modhtmlayout.NewProc("HTMLayoutSelectParentW")
	procHTMLayoutSendEvent = modhtmlayout.NewProc("HTMLayoutSendEvent")
	procHTMLayoutPostEvent = modhtmlayout.NewProc("HTMLayoutPostEvent")
	procHTMLayoutGetChildrenCount = modhtmlayout.NewProc("HTMLayoutGetChildrenCount")
	procHTMLayoutGetNthChild = modhtmlayout.NewProc("HTMLayoutGetNthChild")
	procHTMLayoutGetElementIndex = modhtmlayout.NewProc("HTMLayoutGetElementIndex")
	procHTMLayoutGetParentElement = modhtmlayout.NewProc("HTMLayoutGetParentElement")
	procHTMLayoutInsertElement = modhtmlayout.NewProc("HTMLayoutInsertElement")
	procHTMLayoutDetachElement = modhtmlayout.NewProc("HTMLayoutDetachElement")
	procHTMLayoutDeleteElement = modhtmlayout.NewProc("HTMLayoutDeleteElement")
	procHTMLayoutCloneElement = modhtmlayout.NewProc("HTMLayoutCloneElement")
	procHTMLayoutSwapElements = modhtmlayout.NewProc("HTMLayoutSwapElements")
	procHTMLayoutSortElements = modhtmlayout.NewProc("HTMLayoutSortElements")
	procHTMLayoutSetTimer = modhtmlayout.NewProc("HTMLayoutSetTimer")
	procHTMLayoutGetElementHwnd = modhtmlayout.NewProc("HTMLayoutGetElementHwnd")
	procHTMLayoutGetElementHtml = modhtmlayout.NewProc("HTMLayoutGetElementHtml")
	procHTMLayoutGetElementType = modhtmlayout.NewProc("HTMLayoutGetElementType")
	procHTMLayoutSetElementHtml = modhtmlayout.NewProc("HTMLayoutSetElementHtml")
	procHTMLayoutSetElementInnerText = modhtmlayout.NewProc("HTMLayoutSetElementInnerText")
	procHTMLayoutGetElementInnerText = modhtmlayout.NewProc("HTMLayoutGetElementInnerText")
	procHTMLayoutGetAttributeByName = modhtmlayout.NewProc("HTMLayoutGetAttributeByName")
	procHTMLayoutSetAttributeByName = modhtmlayout.NewProc("HTMLayoutSetAttributeByName")
	procHTMLayoutGetNthAttribute = modhtmlayout.NewProc("HTMLayoutGetNthAttribute")
	procHTMLayoutGetAttributeCount = modhtmlayout.NewProc("HTMLayoutGetAttributeCount")
	procHTMLayoutGetStyleAttribute = modhtmlayout.NewProc("HTMLayoutGetStyleAttribute")
	procHTMLayoutSetStyleAttribute = modhtmlayout.NewProc("HTMLayoutSetStyleAttribute")
	procHTMLayoutGetElementState = modhtmlayout.NewProc("HTMLayoutGetElementState")
	procHTMLayoutSetElementState = modhtmlayout.NewProc("HTMLayoutSetElementState")
	procHTMLayoutMoveElement = modhtmlayout.NewProc("HTMLayoutMoveElement")
	procHTMLayoutMoveElementEx = modhtmlayout.NewProc("HTMLayoutMoveElementEx")
	procHTMLayoutGetElementLocation = modhtmlayout.NewProc("HTMLayoutGetElementLocation")
	procHTMLayoutCallBehaviorMethod = modhtmlayout.NewProc("HTMLayoutCallBehaviorMethod")
)

func HTMLayoutProcND(hwnd HWND, msg uint32, wParam uintptr, lParam uintptr, pbHandled *BOOL) (ret LRESULT) {
	r0, _, _ := syscall.Syscall6(procHTMLayoutProcND.Addr(), 5, uintptr(hwnd), uintptr(msg), uintptr(wParam), uintptr(lParam), uintptr(unsafe.Pointer(pbHandled)), 0)
	ret = LRESULT(r0)
	return
}

func HTMLayoutLoadHtmlEx(hWndHTMLayout HWND, html []byte, htmlSize UINT, baseUrl *uint16) (ret BOOL) {
	var _p0 *byte
	if len(html) > 0 {
		_p0 = &html[0]
	}
	r0, _, _ := syscall.Syscall6(procHTMLayoutLoadHtmlEx.Addr(), 5, uintptr(hWndHTMLayout), uintptr(unsafe.Pointer(_p0)), uintptr(len(html)), uintptr(htmlSize), uintptr(unsafe.Pointer(baseUrl)), 0)
	ret = BOOL(r0)
	return
}

func HTMLayoutLoadFile(hWndHTMLayout HWND, filename *uint16) (ret int, err error) {
	r0, _, e1 := syscall.Syscall(procHTMLayoutLoadFile.Addr(), 2, uintptr(hWndHTMLayout), uintptr(unsafe.Pointer(filename)), 0)
	ret = int(r0)
	if ret == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func HTMLayoutDataReady(hwnd HWND, uri *uint16, data []byte, dataLength int32) (ret int, err error) {
	var _p0 *byte
	if len(data) > 0 {
		_p0 = &data[0]
	}
	r0, _, e1 := syscall.Syscall6(procHTMLayoutDataReady.Addr(), 5, uintptr(hwnd), uintptr(unsafe.Pointer(uri)), uintptr(unsafe.Pointer(_p0)), uintptr(len(data)), uintptr(dataLength), 0)
	ret = int(r0)
	if ret == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func HTMLayoutAttachEventHandler(he HELEMENT, pep uintptr, tag uintptr) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutAttachEventHandler.Addr(), 3, uintptr(he), uintptr(pep), uintptr(tag))
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutDetachEventHandler(he HELEMENT, pep uintptr, tag uintptr) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutDetachEventHandler.Addr(), 3, uintptr(he), uintptr(pep), uintptr(tag))
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutAttachEventHandlerEx(he HELEMENT, pep uintptr, tag uintptr, subscription uint) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall6(procHTMLayoutAttachEventHandlerEx.Addr(), 4, uintptr(he), uintptr(pep), uintptr(tag), uintptr(subscription), 0, 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutWindowAttachEventHandler(hwndLayout HWND, pep uintptr, tag uintptr, subscription uint) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall6(procHTMLayoutWindowAttachEventHandler.Addr(), 4, uintptr(hwndLayout), uintptr(pep), uintptr(tag), uintptr(subscription), 0, 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutWindowDetachEventHandler(hwndLayout HWND, pep uintptr, tag uintptr) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutWindowDetachEventHandler.Addr(), 3, uintptr(hwndLayout), uintptr(pep), uintptr(tag))
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutSetCallback(hWndHTMLayout uintptr, cb uintptr, cbParam uintptr) {
	syscall.Syscall(procHTMLayoutSetCallback.Addr(), 3, uintptr(hWndHTMLayout), uintptr(cb), uintptr(cbParam))
	return
}

func HTMLayout_UseElement(he HELEMENT) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayout_UseElement.Addr(), 1, uintptr(he), 0, 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayout_UnuseElement(he HELEMENT) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayout_UnuseElement.Addr(), 1, uintptr(he), 0, 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutCreateElement(tagname string, textOrNull *uint16, phe *HELEMENT) (ret HLDOM_RESULT, err error) {
	var _p0 *byte
	_p0, err = syscall.BytePtrFromString(tagname)
	if err != nil {
		return
	}
	r0, _, e1 := syscall.Syscall(procHTMLayoutCreateElement.Addr(), 3, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(textOrNull)), uintptr(unsafe.Pointer(phe)))
	ret = HLDOM_RESULT(r0)
	if ret != 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func HTMLayoutGetRootElement(hwnd HWND, pheT *HELEMENT) (ret HLDOM_RESULT, err error) {
	r0, _, e1 := syscall.Syscall(procHTMLayoutGetRootElement.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(pheT)), 0)
	ret = HLDOM_RESULT(r0)
	if ret != 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func HTMLayoutGetFocusElement(hwnd HWND, pheT *HELEMENT) (ret HLDOM_RESULT, err error) {
	r0, _, e1 := syscall.Syscall(procHTMLayoutGetFocusElement.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(pheT)), 0)
	ret = HLDOM_RESULT(r0)
	if ret != 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func HTMLayoutUpdateElementEx(he HELEMENT, flags uint) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutUpdateElementEx.Addr(), 2, uintptr(he), uintptr(flags), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutSetCapture(he HELEMENT) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutSetCapture.Addr(), 1, uintptr(he), 0, 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutSelectElementsW(he HELEMENT, CSS_selectors string, callback uintptr, param uintptr) (ret HLDOM_RESULT, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(CSS_selectors)
	if err != nil {
		return
	}
	r0, _, e1 := syscall.Syscall6(procHTMLayoutSelectElementsW.Addr(), 4, uintptr(he), uintptr(unsafe.Pointer(_p0)), uintptr(callback), uintptr(param), 0, 0)
	ret = HLDOM_RESULT(r0)
	if ret != 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func HTMLayoutSelectParentW(he HELEMENT, selector string, depth uint, heFound *HELEMENT) (ret HLDOM_RESULT, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(selector)
	if err != nil {
		return
	}
	r0, _, e1 := syscall.Syscall6(procHTMLayoutSelectParentW.Addr(), 4, uintptr(he), uintptr(unsafe.Pointer(_p0)), uintptr(depth), uintptr(unsafe.Pointer(heFound)), 0, 0)
	ret = HLDOM_RESULT(r0)
	if ret != 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func HTMLayoutSendEvent(he HELEMENT, appEventCode uint, heSource HELEMENT, reason *uint, handled *BOOL) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall6(procHTMLayoutSendEvent.Addr(), 5, uintptr(he), uintptr(appEventCode), uintptr(heSource), uintptr(unsafe.Pointer(reason)), uintptr(unsafe.Pointer(handled)), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutPostEvent(he HELEMENT, appEventCode uint, heSource HELEMENT, reason *uint) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall6(procHTMLayoutPostEvent.Addr(), 4, uintptr(he), uintptr(appEventCode), uintptr(heSource), uintptr(unsafe.Pointer(reason)), 0, 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutGetChildrenCount(he HELEMENT, count *uint) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutGetChildrenCount.Addr(), 2, uintptr(he), uintptr(unsafe.Pointer(count)), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutGetNthChild(he HELEMENT, n uint, phe *HELEMENT) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutGetNthChild.Addr(), 3, uintptr(he), uintptr(n), uintptr(unsafe.Pointer(phe)))
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutGetElementIndex(he HELEMENT, p_index *uint) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutGetElementIndex.Addr(), 2, uintptr(he), uintptr(unsafe.Pointer(p_index)), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutGetParentElement(he HELEMENT, p_parent_he *HELEMENT) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutGetParentElement.Addr(), 2, uintptr(he), uintptr(unsafe.Pointer(p_parent_he)), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutInsertElement(he HELEMENT, hparent HELEMENT, index uint) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutInsertElement.Addr(), 3, uintptr(he), uintptr(hparent), uintptr(index))
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutDetachElement(he HELEMENT) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutDetachElement.Addr(), 1, uintptr(he), 0, 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutDeleteElement(he HELEMENT) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutDeleteElement.Addr(), 1, uintptr(he), 0, 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutCloneElement(he HELEMENT, phe *HELEMENT) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutCloneElement.Addr(), 2, uintptr(he), uintptr(unsafe.Pointer(phe)), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutSwapElements(he HELEMENT, other HELEMENT) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutSwapElements.Addr(), 2, uintptr(he), uintptr(other), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutSortElements(he HELEMENT, firstIndex uint, lastIndex uint, cmpFunc uintptr, cmpFuncParam uintptr) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall6(procHTMLayoutSortElements.Addr(), 5, uintptr(he), uintptr(firstIndex), uintptr(lastIndex), uintptr(cmpFunc), uintptr(cmpFuncParam), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutSetTimer(he HELEMENT, milliseconds uint) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutSetTimer.Addr(), 2, uintptr(he), uintptr(milliseconds), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutGetElementHwnd(he HELEMENT, p_hwnd *HWND, rootWindow BOOL) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutGetElementHwnd.Addr(), 3, uintptr(he), uintptr(unsafe.Pointer(p_hwnd)), uintptr(rootWindow))
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutGetElementHtml(he HELEMENT, utf8bytes uintptr, outer BOOL) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutGetElementHtml.Addr(), 3, uintptr(he), uintptr(utf8bytes), uintptr(outer))
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutGetElementType(he HELEMENT, p_type uintptr) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutGetElementType.Addr(), 2, uintptr(he), uintptr(p_type), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutSetElementHtml(he HELEMENT, html string, htmlLength int, where uint) (ret HLDOM_RESULT, err error) {
	var _p0 *byte
	_p0, err = syscall.BytePtrFromString(html)
	if err != nil {
		return
	}
	r0, _, e1 := syscall.Syscall6(procHTMLayoutSetElementHtml.Addr(), 4, uintptr(he), uintptr(unsafe.Pointer(_p0)), uintptr(htmlLength), uintptr(where), 0, 0)
	ret = HLDOM_RESULT(r0)
	if ret != 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func HTMLayoutSetElementInnerText(he HELEMENT, text string, length uint) (ret HLDOM_RESULT, err error) {
	var _p0 *byte
	_p0, err = syscall.BytePtrFromString(text)
	if err != nil {
		return
	}
	r0, _, e1 := syscall.Syscall(procHTMLayoutSetElementInnerText.Addr(), 3, uintptr(he), uintptr(unsafe.Pointer(_p0)), uintptr(length))
	ret = HLDOM_RESULT(r0)
	if ret != 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func HTMLayoutGetElementInnerText(he HELEMENT, utf8bytes uintptr) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutGetElementInnerText.Addr(), 2, uintptr(he), uintptr(utf8bytes), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutGetAttributeByName(he HELEMENT, utf8bytes *byte, p_value uintptr) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutGetAttributeByName.Addr(), 3, uintptr(he), uintptr(unsafe.Pointer(utf8bytes)), uintptr(p_value))
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutSetAttributeByName(he HELEMENT, name *byte, value *uint16) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutSetAttributeByName.Addr(), 3, uintptr(he), uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(value)))
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutGetNthAttribute(he HELEMENT, n uint, name uintptr, value uintptr) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall6(procHTMLayoutGetNthAttribute.Addr(), 4, uintptr(he), uintptr(n), uintptr(name), uintptr(value), 0, 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutGetAttributeCount(he HELEMENT, p_count *uint) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutGetAttributeCount.Addr(), 2, uintptr(he), uintptr(unsafe.Pointer(p_count)), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutGetStyleAttribute(he HELEMENT, name uintptr, p_value uintptr) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutGetStyleAttribute.Addr(), 3, uintptr(he), uintptr(name), uintptr(p_value))
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutSetStyleAttribute(he HELEMENT, name uintptr, value uintptr) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutSetStyleAttribute.Addr(), 3, uintptr(he), uintptr(name), uintptr(value))
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutGetElementState(he HELEMENT, pstateBits *uint) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutGetElementState.Addr(), 2, uintptr(he), uintptr(unsafe.Pointer(pstateBits)), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutSetElementState(he HELEMENT, stateBitsToSet uint, stateBitsToClear uint, updateView BOOL) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall6(procHTMLayoutSetElementState.Addr(), 4, uintptr(he), uintptr(stateBitsToSet), uintptr(stateBitsToClear), uintptr(updateView), 0, 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutMoveElement(he HELEMENT, xView int, yView int) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutMoveElement.Addr(), 3, uintptr(he), uintptr(xView), uintptr(yView))
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutMoveElementEx(he HELEMENT, xView int, yView int, width int, height int) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall6(procHTMLayoutMoveElementEx.Addr(), 5, uintptr(he), uintptr(xView), uintptr(yView), uintptr(width), uintptr(height), 0)
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutGetElementLocation(he HELEMENT, p_location *Rect, areas uint) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutGetElementLocation.Addr(), 3, uintptr(he), uintptr(unsafe.Pointer(p_location)), uintptr(areas))
	ret = HLDOM_RESULT(r0)
	return
}

func HTMLayoutCallBehaviorMethod(he HELEMENT, params *METHOD_PARAMS) (ret HLDOM_RESULT) {
	r0, _, _ := syscall.Syscall(procHTMLayoutCallBehaviorMethod.Addr(), 2, uintptr(he), uintptr(unsafe.Pointer(params)), 0)
	ret = HLDOM_RESULT(r0)
	return
}
