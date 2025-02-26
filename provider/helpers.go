package provider

import (
	"fmt"
	"strconv"
)

func IfaceToInt(i interface{}) int {
	// Use for validated values only
	result, _ := strconv.Atoi(i.(string))
	return result
}

func F64ToStr(i interface{}) string {
	return fmt.Sprintf("%.0f", i.(float64))
}
