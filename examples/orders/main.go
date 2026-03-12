package main

import (
	"fmt"
	"os"
	"time"

	sp_api "github.com/fond-of-vertigo/amazon-sp-api"
	"github.com/fond-of-vertigo/amazon-sp-api/constants"
	"github.com/fond-of-vertigo/logger"
)

const PollingDelay = time.Second * 5

func main() {
	log := logger.New(logger.LvlDebug)
	c := sp_api.Config{
		ClientID:     mustGetenv("AMZN_CLIENT_ID"),
		ClientSecret: mustGetenv("AMZN_CLIENT_SECRET"),
		RefreshToken: mustGetenv("AMZN_REFRESH_TOKEN"),
		Endpoint:     constants.Europe,
		Log:          log,
	}

	client, err := sp_api.NewClient(c)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	orderID := "000-0000000-0000000"
	resp, err := client.OrdersAPI.GetOrder(orderID, nil, nil)
	if err != nil {
		log.Errorf("Error while getting order: %w", err)
		return
	}

	log.Infof("Response: %+v", *resp.ResponseBody)
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("missing env var %s", key))
	}
	return v
}
