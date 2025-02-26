package provider

import (
	"fmt"
	"strconv"
)

func StringIsID(i interface{}, k string) ([]string, []error) {
	i, err := strconv.Atoi(i.(string))
	if err != nil {
		return nil, []error{fmt.Errorf("expected %q to be contain an integer ID integer, got %v", k, i)}
	}

	return nil, nil
}
