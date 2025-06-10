package media

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

// HTTPClient define la interfaz para realizar peticiones HTTP
type HTTPClient interface {
	DoRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error
}

// Service implementa MediaService
type Service struct {
	client HTTPClient
}

// NewService crea una nueva instancia del servicio de media
func NewService(client HTTPClient) *Service {
	return &Service{
		client: client,
	}
}

// GetMediaByFileName obtiene un archivo de media por su nombre
func (s *Service) GetMediaByFileName(ctx context.Context, fileName string) (*MediaResponse, error) {
	if fileName == "" {
		return nil, fmt.Errorf("fileName is required")
	}
	
	endpoint := fmt.Sprintf("/api/v1/getMediaByFileName/%s", fileName)
	
	var response MediaResponse
	err := s.client.DoRequest(ctx, "GET", endpoint, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting media file %s: %w", fileName, err)
	}
	
	return &response, nil
}

// UploadMedia sube un archivo de media a WATI
func (s *Service) UploadMedia(ctx context.Context, file io.Reader, fileName string, mediaType string) (*UploadResponse, error) {
	req := &UploadRequest{
		File:      file,
		FileName:  fileName,
		MediaType: mediaType,
	}
	
	return s.UploadMediaWithRequest(ctx, req)
}

// UploadMediaWithRequest sube un archivo de media usando una petición completa
func (s *Service) UploadMediaWithRequest(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	
	// Crear multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// Agregar el archivo
	part, err := writer.CreateFormFile("file", req.FileName)
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %w", err)
	}
	
	_, err = io.Copy(part, req.File)
	if err != nil {
		return nil, fmt.Errorf("error copying file data: %w", err)
	}
	
	// Agregar campos adicionales
	if req.MediaType != "" {
		writer.WriteField("mediaType", req.MediaType)
	}
	
	if req.Caption != "" {
		writer.WriteField("caption", req.Caption)
	}
	
	if req.Description != "" {
		writer.WriteField("description", req.Description)
	}
	
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing multipart writer: %w", err)
	}
	
	// Realizar petición HTTP personalizada para multipart
	response, err := s.doMultipartRequest(ctx, "POST", "/api/v1/uploadMedia", &buf, writer.FormDataContentType())
	if err != nil {
		return nil, fmt.Errorf("error uploading media: %w", err)
	}
	
	return response, nil
}

// DeleteMedia elimina un archivo de media
func (s *Service) DeleteMedia(ctx context.Context, fileName string) error {
	if fileName == "" {
		return fmt.Errorf("fileName is required")
	}
	
	endpoint := fmt.Sprintf("/api/v1/deleteMedia/%s", fileName)
	
	var response BaseResponse
	err := s.client.DoRequest(ctx, "DELETE", endpoint, nil, &response)
	if err != nil {
		return fmt.Errorf("error deleting media file %s: %w", fileName, err)
	}
	
	return nil
}

// GetMediaURL obtiene la URL de un archivo de media
func (s *Service) GetMediaURL(ctx context.Context, fileName string) (string, error) {
	media, err := s.GetMediaByFileName(ctx, fileName)
	if err != nil {
		return "", err
	}
	
	return media.Media.URL, nil
}

// ListMedia obtiene una lista de archivos de media con parámetros opcionales
func (s *Service) ListMedia(ctx context.Context, params *GetMediaParams) (*MediaListResponse, error) {
	if params == nil {
		params = &GetMediaParams{}
	}
	
	params.SetDefaults()
	
	// Construir endpoint con query parameters
	endpoint := "/api/v1/media"
	queryParams := params.ToMap()
	
	if len(queryParams) > 0 {
		var parts []string
		for key, value := range queryParams {
			parts = append(parts, fmt.Sprintf("%s=%s", key, value))
		}
		endpoint += "?" + strings.Join(parts, "&")
	}
	
	var response MediaListResponse
	err := s.client.DoRequest(ctx, "GET", endpoint, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error listing media: %w", err)
	}
	
	return &response, nil
}

