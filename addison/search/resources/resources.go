package resources

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var api_key = "903d1e25ae7a962d667ef53c611e6250"

type request struct {
	Audio string
}

type response struct {
	Status string
	Result struct {
		Title string
	}
}

func searchTrack(w http.ResponseWriter, r *http.Request) {
	reqBody := map[string]interface{}{}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	audio, successful := reqBody["Audio"]
	if !successful {
		w.WriteHeader(http.StatusBadRequest)
	}

	apiRequest := map[string]interface{}{"api_key": api_key, "audio": audio}
	marshalled, err := json.Marshal(apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	apiReq, err := http.Post("https://api.audd.io/recognize", "application/json", bytes.NewBuffer(marshalled))
	if err != nil || apiReq.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
	}

	defer apiReq.Body.Close()

	apiResMarshalled, err := io.ReadAll(apiReq.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	apiRes := response{}
	err = json.Unmarshal(apiResMarshalled, &apiRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if apiRes.Status != "success" {
		w.WriteHeader(http.StatusInternalServerError)
	}

	result := map[string]interface{}{"Id": apiRes.Result.Title}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* Function for Searching Tracks */
	r.HandleFunc("/search", searchTrack).Methods("POST")
	return r
}
