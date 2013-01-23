package recommendsvc

import (
	"fmt"
	"net/http"
	"log"
	"sort"
	"encoding/json"
	"strconv"
	"math"
	"errors"
	"io"
	"bytes"
)

func Build_recommendation_handler(places []Place) (func (response http.ResponseWriter, request *http.Request)) {
	return func(response http.ResponseWriter, request *http.Request) {
		id, count, geo, rng, parseErr := parse_request(request)
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
		scores_map := calculate_scores(query, geo, places, rng)
		scores_sortedmap := NewSortedMap(scores_map)
		sort.Sort(scores_sortedmap)
		var buffer bytes.Buffer
		results := make([]Result, count)
		for i := 0; i < count; i++ {
			key := scores_sortedmap.S[i]
			score := scores_sortedmap.M[key]
			distance := distance(geo, places[key].Geo)
			result := NewResult(places[key], distance, score)
			results[i] = *result
		}
		dat, err := json.Marshal(results)
		if (err == nil) {
			buffer.WriteString(fmt.Sprintf("%s", dat));
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


func calculate_scores(query *Place, geo []float64, places []Place, rng float64) map[int] float64 {
	length := len(places)
	var output = make(map[int] float64, length)
	for i := 0; i < length; i++ {
		switch {
		case places[i].Id == query.Id:
			output[i] = 1000
		case distance(geo, places[i].Geo) > rng:
			output[i] = 1000
		default:
			output[i] = query.Score(&places[i])
		}
	}
	return output
}

func distance(left []float64, right []float64) float64 {
	deg := math.Pi / 180.0
	phi1 := left[0] * deg
	phi2 := right[0] * deg
	lam12 := (right[1] - left[1]) * deg
	d2 := math.Pow(math.Cos(phi1) * math.Sin(phi2) - math.Sin(phi1) * math.Cos(phi2) * math.Cos(lam12), 2.0) + math.Pow(math.Cos(phi2) * math.Sin(lam12), 2.0)
	return 6371.009 * math.Asin(math.Sqrt(d2))
}

func parse_request(request *http.Request) (int, int, []float64, float64, error) {
	id, idErr := strconv.Atoi(request.FormValue("id"))
	count, countErr := strconv.Atoi(request.FormValue("count"))
	lat, latErr := strconv.ParseFloat(request.FormValue("lat"), 64)
	lng, lngErr := strconv.ParseFloat(request.FormValue("lng"), 64)
	rng, rngErr := strconv.ParseFloat(request.FormValue("rng"), 64)
	geo := []float64{lat,lng}
	if idErr == nil && countErr == nil && latErr == nil && lngErr == nil && rngErr == nil {
		return id, count, geo, rng, nil
	}
	return id, count, geo, rng, errors.New("Error parsing request parameters")
}
