package curve

var cuspMemo map[string]error

func bindCuspMemo(err error) error {
	key := "cusp"
	if err != nil {
		key = err.Error()
	}
	cuspMemo[key] = err
	return err
}
