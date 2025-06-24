package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"urlify/model"
	"urlify/service"

	"gofr.dev/pkg/gofr"
)

type URLHandler struct {
	service      service.URLService
	visitService service.VisitService
}

func NewURLHandler(service service.URLService, visitService service.VisitService) *URLHandler {
	return &URLHandler{
		service:      service,
		visitService: visitService,
	}
}

// GetAll godoc
// @Summary Get all URLs
// @Description Fetch all shortened URLs
// @Tags URL
// @Produce json
// @Success 200 {array} model.URL
// @Router /all [get]
func (h *URLHandler) GetAll(ctx *gofr.Context) (interface{}, error) {
	urls, err := h.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

// Shorten godoc
// @Summary Shorten a URL
// @Description Create a shortened URL with optional custom code, visibility, and expiry
// @Tags URL
// @Accept json
// @Produce json
// @Param body body model.ShortenRequest true "URL details"
// @Success 200 {object} model.URL
// @Failure 400 {object} map[string]string
// @Router /shorten [post]
func (h *URLHandler) Shorten(ctx *gofr.Context) (interface{}, error) {
	var req model.ShortenRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, err
	}

	link, err := h.service.Shorten(ctx, req)

	if err != nil {
		return nil, err
	}

	return link, nil
}

// Redirect godoc
// @Summary Redirect to the original long URL
// @Description Redirects using the short code and logs analytics
// @Tags Redirect
// @Produce json
// @Param code path string true "Short code"
// @Success 302 {string} string "Redirects to long URL"
// @Failure 404 {object} map[string]string
// @Router /{code} [get]
func (h *URLHandler) Redirect(ctx *gofr.Context) (interface{}, error) {
	code := ctx.PathParam("code")

	link, err := h.service.GetByCode(ctx, code)
	if err != nil {
		// Return structured JSON 404
		return map[string]interface{}{
			"error": "URL not found",
		}, nil
	}

	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		return map[string]string{"error": "This link has expired."}, nil
	}

	// Extract IP/User-Agent and log visit
	ip, userAgent := getIPAndUserAgentFromHeaders(ctx)
	browser, device := parseUserAgent(userAgent)
	country := getCountryFromIP(ip)

	visit := model.Visit{
		Code:      code,
		Timestamp: time.Now(),
		IP:        ip,
		Country:   country,
		Browser:   browser,
		Device:    device,
	}

	// Asynchronously log
	go func() {
		_ = h.visitService.LogVisit(ctx, visit)
	}()

	fmt.Printf("%#v\n", ctx)

	// ✅ Return redirect instruction (not auto-redirect)
	return map[string]interface{}{
		"redirect": link.LongURL,
	}, nil
}

func getIPAndUserAgentFromHeaders(ctx *gofr.Context) (ip, userAgent string) {
	ip = ctx.Header("X-Forwarded-For")
	if ip == "" {
		ip = ctx.Header("X-Real-IP")
	}
	if ip == "" {
		ip = "Unknown"
	}

	userAgent = ctx.Header("User-Agent")
	if userAgent == "" {
		userAgent = "Unknown"
	}

	return
}

func parseUserAgent(ua string) (browser string, device string) {
	if strings.Contains(ua, "Mobile") {
		device = "Mobile"
	} else {
		device = "Desktop"
	}

	switch {
	case strings.Contains(ua, "Chrome"):
		browser = "Chrome"
	case strings.Contains(ua, "Firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "Safari"):
		browser = "Safari"
	default:
		browser = "Unknown"
	}

	return
}

func getCountryFromIP(ip string) string {
	if ip == "" || ip == "::1" || ip == "127.0.0.1" || ip == "localhost" {
		return "Localhost"
	}

	url := fmt.Sprintf("https://ipwho.is/%s", ip)

	resp, err := http.Get(url)
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Country string `json:"country"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil || result.Country == "" {
		return "Unknown"
	}

	return result.Country
}

// Update godoc
// @Summary Update a short URL
// @Description Modify long URL or visibility of a short URL
// @Tags URL
// @Accept json
// @Produce json
// @Param code path string true "Short code"
// @Param request body model.ShortenRequest true "Updated URL details"
// @Success 200 {object} model.URL
// @Failure 400 {object} map[string]string
// @Router /update/{code} [put]
func (h *URLHandler) Update(ctx *gofr.Context) (interface{}, error) {
	code := ctx.PathParam("code")

	var req model.ShortenRequest
	if err := ctx.Bind(&req); err != nil {
		return nil, err
	}

	updatedURL, err := h.service.Update(ctx, code, req)
	if err != nil {
		return nil, err
	}

	return updatedURL, nil
}

// Delete godoc
// @Summary Delete a short URL
// @Description Delete a short link by its code
// @Tags URL
// @Param code path string true "Short code"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /delete/{code} [delete]
func (h *URLHandler) Delete(ctx *gofr.Context) (interface{}, error) {
	code := ctx.PathParam("code")

	err := h.service.Delete(ctx, code)
	if err != nil {
		return nil, err
	}

	return map[string]string{"message": "deleted successfully"}, nil
}

func HealthHandler(ctx *gofr.Context) (interface{}, error) {
	return map[string]string{"status": "ok"}, nil
}
