package commands

import "strconv"

func ToolIsStringInteger(str string) (int64, bool) {
	isInt64 := true

	// check if value can be contain by int64

	// start with 0
	if len(str) > 1 && str[0] == '0' {
		isInt64 = false
	}

	// has + char
	if len(str) > 0 && str[0] == '+' {
		isInt64 = false
	}

	// check is int64
	i64, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		isInt64 = false
	}
	return i64, isInt64
}
