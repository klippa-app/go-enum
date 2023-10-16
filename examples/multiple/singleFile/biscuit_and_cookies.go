//go:generate go run --mod=mod github.com/klippa-app/go-enum -name=Biscuit -case=snake -gql=full -json -bson -xml -ent
package singlefile

type Biscuit int

const (
	BiscuitDigestive Biscuit = 0 //enum:default
	BiscuitHobnob    Biscuit = 1 << iota
	BiscuitNice
	BiscuitJammieDodger
	BiscuitShortbread
	BiscuitGingerNut
)

//go:generate go run --mod=mod github.com/klippa-app/go-enum -name=Cookie -case=pascal -gql=full -json -bson -xml -ent

type Cookie int

const (
	ChocolateDigestive  Cookie = 0 //enum:default
	ChocolateShortbread Cookie = 1 << iota
	ChocolateFinger
	JaffaCake
	ChocolateHobnob
)
