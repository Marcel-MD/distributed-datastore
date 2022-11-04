package http

import (
	"io"
	"net/http"

	"github.com/Marcel-MD/distributed-datastore/domain"
	"github.com/Marcel-MD/distributed-datastore/presentation/cfg"
	"github.com/Marcel-MD/distributed-datastore/presentation/tcp"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func ListenAndServe() {
	config := cfg.GetConfig()
	port := config.Current.HttpPort

	r := initRouter()

	log.Info().Msg("Listening on port " + port)

	log.Fatal().Err(http.ListenAndServe(":"+port, r)).Msg("Error listening on port " + port)
}

func initRouter() *mux.Router {
	r := mux.NewRouter()

	s := domain.GetStore()

	c := tcp.GetClient()

	r.HandleFunc("/{key}", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key := vars["key"]

		value, err := s.Get(key)
		if err != nil {
			log.Error().Err(err).Msg("Error getting key")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Key not found"))
			return
		}

		w.Write(value)

	}).Methods("GET")

	r.HandleFunc("/{key}", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key := vars["key"]

		value, err := io.ReadAll(r.Body)

		if err != nil {
			log.Error().Err(err).Msg("Error reading body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error reading body"))
			return
		}

		err = s.Set(key, value)
		if err != nil {
			log.Error().Err(err).Msg("Error setting key")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error setting key"))
			return
		}

		c.Set(key, value)

		w.WriteHeader(http.StatusCreated)

	}).Methods("POST")

	r.HandleFunc("/{key}", func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key := vars["key"]

		err := s.Delete(key)
		if err != nil {
			log.Error().Err(err).Msg("Error deleting key")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		c.Delete(key)

		w.WriteHeader(http.StatusNoContent)

	}).Methods("DELETE")

	return r
}
