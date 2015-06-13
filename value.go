package gohl

/**
 * ValueBinaryData - retreive integer data of T_BYTES type
 */
// EXTERN_C UINT VALAPI ValueBinaryData( const VALUE* pval, LPCBYTES* pBytes, UINT* pnBytes );
//sys ValueBinaryData(pval *JsonValue, pBytes *uintptr, pnBytes *uint) (ret uint) = htmlayout.ValueBinaryData

func (v *JsonValue) IsElement() bool {
	return v.T == T_DOM_OBJECT
}

func (v *JsonValue) ToElement() *Element {
	var pv uintptr = 0
	var dummy uint = 0
	r := ValueBinaryData(v, &pv, &dummy)
	if r != HV_OK {
		return nil
	}
	return NewElementFromHandle(HELEMENT(pv))
}
