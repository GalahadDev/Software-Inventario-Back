package database

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

// SubirAStorageSupabase construye una petición  archivo al bucket.
func SubirAStorageSupabase(fileHeader *multipart.FileHeader, bucketName, filePath string) (string, error) {
	supabaseProject := os.Getenv("SUPABASE_PROJECT")
	supabaseServiceKey := os.Getenv("SUPABASE_SERVICE_KEY")

	// 1. Abrir el archivo
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el archivo del form: %v", err)
	}
	defer file.Close()

	// 2. Detectar o adivinar el MIME según la extensión (simple ejemplo)
	extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	var mimeType string
	switch extension {
	case ".png":
		mimeType = "image/png"
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".gif":
		mimeType = "image/gif"
	default:
		mimeType = "application/octet-stream"
	}

	// 3. Crear el multipart en memoria
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", fileHeader.Filename))
	h.Set("Content-Type", mimeType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return "", fmt.Errorf("error al crear part: %v", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("error al copiar contenido del archivo: %v", err)
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	// 4. Construir la URL
	supabaseURL := fmt.Sprintf("https://%s.supabase.co", supabaseProject)
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucketName, filePath)

	req, err := http.NewRequest(http.MethodPost, uploadURL, &body)
	if err != nil {
		return "", fmt.Errorf("error creando request a supabase: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+supabaseServiceKey)

	// 5. Petición
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error al enviar request a supabase: %v", err)
	}
	defer resp.Body.Close()

	// 6. Validar
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("error al subir archivo, status %d, body: %s", resp.StatusCode, readBody(resp))
	}

	// 7. Retornar la URL
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucketName, filePath)
	return publicURL, nil
}

func readBody(resp *http.Response) string {
	var sb strings.Builder
	io.Copy(&sb, resp.Body)
	return sb.String()
}
