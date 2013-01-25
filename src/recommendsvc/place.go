package recommendsvc

import (
	"errors"
	"math"
)

type Place struct {
	Id       int
	Name     string
	Locality string
	Geo      []float64
	Features []float64
}

func FindPlace(id int, places []Place) (*Place, error) {
	for _, place := range places {
		if place.Id == id {
			return &place, nil
		}
	}
	return nil, errors.New("Place not found")
}

func (this *Place) Score(other *Place) float64 {
	length := len(this.Features)
	var value float64
	for i := 0; i < length; i++ {
		value += math.Pow(this.Features[i]-other.Features[i], 2.0)
	}
	return math.Sqrt(value)
}
