package handlers

import (
	"net/http"

	"github.com/s4lat/gokan/database"
	"github.com/s4lat/gokan/log"
)

// Handlers - contains all http handlers methods.
type Handlers struct {
	DB  database.DB
	Log log.Log
}

// IndexHandler - handles index page.
func (h *Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("<strong>Index page</strong>"))
	if err != nil {
		h.Log.Error(err)
	}
}
