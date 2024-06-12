package save

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	resp "github.com/hovanja2011/move/internal/lib/api/response"
	sl "github.com/hovanja2011/move/internal/lib/logger"
)

type Request struct {
	driverPhone string `json:"driverPhone" validate:"required, phone"`
	driverName  string `json:"driverName"  validate:"required, name"`
}

type Response struct {
	resp.Response
	Error string `json:"error,omitempty"`
}

type DriverCreator interface {
	CreateDriver(driverPhone string, driverName string) error
}

func New(log *slog.Logger, driverCreator DriverCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.new"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Error(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Error(validateErr))

			render.JSON(w, r, resp.ValidationErrors(validateErr))

			return
		}

	}
}
