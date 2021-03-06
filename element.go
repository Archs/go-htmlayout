package gohl

// #include <stdlib.h>
// #include <htmlayout.h>
import "C"

import (
	"errors"
	"github.com/lxn/win"
	"reflect"
	"syscall"

	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf16"
	"unsafe"
)

const (
	HLDOM_OK                = C.HLDOM_OK
	HLDOM_INVALID_HWND      = C.HLDOM_INVALID_HWND
	HLDOM_INVALID_HANDLE    = C.HLDOM_INVALID_HANDLE
	HLDOM_PASSIVE_HANDLE    = C.HLDOM_PASSIVE_HANDLE
	HLDOM_INVALID_PARAMETER = C.HLDOM_INVALID_PARAMETER
	HLDOM_OPERATION_FAILED  = C.HLDOM_OPERATION_FAILED
	HLDOM_OK_NOT_HANDLED    = C.HLDOM_OK_NOT_HANDLED

	HV_OK_TRUE           = 0xffffffff
	HV_OK                = C.HV_OK
	HV_BAD_PARAMETER     = C.HV_BAD_PARAMETER
	HV_INCOMPATIBLE_TYPE = C.HV_INCOMPATIBLE_TYPE

	STATE_LINK       = 0x00000001 // selector :link,    any element having href attribute
	STATE_HOVER      = 0x00000002 // selector :hover,   element is under the cursor, mouse hover
	STATE_ACTIVE     = 0x00000004 // selector :active,  element is activated, e.g. pressed
	STATE_FOCUS      = 0x00000008 // selector :focus,   element is in focus
	STATE_VISITED    = 0x00000010 // selector :visited, aux flag - not used internally now.
	STATE_CURRENT    = 0x00000020 // selector :current, current item in collection, e.g. current <option> in <select>
	STATE_CHECKED    = 0x00000040 // selector :checked, element is checked (or selected), e.g. check box or itme in multiselect
	STATE_DISABLED   = 0x00000080 // selector :disabled, element is disabled, behavior related flag.
	STATE_READONLY   = 0x00000100 // selector :read-only, element is read-only, behavior related flag.
	STATE_EXPANDED   = 0x00000200 // selector :expanded, element is in expanded state - nodes in tree view e.g. <options> in <select>
	STATE_COLLAPSED  = 0x00000400 // selector :collapsed, mutually exclusive with EXPANDED
	STATE_INCOMPLETE = 0x00000800 // selector :incomplete, element has images (back/fore/bullet) requested but not delivered.
	STATE_ANIMATING  = 0x00001000 // selector :animating, is currently animating
	STATE_FOCUSABLE  = 0x00002000 // selector :focusable, shall accept focus
	STATE_ANCHOR     = 0x00004000 // selector :anchor, first element in selection (<select miltiple>), STATE_CURRENT is the current.
	STATE_SYNTHETIC  = 0x00008000 // selector :synthetic, synthesized DOM elements - e.g. all missed cells in tables (<td>) are getting this flag
	STATE_OWNS_POPUP = 0x00010000 // selector :owns-popup, anchor(owner) element of visible popup.
	STATE_TABFOCUS   = 0x00020000 // selector :tab-focus, element got focus by tab traversal. engine set it together with :focus.
	STATE_EMPTY      = 0x00040000 // selector :empty - element is empty.
	STATE_BUSY       = 0x00080000 // selector :busy, element is busy. HTMLayoutRequestElementData will set this flag if
	// external data was requested for the element. When data will be delivered engine will reset this flag on the element.

	STATE_DRAG_OVER   = 0x00100000 // drag over the block that can accept it (so is current drop target). Flag is set for the drop target block. At any given moment of time it can be only one such block.
	STATE_DROP_TARGET = 0x00200000 // active drop target. Multiple elements can have this flag when D&D is active.
	STATE_MOVING      = 0x00400000 // dragging/moving - the flag is set for the moving element (copy of the drag-source).
	STATE_COPYING     = 0x00800000 // dragging/copying - the flag is set for the copying element (copy of the drag-source).
	STATE_DRAG_SOURCE = 0x00C00000 // is set in element that is being dragged.

	STATE_POPUP   = 0x40000000 // this element is in popup state and presented to the user - out of flow now
	STATE_PRESSED = 0x04000000 // pressed - close to active but has wider life span - e.g. in MOUSE_UP it
	// is still on, so behavior can check it in MOUSE_UP to discover CLICK condition.
	STATE_HAS_CHILDREN = 0x02000000 // has more than one child.
	STATE_HAS_CHILD    = 0x01000000 // has single child.

	STATE_IS_LTR = 0x20000000 // selector :ltr, the element or one of its nearest container has @dir and that dir has "ltr" value
	STATE_IS_RTL = 0x10000000 // selector :rtl, the element or one of its nearest container has @dir and that dir has "rtl" value

	RESET_STYLE_THIS = 0x0020 // reset styles - this may require if you have styles dependent from attributes,
	RESET_STYLE_DEEP = 0x0010 // use these flags after SetAttribute then. RESET_STYLE_THIS is faster than RESET_STYLE_DEEP.
	MEASURE_INPLACE  = 0x0001 // use this flag if you do not expect any dimensional changes - this is faster than REMEASURE
	MEASURE_DEEP     = 0x0002 // use this flag if changes of some attributes/content may cause change of dimensions of the element
	REDRAW_NOW       = 0x8000
)

var (
	BAD_HELEMENT = HELEMENT(unsafe.Pointer(uintptr(0)))
)

var errorToString = map[HLDOM_RESULT]string{
	HLDOM_OK:                "HLDOM_OK",
	HLDOM_INVALID_HWND:      "HLDOM_INVALID_HWND",
	HLDOM_INVALID_HANDLE:    "HLDOM_INVALID_HANDLE",
	HLDOM_PASSIVE_HANDLE:    "HLDOM_PASSIVE_HANDLE",
	HLDOM_INVALID_PARAMETER: "HLDOM_INVALID_PARAMETER",
	HLDOM_OPERATION_FAILED:  "HLDOM_OPERATION_FAILED",
	HLDOM_OK_NOT_HANDLED:    "HLDOM_OK_NOT_HANDLED",
}

var valueErrorToString = map[VALUE_RESULT]string{
	HV_OK_TRUE:           "HV_OK_TRUE",
	HV_OK:                "HV_OK",
	HV_BAD_PARAMETER:     "HV_BAD_PARAMETER",
	HV_INCOMPATIBLE_TYPE: "HV_INCOMPATIBLE_TYPE",
}

var whitespaceSplitter = regexp.MustCompile(`(\S+)`)

// DomError represents an htmlayout error with an associated
// dom error code
type DomError struct {
	Result  HLDOM_RESULT
	Message string
}

func (e *DomError) Error() string {
	return fmt.Sprintf("%s: %s", errorToString[e.Result], e.Message)
}

func domResultAsString(result HLDOM_RESULT) string {
	return errorToString[result]
}

func domPanic(result C.HLDOM_RESULT, message ...interface{}) {
	panic(&DomError{HLDOM_RESULT(result), fmt.Sprint(message...)})
}

func domPanic2(result HLDOM_RESULT, message ...interface{}) {
	panic(&DomError{result, fmt.Sprint(message...)})
}

type ValueError struct {
	Result  VALUE_RESULT
	Message string
}

func (e *ValueError) Error() string {
	return fmt.Sprintf("%s: %s", valueErrorToString[e.Result], e.Message)
}

func valuePanic(result C.UINT, message ...interface{}) {
	panic(&ValueError{VALUE_RESULT(result), fmt.Sprint(message...)})
}

