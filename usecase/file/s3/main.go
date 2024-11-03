package main

import (
	"fmt"
	"net/http"

	opendevS3 "open-dev/aws/s3"
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

	openS3 := opendevS3.AwsS3{}
	_, _, err = openS3.UploadFile(handler)
	if err == nil {
		fmt.Fprintf(w, "File uploaded successfully.")
		return
	}
	fmt.Fprintf(w, "File uploaded err %s", err)
}

func main() {
	http.HandleFunc("/upload", uploadFileHandler)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
