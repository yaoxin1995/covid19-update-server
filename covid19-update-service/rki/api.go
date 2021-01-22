package rki

import (
	"covid19-update-service/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const rkiBaseUrl = "https://services7.arcgis.com/mOBPykOjAyBO2ZKk/arcgis/rest/services/RKI_Landkreisdaten/FeatureServer/0/query?where=1%3D1&outFields=OBJECTID,cases7_per_100k,last_update&returnGeometry=false&f=json"
const rkiGeoParams = "&geometryType=esriGeometryPoint&inSR=4326&spatialRel=esriSpatialRelWithin"

type response struct {
	Features []feature `json:"features"`
}

type feature struct {
	Attributes attributes `json:"attributes"`
}

type attributes struct {
	ObjectID      uint    `json:"OBJECTID"`
	Cases7Per100k float64 `json:"cases7_per_100k"`
	RawLastUpdate string  `json:"last_update"`
}

func (a *attributes) lastUpdate() (time.Time, error) {
	// RKI API returns timestamp in dd:mm:yyyy HH:mm, but Go does not support 24h time format with leading zeros
	rkiHourRegExp := regexp.MustCompile(`\d{2}:`)

	nonZeroHourTime := rkiHourRegExp.ReplaceAllStringFunc(a.RawLastUpdate, func(entry string) string {
		if strings.Index(entry, "0") == 0 {
			return strings.Replace(entry, "0", "", 1)
		} else {
			return entry
		}
	})
	timeString := strings.Replace(nonZeroHourTime, " Uhr", "", -1)
	t, err := time.Parse("2.1.2006, 15:04", timeString)
	if err != nil {
		return t, fmt.Errorf("could not parse timestamp: %v", err)
	}
	return t, nil
}

type LocationNotFoundError struct{}

func (e *LocationNotFoundError) Error() string {
	return "The provided location is not supported."
}

func requestRKI(url string) (response, error) {
	rkiResponse := response{}
	rawResponse, err := http.Get(url)
	if err != nil {
		return rkiResponse, fmt.Errorf("could not load data fron RKI: %v", err)
	}

	defer rawResponse.Body.Close()
	body, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return rkiResponse, fmt.Errorf("could not read response from RKI: %v", err)
	}
	err = json.Unmarshal(body, &rkiResponse)
	if err != nil {
		return rkiResponse, fmt.Errorf("could not unmarshal JSON response from RKI: %v", err)
	}
	return rkiResponse, nil
}

func getAllRegions() (response, error) {
	return requestRKI(rkiBaseUrl)
}

func GetRegionIDForPosition(position model.GPSPosition) (uint, error) {
	location := fmt.Sprintf("&geometry=%f%%2C%f", position.Longitude, position.Latitude)
	url := fmt.Sprintf("%s%s%s", rkiBaseUrl, rkiGeoParams, location)
	resp, err := requestRKI(url)
	if err != nil {
		return 0, err
	}
	if len(resp.Features) == 0 {
		return 0, &LocationNotFoundError{}
	}
	return resp.Features[0].Attributes.ObjectID, nil
}
