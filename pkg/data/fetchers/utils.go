package fetchers

func unwrap(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}

func orElse(s string, err error) string {
	if err != nil {
		return err.Error()
	}
	return s
}
