package run

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"
)

func Submit(w http.ResponseWriter, r *http.Request) {
    // Parse the multipart form data
    err := r.ParseMultipartForm(1000000) // 100 KB maximum file size
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // Get submission id
	submissionId := r.FormValue("submissionId")
    // Get the file from the form data
    file, _, err := r.FormFile("file")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    // Read the contents of the file into a buffer
    buffer := bytes.NewBuffer(nil)
    if _, err := io.Copy(buffer, file); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // write file to disk
	os.Mkdir(submissionId, 0777);
	err = os.WriteFile(submissionId+"/"+submissionId+".cpp", buffer.Bytes(), 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
    
    // Return a response to the client
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("File uploaded successfully"))

	// run the code
	start := time.Now()
	RunCode(submissionId, "1", "cpp/input")
	end := time.Now()
	println("Time taken: ", end.Sub(start).Milliseconds() , "ms")
}

func test () {
	println("test")
}