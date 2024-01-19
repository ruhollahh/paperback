package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ruhollahh/paperback/api/httputil"
)

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"environment": h.Config.Env,
			"version":     h.Version,
		},
	}

	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		httputil.ServerError(h.Logger, w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonData)
	if err != nil {
		httputil.ServerError(h.Logger, w, r, err)
		return
	}
}
