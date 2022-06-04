package nickandluke

import (
	"encoding/json"
	"net/http"
)

type requestHandler struct {
	guests map[string]string
}

type checkResponse struct {
	Valid bool   `json:"valid"`
	Form  string `json:"form"`
}

func (rh requestHandler) CheckGuest(w http.ResponseWriter, r *http.Request) {
	res := checkResponse{}
	name := r.URL.Query().Get("name")
	if val, ok := rh.guests[name]; ok {
		res.Valid = true
		res.Form = val
	}
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)

}

func RequestHandler() requestHandler {
	guests := make(map[string]string)
	guests["luke"] = "https://tripadvisor.com"
	guests["nick"] = "https://twitter.com"
	return requestHandler{
		guests: guests,
	}
}
