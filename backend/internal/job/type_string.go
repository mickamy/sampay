// Code generated by "stringer -type=Type"; DO NOT EDIT.

package job

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[first-0]
	_ = x[SendEmailJob-1]
	_ = x[last-2]
}

const _Type_name = "firstSendEmailJoblast"

var _Type_index = [...]uint8{0, 5, 17, 21}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
