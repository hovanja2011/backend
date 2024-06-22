package save

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"

	resp "github.com/hovanja2011/move/internal/lib/api/response"

	sl "github.com/hovanja2011/move/internal/lib/logger/sl"
	"github.com/hovanja2011/move/internal/storage"
)

type Request struct {
	idDriver int64  `json:"idDriver" validate:"required, ident"`
	sFrom    string `json:"sFrom" validate:"required, adress"`
	sTo      string `json:"sTo" validate:"required, adress"`
}

type Response struct {
	resp.Response
	Error string `json:"error,omitempty"`
}

type DriveCreator interface {
	CreateDrive(idDriver int64, sFrom string, sTo string) error
}

func New(log *slog.Logger, driveCreator DriveCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

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

		err = driveCreator.CreateDrive(req.idDriver, req.sFrom, req.sTo)
		if errors.Is(err, storage.ErrDriverNotFound) {
			log.Info("driver not found", slog.String("idDriver", fmt.Sprint(req.idDriver)))

			render.JSON(w, r, resp.Error("Driver is not exist"))

			return
		}
	}
}