// GetMediaStats obtiene estadísticas de media
func (s *Service) GetMediaStats(ctx context.Context) (*MediaStatsResponse, error) {
	var response MediaStatsResponse
	err := s.client.DoRequest(ctx, "GET", "/api/v1/media/stats", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("error getting media stats: %w", err)
	}
	
	return &response, nil
}

// UploadImage sube una imagen
func (s *Service) UploadImage(ctx context.Context, file io.Reader, fileName string, caption string) (*UploadResponse, error) {
	req := &UploadRequest{
		File:      file,
		FileName:  fileName,
		MediaType: string(MediaTypeImage),
		Caption:   caption,
	}
	
	return s.UploadMediaWithRequest(ctx, req)
}

// UploadVideo sube un video
func (s *Service) UploadVideo(ctx context.Context, file io.Reader, fileName string, caption string) (*UploadResponse, error) {
	req := &UploadRequest{
		File:      file,
		FileName:  fileName,
		MediaType: string(MediaTypeVideo),
		Caption:   caption,
	}
	
	return s.UploadMediaWithRequest(ctx, req)
}

// UploadAudio sube un archivo de audio
func (s *Service) UploadAudio(ctx context.Context, file io.Reader, fileName string) (*UploadResponse, error) {
	req := &UploadRequest{
		File:      file,
		FileName:  fileName,
		MediaType: string(MediaTypeAudio),
	}
	
	return s.UploadMediaWithRequest(ctx, req)
}

// UploadDocument sube un documento
func (s *Service) UploadDocument(ctx context.Context, file io.Reader, fileName string, caption string) (*UploadResponse, error) {
	req := &UploadRequest{
		File:      file,
		FileName:  fileName,
		MediaType: string(MediaTypeDocument),
		Caption:   caption,
	}
	
	return s.UploadMediaWithRequest(ctx, req)
}

// GetMediaByType obtiene archivos de media filtrados por tipo
func (s *Service) GetMediaByType(ctx context.Context, mediaType MediaType, params *GetMediaParams) (*MediaListResponse, error) {
	if params == nil {
		params = &GetMediaParams{}
	}
	
	params.MediaType = string(mediaType)
	return s.ListMedia(ctx, params)
}

// GetImages obtiene solo imágenes
func (s *Service) GetImages(ctx context.Context, params *GetMediaParams) (*MediaListResponse, error) {
	return s.GetMediaByType(ctx, MediaTypeImage, params)
}

// GetVideos obtiene solo videos
func (s *Service) GetVideos(ctx context.Context, params *GetMediaParams) (*MediaListResponse, error) {
	return s.GetMediaByType(ctx, MediaTypeVideo, params)
}

// GetAudios obtiene solo archivos de audio
func (s *Service) GetAudios(ctx context.Context, params *GetMediaParams) (*MediaListResponse, error) {
	return s.GetMediaByType(ctx, MediaTypeAudio, params)
}

// GetDocuments obtiene solo documentos
func (s *Service) GetDocuments(ctx context.Context, params *GetMediaParams) (*MediaListResponse, error) {
	return s.GetMediaByType(ctx, MediaTypeDocument, params)
}

// SearchMedia busca archivos de media por nombre
func (s *Service) SearchMedia(ctx context.Context, query string, params *GetMediaParams) (*MediaListResponse, error) {
	// Esta funcionalidad dependería de si WATI soporta búsqueda por nombre
	// Por ahora, obtenemos todos y filtramos localmente
	response, err := s.ListMedia(ctx, params)
	if err != nil {
		return nil, err
	}
	
	var filteredMedia []MediaFile
	queryLower := strings.ToLower(query)
	
	for _, media := range response.Media {
		if strings.Contains(strings.ToLower(media.FileName), queryLower) ||
		   strings.Contains(strings.ToLower(media.OriginalName), queryLower) {
			filteredMedia = append(filteredMedia, media)
		}
	}
	
	response.Media = filteredMedia
	response.TotalCount = len(filteredMedia)
	
	return response, nil
}

