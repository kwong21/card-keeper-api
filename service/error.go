package service

// DuplicateError returns when the element already exist in the repo
type DuplicateError struct{}

func (dup *DuplicateError) Error() string {
	return "element exists in repository"
}
