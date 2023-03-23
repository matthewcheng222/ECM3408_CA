package resources

import (
	"encoding/json"
	"net/http"
	"tracks/repository"

	"github.com/gorilla/mux"
)

func createTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var t repository.Track
	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id != t.Id {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	n := repository.Update(t)
	if n > 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		n = repository.Insert(t)
		if n > 0 {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func listTrack(w http.ResponseWriter, r *http.Request) {
	tracks, trackCount := repository.List()

	if trackCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if trackCount == -1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ids := make([]string, 0, len(tracks))

	for _, t := range tracks {
		ids = append(ids, t.Id)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ids)
}

func readTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	track, found := repository.Read(id)

	if found == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else if found > 0 {
		trackData := repository.Track{Id: track.Id, Audio: track.Audio}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(trackData)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func deleteTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result := repository.Delete(id)

	switch {
	case result > 0:
		w.WriteHeader(http.StatusNoContent)
	case result == 0:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Router() http.Handler {
	r := mux.NewRouter()

	/* Function for Creating Tracks */
	r.HandleFunc("/tracks/{id}", createTrack).Methods("PUT")

	/* Function for Listing Tracks */
	r.HandleFunc("/tracks", listTrack).Methods("GET")

	/* Function for Reading Tracks */
	r.HandleFunc("/tracks/{id}", readTrack).Methods("GET")

	/* Function for Deleting Tracks */
	r.HandleFunc("/tracks/{id}", deleteTrack).Methods("DELETE")

	return r
}