// Returns the utf-16 encoding of the utf-8 string s,
// with a terminating NUL added.
func stringToUtf16(s string) []uint16 {
	return utf16.Encode([]rune(s + "\x00"))
}

// Returns the utf-8 encoding of the utf-16 sequence s,
// with a terminating NUL removed.
func utf16ToString(s *uint16) string {
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

func bytePtrToString(s *byte) string {
	if s == nil {
		panic("null cstring")
	}
	bs := make([]byte, 0, 256)
	for p := uintptr(unsafe.Pointer(s)); ; p += 1 {
		b := *(*byte)(unsafe.Pointer(p))
		if b == 0 {
			return string(bs)
		}
		bs = append(bs, b)
	}
	return ""
}

func utf16ToStringLength(s *uint16, length int) string {
	if s == nil {
		panic("null cstring")
	}
	us := make([]uint16, 0, 256)
	for p, i := uintptr(unsafe.Pointer(s)), 0; i < length; p, i = p+2, i+1 {
		u := *(*uint16)(unsafe.Pointer(p))
		us = append(us, u)
	}
	return string(utf16.Decode(us))
}

// Returns pointer to the utf-16 encoding of
// the utf-8 string s, with a terminating NUL added.
func stringToUtf16Ptr(s string) *uint16 {
	return &stringToUtf16(s)[0]
}

/**Marks DOM object as used (a.k.a. AddRef).
 * \param[in] he \b #HELEMENT
 * \return \b #HLDOM_RESULT
 * Application should call this function before using element handle. If the
 * application fails to do that calls to other DOM functions for this handle
 * may result in an error.
 *
 * \sa #HTMLayout_UnuseElement()
 **/
// EXTERN_C HLDOM_RESULT HLAPI HTMLayout_UseElement(HELEMENT he);
//sys HTMLayout_UseElement(he HELEMENT) (ret HLDOM_RESULT) = htmlayout.HTMLayout_UseElement
func use(handle HELEMENT) {
	if dr := HTMLayout_UseElement(handle); dr != HLDOM_OK {
		domPanic2(dr, "UseElement")
	}
}

/**Marks DOM object as unused (a.k.a. Release).
 * Get handle of every element's child element.
 * \param[in] he \b #HELEMENT
 * \return \b #HLDOM_RESULT
 *
 * Application should call this function when it does not need element's
 * handle anymore.
 * \sa #HTMLayout_UseElement()
 **/
// EXTERN_C HLDOM_RESULT HLAPI HTMLayout_UnuseElement(HELEMENT he);
//sys HTMLayout_UnuseElement(he HELEMENT) (ret HLDOM_RESULT) = htmlayout.HTMLayout_UnuseElement
func unuse(handle HELEMENT) {
	if handle != 0 {
		if dr := HTMLayout_UnuseElement(handle); dr != HLDOM_OK {
			domPanic2(dr, "UnuseElement")
		}
	}
}

/*
Element

Represents a single DOM element, owns and manages a Handle
*/
type Element struct {
	handle HELEMENT
}

// Constructors
func NewElementFromHandle(h HELEMENT) *Element {
	if h == BAD_HELEMENT {
		panic("Nil helement")
	}
	e := &Element{BAD_HELEMENT}
	e.setHandle(h)
	runtime.SetFinalizer(e, (*Element).finalize)
	return e
}

/** Create new element, the element is disconnected initially from the DOM.
   Element created with ref_count = 1 thus you \b must call HTMLayout_UnuseElement on returned handler.
* \param[in] tagname \b LPCSTR, html tag of the element e.g. "div", "option", etc.
* \param[in] textOrNull \b LPCWSTR, initial text of the element or NULL. text here is a plain text - method does no parsing.
* \param[out ] phe \b #HELEMENT*, variable to receive handle of the element
 **/
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutCreateElement( LPCSTR tagname, LPCWSTR textOrNull, /*out*/ HELEMENT *phe );
//sys HTMLayoutCreateElement(tagname string, textOrNull *uint16, phe *HELEMENT) (ret HLDOM_RESULT, err error) [failretval != 0] = htmlayout.HTMLayoutCreateElement

func NewElement(tagName string) *Element {
	var handle HELEMENT = BAD_HELEMENT
	if ret, err := HTMLayoutCreateElement(tagName, nil, &handle); err != nil {
		domPanic2(ret, "Failed to create new element")
	}
	return NewElementFromHandle(handle)
}

/**Get root DOM element of HTML document.
 * \param[in] hwnd \b HWND, HTMLayout window for which you need to get root
 * element
 * \param[out ] phe \b #HELEMENT*, variable to receive root element
 * \return \b #HLDOM_RESULT
 *
 * Root DOM object is always a 'HTML' element of the document.
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutGetRootElement(HWND hwnd, HELEMENT *phe);
//sys HTMLayoutGetRootElement(hwnd HWND, pheT *HELEMENT) (ret HLDOM_RESULT, err error) [failretval != 0] = htmlayout.HTMLayoutGetRootElement

func GetRootElement(hwnd win.HWND) *Element {
	var handle HELEMENT = BAD_HELEMENT
	if ret, err := HTMLayoutGetRootElement(HWND(hwnd), &handle); err != nil {
		domPanic2(ret, "Failed to get root element")
	}
	return NewElementFromHandle(handle)
}

/**Get focused DOM element of HTML document.
 * \param[in] hwnd \b HWND, HTMLayout window for which you need to get focus
 * element
 * \param[out ] phe \b #HELEMENT*, variable to receive focus element
 * \return \b #HLDOM_RESULT
 *
 * phe can have null value (0).
 *
 * COMMENT: To set focus on element use HTMLayoutSetElementState(STATE_FOCUS,0)
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutGetFocusElement(HWND hwnd, HELEMENT *phe);
//sys HTMLayoutGetFocusElement(hwnd HWND, pheT *HELEMENT) (ret HLDOM_RESULT, err error) [failretval != 0] = htmlayout.HTMLayoutGetFocusElement
func GetFocusedElement(hwnd win.HWND) *Element {
	var handle HELEMENT = BAD_HELEMENT
	if ret, err := HTMLayoutGetFocusElement(HWND(hwnd), &handle); err != nil {
		domPanic2(ret, "Failed to get focus element")
	}
	if handle != BAD_HELEMENT {
		return NewElementFromHandle(handle)
	}
	return nil
}

// Finalizer method, only to be called from Release or by
// the Go runtime
func (e *Element) finalize() {
	// Detach handlers
	if attachedHandlers, hasHandlers := eventHandlers[e.handle]; hasHandlers {
		for handler := range attachedHandlers {
			tag := uintptr(unsafe.Pointer(handler))
			HTMLayoutDetachEventHandler(e.handle, goElementProc, tag)
		}
		delete(eventHandlers, e.handle)
	}

	// Release the underlying htmlayout handle
	unuse(e.handle)
	e.handle = BAD_HELEMENT
}

func (e *Element) Release() {
	// Unregister the finalizer so that it does not get called by Go
	// and then explicitly finalize this element
	runtime.SetFinalizer(e, nil)
	e.finalize()
}

func (e *Element) setHandle(h HELEMENT) {
	use(h)
	unuse(e.handle)
	e.handle = h
}

func (e *Element) Handle() HELEMENT {
	return e.handle
}

func (e *Element) Equals(other *Element) bool {
	return other != nil && e.handle == other.handle
}

// This is the same as AttachHandler, except that behaviors are singleton instances stored
// in a master map.  They may be shared among many elements since they have no state.
// The only reason we keep a separate set of the behaviors is so that the event handler
// dispatch method can tell if an event handler is a behavior or a regular handler.
func (e *Element) attachBehavior(handler *EventHandler) {
	tag := uintptr(unsafe.Pointer(handler))
	if subscription := handler.Subscription(); subscription == HANDLE_ALL {
		if ret := HTMLayoutAttachEventHandler(e.handle, goElementProc, tag); ret != HLDOM_OK {
			domPanic2(ret, "Failed to attach event handler to element")
		}
	} else {
		if ret := HTMLayoutAttachEventHandlerEx(e.handle, goElementProc, tag, subscription); ret != HLDOM_OK {
			domPanic2(ret, "Failed to attach event handler to element")
		}
	}
}

func (e *Element) AttachHandler(handler *EventHandler) {
	attachedHandlers, hasAttachments := eventHandlers[e.handle]
	if hasAttachments {
		if _, exists := attachedHandlers[handler]; exists {
			// This exact event handler is already attached to this exact element.
			return
		}
	}

	// Don't let the caller disable ATTACH/DETACH events, otherwise we
	// won't know when to throw out our event handler object
	subscription := handler.Subscription()
	subscription &= ^uint(DISABLE_INITIALIZATION & 0xffffffff)

	tag := uintptr(unsafe.Pointer(handler))
	if subscription == HANDLE_ALL {
		if ret := HTMLayoutAttachEventHandler(e.handle, goElementProc, tag); ret != HLDOM_OK {
			domPanic2(ret, "Failed to attach event handler to element")
		}
	} else {
		if ret := HTMLayoutAttachEventHandlerEx(e.handle, goElementProc, tag, subscription); ret != HLDOM_OK {
			domPanic2(ret, "Failed to attach event handler to element")
		}
	}

	if !hasAttachments {
		eventHandlers[e.handle] = make(map[*EventHandler]bool, 8)
	}
	eventHandlers[e.handle][handler] = true
}

func (e *Element) DetachHandler(handler *EventHandler) {
	tag := uintptr(unsafe.Pointer(handler))
	if attachedHandlers, exists := eventHandlers[e.handle]; exists {
		if _, exists := attachedHandlers[handler]; exists {
			if ret := HTMLayoutDetachEventHandler(e.handle, goElementProc, tag); ret != HLDOM_OK {
				domPanic2(ret, "Failed to detach event handler from element")
			}
			delete(attachedHandlers, handler)
			if len(attachedHandlers) == 0 {
				delete(eventHandlers, e.handle)
			}
			return
		}
	}
	panic("cannot detach, handler was not registered")
}

/**Apply changes and refresh element area in its window.
 * \param[in] he \b #HELEMENT
 * \param[in] flags \b UINT, combination of UPDATE_ELEMENT_FLAGS.
 * \return \b #HLDOM_RESULT
 *
 *  Note HTMLayoutUpdateElement is an equivalent of HTMLayoutUpdateElementEx(,RESET_STYLE_DEEP | REMEASURE [| REDRAW_NOW ])
 *
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutUpdateElementEx(HELEMENT he, UINT flags);
//sys HTMLayoutUpdateElementEx(he HELEMENT, flags uint) (ret HLDOM_RESULT) = htmlayout.HTMLayoutUpdateElementEx

func (e *Element) Update(restyle, restyleDeep, remeasure, remeasureDeep, render bool) {
	var flags uint
	if restyle {
		if restyleDeep {
			flags |= RESET_STYLE_DEEP
		} else {
			flags |= RESET_STYLE_THIS
		}
	}
	if remeasure {
		if remeasureDeep {
			flags |= MEASURE_DEEP
		} else {
			flags |= MEASURE_INPLACE
		}
	}
	if render {
		flags |= REDRAW_NOW
	}
	if ret := HTMLayoutUpdateElementEx(e.handle, flags); ret != HLDOM_OK {
		domPanic2(ret, "Failed to update element")
	}
}

/**Set the mouse capture to the specified element.
 * \param[in] he \b #HELEMENT
 * \return \b #HLDOM_RESULT
 *
 * After call to this function all mouse events will be targeted to the element.
 * To remove mouse capture call ::ReleaseCapture() function. It is declared somewhere in <windows.h>.
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutSetCapture(HELEMENT he);
//sys HTMLayoutSetCapture(he HELEMENT) (ret HLDOM_RESULT) = htmlayout.HTMLayoutSetCapture
func (e *Element) Capture() {
	if ret := HTMLayoutSetCapture(e.handle); ret != HLDOM_OK {
		domPanic2(ret, "Failed to set capture for element")
	}
}

// func (e *Element) ReleaseCapture() {
// 	if ok := C.ReleaseCapture(); ok == 0 {
// 		panic("Failed to release capture for element")
// 	}
// }

// Functions for querying elements

// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutSelectElementsW(
//           HELEMENT  he,
//           LPCWSTR   CSS_selectors,
//           HTMLayoutElementCallback*
//                     callback,
//           LPVOID    param);
//sys HTMLayoutSelectElementsW(he HELEMENT, CSS_selectors string, callback uintptr, param uintptr) (ret HLDOM_RESULT, err error) [failretval != 0] = htmlayout.HTMLayoutSelectElementsW
func (e *Element) Select(selector string) []*Element {
	results := make([]*Element, 0, 32)
	if ret, err := HTMLayoutSelectElementsW(e.handle, selector, goSelectCallback, uintptr(unsafe.Pointer(&results))); err != nil {
		domPanic2(ret, "Failed to select dom elements, selector: '", selector, "'")
	}
	return results
}

// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutSelectParentW(
//           HELEMENT  he,
//           LPCWSTR   selector,
//           UINT      depth,
//           /*out*/ HELEMENT* heFound);
//sys HTMLayoutSelectParentW(he HELEMENT, selector string, depth uint, heFound *HELEMENT) (ret HLDOM_RESULT, err error) [failretval != 0] = htmlayout.HTMLayoutSelectParentW

