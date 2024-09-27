package errors

type RefreshTokenNotFoundError struct {
	Message string
}

func (e RefreshTokenNotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "refresh token not found"
}

type ParameterRequiredError struct {
	Message string
}

func (e ParameterRequiredError) Error() string {
	if e.Message != "" {
		return e.Message + " parameter is missing"
	}

	return "required parameter is missing"
}

type IPAddressDoesNotMatchWithTokensIPError struct {
	Message string
}

func (e IPAddressDoesNotMatchWithTokensIPError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "IP address does not match with token IP address"
}

type AccessTokenDoesNotBelongToRefreshTokenError struct {
	Message string
}

func (e AccessTokenDoesNotBelongToRefreshTokenError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "Access token does not belong to refresh token"
}

type HeaderError struct {
	Message string
}

func (e HeaderError) Error() string {
	if e.Message != "" {
		return e.Message + " Header is missing or invalid"
	}

	return "required Header is missing or invalid"
}
