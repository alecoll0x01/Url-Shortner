package service_test

import (
	"errors"
	"testing"

	"github.com/alecoll0x01/url-shortener/internal/repository"
	"github.com/alecoll0x01/url-shortener/internal/service"
)

func TestShorten_ValidURL(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewURLService(repo)

	entry, err := svc.Shorten("https://golang.org", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if entry.Code == "" {
		t.Error("expected a non-empty code")
	}
	if entry.OriginalURL != "https://golang.org" {
		t.Errorf("expected original URL to be preserved, got %s", entry.OriginalURL)
	}
}

func TestShorten_InvalidURL(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewURLService(repo)

	_, err := svc.Shorten("nao-e-uma-url", "")
	if !errors.Is(err, service.ErrInvalidURL) {
		t.Errorf("expected ErrInvalidURL, got %v", err)
	}
}

func TestShorten_CustomCode(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewURLService(repo)

	entry, err := svc.Shorten("https://go.dev", "godev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Code != "godev" {
		t.Errorf("expected code 'godev', got %s", entry.Code)
	}
}

func TestShorten_CustomCode_Conflict(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewURLService(repo)

	_, _ = svc.Shorten("https://go.dev", "godev")

	_, err := svc.Shorten("https://google.com", "godev")
	if !errors.Is(err, service.ErrCodeConflict) {
		t.Errorf("expected ErrCodeConflict, got %v", err)
	}
}

func TestResolve_IncrementsClicks(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewURLService(repo)

	entry, _ := svc.Shorten("https://go.dev", "godev")

	svc.Resolve(entry.Code)
	svc.Resolve(entry.Code)

	stats, err := svc.Stats(entry.Code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Clicks != 2 {
		t.Errorf("expected 2 clicks, got %d", stats.Clicks)
	}
}

func TestResolve_NotFound(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	svc := service.NewURLService(repo)

	_, err := svc.Resolve("naoexiste")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
