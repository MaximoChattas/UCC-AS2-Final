package dto

type AmadeusTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type AmadeusAvailabilityResponse struct {
	Data []struct {
		Available bool `json:"available"`
	}
}
