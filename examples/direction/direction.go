//go:generate go run --mod=mod github.com/klippa-app/go-enum -name=Direction -json
package directions

import "fmt"

type Direction int

const (
	North Direction = iota
	South
	West
	East
)

//go:generate go run --mod=mod github.com/klippa-app/go-enum -name=Cardinal -json

type Cardinal string

const (
	CardUnkown Cardinal = "unknown" //enum:default,invalid
	CardNorth  Cardinal = "north"
	CardSouth  Cardinal = "south"
	CardWest   Cardinal = "west"
	CardEast   Cardinal = "east"
)

func init() {
	direction := North
	cardinal := CardNorth

	fmt.Println(direction, cardinal)
}