// ValidateUpload valida un archivo antes de subirlo
func (s *Service) ValidateUpload(fileName string, size int64, mimeType string) error {
	if fileName == "" {
		return fmt.Errorf("fileName is required")
	}
	
	// Determinar tipo de media basado en MIME type
	mediaType := GetMediaTypeFromMimeType(mimeType)
	
	// Validar tipo MIME
	if !IsSupportedMimeType(mediaType, mimeType) {
		return fmt.Errorf("unsupported MIME type: %s", mimeType)
	}
	
	// Validar tamaño
	if err := ValidateFileSize(mediaType, size); err != nil {
		return err
	}
	
	return nil
}

// GetMediaInfo obtiene información detallada de un archivo
func (s *Service) GetMediaInfo(ctx context.Context, fileName string) (*MediaFile, error) {
	response, err := s.GetMediaByFileName(ctx, fileName)
	if err != nil {
		return nil, err
	}
	
	return &response.Media, nil
}

// IsMediaReady verifica si un archivo está listo para usar
func (s *Service) IsMediaReady(ctx context.Context, fileName string) (bool, error) {
	media, err := s.GetMediaInfo(ctx, fileName)
	if err != nil {
		return false, err
	}
	
	return media.IsReady(), nil
}

// WaitForMediaReady espera hasta que un archivo esté listo
func (s *Service) WaitForMediaReady(ctx context.Context, fileName string, maxWaitSeconds int) (*MediaFile, error) {
	for i := 0; i < maxWaitSeconds; i++ {
		media, err := s.GetMediaInfo(ctx, fileName)
		if err != nil {
			return nil, err
		}
		
		if media.IsReady() {
			return media, nil
		}
		
		if media.Status == string(MediaStatusFailed) {
			return nil, fmt.Errorf("media processing failed for file: %s", fileName)
		}
		
		// Esperar 1 segundo antes del siguiente intento
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(1 * time.Second):
		}
	}
	
	return nil, fmt.Errorf("timeout waiting for media to be ready: %s", fileName)
}

// doMultipartRequest realiza una petición HTTP multipart personalizada
func (s *Service) doMultipartRequest(ctx context.Context, method, endpoint string, body io.Reader, contentType string) (*UploadResponse, error) {
	// Esta función necesitaría acceso directo al cliente HTTP
	// Por simplicidad, usaremos el método estándar con una estructura especial
	
	// Crear una estructura que represente el multipart data
	multipartData := struct {
		ContentType string `json:"contentType"`
		Body        []byte `json:"body"`
	}{
		ContentType: contentType,
	}
	
	// Leer el cuerpo
	if body != nil {
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("error reading multipart body: %w", err)
		}
		multipartData.Body = bodyBytes
	}
	
	var response UploadResponse
	err := s.client.DoRequest(ctx, method, endpoint, multipartData, &response)
	if err != nil {
		return nil, err
	}
	
	return &response, nil
}

// GetFileExtension extrae la extensión de un nombre de archivo
func GetFileExtension(fileName string) string {
	return filepath.Ext(fileName)
}

// GetMimeTypeFromExtension determina el tipo MIME basado en la extensión
func GetMimeTypeFromExtension(extension string) string {
	extension = strings.ToLower(extension)
	
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".mp4":  "video/mp4",
		".mov":  "video/quicktime",
		".avi":  "video/avi",
		".mkv":  "video/mkv",
		".mp3":  "audio/mpeg",
		".aac":  "audio/aac",
		".ogg":  "audio/ogg",
		".opus": "audio/opus",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".txt":  "text/plain",
		".csv":  "text/csv",
	}
	
	if mimeType, exists := mimeTypes[extension]; exists {
		return mimeType
	}
	
	return "application/octet-stream"
}

