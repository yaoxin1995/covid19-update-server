package rki

import (
	"covid19-update-service/model"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const rkiBaseUrl = "https://services7.arcgis.com/mOBPykOjAyBO2ZKk/arcgis/rest/services/RKI_Landkreisdaten/FeatureServer/0/query?where=1%3D1&outFields=OBJECTID,cases7_per_100k&returnGeometry=false&f=json"
const rkiGeoParams = "&geometryType=esriGeometryPoint&inSR=4326&spatialRel=esriSpatialRelWithin"

type Response struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	ObjectID      uint    `json:"OBJECTID"`
	Cases7Per100k float64 `json:"cases7_per_100k"`
}

func requestRKI(url string) (Response, error) {
	rkiResponse := Response{}
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

func GetAllIncidences() (Response, error) {
	return requestRKI(rkiBaseUrl)
}

func GetRkiObjectIdForPosition(position model.GPSPosition) (uint, error) {
	location := fmt.Sprintf("&geometry=%f%%2C%f", position.Longitude, position.Latitude)
	url := fmt.Sprintf("%s%s%s", rkiBaseUrl, rkiGeoParams, location)
	resp, err := requestRKI(url)
	if err != nil {
		return 0, err
	}
	if len(resp.Features) == 0 {
		return 0, errors.New("could not find feature for position")
	}
	return resp.Features[0].Attributes.ObjectID, nil
}
