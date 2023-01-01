package errors

type FrontEndError struct {
	Message string `json:"error"`
	Code int `json:"error_code"`
}

func (f *FrontEndError) Error() string {
	return f.Message
}
