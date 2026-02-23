package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alecoll0x01/url-shortener/internal/repository"
	"github.com/alecoll0x01/url-shortener/internal/service"
	"github.com/alecoll0x01/url-shortener/pkg/response"
	"github.com/go-chi/chi/v5"
)

type URLHandler struct {
	svc     *service.URLService
	baseURL string
}

func NewURLHandler(svc *service.URLService, baseURL string) *URLHandler {
	return &URLHandler{svc: svc, baseURL: baseURL}
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var body struct {
		URL        string `json:"url"`
		CustomCode string `json:"custom_code,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	entry, err := h.svc.Shorten(body.URL, body.CustomCode)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidURL):
			response.Error(w, http.StatusBadRequest, "invalid or unsupported URL")
		case errors.Is(err, service.ErrCodeConflict):
			response.Error(w, http.StatusConflict, "custom code already in use")
		default:
			response.Error(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{
		"code":         entry.Code,
		"original_url": entry.OriginalURL,
		"short_url":    h.baseURL + "/" + entry.Code,
	})
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	entry, err := h.svc.Resolve(code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "short URL not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	http.Redirect(w, r, entry.OriginalURL, http.StatusMovedPermanently)
}

func (h *URLHandler) Stats(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	entry, err := h.svc.Stats(code)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "short URL not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"code":         entry.Code,
		"original_url": entry.OriginalURL,
		"clicks":       entry.Clicks,
	})
}

func (h *URLHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	urls := h.svc.ListAll()
	response.JSON(w, http.StatusOK, urls)
}

func (h *URLHandler) Delete(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	if err := h.svc.Delete(code); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "short URL not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