// Searches up the parent chain to find the first element that matches the given selector.
// Includes the element in the search.  Depth indicates how far the search should progress.
// Depth = 1 means only consider this element.  Depth = 0 means search all the way up to the
// root.  Any other positive value of depth limits the length of the search.
func (e *Element) SelectParentLimit(selector string, depth int) *Element {
	var parent HELEMENT
	if ret, err := HTMLayoutSelectParentW(e.handle, selector, uint(depth), &parent); err != nil {
		domPanic2(ret, "Failed to select parent dom elements, selector: '", selector, "'")
	}
	if parent != 0 {
		return NewElementFromHandle(parent)
	}
	return nil
}

func (e *Element) SelectParent(selector string) *Element {
	return e.SelectParentLimit(selector, 0)
}

/** SendEvent - sends sinking/bubbling event to the child/parent chain of he element.
   First event will be send in SINKING mode (with SINKING flag) - from root to he element itself.
   Then from he element to its root on parents chain without SINKING flag (bubbling phase).

* \param[in] he \b HELEMENT, element to send this event to.
* \param[in] appEventCode \b UINT, event ID, see: #BEHAVIOR_EVENTS
* \param[in] heSource \b HELEMENT, optional handle of the source element, e.g. some list item
* \param[in] reason \b UINT, notification specific event reason code
* \param[out] handled \b BOOL*, variable to receive TRUE if any handler handled it, FALSE otherwise.

**/

// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutSendEvent(
//           HELEMENT he, UINT appEventCode, HELEMENT heSource, UINT_PTR reason, /*out*/ BOOL* handled);
//sys HTMLayoutSendEvent(he HELEMENT, appEventCode uint, heSource HELEMENT, reason *uint, handled *BOOL) (ret HLDOM_RESULT) = htmlayout.HTMLayoutSendEvent

