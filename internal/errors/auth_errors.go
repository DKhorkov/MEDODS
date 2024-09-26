package errors

type RefreshTokenNotFoundError struct {
	message string
}

func (e *RefreshTokenNotFoundError) Error() string {
	if e.message != "" {
		return e.message
	}

	return "refresh token not found"
}
