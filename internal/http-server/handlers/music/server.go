package handlers

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/stepan41k/Testovoe/internal/domain/models"
	"github.com/stepan41k/Testovoe/internal/lib/api/logger/sl"
	resp "github.com/stepan41k/Testovoe/internal/lib/api/response"
	"github.com/stepan41k/Testovoe/internal/service"
)

type Music interface {
	GetSongs(ctx context.Context, song models.Song) (songs []models.Song, err error)
	GetTextSong(ctx context.Context, song models.Song) (verse models.Verse, err error)
	DeleteSong(ctx context.Context, song models.Song) (id int64, err error)
	UpdateSong(ctx context.Context, songDetails models.Song) (id int64, err error)
	AddNewSong(ctx context.Context, song models.Song) (id int64, err error)
}

type MusicHandler struct {
	music Music
	log *slog.Logger
}

func New(music Music, log *slog.Logger) *MusicHandler{
	return &MusicHandler{
		music: music,
		log: log,
	}
}


func (m *MusicHandler) GetSongs(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.music.GetSongs"

		log := m.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.Song

		err := render.Decode(r, &req)
		CheckForErrors(req, w, r, log, err)

		songs, err := m.music.GetSongs(ctx, req)
		if err != nil {
			log.Error("internal error")

			render.JSON(w, r, resp.Response{
				Status: http.StatusInternalServerError,
				Error: "internal error",
			})

			return
		}

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data: songs,
		})	
	}
}


func (m *MusicHandler) GetTextSong(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.music.GetTextSong"

		log := m.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.Song

		err := render.Decode(r, &req)
		CheckForErrors(req, w, r, log, err)

		verse, err := m.music.GetTextSong(ctx, req)

		if err != nil {
			log.Error("internal error")

			render.JSON(w, r, resp.Response{
				Status: http.StatusInternalServerError,
				Error: "internal error",
			})

			return
		}

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data: verse,
		})	
	}
}


func (m *MusicHandler) DeleteSong(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.music.DeleteSong"

		log := m.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.Song

		err := render.Decode(r, &req)
		CheckForErrors(req, w, r, log, err)

		songID, err := m.music.DeleteSong(ctx, req)
		if err != nil {
			log.Error("internal error")

			render.JSON(w, r, resp.Response{
				Status: http.StatusInternalServerError,
				Error: "internal error",
			})

			return
		}

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data: songID,
		})	
	}
}


func (m *MusicHandler) UpdateSong(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.music.UpdateSong"

		log := m.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.Song

		err := render.Decode(r, &req)
		CheckForErrors(req, w, r, log, err)

		songID, err := m.music.UpdateSong(ctx, req)
		if err != nil {
			log.Error("internal error", sl.Err(err))

			render.JSON(w, r, resp.Response{
				Status: http.StatusInternalServerError,
				Error: "internal error",
			})

			return
		}

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data: songID,
		})	
	}
}


func (m *MusicHandler) AddNewSong(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.music.AddNewSong"

		log := m.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.Song

		err := render.Decode(r, &req)
		CheckForErrors(req, w, r, log, err)

		songID, err := m.music.AddNewSong(ctx, req)
		if err != nil {
			if errors.Is(err, service.ErrSongExists) {
				log.Error("song already exists")

				render.JSON(w, r, resp.Response{
					Status: http.StatusConflict,
					Error: "song already exists",
				})

				return
			}

			log.Error("internal error", sl.Err(err))

			render.JSON(w, r, resp.Response{
				Status: http.StatusInternalServerError,
				Error: "internal error",
			})

			return
		}

		render.JSON(w, r, resp.Response{
			Status: http.StatusOK,
			Data: songID,
		})
	}
}


func CheckForErrors(req any, w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {

	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Response{
				Status: http.StatusConflict,
				Error: "empty request",
			})

			return
		}

		log.Error("failed to decode request")
		render.JSON(w, r, resp.Response{
			Status: http.StatusBadRequest,
			Error: "failed to decode request",
		})
		return
	}

	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
				
		log.Error("invalid request", sl.Err(err))

		render.JSON(w, r, resp.ValidationError(validateErr))

		return
	}
}