package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"the-illumiteers-station/backend/internal/sheets"
	"the-illumiteers-station/backend/internal/youtube"
)

type Handler struct {
	youtubeSrv *youtube.Service
	sheetsCli  *sheets.Client
}

func NewHandler(yt *youtube.Service, sh *sheets.Client) *Handler {
	return &Handler{
		youtubeSrv: yt,
		sheetsCli:  sh,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/youtube", h.GetYouTubeFeed)
	r.Get("/api/blog", h.GetBlogFeed)
	r.Post("/api/submit", h.SubmitForm)
}

func (h *Handler) GetYouTubeFeed(w http.ResponseWriter, r *http.Request) {
	videos, err := h.youtubeSrv.GetCachedFeed()
	if err != nil {
		fmt.Printf("GetCachedFeed failed: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

type SubmissionRequest struct {
	Type          string `json:"type"` // "tip", "request", "moment"
	
	// Form A: Submit a Tip
	Name          string `json:"name,omitempty"`
	Contact       string `json:"contact,omitempty"`
	Content       string `json:"content,omitempty"`
	ReferenceLink string `json:"referenceLink,omitempty"`

	// Form B: Content Request
	Requester      string `json:"requester,omitempty"`
	SuggestedTitle string `json:"suggestedTitle,omitempty"`
	RequestDetails string `json:"requestDetails,omitempty"`

	// Form C: Share a Moment
	Submitter        string `json:"submitter,omitempty"`
	EventDescription string `json:"eventDescription,omitempty"`
	MediaURL         string `json:"mediaUrl,omitempty"`
}

func (h *Handler) SubmitForm(w http.ResponseWriter, r *http.Request) {
	var req SubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var err error
	switch req.Type {
	case "tip":
		if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Content) == "" {
			http.Error(w, "name and content details are required", http.StatusBadRequest)
			return
		}
		// Columns: Timestamp, Name, Contact, Content, ReferenceLink
		rowData := []interface{}{timestamp, req.Name, req.Contact, req.Content, req.ReferenceLink}
		err = h.sheetsCli.AppendSubmission("News Tips", rowData)

	case "request":
		if strings.TrimSpace(req.Requester) == "" || strings.TrimSpace(req.SuggestedTitle) == "" || strings.TrimSpace(req.RequestDetails) == "" {
			http.Error(w, "requester name, suggested title, and request details are required", http.StatusBadRequest)
			return
		}
		// Columns: Timestamp, Requester, SuggestedTitle, RequestDetails, ReferenceLink
		rowData := []interface{}{timestamp, req.Requester, req.SuggestedTitle, req.RequestDetails, req.ReferenceLink}
		err = h.sheetsCli.AppendSubmission("Content Requests", rowData)

	case "moment":
		if strings.TrimSpace(req.Submitter) == "" || strings.TrimSpace(req.EventDescription) == "" {
			http.Error(w, "submitter name and milestone details are required", http.StatusBadRequest)
			return
		}
		// Columns: Timestamp, Submitter, EventDescription, MediaURL
		rowData := []interface{}{timestamp, req.Submitter, req.EventDescription, req.MediaURL}
		err = h.sheetsCli.AppendSubmission("Community Moments", rowData)

	default:
		http.Error(w, "invalid submission type", http.StatusBadRequest)
		return
	}

	if err != nil {
		fmt.Printf("Error appending submission to Google Sheets: %v\n", err)
		http.Error(w, fmt.Sprintf("failed to store submission: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handler) GetBlogFeed(w http.ResponseWriter, r *http.Request) {
	if h.sheetsCli == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sheets.GetMockBlogPosts())
		return
	}

	posts, err := h.sheetsCli.GetCachedBlogPosts()
	if err != nil {
		fmt.Printf("GetCachedBlogPosts failed: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
