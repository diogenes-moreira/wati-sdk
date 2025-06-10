package media

import (
	"fmt"
	"io"
	"time"
)

// MediaFile representa un archivo de media en WATI
type MediaFile struct {
	ID          string    `json:"id"`
	FileName    string    `json:"fileName"`
	OriginalName string   `json:"originalName,omitempty"`
	MimeType    string    `json:"mimeType"`
	Size        int64     `json:"size"`
	URL         string    `json:"url"`
	ThumbnailURL string   `json:"thumbnailUrl,omitempty"`
	Duration    int       `json:"duration,omitempty"` // Para audio/video en segundos
	Width       int       `json:"width,omitempty"`    // Para imágenes/videos
	Height      int       `json:"height,omitempty"`   // Para imágenes/videos
	Caption     string    `json:"caption,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Status      string    `json:"status"`
	UploadedBy  string    `json:"uploadedBy,omitempty"`
}

// MediaResponse representa la respuesta de obtener media
type MediaResponse struct {
	BaseResponse
	Media MediaFile `json:"media"`
}

// UploadResponse representa la respuesta de subida de media
type UploadResponse struct {
	BaseResponse
	Media    MediaFile `json:"media"`
	UploadID string    `json:"uploadId,omitempty"`
}

// MediaListResponse representa la respuesta de lista de media
type MediaListResponse struct {
	BaseResponse
	PaginatedResponse
	Media []MediaFile `json:"media"`
}

// UploadRequest representa la petición de subida de media
type UploadRequest struct {
	File        io.Reader `json:"-"`
	FileName    string    `json:"fileName"`
	MediaType   string    `json:"mediaType"`
	Caption     string    `json:"caption,omitempty"`
	Description string    `json:"description,omitempty"`
}

// MediaFilter representa filtros para búsqueda de media
type MediaFilter struct {
	MediaType   string    `json:"mediaType,omitempty"`
	FileName    string    `json:"fileName,omitempty"`
	MinSize     int64     `json:"minSize,omitempty"`
	MaxSize     int64     `json:"maxSize,omitempty"`
	CreatedAfter time.Time `json:"createdAfter,omitempty"`
	CreatedBefore time.Time `json:"createdBefore,omitempty"`
	Status      string    `json:"status,omitempty"`
}

// GetMediaParams representa los parámetros para obtener media
type GetMediaParams struct {
	PageSize   int    `json:"pageSize,omitempty"`
	PageNumber int    `json:"pageNumber,omitempty"`
	MediaType  string `json:"mediaType,omitempty"`
	Status     string `json:"status,omitempty"`
}

// MediaStats representa estadísticas de media
type MediaStats struct {
	TotalFiles    int   `json:"totalFiles"`
	TotalSize     int64 `json:"totalSize"`
	ImageCount    int   `json:"imageCount"`
	VideoCount    int   `json:"videoCount"`
	AudioCount    int   `json:"audioCount"`
	DocumentCount int   `json:"documentCount"`
	OtherCount    int   `json:"otherCount"`
}

// MediaStatsResponse representa la respuesta de estadísticas
type MediaStatsResponse struct {
	BaseResponse
	Stats MediaStats `json:"stats"`
}

// BaseResponse representa la respuesta base de la API
type BaseResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// PaginatedResponse representa una respuesta paginada
type PaginatedResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
	TotalCount int `json:"totalCount"`
}

// MediaType representa los tipos de media soportados
type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeVideo    MediaType = "video"
	MediaTypeAudio    MediaType = "audio"
	MediaTypeDocument MediaType = "document"
	MediaTypeSticker  MediaType = "sticker"
)

// MediaStatus representa los estados de un archivo de media
type MediaStatus string

const (
	MediaStatusUploading MediaStatus = "uploading"
	MediaStatusProcessing MediaStatus = "processing"
	MediaStatusReady     MediaStatus = "ready"
	MediaStatusFailed    MediaStatus = "failed"
	MediaStatusDeleted   MediaStatus = "deleted"
)

// SupportedMimeTypes define los tipos MIME soportados por WATI
var SupportedMimeTypes = map[MediaType][]string{
	MediaTypeImage: {
		"image/jpeg",
		"image/png",
		"image/webp",
		"image/gif",
	},
	MediaTypeVideo: {
		"video/mp4",
		"video/3gpp",
		"video/quicktime",
		"video/avi",
		"video/mkv",
	},
	MediaTypeAudio: {
		"audio/aac",
		"audio/mp4",
		"audio/mpeg",
		"audio/amr",
		"audio/ogg",
		"audio/opus",
	},
	MediaTypeDocument: {
		"application/pdf",
		"application/vnd.ms-powerpoint",
		"application/msword",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"text/plain",
		"text/csv",
	},
}

// MaxFileSizes define los tamaños máximos por tipo de media (en bytes)
var MaxFileSizes = map[MediaType]int64{
	MediaTypeImage:    5 * 1024 * 1024,   // 5MB
	MediaTypeVideo:    16 * 1024 * 1024,  // 16MB
	MediaTypeAudio:    16 * 1024 * 1024,  // 16MB
	MediaTypeDocument: 100 * 1024 * 1024, // 100MB
	MediaTypeSticker:  500 * 1024,        // 500KB
}

// Validate valida la petición de subida
func (r *UploadRequest) Validate() error {
	if r.File == nil {
		return fmt.Errorf("file is required")
	}
	
	if r.FileName == "" {
		return fmt.Errorf("fileName is required")
	}
	
	if r.MediaType == "" {
		return fmt.Errorf("mediaType is required")
	}
	
	// Validar que el tipo de media sea soportado
	mediaType := MediaType(r.MediaType)
	if !IsValidMediaType(mediaType) {
		return fmt.Errorf("unsupported media type: %s", r.MediaType)
	}
	
	return nil
}

// IsValidMediaType verifica si un tipo de media es válido
func IsValidMediaType(mediaType MediaType) bool {
	_, exists := SupportedMimeTypes[mediaType]
	return exists
}

// IsSupportedMimeType verifica si un tipo MIME es soportado para un tipo de media
func IsSupportedMimeType(mediaType MediaType, mimeType string) bool {
	supportedTypes, exists := SupportedMimeTypes[mediaType]
	if !exists {
		return false
	}
	
	for _, supportedType := range supportedTypes {
		if supportedType == mimeType {
			return true
		}
	}
	
	return false
}

// GetMaxFileSize retorna el tamaño máximo permitido para un tipo de media
func GetMaxFileSize(mediaType MediaType) int64 {
	if maxSize, exists := MaxFileSizes[mediaType]; exists {
		return maxSize
	}
	return 5 * 1024 * 1024 // 5MB por defecto
}

// ValidateFileSize valida el tamaño de un archivo
func ValidateFileSize(mediaType MediaType, size int64) error {
	maxSize := GetMaxFileSize(mediaType)
	if size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d bytes for media type %s", 
			size, maxSize, mediaType)
	}
	return nil
}

// GetMediaTypeFromMimeType determina el tipo de media basado en el tipo MIME
func GetMediaTypeFromMimeType(mimeType string) MediaType {
	for mediaType, supportedTypes := range SupportedMimeTypes {
		for _, supportedType := range supportedTypes {
			if supportedType == mimeType {
				return mediaType
			}
		}
	}
	return MediaTypeDocument // Por defecto
}

// IsImage verifica si el archivo es una imagen
func (m *MediaFile) IsImage() bool {
	return MediaType(m.MimeType) == MediaTypeImage || 
		   IsSupportedMimeType(MediaTypeImage, m.MimeType)
}

// IsVideo verifica si el archivo es un video
func (m *MediaFile) IsVideo() bool {
	return MediaType(m.MimeType) == MediaTypeVideo || 
		   IsSupportedMimeType(MediaTypeVideo, m.MimeType)
}

// IsAudio verifica si el archivo es audio
func (m *MediaFile) IsAudio() bool {
	return MediaType(m.MimeType) == MediaTypeAudio || 
		   IsSupportedMimeType(MediaTypeAudio, m.MimeType)
}

// IsDocument verifica si el archivo es un documento
func (m *MediaFile) IsDocument() bool {
	return MediaType(m.MimeType) == MediaTypeDocument || 
		   IsSupportedMimeType(MediaTypeDocument, m.MimeType)
}

// IsReady verifica si el archivo está listo para usar
func (m *MediaFile) IsReady() bool {
	return m.Status == string(MediaStatusReady)
}

// IsProcessing verifica si el archivo está siendo procesado
func (m *MediaFile) IsProcessing() bool {
	return m.Status == string(MediaStatusProcessing) || 
		   m.Status == string(MediaStatusUploading)
}

// HasThumbnail verifica si el archivo tiene thumbnail
func (m *MediaFile) HasThumbnail() bool {
	return m.ThumbnailURL != ""
}

// GetFileExtension retorna la extensión del archivo
func (m *MediaFile) GetFileExtension() string {
	if len(m.FileName) == 0 {
		return ""
	}
	
	for i := len(m.FileName) - 1; i >= 0; i-- {
		if m.FileName[i] == '.' {
			return m.FileName[i:]
		}
	}
	
	return ""
}

// FormatFileSize formatea el tamaño del archivo en formato legible
func (m *MediaFile) FormatFileSize() string {
	const unit = 1024
	if m.Size < unit {
		return fmt.Sprintf("%d B", m.Size)
	}
	
	div, exp := int64(unit), 0
	for n := m.Size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	return fmt.Sprintf("%.1f %cB", float64(m.Size)/float64(div), "KMGTPE"[exp])
}

// ToMap convierte GetMediaParams a un mapa para query parameters
func (p *GetMediaParams) ToMap() map[string]string {
	params := make(map[string]string)
	
	if p.PageSize > 0 {
		params["pageSize"] = fmt.Sprintf("%d", p.PageSize)
	}
	
	if p.PageNumber > 0 {
		params["pageNumber"] = fmt.Sprintf("%d", p.PageNumber)
	}
	
	if p.MediaType != "" {
		params["mediaType"] = p.MediaType
	}
	
	if p.Status != "" {
		params["status"] = p.Status
	}
	
	return params
}

// SetDefaults establece valores por defecto para GetMediaParams
func (p *GetMediaParams) SetDefaults() {
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	
	if p.PageNumber <= 0 {
		p.PageNumber = 1
	}
}

