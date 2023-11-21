package service

import (
	"User-Reservation/client"
	"User-Reservation/dto"
	"User-Reservation/model"
	"User-Reservation/utils"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

type amadeusService struct {
	HTTPClient utils.HttpClientInterface
}

type amadeusServiceInterface interface {
	InsertAmadeusMap(amadeusMapDto dto.AmadeusMapDto) (dto.AmadeusMapDto, error)
	GetAmadeusIdByHotelId(hotelId string) (dto.AmadeusMapDto, error)
	GetAmadeusAvailability(amadeusId string, startDate time.Time, endDate time.Time) (bool, error)
}

var AmadeusService amadeusServiceInterface
var amadeusToken string

func init() {
	AmadeusService = &amadeusService{HTTPClient: &utils.HttpClient{}}
	go getAmadeusToken()
}

func (s *amadeusService) InsertAmadeusMap(amadeusMapDto dto.AmadeusMapDto) (dto.AmadeusMapDto, error) {
	var mapping model.AmadeusMap

	mapping.HotelId = amadeusMapDto.HotelId
	mapping.AmadeusId = amadeusMapDto.AmadeusId

	mapping = client.AmadeusClient.InsertAmadeusMap(mapping)

	if mapping.HotelId == "" {
		return amadeusMapDto, errors.New("error creating mapping")
	}

	return amadeusMapDto, nil
}

func (s *amadeusService) GetAmadeusIdByHotelId(hotelId string) (dto.AmadeusMapDto, error) {

	var mapping model.AmadeusMap = client.AmadeusClient.GetAmadeusIdByHotelId(hotelId)

	var amadeusMapDto dto.AmadeusMapDto

	if mapping.HotelId == "" {
		return amadeusMapDto, errors.New("no amadeus id set")
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

	for true {

		req, err := http.NewRequest("POST", url, strings.NewReader(data))
		if err != nil {
			log.Error("Error creating request:", err)
			return
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

func (s *amadeusService) GetAmadeusAvailability(amadeusId string, startDate time.Time, endDate time.Time) (bool, error) {

	formatedStartDate := startDate.Format("2006-01-02")
	formatedEndDate := endDate.Format("2006-01-02")

	url := fmt.Sprintf("https://test.api.amadeus.com/v3/shopping/hotel-offers?hotelIds=%s&checkInDate=%s&checkOutDate=%s", amadeusId, formatedStartDate, formatedEndDate)

	req, err := http.NewRequest("GET", url, strings.NewReader(""))

	if err != nil {
		log.Error("Error creating request:", err)
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+amadeusToken)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		log.Error("Error making request:", err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("Error during request:", resp.StatusCode)
		return false, err
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Error("Error reading response ", err)
		return false, err
	}

	var availabilityResponse dto.AmadeusAvailabilityResponse
	err = json.Unmarshal(body, &availabilityResponse)

	if err != nil {
		log.Error("Error parsing JSON ", err)
		return false, err
	}

	// Amadeus not working correctly
	if len(availabilityResponse.Data) == 0 {
		return true, nil
	}

	return availabilityResponse.Data[0].Available, nil

}
