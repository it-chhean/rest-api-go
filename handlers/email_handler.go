package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"email-api/models"
)

// EmailRepository is the contract the handlers depend on.
type EmailRepository interface {
	GetAll() ([]models.Email, error)
	GetByID(id int) (*models.Email, error)
	Create(e *models.Email) error
	Update(e *models.Email) error
	Delete(id int) error
}

// Handler holds shared dependencies for all email endpoints.
type Handler struct {
	store EmailRepository
}

// NewHandler constructs a Handler with the given store.
func NewHandler(store EmailRepository) *Handler {
	return &Handler{store: store}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "invalid id: must be an integer", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func decodeBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}


func (h *Handler) GetEmails(w http.ResponseWriter, r *http.Request) {
	emails, err := h.store.GetAll()
	if err != nil {
		http.Error(w, "could not retrieve emails", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, emails)
}

func (h *Handler) GetEmailByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	email, err := h.store.GetByID(id)
	if errors.Is(err, models.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "could not retrieve email", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, email)
}

func (h *Handler) CreateEmail(w http.ResponseWriter, r *http.Request) {
	var email models.Email
	if !decodeBody(w, r, &email) {
		return
	}
	if err := email.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err := h.store.Create(&email); err != nil {
		http.Error(w, "could not create email", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, email)
}

func (h *Handler) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var updated models.Email
	if !decodeBody(w, r, &updated) {
		return
	}
	if err := updated.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	updated.ID = id
	if err := h.store.Update(&updated); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "could not update email", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) DeleteEmail(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	if err := h.store.Delete(id); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "could not delete email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
