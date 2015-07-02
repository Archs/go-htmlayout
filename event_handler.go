package gohl

type EventHandler struct {
	OnAttached      func(he *Element)
	OnDetached      func(he *Element)
	OnMouse         func(he *Element, params *MouseParams) bool
	OnKey           func(he *Element, params *KeyParams) bool
	OnFocus         func(he *Element, params *FocusParams) bool
	OnDraw          func(he *Element, params *DrawParams) bool
	OnTimer         func(he *Element, params *TimerParams) bool
	OnBehaviorEvent func(he *Element, params *BehaviorEventParams) bool
	OnMethodCall    func(he *Element, params *MethodParams) bool
	OnDataArrived   func(he *Element, params *DataArrivedParams) bool
	OnSize          func(he *Element)
	OnScroll        func(he *Element, params *ScrollParams) bool
	OnExchange      func(he *Element, params *ExchangeParams) bool
	OnGesture       func(he *Element, params *GestureParams) bool
}

func (e *EventHandler) Subscription() uint {
	var subscription uint = 0
	add := func(f interface{}, flag uint) {
		if f != nil {
			subscription |= flag
		}
	}

	// OnAttached and OnDetached purposely omitted, since we must receive these events
	add(e.OnMouse, HANDLE_MOUSE)
	add(e.OnKey, HANDLE_KEY)
	add(e.OnFocus, HANDLE_FOCUS)
	add(e.OnDraw, HANDLE_DRAW)
	add(e.OnTimer, HANDLE_TIMER)
	add(e.OnBehaviorEvent, HANDLE_BEHAVIOR_EVENT)
	add(e.OnMethodCall, HANDLE_METHOD_CALL)
	add(e.OnDataArrived, HANDLE_DATA_ARRIVED)
	add(e.OnSize, HANDLE_SIZE)
	add(e.OnScroll, HANDLE_SCROLL)
	add(e.OnExchange, HANDLE_EXCHANGE)
	add(e.OnGesture, HANDLE_GESTURE)

	return subscription
}
