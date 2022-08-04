package nickandluke

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const oneGuestUrl = "https://docs.google.com/forms/d/e/1FAIpQLSdXF80AevtDqkC7ZTynrzXRuwfZCjQPTpsLhCEfuRPSOCCgww/viewform?usp=sf_link"
const twoGuestsUrl = "https://docs.google.com/forms/d/e/1FAIpQLSevxS_HMScw6Nhcru3ke8GeqWfJnBAA_AdWPc-1eRmgS4G6LQ/viewform?usp=sf_link"
const guestFile = "staging/guests.csv"

type guests map[string]string

type requestHandler struct {
	guests guests
}

func (rh requestHandler) String() string {
	var sb strings.Builder
	for guest, url := range rh.guests {
		sb.WriteString(fmt.Sprintf("%-20s%s\n", guest, url))
	}

	return sb.String()
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
func cleanupGuest(guest string) string {
	return strings.ToLower(strings.TrimSpace(guest))
}

func parseGuests(rows [][]string) guests {
	guests := make(map[string]string)
	numOneGuest := 0
	numTwoGuests := 0
	for i := 0; i < len(rows); i++ {
		row := rows[i]
		if len(row) != 2 {
			panic(fmt.Sprintf("Row %s does not have two records", row))
		}
		guest1 := cleanupGuest(row[0])
		guest2 := cleanupGuest(row[1])

		if guest1 == "" {
			panic(fmt.Sprintf("Row %s has empty first guest", row))
		}
		var url string
		// Has one guests
		if guest2 == "" {
			url = oneGuestUrl
			numOneGuest += 1
		} else {
			url = twoGuestsUrl
			numTwoGuests += 1
			if _, ok := guests[guest2]; ok {
				panic(fmt.Sprintf("Found duplicate guest %s", guest2))
			}
			guests[guest2] = url
		}

		if _, ok := guests[guest1]; ok {
			panic(fmt.Sprintf("Found duplicate guest %s", guest1))
		}

		guests[guest1] = url

	}
	if numOneGuest == 0 {
		panic("Found no one-guests!")
	}
	if numTwoGuests == 0 {
		panic("Found no two-guests!")
	}

	return guests
}

func loadGuests() guests {
	f, err := os.Open(guestFile)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	csvReader.Comma = ','
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return parseGuests(data)
}

func RequestHandler() requestHandler {
	guests := loadGuests()
	//guests["luke massa"] = "https://tripadvisor.com"
	//guests["nick andersen"] = "https://twitter.com"
	return requestHandler{
		guests: guests,
	}
}
