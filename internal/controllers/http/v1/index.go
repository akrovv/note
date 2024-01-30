package v1

import (
	"net/http"

	"go.uber.org/zap"
)

type IndexHandler struct {
	Logger *zap.Logger
}

func NewIndexHandler(logger *zap.Logger) *IndexHandler {
	return &IndexHandler{Logger: logger}
}

func (ih *IndexHandler) Index(w http.ResponseWriter, r *http.Request) {
	ih.Logger.Info("listening index.html")
	http.ServeFile(w, r, "front/index.html")
}
