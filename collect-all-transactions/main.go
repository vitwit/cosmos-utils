package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RequestPrefix = "https://api-osmosis.cosmostation.io/v1/account/new_txs/"

//"https://api.cosmostation.io/v1/account/new_txs/"

type TxHeader struct {
	ID        int    `json:"id"`
	ChainID   string `json:"chain_id"`
	BlockID   int    `json:"block_id"`
	Timestamp string `json:"timestamp"`
}

type Body struct {
	Messages []map[string]interface{} `json:"messages"`
}

type Txn struct {
	Body Body `json:"body"`
}

type TxData struct {
	Height    string              `json:"height,omitempty"`
	TxHash    string              `json:"txhash,omitempty"`
	Codespace string              `json:"codespace,omitempty"`
	Code      int                 `json:"code,omitempty"`
	Data      string              `json:"data,omitempty"`
	RawLog    string              `json:"raw_log,omitempty"`
	Logs      sdk.ABCIMessageLogs `json:"logs"`
	Info      string              `json:"info,omitempty"`
	GasWanted string              `json:"gas_wanted,omitempty"`
	GasUsed   string              `json:"gas_used,omitempty"`
	Tx        Txn                 `json:"tx,omitempty"`
	Timestamp string              `json:"timestamp,omitempty"`
	// Events    []tmTypes.Event `json:"events"`
}

type Tx struct {
	Header TxHeader `json:"header"`
	Data   TxData   `json:"data"`
}

