package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nursu79/go-production-api/internal/config"
	"github.com/nursu79/go-production-api/internal/delivery/http/response"
	"github.com/nursu79/go-production-api/internal/domain"
	"github.com/nursu79/go-production-api/internal/usecase"
)

type AIHandler struct {
	cfg        *config.Config
	httpClient *http.Client
	usecase    usecase.ChatHistoryUsecase
}

func NewAIHandler(cfg *config.Config, uc usecase.ChatHistoryUsecase) *AIHandler {
	return &AIHandler{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // 60s timeout for LLM / file upload tasks!
		},
		usecase: uc,
	}
}

func (h *AIHandler) RegisterRoutes(r chi.Router) {
	r.Post("/ingest", h.IngestDocument)
	r.Post("/query", h.QueryAssistant)
	r.Get("/history", h.GetHistory)
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
	// We must read the request body first so we can parse out the original user query
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		response.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "Failed to read request body"})
		return
	}

	var userReq struct {
		Query string `json:"query"`
	}
	// We ignore the error just in case, but keep what we can
	json.Unmarshal(reqBody, &userReq)

	// Recreate proxy request
	proxyReq, err := http.NewRequest("POST", h.cfg.AIServiceUrl+"/query", bytes.NewBuffer(reqBody))
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
	// Read response payload to get the answer and context
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		response.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to read response"})
		return
	}

	var result struct {
		Answer  string   `json:"answer"`
		Context []string `json:"context"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		response.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Invalid response from AI service"})
		return
	}

	// Persist the request and response in PostgreSQL
	userID, ok := r.Context().Value(domain.UserIDKey).(string)
	if ok {
		var pgUserID pgtype.UUID
		if err := pgUserID.Scan(userID); err == nil {
			contextBytes, _ := json.Marshal(result.Context)
			h.usecase.SaveChatHistory(r.Context(), pgUserID, userReq.Query, result.Answer, contextBytes)
		}
	}

	// Send back to client
	var originalResult map[string]interface{}
	json.Unmarshal(respBody, &originalResult)

	response.RespondJSON(w, resp.StatusCode, originalResult)
}

func (h *AIHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(domain.UserIDKey).(string)
	if !ok {
		response.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	var pgUserID pgtype.UUID
	if err := pgUserID.Scan(userIDStr); err != nil {
		response.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid User ID"})
		return
	}

	// Simplistic pagination wrapper mapping defaults explicitly without manual querystring injection parses
	// E.g get limit offset safely.
	history, err := h.usecase.GetChatHistory(r.Context(), pgUserID, 50, 0)
	if err != nil {
		response.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve history"})
		return
	}

	// Safely return empty arrays
	if len(history) == 0 {
		response.RespondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	response.RespondJSON(w, http.StatusOK, history)
}
