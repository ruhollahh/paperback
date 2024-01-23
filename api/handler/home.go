package handler

import (
	"context"
	"net/http"

	"github.com/ruhollahh/paperback/web/views/pages"
)

func (h *Handler) Home(w http.ResponseWriter, _ *http.Request) {
	pages.Home().Render(context.Background(), w)
}
