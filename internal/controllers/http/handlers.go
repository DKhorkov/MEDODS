package httpcontroller

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/DKhorkov/hmtm-sso/pkg/logging"
	customerrors "github.com/DKhorkov/medods/internal/errors"

	"github.com/DKhorkov/medods/internal/entities"
	"github.com/DKhorkov/medods/internal/interfaces"
	"github.com/DKhorkov/medods/internal/security"
)

type TokensHandler struct {
	UseCases interfaces.UseCases
	Logger   *slog.Logger
}

func (handler TokensHandler) GetHandleFunc() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		handler.Logger.Info(
			"Tokens request received",
			"Method", request.Method,
			"URL", request.URL,
			"RequestURI", request.RequestURI,
			"UserAgent", request.UserAgent(),
			"RemoteAddr", request.RemoteAddr,
		)

		if request.Method != http.MethodPost && request.Method != http.MethodPut {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if request.Method == http.MethodPost {
			handler.createTokensHandler(
				writer,
				request,
			)
		} else if request.Method == http.MethodPut {
			handler.refreshTokensHandler(
				writer,
				request,
			)
		}
	}
}

func (handler TokensHandler) createTokensHandler(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var requestBody map[string]string
	if err := getRequestBody(request, handler.Logger, &requestBody); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	guid, found := requestBody["GUID"]
	if !found {
		err := customerrors.ParameterRequiredError{Message: "GUID"}
		handler.Logger.Error(
			"Parameter required",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	data := entities.CreateTokensDTO{
		GUID: guid,
		IP:   getUserIP(request),
	}

	tokens, err := handler.UseCases.CreateTokens(data)
	if err != nil {
		handler.Logger.Error(
			"JSON decoding error",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	tokens.RefreshToken = security.Encode([]byte(tokens.RefreshToken))
	writer.Header().Set("Authorization", "Bearer "+tokens.AccessToken)
	renderJSON(writer, map[string]string{"refreshToken": tokens.RefreshToken})
}

func (handler TokensHandler) refreshTokensHandler(
	writer http.ResponseWriter,
	request *http.Request,
) {
	var requestBody map[string]string
	if err := getRequestBody(request, handler.Logger, &requestBody); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	authorizationHeader := request.Header.Get("Authorization")
	if authorizationHeader == "" {
		err := customerrors.HeaderError{Message: "Authorization"}
		handler.Logger.Error(
			"Authorization header required",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	authorizationHeaderValues := strings.Split(authorizationHeader, " ")
	if len(authorizationHeaderValues) != 2 || authorizationHeaderValues[0] != "Bearer" {
		err := customerrors.HeaderError{Message: "Authorization"}
		handler.Logger.Error(
			"Authorization header is invalid",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken := authorizationHeaderValues[1]

	encodedRefreshToken, found := requestBody["refreshToken"]
	if !found {
		err := customerrors.ParameterRequiredError{Message: "refreshToken"}
		handler.Logger.Error(
			"Parameter required",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	refreshToken, err := security.Decode(encodedRefreshToken)
	if err != nil {
		handler.Logger.Error(
			"Refresh token decoding error",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		http.Error(writer, customerrors.InvalidJWTError{}.Error(), http.StatusBadRequest)
		return
	}

	data := entities.RefreshTokensDTO{
		Tokens: entities.Tokens{
			AccessToken:  accessToken,
			RefreshToken: string(refreshToken),
		},
		IP: getUserIP(request),
	}

	tokens, err := handler.UseCases.RefreshTokens(data)
	if err != nil {
		handler.Logger.Error(
			"Refreshing tokens error",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		var possibleError customerrors.AccessTokenDoesNotBelongToRefreshTokenError
		if errors.As(err, &possibleError) {
			http.Error(writer, possibleError.Error(), http.StatusBadRequest)
		} else {
			http.Error(writer, customerrors.InvalidJWTError{}.Error(), http.StatusBadRequest)
		}

		return
	}

	tokens.RefreshToken = security.Encode([]byte(tokens.RefreshToken))
	writer.Header().Set("Authorization", "Bearer "+tokens.AccessToken)
	renderJSON(writer, map[string]string{"refreshToken": tokens.RefreshToken})
}