// For delivering programmatic events to this element.
// Returns true if the event was handled, false otherwise
func (e *Element) SendEvent(eventCode uint, source *Element, reason uint) bool {
	var handled BOOL = 0
	if ret := HTMLayoutSendEvent(e.handle, eventCode, source.handle, &reason, &handled); ret != HLDOM_OK {
		domPanic2(ret, "Failed to send event")
	}
	return handled != 0
}

/** PostEvent - post sinking/bubbling event to the child/parent chain of he element.
 *  Function will return immediately posting event into input queue of the application.
 *
 * \param[in] he \b HELEMENT, element to send this event to.
 * \param[in] appEventCode \b UINT, event ID, see: #BEHAVIOR_EVENTS
 * \param[in] heSource \b HELEMENT, optional handle of the source element, e.g. some list item
 * \param[in] reason \b UINT, notification specific event reason code

 **/
//
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutPostEvent(
//           HELEMENT he, UINT appEventCode, HELEMENT heSource, UINT reason);
//sys HTMLayoutPostEvent(he HELEMENT, appEventCode uint, heSource HELEMENT, reason *uint) (ret HLDOM_RESULT) = htmlayout.HTMLayoutPostEvent

// For asynchronously delivering programmatic events to this element.
func (e *Element) PostEvent(eventCode uint, source *Element, reason uint) {
	if ret := HTMLayoutPostEvent(e.handle, eventCode, source.handle, &reason); ret != HLDOM_OK {
		domPanic2(ret, "Failed to post event")
	}
}

//
// DOM structure accessors/modifiers:
//

/**Get number of child elements.
 * \param[in] he \b #HELEMENT, element which child elements you need to count
 * \param[out] count \b UINT*, variable to receive number of child elements
 * \return \b #HLDOM_RESULT
 *
 * \par Example:
 * for paragraph defined as
 * \verbatim <p>Hello <b>wonderfull</b> world!</p> \endverbatim
 * count will be set to 1 as the paragraph has only one sub element:
 * \verbatim <b>wonderfull</b> \endverbatim
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutGetChildrenCount(HELEMENT he, UINT* count);
//sys HTMLayoutGetChildrenCount(he HELEMENT, count *uint) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetChildrenCount
func (e *Element) ChildCount() uint {
	var count uint
	if ret := HTMLayoutGetChildrenCount(e.handle, &count); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get child count")
	}
	return count
}

/**Get handle of every element's child element.
 * \param[in] he \b #HELEMENT
 * \param[in] n \b UINT, number of the child element
 * \param[out] phe \b #HELEMENT*, variable to receive handle of the child
 * element
 * \return \b #HLDOM_RESULT
 *
 * \par Example:
 * for paragraph defined as
 * \verbatim <p>Hello <b>wonderfull</b> world!</p> \endverbatim
 * *phe will be equal to handle of &lt;b&gt; element:
 * \verbatim <b>wonderfull</b> \endverbatim
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutGetNthChild(HELEMENT he, UINT n, HELEMENT* phe);
//sys HTMLayoutGetNthChild(he HELEMENT, n uint, phe *HELEMENT) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetNthChild
func (e *Element) Child(index uint) *Element {
	var child HELEMENT
	if ret := HTMLayoutGetNthChild(e.handle, index, &child); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get child at index: ", index)
	}
	return NewElementFromHandle(child)
}

func (e *Element) Children() []*Element {
	slice := make([]*Element, 0, 32)
	for i := uint(0); i < e.ChildCount(); i++ {
		slice = append(slice, e.Child(i))
	}
	return slice
}

/**Get element index.
 * \param[in] he \b #HELEMENT
 * \param[out] p_index \b LPUINT, variable to receive number of the element
 * among parent element's subelements.
 * \return \b #HLDOM_RESULT
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutGetElementIndex(HELEMENT he, LPUINT p_index);
//sys HTMLayoutGetElementIndex(he HELEMENT, p_index *uint) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetElementIndex
func (e *Element) Index() uint {
	var index uint
	if ret := HTMLayoutGetElementIndex(e.handle, &index); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get element's index")
	}
	return index
}

/**Get parent element.
 * \param[in] he \b #HELEMENT, element which parent you need
 * \param[out] p_parent_he \b #HELEMENT*, variable to recieve handle of the
 * parent element
 * \return \b #HLDOM_RESULT
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutGetParentElement(HELEMENT he, HELEMENT* p_parent_he);
//sys HTMLayoutGetParentElement(he HELEMENT, p_parent_he *HELEMENT) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetParentElement
func (e *Element) Parent() *Element {
	var parent HELEMENT
	if ret := HTMLayoutGetParentElement(e.handle, &parent); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get parent")
	}
	if parent != 0 {
		return NewElementFromHandle(parent)
	}
	return nil
}

/** Insert element at \i index position of parent.
   It is not an error to insert element which already has parent - it will be disconnected first, but
   you need to update elements parent in this case.
* \param index \b UINT, position of the element in parent collection.
  It is not an error to provide index greater than elements count in parent -
  it will be appended.
**/
//
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutInsertElement( HELEMENT he, HELEMENT hparent, UINT index );
//sys HTMLayoutInsertElement(he HELEMENT,hparent HELEMENT, index uint) (ret HLDOM_RESULT) = htmlayout.HTMLayoutInsertElement

func (e *Element) InsertChild(child *Element, index uint) {
	if ret := HTMLayoutInsertElement(child.handle, e.handle, index); ret != HLDOM_OK {
		domPanic2(ret, "Failed to insert child element at index: ", index)
	}
}

func (e *Element) AppendChild(child *Element) {
	count := e.ChildCount()
	if ret := HTMLayoutInsertElement(child.handle, e.handle, count); ret != HLDOM_OK {
		domPanic2(ret, "Failed to append child element")
	}
}

/** Take element out of its container (and DOM tree).
   Element will be destroyed when its reference counter will become zero
**/
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutDetachElement( HELEMENT he );
//sys HTMLayoutDetachElement(he HELEMENT) (ret HLDOM_RESULT) = htmlayout.HTMLayoutDetachElement
func (e *Element) Detach() {
	if ret := HTMLayoutDetachElement(e.handle); ret != HLDOM_OK {
		domPanic2(ret, "Failed to detach element from dom")
	}
}

/**Delete element.
 * \param[in] he \b #HELEMENT
 * \return \b #HLDOM_RESULT
 *
 * This function removes element from the DOM tree and then deletes it.
 *
 * \warning After call to this function \c he will become invalid.
 **/
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutDeleteElement(HELEMENT he);
//sys HTMLayoutDeleteElement(he HELEMENT) (ret HLDOM_RESULT) = htmlayout.HTMLayoutDeleteElement
func (e *Element) Delete() {
	if ret := HTMLayoutDeleteElement(e.handle); ret != HLDOM_OK {
		domPanic2(ret, "Failed to delete element from dom")
	}
	e.finalize()
}

/** Create new element as copy of existing element, new element is a full (deep) copy of the element and
   is disconnected initially from the DOM.
   Element created with ref_count = 1 thus you \b must call HTMLayout_UnuseElement on returned handler.
* \param[in] he \b #HELEMENT, source element.
* \param[out ] phe \b #HELEMENT*, variable to receive handle of the new element.
 **/
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutCloneElement( HELEMENT he, /*out*/ HELEMENT *phe );
//sys HTMLayoutCloneElement(he HELEMENT, phe *HELEMENT) (ret HLDOM_RESULT) = htmlayout.HTMLayoutCloneElement

// Makes a deep clone of the receiver, the resulting subtree is not attached to the dom.
func (e *Element) Clone() *Element {
	var clone HELEMENT
	if ret := HTMLayoutCloneElement(e.handle, &clone); ret != HLDOM_OK {
		domPanic2(ret, "Failed to clone element")
	}
	return NewElementFromHandle(clone)
}

/** HTMLayoutSwapElements - swap element positions.
 * Function changes "insertion points" of two elements. So it swops indexes and parents of two elements.
 * \param[in] he1 \b HELEMENT, first element.
 * \param[in] he2 \b HELEMENT, second element.
 **/
//
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutSwapElements(
//          HELEMENT he1, HELEMENT he2 );
//sys HTMLayoutSwapElements(he HELEMENT, other HELEMENT) (ret HLDOM_RESULT) = htmlayout.HTMLayoutSwapElements
func (e *Element) Swap(other *Element) {
	if ret := HTMLayoutSwapElements(e.handle, other.handle); ret != HLDOM_OK {
		domPanic2(ret, "Failed to swap elements")
	}
}

/** HTMLayoutSortElements - sort children of the element.
 * \param[in] he \b HELEMENT, element which children to be sorted.
 * \param[in] firstIndex \b UINT, first child index to start sorting from.
 * \param[in] lastIndex \b UINT, last index of the sorting range, element with this index will not be included in the sorting.
 * \param[in] cmpFunc \b ELEMENT_COMPARATOR, comparator function.
 * \param[in] cmpFuncParam \b LPVOID, parameter to be passed in comparator function.
 **/

// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutSortElements(
//          HELEMENT he, UINT firstIndex, UINT lastIndex,
//          ELEMENT_COMPARATOR* cmpFunc, LPVOID cmpFuncParam );
//sys HTMLayoutSortElements(he HELEMENT, firstIndex uint, lastIndex uint, cmpFunc uintptr, cmpFuncParam uintptr) (ret HLDOM_RESULT) = htmlayout.HTMLayoutSortElements

// Sorts 'count' child elements starting at index 'start'.  Uses comparator to define the
// order.  Comparator should return -1, or 0, or 1 to indicate less, equal or greater
func (e *Element) SortChildrenRange(start, count uint, comparator func(*Element, *Element) int) {
	end := start + count
	arg := uintptr(unsafe.Pointer(&comparator))
	if ret := HTMLayoutSortElements(e.handle, start, end, goElementComparator, arg); ret != HLDOM_OK {
		domPanic2(ret, "Failed to sort elements")
	}
}

func (e *Element) SortChildren(comparator func(*Element, *Element) int) {
	e.SortChildrenRange(0, e.ChildCount(), comparator)
}

/** Start Timer for the element.
   Element will receive on_timer event
   To stop timer call HTMLayoutSetTimer( he, 0 );
**/
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutSetTimer( HELEMENT he, UINT milliseconds );
//sys HTMLayoutSetTimer(he HELEMENT, milliseconds uint) (ret HLDOM_RESULT) = htmlayout.HTMLayoutSetTimer
func (e *Element) SetTimer(ms int) {
	if ret := HTMLayoutSetTimer(e.handle, uint(ms)); ret != HLDOM_OK {
		domPanic2(ret, "Failed to set timer")
	}
}

func (e *Element) CancelTimer() {
	e.SetTimer(0)
}

/**Get HWND of containing window.
 * \param[in] he \b #HELEMENT
 * \param[out] p_hwnd \b HWND*, variable to receive window handle
 * \param[in] rootWindow \b BOOL, handle of which window to get:
 * - TRUE - HTMLayout window
 * - FALSE - nearest parent element having overflow:auto or :scroll
 * \return \b #HLDOM_RESULT
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutGetElementHwnd(HELEMENT he, HWND* p_hwnd, BOOL rootWindow);
//sys HTMLayoutGetElementHwnd(he HELEMENT, p_hwnd *HWND, rootWindow BOOL) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetElementHwnd
func (e *Element) Hwnd() win.HWND {
	var hwnd win.HWND
	if ret := HTMLayoutGetElementHwnd(e.handle, (*HWND)(&hwnd), 0); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get element's hwnd")
	}
	return hwnd
}

func (e *Element) RootHwnd() win.HWND {
	var hwnd win.HWND
	if ret := HTMLayoutGetElementHwnd(e.handle, (*HWND)(&hwnd), 1); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get element's root hwnd")
	}
	return hwnd
}

/**Get text of the element and information where child elements are placed.
 * \param[in] he \b #HELEMENT
 * \param[out] utf8bytes \b pointer to byte address receiving UTF8 encoded HTML
 * \param[in] outer \b BOOL, if TRUE will retunr outer HTML otherwise inner.
 * \return \b #HLDOM_RESULT
 */
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutGetElementHtml(HELEMENT he, LPBYTE* utf8bytes, BOOL outer);
//sys HTMLayoutGetElementHtml(he HELEMENT, utf8bytes uintptr, outer BOOL) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetElementHtml
func (e *Element) Html() string {
	var data *C.char
	p := unsafe.Pointer(&data)
	if ret := HTMLayoutGetElementHtml(e.handle, uintptr(p), BOOL(0)); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get inner html")
	}
	str := C.GoString(data)
	// TODO need free??
	// C.free(unsafe.Pointer(data))
	return str
}

func (e *Element) OuterHtml() string {
	var data *C.char
	p := unsafe.Pointer(&data)
	if ret := HTMLayoutGetElementHtml(e.handle, uintptr(p), BOOL(1)); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get inner html")
	}
	str := C.GoString(data)
	// C.free(unsafe.Pointer(data))
	return str
}

