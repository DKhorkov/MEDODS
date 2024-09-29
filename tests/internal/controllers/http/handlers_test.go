package controllers__test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/DKhorkov/hmtm-sso/pkg/logging"
	httpcontroller "github.com/DKhorkov/medods/internal/controllers/http"
	"github.com/DKhorkov/medods/internal/entities"
	customerrors "github.com/DKhorkov/medods/internal/errors"
	mocks "github.com/DKhorkov/medods/internal/mocks/repositories"
	"github.com/DKhorkov/medods/internal/security"
	"github.com/DKhorkov/medods/internal/services"
	"github.com/DKhorkov/medods/internal/usecases"
	testconfig "github.com/DKhorkov/medods/tests/config"
	"github.com/stretchr/testify/assert"
)

var testsConfig = testconfig.New()

func TestControllersHTTPTokensHandlerInvalidMethod(t *testing.T) {
	t.Run("invalid HTTP method", func(t *testing.T) {
		logger := logging.GetInstance(testsConfig.Logging.Level, testsConfig.Logging.LogFilePath)
		useCases := &usecases.CommonUseCases{Logger: logger}

		request := httptest.NewRequest(
			http.MethodGet,
			"/tokens",
			nil,
		)

		writer := httptest.NewRecorder()
		handleFunc := httpcontroller.TokensHandler{UseCases: useCases, Logger: logger}.GetHandleFunc()
		handleFunc(writer, request)

		result := writer.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)
	})
}

func TestControllersHTTPTokensHandlerCreateTokens(t *testing.T) {
	t.Run("successfully create tokens", func(t *testing.T) {
		authRepository := &mocks.MockedAuthRepository{RefreshTokensStorage: map[int]*entities.RefreshToken{}}
		usersRepository := &mocks.MockedUsersRepository{}
		authService := &services.CommonAuthService{AuthRepository: authRepository}
		usersService := &services.CommonUsersService{UsersRepository: usersRepository}
		logger := logging.GetInstance(testsConfig.Logging.Level, testsConfig.Logging.LogFilePath)
		useCases := &usecases.CommonUseCases{
			AuthService:  authService,
			UsersService: usersService,
			HashCost:     testsConfig.HashCost,
			JWTConfig:    testsConfig.JWT,
			SMTPConfig:   testsConfig.SMTP,
			Logger:       logger,
		}

		bodyData := map[string]interface{}{"GUID": testsConfig.RefreshToken.GUID}
		body, err := json.Marshal(bodyData)
		if err != nil {
			t.Fatal(err)
		}

		request := httptest.NewRequest(
			http.MethodPost,
			"/tokens",
			strings.NewReader(string(body)),
		)

		writer := httptest.NewRecorder()
		handleFunc := httpcontroller.TokensHandler{UseCases: useCases, Logger: logger}.GetHandleFunc()
		handleFunc(writer, request)

		result := writer.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusOK, result.StatusCode)

		authorizationHeader := result.Header.Get("Authorization")
		assert.NotEqual(t, "", authorizationHeader)

		authorizationHeaderValues := strings.Split(authorizationHeader, " ")
		assert.Len(t, authorizationHeaderValues, 2)
		assert.Equal(t, "Bearer", authorizationHeaderValues[0])
		assert.NotEqual(t, "", authorizationHeaderValues[1])

		responseBodyData, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		var responseBody = struct {
			RefreshToken string `json:"refreshToken"`
		}{}

		if err = json.Unmarshal(responseBodyData, &responseBody); err != nil {
			t.Fatal(err)
		}

		assert.NotEqual(t, "", responseBody.RefreshToken)
	})

	t.Run("GUID parameter required", func(t *testing.T) {
		authRepository := &mocks.MockedAuthRepository{RefreshTokensStorage: map[int]*entities.RefreshToken{}}
		usersRepository := &mocks.MockedUsersRepository{}
		authService := &services.CommonAuthService{AuthRepository: authRepository}
		usersService := &services.CommonUsersService{UsersRepository: usersRepository}
		logger := logging.GetInstance(testsConfig.Logging.Level, testsConfig.Logging.LogFilePath)
		useCases := &usecases.CommonUseCases{
			AuthService:  authService,
			UsersService: usersService,
			HashCost:     testsConfig.HashCost,
			JWTConfig:    testsConfig.JWT,
			SMTPConfig:   testsConfig.SMTP,
			Logger:       logger,
		}

		request := httptest.NewRequest(
			http.MethodPost,
			"/tokens",
			nil,
		)

		writer := httptest.NewRecorder()
		handleFunc := httpcontroller.TokensHandler{UseCases: useCases, Logger: logger}.GetHandleFunc()
		handleFunc(writer, request)

		result := writer.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
		responseBodyData, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		errorMessage := strings.Split(string(responseBodyData), "\n")[1]
		assert.Equal(t, customerrors.ParameterRequiredError{Parameter: "GUID"}.Error(), errorMessage)
	})
}

