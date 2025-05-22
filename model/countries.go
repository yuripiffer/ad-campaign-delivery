package model

type (
	Country string
)

// REMINDER: also insert the country in map Countries whenever
// a new country is added as a constant.
const (
	France       Country = "FR"
	Spain        Country = "ES"
	UK           Country = "UK"
	UnitedStates Country = "US"
)

var Countries = map[string]Country{
	"FR": France,
	"ES": Spain,
	"UK": UK,
	"US": UnitedStates,
}
