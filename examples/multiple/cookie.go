//go:generate go run --mod=mod github.com/klippa-app/go-enum -case=pascal -gql=full -json -bson -xml -ent
package multiple

type Cookie int

const (
	ChocolateDigestive  Cookie = 0 //enum:default
	ChocolateShortbread Cookie = 1 << iota
	ChocolateFinger
	JaffaCake
	ChocolateHobnob
)
