package gohl

/*
#include <htmlayout.h>
*/
import "C"

/**
 * ValueBinaryData - retreive integer data of T_BYTES type
 */
// EXTERN_C UINT VALAPI ValueBinaryData( const VALUE* pval, LPCBYTES* pBytes, UINT* pnBytes );
//sys ValueBinaryData(pval *JsonValue, pBytes *uintptr, pnBytes *uint) (ret uint) = htmlayout.ValueBinaryData

func (v *JsonValue) IsElement() bool {
	return v.t == T_DOM_OBJECT
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

func (v *JsonValue) IsString() bool {
	return v.t == T_STRING
}

/**
 * ValueToString - converts value to T_STRING inplace:
 * - CVT_SIMPLE - parse/emit terminal values (T_INT, T_FLOAT, T_LENGTH, T_STRING)
 * - CVT_JSON_LITERAL - parse/emit value using JSON literal rules: {}, [], "string", true, false, null
 * - CVT_JSON_MAP - parse/emit MAP value without enclosing '{' and '}' brackets.
 */
// EXTERN_C UINT VALAPI ValueToString( VALUE* pval, /*VALUE_STRING_CVT_TYPE*/ UINT how );
//sys ValueToString(pval *JsonValue, how uint) (ret uint) = htmlayout.ValueToString

/**
 * ValueStringData - returns string data for T_STRING type
 * For T_FUNCTION returns name of the fuction.
 */
// EXTERN_C UINT VALAPI ValueStringData( const VALUE* pval, LPCWSTR* pChars, UINT* pNumChars );
//sys ValueStringData(pval *JsonValue, pChars **uint16, pNumChars *uint) (ret uint) = htmlayout.ValueStringData

func (v *JsonValue) String() string {
	how := uint(C.CVT_SIMPLE)
	if v.IsString() {
		var pChars *uint16
		num := uint(0)
		ValueStringData(v, &pChars, &num)
		return utf16ToString(pChars)
	}
	t := *v
	ValueToString(&t, how)
	tt := &t
	return tt.String()
}
