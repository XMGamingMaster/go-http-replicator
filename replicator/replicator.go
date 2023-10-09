package replicator

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Replicator struct {
	targets []string
}

func (r *Replicator) SetTargets(targets []string) {
	r.targets = targets
}

func (r *Replicator) Replicate(request *http.Request, target string, body *[]byte) {
	httpClient := http.Client{}

	replicateRequest, err := http.NewRequest(request.Method, target, bytes.NewReader(*body))

	replicateRequest.Header = request.Header
	replicateRequest.Header.Set("Host", request.Host)
    replicateRequest.Header.Set("X-Forwarded-For", request.RemoteAddr)

	replicateResponse, err := httpClient.Do(replicateRequest)
	if err != nil {
		fmt.Printf("Failed to replicate request to %s: %v", target, err)
		return
	}
	defer replicateResponse.Body.Close()
}

func (r *Replicator) Handler(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, target := range r.targets {
		go r.Replicate(request, target, &body)
	}

	writer.WriteHeader(204)
}
