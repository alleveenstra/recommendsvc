package recommendsvc

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
)

func Http_json(output string, response http.ResponseWriter, request *http.Request) {
	var buffer bytes.Buffer
	buffer.WriteString(output)
	log.Printf("Succesfully handled request %s %s", request.URL.Path, request.URL.RawQuery)
	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	response.Header().Set("Content-Length", strconv.Itoa(buffer.Len()))
	io.WriteString(response, buffer.String())
}

func Http_error(err error, code int, response http.ResponseWriter, request *http.Request) {
	if err != nil {
		log.Panicf("%v %s %s", err, request.URL.Path, request.URL.RawQuery)
		switch code {
		case 400:
			http.Error(response, "400 Malformed request.", 400)
		case 500:
			http.Error(response, "500 Internal server error.", 500)
		}
	}
}

func Log_error(err error) {
	if err != nil {
		log.Fatal(fmt.Sprintf("%v", err))
	}
}

func Fatal_error(err error) {
	if err != nil {
		log.Fatal(fmt.Sprintf("%v", err))
		os.Exit(-1)
	}
}

func Scores_to_result(scores map[int]float64, places []Place, count int) []Result {
	scores_sortedmap := NewSortedMap(scores)
	sort.Sort(scores_sortedmap)
	results := make([]Result, count)
	for i := 0; i < count; i++ {
		key := scores_sortedmap.S[i]
		score := scores_sortedmap.M[key]
		result := NewResult(places[key], score)
		results[i] = *result
	}
	return results
}

func Distance(left []float64, right []float64) float64 {
	deg := math.Pi / 180.0
	phi1 := left[0] * deg
	phi2 := right[0] * deg
	lam12 := (right[1] - left[1]) * deg
	d2 := math.Pow(math.Cos(phi1)*math.Sin(phi2)-math.Sin(phi1)*math.Cos(phi2)*math.Cos(lam12), 2.0) + math.Pow(math.Cos(phi2)*math.Sin(lam12), 2.0)
	return 6371.009 * math.Asin(math.Sqrt(d2))
}
