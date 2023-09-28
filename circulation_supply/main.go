package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/passage/supply/circulating", GetCircSupply)
	fmt.Println("Listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

type Response struct {
	Amount Amount `json:"amount"`
}

type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

func GetCircSupply(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("vestings.csv")
	if err != nil {
		http.Error(w, "failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Assume header row exists
	if _, err := reader.Read(); err != nil {
		http.Error(w, "error reading record", http.StatusInternalServerError)
		return
	}

	sum := 0
	currentTime := time.Now().Unix() // Current UNIX time

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "error reading record", http.StatusInternalServerError)
			return
		}

		unixTime, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			continue
		}

		if unixTime > currentTime {
			value, err := strconv.Atoi(record[1])
			if err != nil {
				http.Error(w, "Failed to make request:", http.StatusInternalServerError)
				return
			}
			sum += value
		}
	}

	url := "https://api.passage.vitwit.com/cosmos/bank/v1beta1/supply/upasg"
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to make request:", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Received non-200 response code:", http.StatusInternalServerError)
		return
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	originalAmount, err := strconv.ParseInt(response.Amount.Amount, 10, 64)
	if err != nil {
		http.Error(w, "Failed to convert amount to int:", http.StatusInternalServerError)
		return
	}

	amt := originalAmount - int64(sum)
	response.Amount.Amount = strconv.Itoa(int(amt))

	json.NewEncoder(w).Encode(amt)
}
