package targets

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/vitwit/cosmos-utils/relayer-alerter/config"
	"github.com/vitwit/cosmos-utils/relayer-alerter/db"
)

// GetEndpointsStatus to get alert about endpoints status
func GetEndpointsStatus(cfg *config.Config) error {
	var ops HTTPOptions

	addresses, err := db.GetAllAddress(bson.M{}, bson.M{}, cfg.MongoDB.Database)
	if err != nil {
		log.Printf("Error while getting addresses list from db : %v", err)
		return err
	}
	var msg string

	for _, value := range addresses {
		if value.NetworkName != "persistence" && value.NetworkName != "iris" && value.NetworkName != "Crypto.com" { // for now ignore endpoint alerts for persistence
			ops = HTTPOptions{
				Endpoint: value.RPC + "/status",
				Method:   http.MethodGet,
			}

			fmt.Println(value.NetworkName)

			_, err := HitHTTPTarget(ops)
			if err != nil {
				log.Printf("Error in rpc: %v", err)
				msg = msg + fmt.Sprintf("⛔ Unreachable to %s RPC :: %s and the ERROR is : %v\n", value.NetworkName, ops.Endpoint, err.Error())
			}

			ops = HTTPOptions{
				Endpoint: value.LCD + "/node_info",
				Method:   http.MethodGet,
			}

			_, err = HitHTTPTarget(ops)
			if err != nil {
				log.Printf("Error in lcd endpoint: %v", err)
				msg = msg + fmt.Sprintf("⛔ Unreachable to %s LCD :: %s and the ERROR is : %v\n\n", value.NetworkName, ops.Endpoint, err.Error())
			}
		}
	}

	if msg != "" {
		err = SendTelegramAlert(msg, cfg)
		if err != nil {
			log.Printf("Error while sending telegram alert : %v", err)
		}
	}

	return nil
}