func TestControllersHTTPTokensHandlerRefreshTokens(t *testing.T) {
	t.Run("successfully refresh tokens", func(t *testing.T) {
		hashedRefreshTokenValue, err := security.HashRefreshToken(testsConfig.RefreshToken.Value, testsConfig.HashCost)
		if err != nil {
			t.Fatal(err)
		}

		dbRefreshToken := &entities.RefreshToken{
			ID:    1,
			Value: hashedRefreshTokenValue,
			TTL:   time.Now().Add(time.Hour),
			GUID:  testsConfig.RefreshToken.GUID,
		}

		authRepository := &mocks.MockedAuthRepository{
			RefreshTokensStorage: map[int]*entities.RefreshToken{
				dbRefreshToken.ID: dbRefreshToken,
			},
		}

		usersRepository := &mocks.MockedUsersRepository{}
		authService := &services.CommonAuthService{AuthRepository: authRepository}
		usersService := &services.CommonUsersService{UsersRepository: usersRepository}
		logger := logging.GetInstance(testsConfig.Logging.Level, testsConfig.Logging.LogFilePath)
		useCases := &usecases.CommonUseCases{
			AuthService:  authService,
			UsersService: usersService,
			HashCost:     testsConfig.HashCost,
			JWTConfig:    testsConfig.JWT,
			SMTPConfig:   testsConfig.SMTP,
			Logger:       logger,
		}

		refreshToken, err := security.GenerateJWT(
			security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     testsConfig.RefreshToken.Value,
				GUID:      testsConfig.RefreshToken.GUID,
			},
		)

		if err != nil {
			t.Fatal(err)
		}

		encodedRefreshToken := security.Encode([]byte(refreshToken))

		accessToken, err := security.GenerateJWT(
			security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     strconv.Itoa(dbRefreshToken.ID),
				GUID:      testsConfig.RefreshToken.GUID,
			},
		)

		if err != nil {
			t.Fatal(err)
		}

		bodyData := map[string]interface{}{"refreshToken": encodedRefreshToken}
		body, err := json.Marshal(bodyData)
		if err != nil {
			t.Fatal(err)
		}

		request := httptest.NewRequest(
			http.MethodPut,
			"/tokens",
			strings.NewReader(string(body)),
		)

		request.Header.Set("Authorization", "Bearer "+accessToken)
		request.Header.Set("X-Real-Ip", testsConfig.IP)
		writer := httptest.NewRecorder()
		handleFunc := httpcontroller.TokensHandler{UseCases: useCases, Logger: logger}.GetHandleFunc()
		handleFunc(writer, request)

		result := writer.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusOK, result.StatusCode)

		authorizationHeader := result.Header.Get("Authorization")
		assert.NotEqual(t, "", authorizationHeader)

		authorizationHeaderValues := strings.Split(authorizationHeader, " ")
		assert.Len(t, authorizationHeaderValues, 2)
		assert.Equal(t, "Bearer", authorizationHeaderValues[0])
		assert.NotEqual(t, "", authorizationHeaderValues[1])

		responseBodyData, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		var responseBody = struct {
			RefreshToken string `json:"refreshToken"`
		}{}

		if err = json.Unmarshal(responseBodyData, &responseBody); err != nil {
			t.Fatal(err)
		}

		assert.NotEqual(t, "", responseBody.RefreshToken)
	})

	t.Run("accessToken Header required", func(t *testing.T) {
		logger := logging.GetInstance(testsConfig.Logging.Level, testsConfig.Logging.LogFilePath)
		useCases := &usecases.CommonUseCases{Logger: logger}

		request := httptest.NewRequest(
			http.MethodPut,
			"/tokens",
			nil,
		)

		writer := httptest.NewRecorder()
		handleFunc := httpcontroller.TokensHandler{UseCases: useCases, Logger: logger}.GetHandleFunc()
		handleFunc(writer, request)

		result := writer.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
		responseBodyData, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		errorMessage := strings.Split(string(responseBodyData), "\n")[1]
		assert.Equal(t, customerrors.HeaderError{Header: "Authorization"}.Error(), errorMessage)
	})

	t.Run("Authorization header Transport error", func(t *testing.T) {
		logger := logging.GetInstance(testsConfig.Logging.Level, testsConfig.Logging.LogFilePath)
		useCases := &usecases.CommonUseCases{Logger: logger}

		accessToken, err := security.GenerateJWT(
			security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     strconv.Itoa(1),
				GUID:      testsConfig.RefreshToken.GUID,
			},
		)

		if err != nil {
			t.Fatal(err)
		}

		request := httptest.NewRequest(
			http.MethodPut,
			"/tokens",
			nil,
		)

		request.Header.Set("Authorization", "InvalidTransport "+accessToken)
		writer := httptest.NewRecorder()
		handleFunc := httpcontroller.TokensHandler{UseCases: useCases, Logger: logger}.GetHandleFunc()
		handleFunc(writer, request)

		result := writer.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
		responseBodyData, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		errorMessage := strings.Split(string(responseBodyData), "\n")[1]
		assert.Equal(t, customerrors.HeaderError{Header: "Authorization"}.Error(), errorMessage)
	})

	t.Run("refreshToken parameter required", func(t *testing.T) {
		logger := logging.GetInstance(testsConfig.Logging.Level, testsConfig.Logging.LogFilePath)
		useCases := &usecases.CommonUseCases{Logger: logger}

		accessToken, err := security.GenerateJWT(
			security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     strconv.Itoa(1),
				GUID:      testsConfig.RefreshToken.GUID,
			},
		)

		if err != nil {
			t.Fatal(err)
		}

		request := httptest.NewRequest(
			http.MethodPut,
			"/tokens",
			nil,
		)

		request.Header.Set("Authorization", "Bearer "+accessToken)
		writer := httptest.NewRecorder()
		handleFunc := httpcontroller.TokensHandler{UseCases: useCases, Logger: logger}.GetHandleFunc()
		handleFunc(writer, request)

		result := writer.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
		responseBodyData, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		errorMessage := strings.Split(string(responseBodyData), "\n")[1]
		assert.Equal(t, customerrors.ParameterRequiredError{Parameter: "refreshToken"}.Error(), errorMessage)
	})

	t.Run("refreshToken parameter required", func(t *testing.T) {
		logger := logging.GetInstance(testsConfig.Logging.Level, testsConfig.Logging.LogFilePath)
		useCases := &usecases.CommonUseCases{Logger: logger}

		accessToken, err := security.GenerateJWT(
			security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     strconv.Itoa(1),
				GUID:      testsConfig.RefreshToken.GUID,
			},
		)

		if err != nil {
			t.Fatal(err)
		}

		bodyData := map[string]interface{}{"refreshToken": "invalid encoded refreshToken"}
		body, err := json.Marshal(bodyData)
		if err != nil {
			t.Fatal(err)
		}

		request := httptest.NewRequest(
			http.MethodPut,
			"/tokens",
			strings.NewReader(string(body)),
		)

		request.Header.Set("Authorization", "Bearer "+accessToken)
		writer := httptest.NewRecorder()
		handleFunc := httpcontroller.TokensHandler{UseCases: useCases, Logger: logger}.GetHandleFunc()
		handleFunc(writer, request)

		result := writer.Result()
		defer result.Body.Close()

		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
		responseBodyData, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		errorMessage := strings.Split(string(responseBodyData), "\n")[0]
		assert.Equal(t, customerrors.InvalidJWTError{}.Error(), errorMessage)
	})
}