/**Get element's type.
 * \param[in] he \b #HELEMENT
 * \param[out] p_type \b LPCSTR*, receives name of the element type.
 * \return \b #HLDOM_RESULT
 *
 * \par Example:
 * For &lt;div&gt; tag p_type will be set to "div".
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutGetElementType(HELEMENT he, LPCSTR* p_type);
//sys HTMLayoutGetElementType(he HELEMENT, p_type uintptr) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetElementType
func (e *Element) Type() string {
	var data *C.char
	if ret := HTMLayoutGetElementType(e.handle, (uintptr)(unsafe.Pointer(&data))); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get element type")
	}
	return C.GoString(data)
}

/**Set inner or outer html of the element.
 * \param[in] he \b #HELEMENT
 * \param[in] html \b LPCBYTE, UTF-8 encoded string containing html text
 * \param[in] htmlLength \b DWORD, length in bytes of \c html.
 * \param[in] where \b UINT, possible values are:
 * - SIH_REPLACE_CONTENT - replace content of the element
 * - SIH_INSERT_AT_START - insert html before first child of the element
 * - SIH_APPEND_AFTER_LAST - insert html after last child of the element
 *
 * - SOH_REPLACE - replace element by html, a.k.a. element.outerHtml = "something"
 * - SOH_INSERT_BEFORE - insert html before the element
 * - SOH_INSERT_AFTER - insert html after the element
 *   ATTN: SOH_*** operations do not work for inline elements like <SPAN>
 *
 * \return /b #HLDOM_RESULT
  **/
// EXTERN_C HLDOM_RESULT HLAPI
//       HTMLayoutSetElementHtml(HELEMENT he, LPCBYTE html, DWORD htmlLength, UINT where);
//sys HTMLayoutSetElementHtml(he HELEMENT, html string, htmlLength int, where uint) (ret HLDOM_RESULT, err error) [failretval != 0] = htmlayout.HTMLayoutSetElementHtml
func (e *Element) SetHtml(html string) {
	if ret, err := HTMLayoutSetElementHtml(e.handle, html, len(html), SIH_REPLACE_CONTENT); err != nil {
		domPanic2(ret, "Failed to replace element's html")
	}
}

