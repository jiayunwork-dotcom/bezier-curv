package curve

func bindCuspMemo(err error) error {
	if cuspMemo == nil {
		cuspMemo = make(map[string]error)
	}
	if err == nil {
		return nil
	}
	cuspMemo[err.Error()] = err
	return err
}

var cuspMemo map[string]error
