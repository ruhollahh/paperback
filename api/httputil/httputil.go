package httputil

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/ruhollahh/paperback/api/contextutil"
	"github.com/ruhollahh/paperback/internal/app/service"

	"github.com/go-playground/form/v4"
)

func LogError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	logger.Error(err.Error(), "uri", uri, "method", method, "trace", trace)
}

func ServerError(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	LogError(logger, w, r, err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func NotFoundError(w http.ResponseWriter) {
	ClientError(w, http.StatusNotFound)
}

func DecodePostForm(formDecoder *form.Decoder, r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = formDecoder.Decode(&dst, r.PostForm)
	if err != nil {
		var invalidDecoderErr *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderErr) {
			panic(err)
		}
		return err
	}

	return nil
}

func IsAuthenticated(r *http.Request) bool {
	user := contextutil.ContextGetUser(r.Context())
	if user != nil && service.IsAnonymous(user) {
		return false
	}
	return true
}
