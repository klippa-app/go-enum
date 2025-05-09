//go:generate go run --mod=mod github.com/klippa-app/go-enum -no-stringer -gql=full -json -bson -xml -ent
package day

type Day string

const (
	Unknown   Day = "UnKnOwN" //enum:invalid
	Monday    Day = "MoN dAy"
	Tuesday   Day = "TuEs DaY"
	Wednesday Day = "WeDnEs DaY"
	Thursday  Day = "ThUrS dAy"
	Friday    Day = "FrI dAy"
	Saturday  Day = "SaTuR dAy"
	Sunday    Day = "SuN dAy"
)

func (day Day) String() string {
	return string(day)
}