func (e *Element) PrependHtml(prefix string) {
	if ret, err := HTMLayoutSetElementHtml(e.handle, prefix, len(prefix), SIH_INSERT_AT_START); err != nil {
		domPanic2(ret, "Failed to replace element's html")
	}
}

func (e *Element) AppendHtml(suffix string) {
	if ret, err := HTMLayoutSetElementHtml(e.handle, suffix, len(suffix), SIH_APPEND_AFTER_LAST); err != nil {
		domPanic2(ret, "Failed to replace element's html")
	}
}

/**Set inner text of the element.
 * \param[in] he \b #HELEMENT
 * \param[in] utf8bytes \b pointer, UTF8 encoded plain text
 * \param[in] length \b UINT, number of bytes in utf8bytes sequence
 * \return \b #HLDOM_RESULT
 */
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutSetElementInnerText(HELEMENT he, LPCBYTE utf8bytes, UINT length);
//sys HTMLayoutSetElementInnerText(he HELEMENT, text string, length uint) (ret HLDOM_RESULT, err error) [failretval != 0] = htmlayout.HTMLayoutSetElementInnerText
func (e *Element) SetInnerText(text string) {
	if ret, err := HTMLayoutSetElementInnerText(e.handle, text, uint(len(text))); err != nil {
		domPanic2(ret, "Failed to replace element's text")
	}
}

/**Get inner text of the element.
 * \param[in] he \b #HELEMENT
 * \param[out] utf8bytes \b pointer to byte address receiving UTF8 encoded plain text
 * \return \b #HLDOM_RESULT
 */
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutGetElementInnerText(HELEMENT he, LPBYTE* utf8bytes);
//sys HTMLayoutGetElementInnerText(he HELEMENT, utf8bytes uintptr) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetElementInnerText

func (e *Element) GetInnerText() string {
	var data *C.char
	if ret := HTMLayoutGetElementInnerText(e.handle, (uintptr)(unsafe.Pointer(&data))); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get text")
	}
	return C.GoString(data)
}

//
// HTML attribute accessors/modifiers:
//

/**Get value of any element's attribute by name.
 * \param[in] he \b #HELEMENT
 * \param[in] name \b LPCSTR, attribute name
 * \param[out] p_value \b LPCWSTR*, will be set to address of the string
 * containing attribute value
 * \return \b #HLDOM_RESULT
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutGetAttributeByName(HELEMENT he, LPCSTR name, LPCWSTR* p_value);
//sys HTMLayoutGetAttributeByName(he HELEMENT, utf8bytes *byte, p_value uintptr) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetAttributeByName

// Returns the value of attr and a boolean indicating whether or not that attr exists.
// If the boolean is true, then the returned string is valid.
func (e *Element) Attr(key string) (string, bool) {
	szValue := (*C.WCHAR)(nil)
	szKey := C.CString(key)
	defer C.free(unsafe.Pointer(szKey))
	if ret := HTMLayoutGetAttributeByName(e.handle, (*byte)(unsafe.Pointer(szKey)), (uintptr)(unsafe.Pointer(&szValue))); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get attribute: ", key)
	}
	if szValue != nil {
		return utf16ToString((*uint16)(szValue)), true
	}
	return "", false
}

func (e *Element) AttrAsFloat(key string) (float64, bool, error) {
	var f float64
	var err error
	if s, exists := e.Attr(key); !exists {
		return 0.0, false, nil
	} else if f, err = strconv.ParseFloat(s, 64); err != nil {
		return 0.0, true, err
	}
	return float64(f), true, nil
}

func (e *Element) AttrAsInt(key string) (int, bool, error) {
	var i int
	var err error
	if s, exists := e.Attr(key); !exists {
		return 0, false, nil
	} else if i, err = strconv.Atoi(s); err != nil {
		return 0, true, err
	}
	return i, true, nil
}

/**Set attribute's value.
 * \param[in] he \b #HELEMENT
 * \param[in] name \b LPCSTR, attribute name
 * \param[in] value \b LPCWSTR, new attribute value or 0 if you want to remove attribute.
 * \return \b #HLDOM_RESULT
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutSetAttributeByName(HELEMENT he, LPCSTR name, LPCWSTR value);
//sys HTMLayoutSetAttributeByName(he HELEMENT, name *byte, value *uint16) (ret HLDOM_RESULT) = htmlayout.HTMLayoutSetAttributeByName

func (e *Element) SetAttr(key string, value interface{}) {
	szKey := syscall.StringBytePtr(key)
	var ret HLDOM_RESULT = HLDOM_OK
	switch v := value.(type) {
	case string:
		ret = HTMLayoutSetAttributeByName(e.handle, szKey, stringToUtf16Ptr(v))
	case float32:
		ret = HTMLayoutSetAttributeByName(e.handle, szKey, stringToUtf16Ptr(strconv.FormatFloat(float64(v), 'g', -1, 64)))
	case float64:
		ret = HTMLayoutSetAttributeByName(e.handle, szKey, stringToUtf16Ptr(strconv.FormatFloat(float64(v), 'g', -1, 64)))
	case int:
		ret = HTMLayoutSetAttributeByName(e.handle, szKey, stringToUtf16Ptr(strconv.Itoa(v)))
	case int32:
		ret = HTMLayoutSetAttributeByName(e.handle, szKey, stringToUtf16Ptr(strconv.FormatInt(int64(v), 10)))
	case int64:
		ret = HTMLayoutSetAttributeByName(e.handle, szKey, stringToUtf16Ptr(strconv.FormatInt(v, 10)))
	case nil:
		ret = HTMLayoutSetAttributeByName(e.handle, szKey, nil)
	default:
		panic(fmt.Sprintf("Don't know how to format this argument type: %s", reflect.TypeOf(v)))
	}
	if ret != HLDOM_OK {
		domPanic2(ret, "Failed to set attribute: "+key)
	}
}

func (e *Element) RemoveAttr(key string) {
	e.SetAttr(key, nil)
}

/**Get value of any element's attribute by attribute's number.
 * \param[in] he \b #HELEMENT
 * \param[in] n \b UINT, number of desired attribute
 * \param[out] p_name \b LPCSTR*, will be set to address of the string
 * containing attribute name
 * \param[out] p_value \b LPCWSTR*, will be set to address of the string
 * containing attribute value
 * \return \b #HLDOM_RESULT
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutGetNthAttribute(HELEMENT he, UINT n, LPCSTR* p_name, LPCWSTR* p_value);
//sys HTMLayoutGetNthAttribute(he HELEMENT, n uint, name uintptr, value uintptr) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetNthAttribute

func (e *Element) AttrByIndex(index int) (string, string) {
	szName := (*C.CHAR)(nil)
	szValue := (*C.WCHAR)(nil)
	if ret := HTMLayoutGetNthAttribute(e.handle, uint(index), uintptr(unsafe.Pointer(&szName)), uintptr(unsafe.Pointer(&szValue))); ret != HLDOM_OK {
		domPanic2(ret, fmt.Sprintf("Failed to get attribute by index: %u", index))
	}
	return C.GoString((*C.char)(szName)), utf16ToString((*uint16)(szValue))
}

/**Get number of element's attributes.
 * \param[in] he \b #HELEMENT
 * \param[out] p_count \b LPUINT, variable to receive number of element
 * attributes.
 * \return \b #HLDOM_RESULT
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutGetAttributeCount(HELEMENT he, LPUINT p_count);
//sys HTMLayoutGetAttributeCount(he HELEMENT, p_count *uint) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetAttributeCount

func (e *Element) AttrCount() uint {
	var count uint = 0
	if ret := HTMLayoutGetAttributeCount(e.handle, &count); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get attribute count")
	}
	return count
}

//
// CSS style attribute accessors/mutators
//

/**Get element's style attribute.
 * \param[in] he \b #HELEMENT
 * \param[in] name \b LPCSTR, name of the style attribute
 * \param[out] p_value \b LPCWSTR*, variable to receive value of the style attribute.
 *
 * Style attributes are those that are set using css. E.g. "font-face: arial" or "display: block".
 *
 * \sa #HTMLayoutSetStyleAttribute()
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutGetStyleAttribute(HELEMENT he, LPCSTR name, LPCWSTR* p_value);
//sys HTMLayoutGetStyleAttribute(he HELEMENT, name uintptr,  p_value uintptr) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetStyleAttribute

// Returns the value of the style and a boolean indicating whether or not that style exists.
// If the boolean is true, then the returned string is valid.
func (e *Element) Style(key string) (string, bool) {
	szValue := (*C.WCHAR)(nil)
	szKey := C.CString(key)
	defer C.free(unsafe.Pointer(szKey))
	if ret := HTMLayoutGetStyleAttribute(e.handle, uintptr(unsafe.Pointer(szKey)), uintptr(unsafe.Pointer(&szValue))); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get style: "+key)
	}
	if szValue != nil {
		return utf16ToString((*uint16)(szValue)), true
	}
	return "", false
}

/**Get element's style attribute.
 * \param[in] he \b #HELEMENT
 * \param[in] name \b LPCSTR, name of the style attribute
 * \param[out] value \b LPCWSTR, value of the style attribute or NULL for clearing the attribute
 *
 * Style attributes are those that are set using css. E.g. "font-face: arial" or "display: block".
 *
 * \sa #HTMLayoutGetStyleAttribute()
 **/
// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutSetStyleAttribute(HELEMENT he, LPCSTR name, LPCWSTR value);
//sys HTMLayoutSetStyleAttribute(he HELEMENT, name uintptr,  value uintptr) (ret HLDOM_RESULT) = htmlayout.HTMLayoutSetStyleAttribute

func (e *Element) SetStyle(key string, value interface{}) {
	szKey := C.CString(key)
	defer C.free(unsafe.Pointer(szKey))
	var valuePtr *uint16 = nil

	switch v := value.(type) {
	case string:
		valuePtr = stringToUtf16Ptr(v)
	case float32:
		valuePtr = stringToUtf16Ptr(strconv.FormatFloat(float64(v), 'g', -1, 64))
	case float64:
		valuePtr = stringToUtf16Ptr(strconv.FormatFloat(float64(v), 'g', -1, 64))
	case int:
		valuePtr = stringToUtf16Ptr(strconv.Itoa(v))
	case int32:
		valuePtr = stringToUtf16Ptr(strconv.FormatInt(int64(v), 10))
	case int64:
		valuePtr = stringToUtf16Ptr(strconv.FormatInt(v, 10))
	case nil:
		valuePtr = nil
	default:
		panic(fmt.Sprintf("Don't know how to format this argument type: %s", reflect.TypeOf(v)))
	}

	if ret := HTMLayoutSetStyleAttribute(e.handle, uintptr(unsafe.Pointer(szKey)), uintptr(unsafe.Pointer(valuePtr))); ret != HLDOM_OK {
		domPanic2(ret, "Failed to set style: "+key)
	}
}

func (e *Element) RemoveStyle(key string) {
	e.SetStyle(key, nil)
}

func (e *Element) ClearStyles(key string) {
	if ret := HTMLayoutSetStyleAttribute(e.handle, 0, 0); ret != HLDOM_OK {
		domPanic2(ret, "Failed to clear all styles")
	}
}

//
// Element state manipulation
//

/** Get/set state bits, stateBits*** accept or'ed values above
 **/
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutGetElementState( HELEMENT he, UINT* pstateBits);
//sys HTMLayoutGetElementState(he HELEMENT, pstateBits *uint) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetElementState

// Gets the whole set of state flags for this element
func (e *Element) StateFlags() uint {
	var state uint
	if ret := HTMLayoutGetElementState(e.handle, &state); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get element state flags")
	}
	return state
}

// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutSetElementState( HELEMENT he, UINT stateBitsToSet, UINT stateBitsToClear, BOOL updateView);
//sys HTMLayoutSetElementState(he HELEMENT, stateBitsToSet uint, stateBitsToClear uint, updateView BOOL) (ret HLDOM_RESULT) = htmlayout.HTMLayoutSetElementState

// Replaces the whole set of state flags with the specified value
func (e *Element) SetStateFlags(flags uint) {
	shouldUpdate := BOOL(1)
	if ret := HTMLayoutSetElementState(e.handle, flags, ^flags, shouldUpdate); ret != HLDOM_OK {
		domPanic2(ret, "Failed to set element state flags")
	}
}

// Returns true if the specified flag is "on"
func (e *Element) State(flag uint) bool {
	return e.StateFlags()&flag != 0
}

// Sets the specified flag to "on" or "off" according to the value of the provided boolean
func (e *Element) SetState(flag uint, on bool) {
	addBits := uint(0)
	clearBits := uint(0)
	if on {
		addBits = flag
	} else {
		clearBits = flag
	}
	shouldUpdate := BOOL(1)
	if ret := HTMLayoutSetElementState(e.handle, addBits, clearBits, shouldUpdate); ret != HLDOM_OK {
		domPanic2(ret, "Failed to set element state flag")
	}
}

//
// Functions for retrieving/setting the various dimensions of an element
//

/** HTMLayoutMoveElement - moves element from its normal place to the position defined by xView, yView.
 *
 * \param[in] he \b HELEMENT, element.
 * \param[in] xView \b INT, new x coordinate of content box of the element relative to the view - htmlayout window.
 * \param[in] yView \b INT, new y coordinate of content box of the element relative to the view - htmlayout window.
 *
 * If element is moved outside of the view then HTMLayoutMoveElement will create popup window for it.
 *
 **/
//
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutMoveElement( HELEMENT he, INT xView, INT yView);
//sys HTMLayoutMoveElement(he HELEMENT, xView int, yView int) (ret HLDOM_RESULT) = htmlayout.HTMLayoutMoveElement
func (e *Element) Move(x, y int) {
	if ret := HTMLayoutMoveElement(e.handle, x, y); ret != HLDOM_OK {
		domPanic2(ret, "Failed to move element")
	}
}

/** HTMLayoutMoveElementEx - moves and resizes the element from its normal place to the position defined by xView, yView.
 *
 * \param[in] he \b HELEMENT, element.
 * \param[in] xView \b INT, new x coordinate of content box of the element relative to the view - htmlayout window.
 * \param[in] yView \b INT, new y coordinate of content box of the element relative to the view - htmlayout window.
 * \param[in] width \b INT, new width of content box of the element relative to the view - htmlayout window.
 * \param[in] height \b INT, new height of content box of the element relative to the view - htmlayout window.
 *
 * If element is moved outside of the view then HTMLayoutMoveElement will create popup window for it.
 *
 **/
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutMoveElementEx( HELEMENT he, INT xView, INT yView,INT width, INT height);
//sys HTMLayoutMoveElementEx(he HELEMENT, xView int, yView int, width int, height int) (ret HLDOM_RESULT) = htmlayout.HTMLayoutMoveElementEx

func (e *Element) Resize(x, y, w, h int) {
	if ret := HTMLayoutMoveElementEx(e.handle, x, y, w, h); ret != HLDOM_OK {
		domPanic2(ret, "Failed to resize element")
	}
}

