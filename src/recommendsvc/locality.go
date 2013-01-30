package recommendsvc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func Build_locality_handler(places []Place) func(response http.ResponseWriter, request *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		id, count, locality, parseErr := parse_locality_request(request)
		if parseErr != nil {
			log.Println(fmt.Sprintf("%v %s %s", parseErr, request.URL.Path, request.URL.RawQuery))
			http.Error(response, "400 Malformed request.", 400)
			return
		}
		query, findErr := FindPlace(id, places)
		if findErr != nil {
			log.Println(fmt.Sprintf("%v %s %s", parseErr, request.URL.Path, request.URL.RawQuery))
			http.Error(response, "400 Malformed request.", 400)
			return
		}
		scores := calculate_locality_scores(query, locality, places)
		results := Scores_to_result(scores, places, count)
		dat, err := json.Marshal(results)
		var buffer bytes.Buffer
		if err == nil {
			buffer.WriteString(fmt.Sprintf("%s", dat))
		} else {
			log.Panicf("Marshalling error %v", err)
			http.Error(response, "500 Internal server error.", 500)
			return
		}
		log.Println(fmt.Sprintf("Succesfully handled request %s %s", request.URL.Path, request.URL.RawQuery))
		response.Header().Set("Content-Type", "application/json; charset=utf-8")
		response.Header().Set("Content-Length", strconv.Itoa(buffer.Len()))
		io.WriteString(response, buffer.String())
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
		case places[i].Locality != locality:
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
