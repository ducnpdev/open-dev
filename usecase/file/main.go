package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form to retrieve the file
	err := r.ParseMultipartForm(10 << 20) // Max memory 10 MB
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	// Retrieve file from posted form-data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Use the *multipart.FileHeader (handler) to access file metadata
	fmt.Fprintf(w, "Uploaded File: %+v\n", handler.Filename)
	fmt.Fprintf(w, "File Size: %+v\n", handler.Size)
	fmt.Fprintf(w, "MIME Header: %+v\n", handler.Header)

	// Save the file to the server
	dst, err := os.Create("/Users/ducnp/Downloads/open-dev/usecase/file" + handler.Filename)
	if err != nil {
		http.Error(w, "Could not create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file's data to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully.")
}

func main() {
	http.HandleFunc("/upload", uploadFileHandler)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