// EXTERN_C  HLDOM_RESULT HLAPI HTMLayoutGetElementLocation(HELEMENT he, LPRECT p_location, UINT areas /*ELEMENT_AREAS*/);
//sys HTMLayoutGetElementLocation(he HELEMENT, p_location *Rect, areas uint) (ret HLDOM_RESULT) = htmlayout.HTMLayoutGetElementLocation
func (e *Element) getRect(rectTypeFlags uint) (left, top, right, bottom int) {
	r := Rect{}
	if ret := HTMLayoutGetElementLocation(e.handle, &r, rectTypeFlags); ret != HLDOM_OK {
		domPanic2(ret, "Failed to get element rect")
	}
	return int(r.Left), int(r.Top), int(r.Right), int(r.Bottom)
}

func (e *Element) ContentBox() (left, top, right, bottom int) {
	return e.getRect(CONTENT_BOX)
}

func (e *Element) ContentViewBox() (left, top, right, bottom int) {
	return e.getRect(CONTENT_BOX | VIEW_RELATIVE)
}

func (e *Element) ContentBoxSize() (width, height int) {
	l, t, r, b := e.getRect(CONTENT_BOX)
	return int(r - l), int(b - t)
}

func (e *Element) PaddingBox() (left, top, right, bottom int) {
	return e.getRect(PADDING_BOX)
}

func (e *Element) PaddingViewBox() (left, top, right, bottom int) {
	return e.getRect(PADDING_BOX | VIEW_RELATIVE)
}

func (e *Element) PaddingBoxSize() (width, height int) {
	l, t, r, b := e.getRect(PADDING_BOX)
	return int(r - l), int(b - t)
}

func (e *Element) BorderBox() (left, top, right, bottom int) {
	return e.getRect(BORDER_BOX)
}

func (e *Element) BorderViewBox() (left, top, right, bottom int) {
	return e.getRect(BORDER_BOX | VIEW_RELATIVE)
}

func (e *Element) BorderBoxSize() (width, height int) {
	l, t, r, b := e.getRect(BORDER_BOX)
	return int(r - l), int(b - t)
}

func (e *Element) MarginBox() (left, top, right, bottom int) {
	return e.getRect(MARGIN_BOX)
}

func (e *Element) MarginViewBox() (left, top, right, bottom int) {
	return e.getRect(MARGIN_BOX | VIEW_RELATIVE)
}

func (e *Element) MarginBoxSize() (width, height int) {
	l, t, r, b := e.getRect(MARGIN_BOX)
	return int(r - l), int(b - t)
}

//
// Functions for retrieving/setting the value in widget input controls
//

type METHOD_PARAMS struct {
	MethodId uint32
	Text     *uint16
	Length   uint32
}

// type METHOD_PARAMS C.METHOD_PARAMS

// typedef struct _METHOD_PARAMS METHOD_PARAMS;
//
/** HTMLayoutCallMethod - calls behavior specific method.
 * \param[in] he \b HELEMENT, element - source of the event.
 * \param[in] params \b METHOD_PARAMS, pointer to method param block
 **/
//
// EXTERN_C HLDOM_RESULT HLAPI HTMLayoutCallBehaviorMethod(
//           HELEMENT he, METHOD_PARAMS* params);
//sys HTMLayoutCallBehaviorMethod(he HELEMENT, params *METHOD_PARAMS) (ret HLDOM_RESULT) = htmlayout.HTMLayoutCallBehaviorMethod

// calls behavior specific method.
func (e *Element) ValueAsString() (string, error) {
	args := &METHOD_PARAMS{MethodId: GET_TEXT_VALUE}
	ret := HTMLayoutCallBehaviorMethod(e.handle, args)
	if ret == HLDOM_OK_NOT_HANDLED {
		domPanic2(ret, "This type of element does not provide data in this way.  Try a <widget>.")
	} else if ret != HLDOM_OK {
		domPanic2(ret, "Could not get text value")
	}
	if args.Text == nil {
		return "", errors.New("Nil string pointer")
	}
	return utf16ToStringLength(args.Text, int(args.Length)), nil
}

func (e *Element) SetValue(value interface{}) {
	switch v := value.(type) {
	case string:
		args := &METHOD_PARAMS{
			MethodId: SET_TEXT_VALUE,
			Text:     stringToUtf16Ptr(v),
			Length:   uint32(len(v)),
		}
		ret := HTMLayoutCallBehaviorMethod(e.handle, args)
		if ret == HLDOM_OK_NOT_HANDLED {
			domPanic2(ret, "This type of element does not accept data in this way.  Try a <widget>.")
		} else if ret != HLDOM_OK {
			domPanic2(ret, "Could not set text value")
		}
	default:
		panic("Don't know how to set values of this type")
	}
}

//
// The following are not strictly wrappers of htmlayout functions, but rather convenience
// functions that are helpful in common use cases
//

func (e *Element) Describe() string {
	s := e.Type()
	if value, exists := e.Attr("id"); exists {
		s += "#" + value
	}
	if value, exists := e.Attr("class"); exists {
		values := strings.Split(value, " ")
		for _, v := range values {
			s += "." + v
		}
	}
	return s
}

// Returns the first of the child elements matching the selector.  If no elements
// match, the function panics
func (e *Element) SelectFirst(selector string) *Element {
	results := e.Select(selector)
	if len(results) == 0 {
		panic(fmt.Sprintf("No elements match selector '%s'", selector))
	}
	return results[0]
}

// Returns the only child element that matches the selector.  If no elements match
// or more than one element matches, the function panics
func (e *Element) SelectUnique(selector string) *Element {
	results := e.Select(selector)
	if len(results) == 0 {
		panic(fmt.Sprintf("No elements match selector '%s'", selector))
	} else if len(results) > 1 {
		panic(fmt.Sprintf("More than one element match selector '%s'", selector))
	}
	return results[0]
}

// A wrapper of SelectUnique that auto-prepends a hash to the provided id.
// Useful when selecting elements base on a programmatically retrieved id (which does
// not already have the hash on it)
func (e *Element) SelectId(id string) *Element {
	return e.SelectUnique(fmt.Sprintf("#%s", id))
}

//
// Functions for manipulating the set of classes applied to this element:
//

// Returns true if the specified class is among those listed in the "class" attribute.
func (e *Element) HasClass(class string) bool {
	if classList, exists := e.Attr("class"); !exists {
		return false
	} else if classes := whitespaceSplitter.FindAllString(classList, -1); classes == nil {
		return false
	} else {
		for _, item := range classes {
			if class == item {
				return true
			}
		}
	}
	return false
}

// Adds the specified class to the classes listed in the "class" attribute, or does nothing
// if this class is already included in the list.
func (e *Element) AddClass(class string) {
	if classList, exists := e.Attr("class"); !exists {
		e.SetAttr("class", class)
	} else if classes := whitespaceSplitter.FindAllString(classList, -1); classes == nil {
		e.SetAttr("class", class)
	} else {
		for _, item := range classes {
			if class == item {
				return
			}
		}
		classes = append(classes, class)
		e.SetAttr("class", strings.Join(classes, " "))
	}
}

// Removes the specified class from the classes listed in the "class" attribute, or does nothing
// if this class is not included in the list.
func (e *Element) RemoveClass(class string) {
	if classList, exists := e.Attr("class"); exists {
		if classes := whitespaceSplitter.FindAllString(classList, -1); classes != nil {
			for i, item := range classes {
				if class == item {
					// Delete the item from the list
					classes = append(classes[:i], classes[i+1:]...)
					e.SetAttr("class", strings.Join(classes, " "))
					return
				}
			}
		}
	}
}
