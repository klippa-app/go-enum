//go:generate go run --mod=mod github.com/klippa-app/go-enum -case=upper_snake -gql=full -json -bson -xml -ent
package day

type Day int

const (
	Unknown Day = 0 //enum:invalid,default
	Monday  Day = 1 << iota
	Tuesday
	Wednesday
	Thursday
	Friday
)
const (
	Saturday = Friday<<iota + 1
	Sunday
)
