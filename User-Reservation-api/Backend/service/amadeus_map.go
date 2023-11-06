package service

import (
	"User-Reservation/client"
	"User-Reservation/dto"
	"User-Reservation/model"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

type amadeusService struct{}

type amadeusServiceInterface interface {
	InsertAmadeusMap(amadeusMapDto dto.AmadeusMapDto) (dto.AmadeusMapDto, error)
	GetAmadeusIdByHotelId(hotelId string) (dto.AmadeusMapDto, error)
}

var AmadeusService amadeusServiceInterface
var amadeusToken string

func init() {
	AmadeusService = &amadeusService{}
	go getAmadeusToken()
}

func (s *amadeusService) InsertAmadeusMap(amadeusMapDto dto.AmadeusMapDto) (dto.AmadeusMapDto, error) {
	var mapping model.AmadeusMap

	mapping.HotelId = amadeusMapDto.HotelId
	mapping.AmadeusId = amadeusMapDto.AmadeusId

	mapping = client.InsertAmadeusMap(mapping)

	if mapping.HotelId == "" {
		return amadeusMapDto, errors.New("error creating mapping")
	}

	return amadeusMapDto, nil
}

func (s *amadeusService) GetAmadeusIdByHotelId(hotelId string) (dto.AmadeusMapDto, error) {

	var mapping model.AmadeusMap = client.GetAmadeusIdByHotelId(hotelId)

	var amadeusMapDto dto.AmadeusMapDto

	if mapping.HotelId == "" {
		return amadeusMapDto, errors.New("hotel not found")
	}

	amadeusMapDto.HotelId = mapping.HotelId
	amadeusMapDto.AmadeusId = mapping.AmadeusId

	return amadeusMapDto, nil
}

func getAmadeusToken() {

	var tokenResponse dto.AmadeusTokenResponse

	amadeusKey := "50w5hQBPoihvKXUXvAV8LOxRto3rRxCD"
	amadeusSecret := "e5BnlfA2bZYDiaSl"

	url := "https://test.api.amadeus.com/v1/security/oauth2/token"
	data := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", amadeusKey, amadeusSecret)

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		log.Error("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for true {
		httpClient := http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Error("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("Error during request:", resp.StatusCode)
			return
		}

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Error("Error reading response ", err)
			return
		}

		err = json.Unmarshal(body, &tokenResponse)

		if err != nil {
			log.Error("Error parsing JSON ", err)
			return
		}

		amadeusToken = tokenResponse.AccessToken
		log.Info("New token: ", amadeusToken)

		time.Sleep(1790 * time.Second)
	}

}
