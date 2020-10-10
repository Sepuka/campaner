package errors

func NewStorageError(repo string, err error) error {
	return campanerError{
		errorType:     Storage,
		originalError: err,
		context: map[string]string{
			`repository`: repo,
		},
	}
}
