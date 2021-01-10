package cardservice

// DuplicateError returns when the element already exist in the repo
type DuplicateError struct {
	Message string
}

// UnknownError returns for error not handled yet by the API
type UnknownError struct {
	Message string
}

func (duplicate *DuplicateError) Error() string {
	return "element exists in repository. original message = " + duplicate.Message
}

func (unknown *UnknownError) Error() string {
	message := "Unhandled error, mesage from client = " + unknown.Message
	return message
}
