package handler

import (
	"urlify/service"

	"gofr.dev/pkg/gofr"
)

type VisitHandler struct {
	service service.VisitService
}

func NewVisitHandler(s service.VisitService) *VisitHandler {
	return &VisitHandler{service: s}
}

// GetAnalytics godoc
// @Summary Get analytics for a short URL
// @Description Fetch real-time visit logs for a short link
// @Tags Analytics
// @Produce json
// @Param code path string true "Short code"
// @Success 200 {array} model.Visit
// @Failure 400 {object} map[string]string
// @Router /analytics/{code} [get]
func (h *VisitHandler) GetAnalytics(ctx *gofr.Context) (interface{}, error) {
    code := ctx.PathParam("code")
    visits, err := h.service.GetAnalytics(ctx.Context, code)
    if err != nil {
        return nil, err
    }
    return visits, nil
}
