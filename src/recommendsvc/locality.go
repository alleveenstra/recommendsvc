package recommendsvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func Build_locality_handler(places []Place) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		id, count, locality, parseErr := parse_locality_request(request)
		if parseErr != nil {
			Http_error(parseErr, 400, response, request)
			return
		}
		query, findErr := FindPlace(id, places)
		if findErr != nil {
			Http_error(findErr, 400, response, request)
			return
		}
		scores := calculate_locality_scores(query, locality, places)
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

func calculate_locality_scores(query *Place, locality string, places []Place) map[int]float64 {
	length := len(places)
	var output = make(map[int]float64, length)
	for i := 0; i < length; i++ {
		output[i] = 0
		switch {
		case places[i].Id == query.Id:
			output[i] += 100
		case strings.ToLower(places[i].Locality) != locality:
			output[i] += 10
		default:
			output[i] += query.Score(&places[i])
		}
	}
	return output
}

func parse_locality_request(request *http.Request) (int, int, string, error) {
	id, idErr := strconv.Atoi(request.FormValue("id"))
	count, countErr := strconv.Atoi(request.FormValue("count"))
	locality := request.FormValue("locality")
	if idErr == nil && countErr == nil {
		return id, count, locality, nil
	}
	return id, count, locality, errors.New("Error parsing request parameters")
}
