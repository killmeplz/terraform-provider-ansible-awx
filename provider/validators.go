package provider

import (
	"fmt"
	"strconv"
)

func StringIsNotEmpty(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if v == "" {
		return nil, []error{fmt.Errorf("expected %q to not be an empty string, got %v", k, i)}
	}

	return nil, nil
}

func StringIsID(i interface{}, k string) ([]string, []error) {
	i, err := strconv.Atoi(i.(string))
	if err != nil {
		return nil, []error{fmt.Errorf("expected %q to be contain an integer ID integer, got %v", k, i)}
	}

	return nil, nil
}
