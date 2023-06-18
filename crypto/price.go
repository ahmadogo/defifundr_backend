package crypto

import (
	"io/ioutil"
	"net/http"

	"github.com/demola234/defiraise/interfaces"
)

func GetEthPrice() (string, error) {
	//? Get Price from api
	resp, err := http.Get("https://api.coinbase.com/v2/prices/ETH-USD/spot")
	if err != nil {
		return "", err
	}

	//? Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	price, err := interfaces.UnmarshalCurrentPrice(body)
	if err != nil {
		return "", err
	}

	amount := price.Data.Amount

	return amount, nil

}
