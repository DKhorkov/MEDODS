package httpcontroller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/DKhorkov/hmtm-sso/pkg/logging"
)

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
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}

	if ip == "" {
		ip = r.RemoteAddr
	}

	return ip
}

func getRequestBody[T any](request *http.Request, logger *slog.Logger, storage T) error {
	err := json.NewDecoder(request.Body).Decode(storage)
	if err != nil {
		logger.Error(
			"JSON decoding error",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return err
	}

	return nil
}
