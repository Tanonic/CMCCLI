package gv

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readKey() (string){
	dat, err := ioutil.ReadFile(".apikey")
	check(err)
	return string(dat)
}

func ReadConversion() (string){
	dat, err := ioutil.ReadFile(".conversion")
	check(err)
	return string(dat)
}

//TODO implement error checking for HTTP error codes
//give the part from the url after "/v1" as a parameter to this function and it will get the information from the api and return the raw json data
func GetFromApi(directory string) ([]byte) {

	const baseURL = "https://pro-api.coinmarketcap.com/v1"
	url := fmt.Sprintf(baseURL + directory)

	cmcClient := http.Client{
		Timeout: time.Second * 2,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	request.Header.Set("X-CMC_PRO_API_KEY", readKey())

	response, err := cmcClient.Do(request)
	if err != nil {
		if v, ok := err.(*net.OpError); ok {
			if v.Timeout() {
				log.Fatal("Timeout reached")
			}
		}
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		return data
	}
	return nil
}