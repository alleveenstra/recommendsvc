package recommendsvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func Build_geo_handler(places []Place) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		id, count, geo, rng, parseErr := parse_geo_request(request)
		if parseErr != nil {
			Http_error(parseErr, 400, response, request)
			return
		}
		query, findErr := FindPlace(id, places)
		if findErr != nil {
			Http_error(findErr, 400, response, request)
			return
		}
		scores := calculate_geo_scores(query, geo, places, rng)
		results := Scores_to_result(scores, places, count)
		dat, err := json.MarshalIndent(results, " ", "  ")
		if err == nil {
			Http_json(fmt.Sprintf("%s", dat), response, request)
		} else {
			Http_error(findErr, 500, response, request)
			return
		}
	}
}

func calculate_geo_scores(query *Place, geo []float64, places []Place, rng float64) map[int]float64 {
	length := len(places)
	var output = make(map[int]float64, length)
	for i := 0; i < length; i++ {
		output[i] = 0
		switch {
		case places[i].Id == query.Id:
			output[i] += 100
		case Distance(geo, places[i].Geo) > rng:
			output[i] += 10
		default:
			output[i] += query.Score(&places[i])
		}
	}
	return output
}

func parse_geo_request(request *http.Request) (int, int, []float64, float64, error) {
	id, idErr := strconv.Atoi(request.FormValue("id"))
	count, countErr := strconv.Atoi(request.FormValue("count"))
	lat, latErr := strconv.ParseFloat(request.FormValue("lat"), 64)
	lng, lngErr := strconv.ParseFloat(request.FormValue("lng"), 64)
	rng, rngErr := strconv.ParseFloat(request.FormValue("rng"), 64)
	geo := []float64{lat, lng}
	if idErr == nil && countErr == nil && latErr == nil && lngErr == nil && rngErr == nil {
		return id, count, geo, rng, nil
	}
	return id, count, geo, rng, errors.New("Error parsing request parameters")
}
