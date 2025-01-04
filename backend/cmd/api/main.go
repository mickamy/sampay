package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	fmt.Println("API server is running...")
	http.HandleFunc("/echo", echoHandler)

	fmt.Println("listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	_, err = fmt.Fprintf(w, "Echo: %s", string(body))
	if err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}
