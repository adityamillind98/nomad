// Code generated by "stringer -trimprefix=UiColor -output=colored_ui_string_uicolorwhen.go -linecomment -type=UiColorWhen"; DO NOT EDIT.

package ui

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UiColorAuto-0]
	_ = x[UiColorNever-1]
	_ = x[UiColorAlways-2]
}

const _UiColorWhen_name = "AutoNeverAlways"

var _UiColorWhen_index = [...]uint8{0, 4, 9, 15}

func (i UiColorWhen) String() string {
	if i < 0 || i >= UiColorWhen(len(_UiColorWhen_index)-1) {
		return "UiColorWhen(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _UiColorWhen_name[_UiColorWhen_index[i]:_UiColorWhen_index[i+1]]
}
