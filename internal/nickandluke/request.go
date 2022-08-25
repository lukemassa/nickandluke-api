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

type guestConfiguration struct {
	guests map[string]string
	email  string
}

type requestHandler struct {
	guestConfiguration guestConfiguration
}

func (rh requestHandler) String() string {
	var sb strings.Builder
	for guest, url := range rh.guestConfiguration.guests {
		sb.WriteString(fmt.Sprintf("%-20s%s\n", guest, url))
	}
	sb.WriteString(fmt.Sprintf("Email: %s", rh.guestConfiguration.email))

	return sb.String()
}

type checkResponse struct {
	Valid bool   `json:"valid"`
	Form  string `json:"form"`
	Email string `json:"email"`
}

func (rh requestHandler) CheckGuest(w http.ResponseWriter, r *http.Request) {
	res := checkResponse{}
	name := r.URL.Query().Get("name")
	if val, ok := rh.guestConfiguration.guests[name]; ok {
		res.Valid = true
		res.Form = val
		res.Email = rh.guestConfiguration.email
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

func parseGuests(rows [][]string) guestConfiguration {
	ret := guestConfiguration{}
	guests := make(map[string]string)
	numOneGuest := 0
	numTwoGuests := 0
	email := rows[0][0]
	if !strings.Contains(email, "@") {
		panic("Row 1 does not look like an email")
	}
	ret.email = email

	for i := 1; i < len(rows); i++ {
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
	ret.guests = guests

	return ret
}

func loadGuestConfiguration() guestConfiguration {
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
	guestConfiguration := loadGuestConfiguration()

	return requestHandler{
		guestConfiguration: guestConfiguration,
	}
}
