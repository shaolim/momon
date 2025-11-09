package messaging

import (
	"fmt"
	"net/http"
)

func (m *messaging) Callback(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome to Momon HTTP Server!\n")
	fmt.Fprintf(w, "Request path: %s\n", r.URL.Path)
}
