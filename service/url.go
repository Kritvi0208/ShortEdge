package service

import (
	"context"
	"fmt"
	"github.com/Kritvi0208/ShortEdge/model"
	"github.com/Kritvi0208/ShortEdge/store"
	"math/rand"
	"strings"
	"time"
)

type URLService interface {
	GetAll(ctx context.Context) ([]model.URL, error)
	Shorten(ctx context.Context, req model.ShortenRequest) (model.URL, error)
	GetByCode(ctx context.Context, code string) (model.URL, error)
	Update(ctx context.Context, code string, req model.ShortenRequest) (model.URL, error)
	Delete(ctx context.Context, code string) error
}

type urlService struct {
	store store.URL
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func New(s store.URL) URLService {
	return &urlService{store: s}
}

func (u *urlService) GetAll(ctx context.Context) ([]model.URL, error) {
	all, err := u.store.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var valid []model.URL

	for _, url := range all {
		// Include only if:
		// - There's no expiry (ExpiresAt == nil)
		// - OR expiry is in the future
		if url.ExpiresAt == nil || url.ExpiresAt.After(now) {
			valid = append(valid, url)
		}
	}

	return valid, nil
}

func (u *urlService) Shorten(ctx context.Context, req model.ShortenRequest) (model.URL, error) {
	// Validation
	if req.LongURL == "" {
		return model.URL{}, fmt.Errorf("long URL is required")
	}

	code := req.CustomCode

	// If custom code provided, check if it exists
	if code != "" {
		_, err := u.store.GetByCode(ctx, code)
		if err == nil {
			return model.URL{}, fmt.Errorf("custom code '%s' already exists", code)
		}
	}

	// If no custom code, generate unique code with retries
	if code == "" {
		const maxRetries = 5
		for i := 0; i < maxRetries; i++ {
			code = generateCode()
			_, err := u.store.GetByCode(ctx, code)
			if err != nil {
				// Not found, safe to use this code
				break
			}
			code = "" // reset to retry
		}
		if code == "" {
			return model.URL{}, fmt.Errorf("failed to generate unique short code after %d attempts", maxRetries)
		}
	}

	// Normalize visibility
	visibility := strings.ToLower(req.Visibility)
	if visibility != "private" {
		visibility = "public"
	}
	link := model.URL{
		Code:       code,
		LongURL:    req.LongURL,
		Visibility: visibility,
		CreatedAt:  time.Now(),
		ExpiresAt:  req.ExpiresAt, // ⏳
	}

	err := u.store.Create(ctx, link)
	return link, err
}

func (u *urlService) GetByCode(ctx context.Context, code string) (model.URL, error) {
	return u.store.GetByCode(ctx, code)
}

func (u *urlService) Update(ctx context.Context, code string, req model.ShortenRequest) (model.URL, error) {
	existing, err := u.store.GetByCode(ctx, code)
	if err != nil {
		return model.URL{}, fmt.Errorf("short code not found")
	}

	// Update fields
	existing.LongURL = req.LongURL
	existing.Visibility = req.Visibility
	existing.CreatedAt = time.Now()

	err = u.store.Update(ctx, code, existing)
	if err != nil {
		return model.URL{}, err
	}

	return existing, nil
}

func (u *urlService) Delete(ctx context.Context, code string) error {
	return u.store.Delete(ctx, code)
}

func generateCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
