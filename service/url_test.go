package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Kritvi0208/ShortEdge/model"
	"github.com/Kritvi0208/ShortEdge/service"
	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	urls map[string]model.URL
}

func newMockStore() *mockStore {
	return &mockStore{urls: make(map[string]model.URL)}
}

func (m *mockStore) GetAll(ctx context.Context) ([]model.URL, error) {
	var all []model.URL
	for _, v := range m.urls {
		all = append(all, v)
	}
	return all, nil
}

func (m *mockStore) GetByCode(ctx context.Context, code string) (model.URL, error) {
	val, ok := m.urls[code]
	if !ok {
		return model.URL{}, errors.New("not found")
	}
	return val, nil
}

func (m *mockStore) Create(ctx context.Context, url model.URL) error {
	if _, exists := m.urls[url.Code]; exists {
		return errors.New("code already exists")
	}
	m.urls[url.Code] = url
	return nil
}

func (m *mockStore) Update(ctx context.Context, code string, updated model.URL) error {
	if _, exists := m.urls[code]; !exists {
		return errors.New("not found")
	}
	m.urls[code] = updated
	return nil
}

func (m *mockStore) Delete(ctx context.Context, code string) error {
	if _, exists := m.urls[code]; !exists {
		return errors.New("not found")
	}
	delete(m.urls, code)
	return nil
}

func TestShorten_WithCustomCode(t *testing.T) {
	mock := newMockStore()
	svc := service.New(mock)

	req := model.ShortenRequest{
		LongURL:    "https://example.com",
		CustomCode: "custom123",
		Visibility: "public",
	}

	result, err := svc.Shorten(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, req.CustomCode, result.Code)
	assert.Equal(t, req.LongURL, result.LongURL)
}

func TestShorten_CustomCodeExists(t *testing.T) {
	mock := newMockStore()
	mock.urls["custom123"] = model.URL{Code: "custom123"}
	svc := service.New(mock)

	req := model.ShortenRequest{
		LongURL:    "https://example.com",
		CustomCode: "custom123",
	}

	_, err := svc.Shorten(context.Background(), req)
	assert.Error(t, err)
}

func TestGetByCode(t *testing.T) {
	mock := newMockStore()
	mock.urls["abc123"] = model.URL{Code: "abc123", LongURL: "https://x.com"}
	svc := service.New(mock)

	url, err := svc.GetByCode(context.Background(), "abc123")
	assert.NoError(t, err)
	assert.Equal(t, "https://x.com", url.LongURL)
}

func TestUpdate(t *testing.T) {
	mock := newMockStore()
	mock.urls["abc123"] = model.URL{Code: "abc123", LongURL: "old"}
	svc := service.New(mock)

	req := model.ShortenRequest{
		LongURL:    "https://new.com",
		Visibility: "private",
	}

	updated, err := svc.Update(context.Background(), "abc123", req)
	assert.NoError(t, err)
	assert.Equal(t, "https://new.com", updated.LongURL)
	assert.Equal(t, "private", updated.Visibility)
}

func TestDelete(t *testing.T) {
	mock := newMockStore()
	mock.urls["to-delete"] = model.URL{Code: "to-delete"}
	svc := service.New(mock)

	err := svc.Delete(context.Background(), "to-delete")
	assert.NoError(t, err)

	_, err = svc.GetByCode(context.Background(), "to-delete")
	assert.Error(t, err)
}

func TestGetAll_FilterExpired(t *testing.T) {
	mock := newMockStore()

	now := time.Now()

	mock.urls["active"] = model.URL{
		Code:      "active",
		LongURL:   "https://valid.com",
		ExpiresAt: nil,
	}
	mock.urls["future"] = model.URL{
		Code:      "future",
		LongURL:   "https://future.com",
		ExpiresAt: ptrTime(now.Add(1 * time.Hour)),
	}
	mock.urls["expired"] = model.URL{
		Code:      "expired",
		LongURL:   "https://expired.com",
		ExpiresAt: ptrTime(now.Add(-1 * time.Hour)),
	}

	svc := service.New(mock)

	all, err := svc.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, all, 2)
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
