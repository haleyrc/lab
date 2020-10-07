package sdk

type ValidationError struct {
	Message string
}

func (ve ValidationError) Error() string {
	return "validation failed"
}
