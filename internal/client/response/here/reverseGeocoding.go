package here

type ReverseGeocode struct {
	// MetaInfo struct {
	// 	TimeStamp string `json:"TimeStamp"`
	// } `json:"MetaInfo"`
	Items []struct {
		Title           string `json:"heading"`
		ID              string `json:"id"`
		ResultType      string `json:"resultType"`
		HouseNumberType string `json:"houseNumberType"`
		Address         struct {
			Label       string `json:"label"`
			CountryCode string `json:"countryCode"`
			CountryName string `json:"countryName"`
			State       string `json:"state"`
			County      string `json:"county"`
			City        string `json:"city"`
			District    string `json:"district"`
			Street      string `json:"street"`
			HouseNumber string `json:"houseNumber"`
			PostalCode  string `json:"postalCode"`
		} `json:"Address"`
		Position Position   `json:"position"`
		Access   []Position `json:"access"`
		Distance int32      `json:"distance"`
	} `json:"items"`
}

type Position struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type ReverseGeocoder struct {
	Response struct {
		MetaInfo struct {
			TimeStamp string `json:"TimeStamp"`
		} `json:"MetaInfo"`
		View []struct {
			Result []struct {
				MatchLevel string `json:"MatchLevel"`
				Location   struct {
					Address struct {
						Label       string `json:"Label"`
						Country     string `json:"Country"`
						State       string `json:"State"`
						County      string `json:"County"`
						City        string `json:"City"`
						District    string `json:"District"`
						Street      string `json:"Street"`
						HouseNumber string `json:"HouseNumber"`
						PostalCode  string `json:"PostalCode"`
					} `json:"Address"`
				} `json:"Location"`
			} `json:"Result"`
		} `json:"View"`
	} `json:"Response"`
}

type Geocoder struct {
	AppID   string `json:"app_id"`
	AppCode string `json:"app_code"`
}
