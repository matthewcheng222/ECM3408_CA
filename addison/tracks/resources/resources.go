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
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if id == t.Id {
			if n := repository.Insert(t); n > 0 {
				w.WriteHeader(204)
			} else if n := repository.Update(t); n > 0 {
				w.WriteHeader(201)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func listTrack(w http.ResponseWriter, r *http.Request) {
	tracks, trackCount := repository.List()

	if trackCount == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	if trackCount == -1 {
		w.WriteHeader(http.StatusInternalServerError)
	}

	Ids := []string{}

	for _, t := range tracks {
		Ids = append(Ids, t.Id)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Ids)
}

func readTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if c, n := repository.Read(id); n > 0 {
		d := repository.Track{Id: c.Id, Audio: c.Audio}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(d)
	} else if n == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func deleteTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, successful := vars["id"]
	if !successful {
		w.WriteHeader(http.StatusBadRequest)
	}

	result := repository.Delete(id)

	if result > 0 {
		w.WriteHeader(http.StatusNoContent)
	} else if result == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
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
