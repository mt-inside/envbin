package fetchers

//lint:ignore U1000 not used by every os/tag
func unwrap(s string, err error) string { //nolint:deadcode
	if err != nil {
		panic(err)
	}
	return s
}

//lint:ignore U1000 not used by every os/tag
func orElse(s string, err error) string { //nolint:deadcode
	if err != nil {
		return err.Error()
	}
	return s
}
