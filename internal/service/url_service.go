package service

import (
	"errors"
	"math/rand"
	"net/url"
	"time"

	"github.com/alecoll0x01/internal/repository"
)

var (
	ErrInvalidURL   = errors.New("invalid url")
	ErrCodeConflict = errors.New("code already exists")
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type URLService struct {
	repo repository.URLRepository
	rng  *rand.Rand
}

func NewURLService(repo repository.URLRepository) *URLService {
	return &URLService{
		repo: repo,
		rng:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Shorten cria um código curto para a URL fornecida
func (s *URLService) Shorten(rawURL string, customCode string) (repository.URL, error) {
	// Valida a URL
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return repository.URL{}, ErrInvalidURL
	}

	code := customCode
	if code == "" {
		code = s.generateCode(6)
	} else {
		// Verifica se o código customizado já existe
		if _, err := s.repo.FindByCode(code); err == nil {
			return repository.URL{}, ErrCodeConflict
		}
	}

	entry := repository.URL{
		Code:        code,
		OriginalURL: rawURL,
	}

	if err := s.repo.Save(entry); err != nil {
		return repository.URL{}, err
	}

	return entry, nil
}

// Resolve retorna a URL original e incrementa os cliques
func (s *URLService) Resolve(code string) (repository.URL, error) {
	entry, err := s.repo.FindByCode(code)
	if err != nil {
		return repository.URL{}, err
	}

	_ = s.repo.IncrementClicks(code)

	return entry, nil
}

// Stats retorna os dados de uma URL pelo código
func (s *URLService) Stats(code string) (repository.URL, error) {
	return s.repo.FindByCode(code)
}

// ListAll retorna todas as URLs cadastradas
func (s *URLService) ListAll() []repository.URL {
	return s.repo.ListAll()
}

func (s *URLService) generateCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[s.rng.Intn(len(charset))]
	}
	return string(b)
}
