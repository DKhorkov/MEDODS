package httpcontroller

import (
	"encoding/json"
	"github.com/DKhorkov/medods/internal/entities"
	"github.com/DKhorkov/medods/internal/interfaces"
	"github.com/DKhorkov/medods/internal/security"
	"net/http"
)

func GetTokensHandler(useCases interfaces.UseCases) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost && request.Method != http.MethodPut {
			http.Error(writer, "Method not allowed", http.StatusInternalServerError)
			return
		}

		IP := getUserIP(request)

		var requestBody map[string]string
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if request.Method == http.MethodPost {
			createTokensHandlerProcessor(
				writer,
				useCases,
				IP,
				requestBody,
			)
		} else if request.Method == http.MethodPut {
			refreshTokensHandlerProcessor(
				writer,
				useCases,
				IP,
				requestBody,
			)
		}
	}
}

func createTokensHandlerProcessor(
	writer http.ResponseWriter,
	useCases interfaces.UseCases,
	IP string,
	requestBody map[string]string,
) {

	GUID, found := requestBody["GUID"]
	if !found {
		http.Error(writer, "Missing GUID", http.StatusBadRequest)
		return
	}

	data := entities.CreateTokensDTO{
		GUID: GUID,
		IP:   IP,
	}

	tokens, err := useCases.CreateTokens(data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	tokens.RefreshToken = security.Encode([]byte(tokens.RefreshToken))
	renderJSON(writer, tokens)
	return
}

func refreshTokensHandlerProcessor(
	writer http.ResponseWriter,
	useCases interfaces.UseCases,
	IP string,
	requestBody map[string]string,
) {

	accessToken, found := requestBody["accessToken"]
	if !found {
		http.Error(writer, "Missing access token", http.StatusBadRequest)
		return
	}

	decodedRefreshToken, found := requestBody["refreshToken"]
	if !found {
		http.Error(writer, "Missing refresh token", http.StatusBadRequest)
		return
	}

	refreshToken, err := security.Decode(decodedRefreshToken)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	data := entities.RefreshTokensDTO{
		Tokens: entities.Tokens{
			AccessToken:  accessToken,
			RefreshToken: string(refreshToken),
		},
		IP: IP,
	}

	tokes, err := useCases.RefreshTokens(data)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}

	renderJSON(writer, tokes)
	return
}

// renderJSON преобразует 'v' в формат JSON и записывает результат, в виде ответа, в w.
func renderJSON(writer http.ResponseWriter, value interface{}) {
	jsonResponse, err := json.Marshal(value)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if _, err = writer.Write(jsonResponse); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

// getUserIP retrieves IP address from request.
func getUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}

	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	return IPAddress
}
