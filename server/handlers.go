package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// requestData struct is used to decode the JSON request and response body
type requestData struct {
	URL string `json:"url"`
}

// getHandler handles GET requests to the /{param} endpoint
func (s *HTTPServer) getHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the parameter from the URL
	vars := mux.Vars(r)
	paramValue, ok := vars["param"]
	if !ok {
		http.Error(w, "Missing parameter", http.StatusBadRequest)
		return
	}

	// Retrieve the link associated with the parameter
	link, err := s.URLStore.GetLink(paramValue)
	if err != nil {
		return
	}

	// Redirect to the retrieved URL
	http.Redirect(w, r, link, http.StatusFound)
}

// postHandler handles POST requests to the server
func (s *HTTPServer) postHandler(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	var data requestData

	// Decode the JSON request body into the requestData struct
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Generate a short link for the provided URL
	shortLink, err := s.ShortLinkGenerator.GenerateLink(data.URL, s.URLStore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Prepare the response data
	responseData := requestData{URL: shortLink}

	// Encode the response data as JSON
	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	// Write the JSON response to the response writer
	w.Write(responseJSON)
}
