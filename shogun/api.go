package shogun

import (
	"fmt"
	"net/http"

	"github.com/nwillems/shogun-graph/shogun/store"
)

type ShogunAPI struct {
	store store.Storage
}

func NewAPI(store store.Storage) *ShogunAPI {
	return &ShogunAPI{
		store: store,
	}
}

func (api *ShogunAPI) Register(mux *http.ServeMux) error {
	mux.HandleFunc("POST /query", api.Query)
	mux.HandleFunc("GET /query/{node_id}", api.NodeQuery)

	mux.HandleFunc("POST /storage", api.Store)

	return nil
}

func (api *ShogunAPI) Query(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)

	rw.Write([]byte("{Not json}"))
}

func (api *ShogunAPI) NodeQuery(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)

	val := fmt.Sprintf("{Not json: %s}", req.PathValue("node_id"))
	rw.Write([]byte(val))
}

func (api *ShogunAPI) Store(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)

	rw.Write([]byte("{Nothing stored}"))
}
