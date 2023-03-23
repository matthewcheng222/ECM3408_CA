package resources

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var api_key = "f153bf06d224509250503c04716b331d"

type Request struct {
	Audio string
}

type APIResponse struct {
	Status string
	Result struct {
		Title string
	}
}

func searchTrack(w http.ResponseWriter, r *http.Request) {
	var reqBody Request

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	apiRequest := map[string]interface{}{"api_token": api_key, "audio": reqBody.Audio}
	marshalled, err := json.Marshal(apiRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	apiReq, err := http.Post("https://api.audd.io/recognize", "application/json", bytes.NewBuffer(marshalled))
	if err != nil || apiReq.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer apiReq.Body.Close()

	apiResMarshalled, err := io.ReadAll(apiReq.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var apiRes APIResponse
	err = json.Unmarshal(apiResMarshalled, &apiRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if apiRes.Status != "success" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{"Id": apiRes.Result.Title}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* Function for Searching Tracks */
	r.HandleFunc("/search", searchTrack).Methods("POST")
	return r
}
