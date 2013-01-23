package recommendsvc

type Result struct {
	Id int
	Name string
	Locality string
	Geo []float64
	Distance float64
	Score float64
}

func NewResult(place Place, distance float64, score float64) *Result {
	result := new(Result)
	result.Id = place.Id
	result.Name = place.Name
	result.Locality = place.Locality
	result.Geo = place.Geo
	result.Distance = distance
	result.Score = score
	return result
}


