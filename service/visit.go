package service

import (
	"context"
	"github.com/Kritvi0208/ShortEdge/model"
	"github.com/Kritvi0208/ShortEdge/store"
)

type VisitService interface {
	GetAnalytics(ctx context.Context, code string) ([]model.Visit, error)
	LogVisit(ctx context.Context, visit model.Visit) error
}

type visitService struct {
	store store.Visit
}

func NewVisitService(s store.Visit) VisitService {
	return &visitService{store: s}
}

func (v *visitService) GetAnalytics(ctx context.Context, code string) ([]model.Visit, error) {
	return v.store.GetAnalytics(ctx, code)
}

func (v *visitService) LogVisit(ctx context.Context, visit model.Visit) error {
	return v.store.LogVisit(ctx, visit)
}
