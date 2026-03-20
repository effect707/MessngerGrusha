package http

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	jwtpkg "github.com/effect707/MessngerGrusha/internal/pkg/jwt"
	fileuc "github.com/effect707/MessngerGrusha/internal/usecase/file"
)

type FileHandler struct {
	fileUC       *fileuc.UseCase
	tokenManager *jwtpkg.TokenManager
	logger       *slog.Logger
}

func NewFileHandler(fileUC *fileuc.UseCase, tokenManager *jwtpkg.TokenManager, logger *slog.Logger) *FileHandler {
	return &FileHandler{
		fileUC:       fileUC,
		tokenManager: tokenManager,
		logger:       logger,
	}
}

const maxUploadSize = 50 << 20 

func (h *FileHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.authenticate(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "file too large", http.StatusRequestEntityTooLarge)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	messageID, err := uuid.Parse(r.FormValue("message_id"))
	if err != nil {
		http.Error(w, "invalid message_id", http.StatusBadRequest)
		return
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	var durationMs *int
	if d := r.FormValue("duration_ms"); d != "" {
		v, err := strconv.Atoi(d)
		if err == nil {
			durationMs = &v
		}
	}

	attachment, err := h.fileUC.Upload(r.Context(), fileuc.UploadInput{
		MessageID:  messageID,
		UserID:     userID,
		FileName:   header.Filename,
		FileSize:   header.Size,
		MimeType:   mimeType,
		DurationMs: durationMs,
		Reader:     file,
	})
	if err != nil {
		h.logger.Error("upload failed", slog.String("error", err.Error()))
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id":"%s","file_name":"%s","file_size":%d,"mime_type":"%s"}`,
		attachment.ID, attachment.FileName, attachment.FileSize, attachment.MimeType)
}

func (h *FileHandler) Download(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := h.authenticate(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	attachmentID, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "invalid attachment id", http.StatusBadRequest)
		return
	}

	reader, attachment, err := h.fileUC.Download(r.Context(), attachmentID, userID)
	if err != nil {
		h.logger.Error("download failed", slog.String("error", err.Error()))
		http.Error(w, "download failed", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	w.Header().Set("Content-Type", attachment.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, attachment.FileName))
	w.Header().Set("Content-Length", strconv.FormatInt(attachment.FileSize, 10))

	io.Copy(w, reader)
}

func (h *FileHandler) authenticate(r *http.Request) (uuid.UUID, error) {
	token := r.Header.Get("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if token == "" {
		token = r.URL.Query().Get("token")
	}

	if token == "" {
		return uuid.UUID{}, fmt.Errorf("missing token")
	}

	claims, err := h.tokenManager.ParseToken(token)
	if err != nil {
		return uuid.UUID{}, err
	}

	return claims.UserID, nil
}
