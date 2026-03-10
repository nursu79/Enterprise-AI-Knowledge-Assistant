package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/nursu79/go-production-api/internal/config"
	"github.com/nursu79/go-production-api/internal/delivery/http/response"
)

type AIHandler struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewAIHandler(cfg *config.Config) *AIHandler {
	return &AIHandler{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // 60s timeout for LLM / file upload tasks!
		},
	}
}

func (h *AIHandler) RegisterRoutes(r chi.Router) {
	r.Post("/ingest", h.IngestDocument)
	r.Post("/query", h.QueryAssistant)
}

func (h *AIHandler) IngestDocument(w http.ResponseWriter, r *http.Request) {
	// Parse max 10MB file
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		response.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		response.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	io.Copy(part, file)
	writer.Close()
	
	proxyReq, err := http.NewRequest("POST", h.cfg.AIServiceUrl+"/ingest", body)
	if err != nil {
		response.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	
	proxyReq.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := h.httpClient.Do(proxyReq)
	if err != nil {
		response.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	response.RespondJSON(w, resp.StatusCode, result)
}

func (h *AIHandler) QueryAssistant(w http.ResponseWriter, r *http.Request) {
	// Forward JSON payload safely
	proxyReq, err := http.NewRequest("POST", h.cfg.AIServiceUrl+"/query", r.Body)
	if err != nil {
		response.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	proxyReq.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(proxyReq)
	if err != nil {
		response.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	response.RespondJSON(w, resp.StatusCode, result)
}
