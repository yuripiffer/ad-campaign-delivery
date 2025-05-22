package model

type (
	Country string
)

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
