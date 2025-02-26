package provider

import "strconv"

func IfaceToInt(i interface{}) int {
	// Use for validated values only
	result, _ := strconv.Atoi(i.(string))
	return result
}
