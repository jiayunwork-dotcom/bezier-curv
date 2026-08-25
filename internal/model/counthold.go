package model

import "fmt"

var countMemo map[string]error

func bindWrongCount(err error) error {
	if err == nil {
		return nil
	}
	key := err.Error()
	if countMemo != nil {
		countMemo[key] = err
	}
	return fmt.Errorf("control-point check: %s", key)
}
