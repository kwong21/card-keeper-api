package service

// DuplicateError returns when the element already exist in the repo
type DuplicateError struct{}

// UnknownError returns for error not handled yet by the API
type UnknownError struct {
	Message string
}

func (dup *DuplicateError) Error() string {
	return "element exists in repository"
}

func (unknown *UnknownError) Error() string {
	message := "Unhandled error, mesage from client = " + unknown.Message
	return message
}