type AmountObj struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Event struct {
	Type       string      `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

type Log struct {
	Events []Event `json:"events"`
}

func collectAllTxns(address string) error {
	fromId := "0"

	dataFile, err := os.Create("data.csv")
	if err != nil {
		return fmt.Errorf("Error in creating the CSV file: %s", err.Error())
	}

	defer dataFile.Close()

	writer := csv.NewWriter(dataFile)
	defer writer.Flush()

	header := []string{
		"ID",
		"ChainID",
		"BlockID",
		"Height",
		"TxHash",
		"Codespace",
		"Code",
		"Data",
		"RawLog",
		"Logs",
		"Info",
		"GasWanted",
		"GasUsed",
		"Tx",
		"Timestamp",
		"Type",
		"FromAddress",
		"ToAddress",
		"Amount",
		"Withdraw commission",
		"Auto claim rewards",
	}
	if err = writer.Write(header); err != nil {
		return fmt.Errorf("Error in writing header info: %s", err.Error())
	}

	for {
		url := RequestPrefix + address + "?limit=20&from=" + fromId

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		// ...
		req.Header.Add("Origin", `https://www.mintscan.io`)
		req.Header.Add("Referer", `https://www.mintscan.io/`)
		resp, err := client.Do(req)
		//resp, err := http.Get(req)
		if err != nil {
			return fmt.Errorf("Error in fetching the txn data: %s", err.Error())
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error in reading response body: %s", err.Error())
		}

		var txs []Tx
		if err = json.Unmarshal(body, &txs); err != nil {
			return fmt.Errorf("Error in unmarshaling the JSON tx data: %s", err.Error())
		}

		if len(txs) == 0 {
			break
		}

		for _, tx := range txs {
			logs, err := json.Marshal(tx.Data.Logs)
			if err != nil {
				return fmt.Errorf("Error in marshaling logs: %s", err.Error())
			}

			txn, err := json.Marshal(tx.Data.Tx)
			if err != nil {
				return fmt.Errorf("Error in marshaling Tx: %s", err.Error())
			}

			data := []string{
				strconv.Itoa(tx.Header.ID),
				tx.Header.ChainID,
				strconv.Itoa(tx.Header.BlockID),
				tx.Data.Height,
				tx.Data.TxHash,
				tx.Data.Codespace,
				strconv.Itoa(tx.Data.Code),
				tx.Data.Data,
				tx.Data.RawLog,
				string(logs),
				tx.Data.Info,
				tx.Data.GasWanted,
				tx.Data.GasUsed,
				string(txn),
				tx.Data.Timestamp,
			}

			messages := tx.Data.Tx.Body.Messages

			for i, msg := range messages {
				if msg["@type"] == "/cosmos.bank.v1beta1.MsgSend" {
					amount := msg["amount"]
					amountBytes, err := json.Marshal(amount)
					if err != nil {
						return fmt.Errorf("Error in marshaling Tx: %s", err.Error())
					}

					var amounts []AmountObj

					if err = json.Unmarshal(amountBytes, &amounts); err != nil {
						return fmt.Errorf("Error in unmarshaling the JSON msg send data: %s", err.Error())
					}

					if len(messages) > 1 && i != 0 {
						data1 := []string{"", "", "", "", "", "",
							"", "", "", "", "", "", "", "", "",
							fmt.Sprintf("%v", msg["@type"]),
							fmt.Sprintf("%v", msg["from_address"]),
							fmt.Sprintf("%v", msg["to_address"]),
							amounts[0].Amount,
						}

						if err = writer.Write(data1); err != nil {
							return fmt.Errorf("Error in writing data: %s", err.Error())
						}
					} else {
						data = append(data, []string{
							fmt.Sprintf("%v", msg["@type"]),
							fmt.Sprintf("%v", msg["from_address"]),
							fmt.Sprintf("%v", msg["to_address"]),
							amounts[0].Amount,
						}...)
					}
				}

				if msg["@type"] == "/cosmos.staking.v1beta1.MsgDelegate" ||
					msg["@type"] == "/cosmos.staking.v1beta1.MsgUndelegate" ||
					msg["@type"] == "/cosmos.staking.v1beta1.MsgBeginRedelegate" {
					var arrLog []Log
					amount := msg["amount"]
					amountBytes, err := json.Marshal(amount)
					if err != nil {
						return fmt.Errorf("Error in marshaling Tx: %s", err.Error())
					}

					var amounts AmountObj

					if err = json.Unmarshal(amountBytes, &amounts); err != nil {
						return fmt.Errorf("Error in unmarshaling the JSON delete data: %s", err.Error())
					}

					if len(messages) > 1 && i != 0 {
						var rewards string

						if err = json.Unmarshal(logs, &arrLog); err != nil {
							return fmt.Errorf("Error in unmarshaling the JSON delegate logs data: %s", err.Error())
						}

						log := arrLog[i]
						var found bool
						var index int

						for lI, e := range log.Events {
							if e.Type == "transfer" {
								for _, a := range e.Attributes {
									if a.Key == "recipient" && a.Value == address {
										index = lI
										found = true
									}
								}
							}
						}

						if found {
							attrs := log.Events[index]
							for _, atV := range attrs.Attributes {
								if atV.Key == "amount" {
									rewards = atV.Value
								}
							}
						}

						data1 := []string{"", "", "", "", "", "",
							"", "", "", "", "", "", "", "", "",
							fmt.Sprintf("%v", msg["@type"]),
							fmt.Sprintf("%v", msg["delegator_address"]),
							fmt.Sprintf("%v", msg["validator_address"]),
							amounts.Amount,
							"",
							rewards,
						}

						if err = writer.Write(data1); err != nil {
							return fmt.Errorf("Error in writing data: %s", err.Error())
						}
					} else {
						var rewards string

						if err = json.Unmarshal(logs, &arrLog); err != nil {
							return fmt.Errorf("Error in unmarshaling the JSON delegate logs single data: %s", err.Error())
						}

						var found bool
						var index int
						if len(arrLog) > 0 {
							log := arrLog[i]
							for lI, e := range log.Events {
								if e.Type == "transfer" {
									for _, a := range e.Attributes {
										if a.Key == "recipient" && a.Value == address {
											index = lI
											found = true
										}
									}
								}
							}

							if found {
								attrs := log.Events[index]
								for _, atV := range attrs.Attributes {
									if atV.Key == "amount" {
										rewards = atV.Value
									}
								}
							}
						}

						data = append(data, []string{
							fmt.Sprintf("%v", msg["@type"]),
							fmt.Sprintf("%v", msg["delegator_address"]),
							fmt.Sprintf("%v", msg["validator_address"]),
							amounts.Amount,
							"",
							rewards,
						}...)
					}
				}

				if msg["@type"] == "/ibc.applications.transfer.v1.MsgTransfer" {
					amount := msg["token"]
					amountBytes, err := json.Marshal(amount)
					if err != nil {
						return fmt.Errorf("Error in marshaling Tx: %s", err.Error())
					}

					var amounts AmountObj

					if err = json.Unmarshal(amountBytes, &amounts); err != nil {
						return fmt.Errorf("Error in unmarshaling the JSON transfer data: %s", err.Error())
					}

					if len(messages) > 1 && i != 0 {
						data1 := []string{"", "", "", "", "", "",
							"", "", "", "", "", "", "", "", "",
							fmt.Sprintf("%v", msg["@type"]),
							fmt.Sprintf("%v", msg["sender"]),
							fmt.Sprintf("%v", msg["receiver"]),
							amounts.Amount,
						}

						if err = writer.Write(data1); err != nil {
							return fmt.Errorf("Error in writing data: %s", err.Error())
						}
					} else {
						data = append(data, []string{
							fmt.Sprintf("%v", msg["@type"]),
							fmt.Sprintf("%v", msg["sender"]),
							fmt.Sprintf("%v", msg["receiver"]),
							amounts.Amount,
						}...)
					}
				}

				if msg["@type"] == "/cosmos.gov.v1beta1.MsgVote" {
					if len(messages) > 1 && i != 0 {
						data1 := []string{"", "", "", "", "", "",
							"", "", "", "", "", "", "", "", "",
							fmt.Sprintf("%v", msg["@type"]),
							fmt.Sprintf("%v", msg["voter"]),
						}

						if err = writer.Write(data1); err != nil {
							return fmt.Errorf("Error in writing data: %s", err.Error())
						}
					} else {
						data = append(data, []string{
							fmt.Sprintf("%v", msg["@type"]),
							fmt.Sprintf("%v", msg["voter"]),
						}...)
					}
				}

				if msg["@type"] == "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward" {
					if len(messages) > 1 && i != 0 {
						amount, err := GetWithdrawAmount(logs, i, "withdraw_rewards", "amount")

						if err != nil {
							return fmt.Errorf("Error in unmarshaling the JSON data: %s", err.Error())
						}

						data1 := []string{"", "", "", "", "", "",
							"", "", "", "", "", "", "", "", "",
							fmt.Sprintf("%v", msg["@type"]),
							fmt.Sprintf("%v", msg["delegator_address"]),
							fmt.Sprintf("%v", msg["validator_address"]),
							amount,
						}

						if err = writer.Write(data1); err != nil {
							return fmt.Errorf("Error in writing data: %s", err.Error())
						}
					} else {
						amount, err := GetWithdrawAmount(logs, i, "withdraw_rewards", "amount")

						if err != nil {
							return fmt.Errorf("Error in unmarshaling the JSON data: %s", err.Error())
						}

						data = append(data, []string{
							fmt.Sprintf("%v", msg["@type"]),
							fmt.Sprintf("%v", msg["delegator_address"]),
							fmt.Sprintf("%v", msg["validator_address"]),
							amount,
						}...)
					}
				}

				if msg["@type"] == "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission" {

					if len(messages) > 1 && i != 0 {
						amount, err := GetWithdrawAmount(logs, i, "withdraw_commission", "amount")

						if err != nil {
							return fmt.Errorf("Error in unmarshaling the JSON data: %s", err.Error())
						}

						data1 := []string{"", "", "", "", "", "",
							"", "", "", "", "", "", "", "", "",
							fmt.Sprintf("%v", msg["@type"]),
							fmt.Sprintf("%v", msg["validator_address"]),
							"",
							"",
							amount,
						}

						if err = writer.Write(data1); err != nil {
							return fmt.Errorf("Error in writing data: %s", err.Error())
						}
					} else {
						amount, err := GetWithdrawAmount(logs, i, "withdraw_commission", "amount")

						if err != nil {
							return fmt.Errorf("Error in unmarshaling the JSON data: %s", err.Error())
						}

						data = append(data, []string{
							fmt.Sprintf("%v", msg["@type"]),
							fmt.Sprintf("%v", msg["validator_address"]),
							"",
							"",
							amount,
						}...)
					}
				}

			}

			if err = writer.Write(data); err != nil {
				return fmt.Errorf("Error in writing data: %s", err.Error())
			}
		}

		fromId = strconv.Itoa(txs[len(txs)-1].Header.ID)
	}

	return nil
}

func GetWithdrawAmount(logs []byte, i int, input1 string,
	input2 string) (string, error) {
	var arrLog []Log

	var amount string
	if err := json.Unmarshal(logs, &arrLog); err != nil {
		return "", fmt.Errorf("Error in unmarshaling the JSON withdraw amount data: %s", err.Error())
	}

	if len(arrLog) > 0 {
		log := arrLog[i]

		for _, e := range log.Events {
			if e.Type == input1 {
				for _, a := range e.Attributes {
					if a.Key == input2 {
						amount = a.Value
					}
				}
			}
		}

	}

	return amount, nil
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Failed to collect txs: address not provided")
	}
	address := args[1]
	log.Infof("Fetching the details of the txs for the address: %s", address)
	err := collectAllTxns(address)
	if err != nil {
		log.Errorf("collecting all transactions failed: %s", err.Error())
	} else {
		log.Info("fetching txs complete!")
	}
}
