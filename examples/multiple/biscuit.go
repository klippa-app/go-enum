//go:generate go run --mod=mod github.com/klippa-app/go-enum -case=snake -gql=full -json -bson -xml -ent
package multiple

type Biscuit int

const (
	BiscuitDigestive Biscuit = 0 //enum:default
	BiscuitHobnob    Biscuit = 1 << iota
	BiscuitNice
	BiscuitJammieDodger
	BiscuitShortbread
	BiscuitGingerNut
)
