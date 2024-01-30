package v1

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"note/internal/models"
	"note/internal/service"
	"note/internal/tools"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Service interface {
	Get(context.Context, service.GetNote) (*models.Note, error)
	Create(context.Context, service.CreateNote) error
	Update(context.Context, service.UpdateNote) error
	Delete(context.Context, service.DeleteNote) error
	GetAll(context.Context, service.GetNotes) ([]*models.Note, error)
}

type handlers struct {
	noteService Service
	Logger      *zap.Logger
}

func NewNoteHandler(service Service, logger *zap.Logger) handlers {
	return handlers{noteService: service, Logger: logger}
}

const contentType = "application/json"

func (h *handlers) Create(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != contentType {
		h.Logger.Warn("not found application/json header")
		err := tools.ErrorJSON(w, errors.New("not found application/json header"), http.StatusBadRequest)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		h.Logger.Warn(err.Error())
		err = tools.ErrorJSON(w, errors.New("can't read from body"), http.StatusInternalServerError)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	cr := service.CreateNote{}
	err = json.Unmarshal(data, &cr)

	if err != nil {
		h.Logger.Warn(err.Error())
		err = tools.ErrorJSON(w, errors.New("can't read json"), http.StatusInternalServerError)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	err = h.noteService.Create(r.Context(), cr)

	if err != nil {
		h.Logger.Warn(err.Error())
		err = tools.ErrorJSON(w, errors.New("can't create a note"), http.StatusInternalServerError)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	err = tools.WriteJSON(w, struct {
		Response string `json:"response"`
	}{"successfully created"})

	if err != nil {
		h.Logger.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *handlers) UpdateByID(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-type") != contentType {
		h.Logger.Warn("not found application/json header")
		err := tools.ErrorJSON(w, errors.New("not found application/json header"), http.StatusBadRequest)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	vars := mux.Vars(r)

	id, in := vars["id"]

	if !in {
		h.Logger.Warn("id not found")
		err := tools.ErrorJSON(w, errors.New("id not found"), http.StatusBadRequest)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		h.Logger.Warn(err.Error())
		err = tools.ErrorJSON(w, errors.New("can't read from body"), http.StatusInternalServerError)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	up := service.UpdateNote{ID: id}
	err = json.Unmarshal(data, &up)

	if err != nil {
		h.Logger.Warn(err.Error())
		err = tools.ErrorJSON(w, errors.New("can't read the json"), http.StatusInternalServerError)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	err = h.noteService.Update(r.Context(), up)

	if err != nil {
		h.Logger.Warn(err.Error())
		err = tools.ErrorJSON(w, errors.New("can't update a note"), http.StatusInternalServerError)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	err = tools.WriteJSON(w, struct {
		Response string `json:"response"`
	}{"successfully updated"})

	if err != nil {
		h.Logger.Warn("server can't write")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *handlers) DeleteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, in := vars["id"]

	if !in {
		h.Logger.Warn("id not found")
		err := tools.ErrorJSON(w, errors.New("id not found"), http.StatusBadRequest)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	dn := service.DeleteNote{ID: id}
	err := h.noteService.Delete(r.Context(), dn)

	if err != nil {
		h.Logger.Warn(err.Error())
		err = tools.ErrorJSON(w, errors.New("can't delete a note"), http.StatusInternalServerError)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	err = tools.WriteJSON(w, struct {
		Response string `json:"response"`
	}{"successfully deleted"})

	if err != nil {
		h.Logger.Warn("server can't write")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, in := vars["id"]

	if !in {
		h.Logger.Warn("id not found")
		err := tools.ErrorJSON(w, errors.New("id not found"), http.StatusBadRequest)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	gn := service.GetNote{ID: id}
	note, err := h.noteService.Get(r.Context(), gn)

	if err != nil {
		h.Logger.Warn(err.Error())
		err = tools.ErrorJSON(w, errors.New("can't get a note"), http.StatusInternalServerError)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	err = tools.WriteJSON(w, note)

	if err != nil {
		h.Logger.Warn("server can't write")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *handlers) GetAll(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query().Get("order_by")
	notes, err := h.noteService.GetAll(r.Context(), service.GetNotes{OrderBy: queries})

	if err != nil {
		h.Logger.Warn(err.Error())
		err = tools.ErrorJSON(w, errors.New("can't get notes"), http.StatusInternalServerError)

		if err != nil {
			h.Logger.Warn(err.Error())
		}

		return
	}

	err = tools.WriteJSON(w, notes)

	if err != nil {
		h.Logger.Warn("server can't write")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
