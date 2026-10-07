package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (server *Server) uploadImage(w http.ResponseWriter, r *http.Request) {
	// Maximum upload size 10MB rakhte hain multi-photos ke liye
	r.ParseMultipartForm(10 << 20)

	// "images" naam ka array extract karna
	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		http.Error(w, "No images found in request", http.StatusBadRequest)
		return
	}

	os.MkdirAll("uploads", os.ModePerm)
	var uploadedUrls []string

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// Naya unique naam (time + index)
		fileName := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), i, filepath.Ext(fileHeader.Filename))
		filePath := filepath.Join("uploads", fileName)

		dst, err := os.Create(filePath)
		if err != nil {
			continue
		}

		io.Copy(dst, file)
		dst.Close()

		uploadedUrls = append(uploadedUrls, fmt.Sprintf("http://localhost:8080/uploads/%s", fileName))
	}

	// Array of URLs return karenge
	response := map[string]interface{}{
		"message": "Images uploaded successfully",
		"urls":    uploadedUrls,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
