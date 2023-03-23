package resources

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func retrieveTrack(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	audio, successful := t["Audio"]
	if audio == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	if !successful {
		w.WriteHeader(http.StatusBadRequest)
	}

	reqBody, err := json.Marshal(map[string]interface{}{"Audio": audio})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	searchResult, err := http.Post("http://127.0.0.1:3001/search", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer searchResult.Body.Close()
	searchBody := map[string]interface{}{}
	err = json.NewDecoder(searchResult.Body).Decode(&searchBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	targetId, successful := searchBody["Id"]
	if !successful {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	trackResult, err := http.Get("http://127.0.0.1:3000/tracks/" + strings.Replace(targetId.(string), " ", "+", -1))
	defer trackResult.Body.Close()

	trackBody := map[string]interface{}{}
	err = json.NewDecoder(trackResult.Body).Decode(&trackBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tAudio, successful := trackBody["Audio"]
	if !successful {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{"Audio": tAudio}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* Function for Searching Tracks */
	r.HandleFunc("/cooltown", retrieveTrack).Methods("POST")
	return r
}
