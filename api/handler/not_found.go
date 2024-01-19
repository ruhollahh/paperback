package handler

import (
	"github.com/ruhollahh/paperback/api/httputil"
	"net/http"
)

func (h *Handler) NotFound(w http.ResponseWriter, _ *http.Request) {
	httputil.NotFoundError(w)
}
