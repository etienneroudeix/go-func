package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"shipb/lib/linearsearch"
	"shipb/lib/map"
	"math"
	"math/rand"
	"strconv"
)

const gridSize = 8
const xOffset = 97
const yOffset = 49

var shipTypes = []int{2, 3, 3, 4, 5}
var xList = []rune("abcdefgh")
var yList = []rune("12345678")

var reader *bufio.Reader

var player Player
var computer Player

type Player struct {
	Name string
	ShipPosList []Map.ShipPos
}

var ended = false

var playerShipGrid = [8][8]string{
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
}

var strikeGrid = [8][8]string{
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
	[8]string{" ", " ", " ", " ", " ", " ", " ", " "},
}

func main() {
	reader = bufio.NewReader(os.Stdin)

	computer.Name = "Computer"

	fmt.Println("Welcome to Ship Battle !")

	player.Name = getPlayerName()
	fmt.Printf("Welcome %s !\n", player.Name)

	getShips(shipTypes...)
	fmt.Println("Your armada is in position ! Prepare for fight !")

	fmt.Printf("Generating the oponant fleet...")
	generateComputerFleet()
	fmt.Printf("Done\n")

	fmt.Println("Time for your first strike !")
	for !ended {
		strike()
	}
}

func strike() {
	fmt.Printf("Enter strike position : (format xy) :\n")

	read, _ := reader.ReadString('\n')
	read = strings.Replace(read, "\n", "", -1)

	if len(read) != 2 {
		fmt.Println("bad coord")
		strike()
		return
	}

	coord := strings.Split(read, "")
	pos := Map.Pos{int([]rune(coord[0])[0]), int([]rune(coord[1])[0])}

	if !checkPosInMap(pos) {
		fmt.Println("coord out of the map")
		strike()
		return
	}

	hit := false
	for i, shipPos := range computer.ShipPosList {
		fmt.Println(i, Map.Explode(shipPos), []Map.Pos{pos})
		fmt.Println(i, Map.Intersect(Map.Explode(shipPos), []Map.Pos{pos}))
		if len(Map.Intersect(Map.Explode(shipPos), []Map.Pos{pos})) > 0 {
			hit = true
		}
	}

	setStrikeOnGrid(pos, hit)

	if (hit) {
		fmt.Println("Hit !")
	} else {
		fmt.Println("Water :(")
	}

	return
}

func getShips(shipTypes ...int) {

	for i, v := range shipTypes {
		shipPos := getShip(v)

		player.ShipPosList = append(player.ShipPosList, shipPos)

		setShipOnGrid(shipPos, i)

		fmt.Println("Ship is ready for battle !")
	}

	return
}

func generateComputerFleet() {

	for _, v := range shipTypes {
		shipPos := randomShip(v, computer)
		computer.ShipPosList = append(computer.ShipPosList, shipPos)
	}

	return
}

func hasConflict(newPos Map.ShipPos, list []Map.ShipPos) bool {

	for _, pos := range list {
		if len(Map.Intersect(Map.Explode(newPos), Map.Explode(pos))) > 0 {
			return true
		}
	}

	return false
}

func randomShip(length int, knownPos Player) (shipPos Map.ShipPos) {
	// pick random start pos
	shipPos.Start = Map.Pos{
		X: rand.Intn(8) + xOffset,
		Y: rand.Intn(8) + yOffset,
	}

	// pick random direction <0 ^1 >2 v3
	dir := rand.Intn(4)
	switch dir {
	case 0:
		shipPos.End.X = shipPos.Start.X - length + 1
		shipPos.End.Y = shipPos.Start.Y
	case 1:
		shipPos.End.X = shipPos.Start.X
		shipPos.End.Y = shipPos.Start.Y - length + 1
	case 2:
		shipPos.End.X = shipPos.Start.X + length - 1
		shipPos.End.Y = shipPos.Start.Y
	case 3:
		shipPos.End.X = shipPos.Start.X
		shipPos.End.Y = shipPos.Start.Y + length - 1
	}

	if !checkPosInMap(shipPos.Start) || !checkPosInMap(shipPos.End) || hasConflict(shipPos, knownPos.ShipPosList) {
		fmt.Printf(".")
		return randomShip(length, knownPos)
	}

	fmt.Printf("\nShip generated at pos %s%s:%s%s\n", string(shipPos.Start.X), string(shipPos.Start.Y), string(shipPos.End.X), string(shipPos.End.Y))
	return
}

func getShip(length int) Map.ShipPos {
	fmt.Printf("Enter ship position. Expected lenght: %d (format xy:xy) [empty for random] :\n", length)

	read, _ := reader.ReadString('\n')
	read = strings.Replace(read, "\n", "", -1)

	if read == "" {
		fmt.Printf("Generating random ship...")
		return randomShip(length, player)
	}

	ok, err, shipPos := extractShipPos(read)
	if !ok {
		fmt.Printf("Bad position : %s\n", err)
		return getShip(length)
	}

	// check length
	if !shipPos.IsHorizontal() && !shipPos.IsVertical() {
		fmt.Println("Bad position : ship not straight")
		return getShip(length)
	}

	var l int
	if shipPos.IsHorizontal() {
		l = int(math.Abs(float64(shipPos.End.X - shipPos.Start.X)) + 1)
	} else {
		l = int(math.Abs(float64(shipPos.End.Y - shipPos.Start.Y)) + 1)
	}

	if l != length {
		fmt.Printf("Bad position : bad length, %d instead of %d\n", l, length)
		return getShip(length)
	}

	if hasConflict(shipPos, player.ShipPosList) {
		fmt.Println("Bad position : conflict with already set ships")
		return getShip(length)
	}

	return shipPos
}

func extractShipPos(data string) (result bool, err string, shipPos Map.ShipPos) {
	result = false

	pos := strings.Split(data, ":")
	if len(pos) != 2 {
		err = "bad format"
		return
	}

	if len(pos[0]) != 2 {
		err = "bad coord"
		return
	}
	if len(pos[1]) != 2 {
		err = "bad coord"
		return
	}

	coord1 := strings.Split(pos[0], "")
	shipPos.Start = Map.Pos{int([]rune(coord1[0])[0]), int([]rune(coord1[1])[0])}
	coord2 := strings.Split(pos[1], "")
	shipPos.End = Map.Pos{int([]rune(coord2[0])[0]), int([]rune(coord2[1])[0])}

	if !checkPosInMap(shipPos.Start) || !checkPosInMap(shipPos.End) {
		err = "coord out of the map"
		return
	}

	result = true
	return
}

func checkPosInMap(pos Map.Pos) bool {

	if !Linearsearch.Contains(xList, rune(pos.X)) {
		return false
	}

	if !Linearsearch.Contains(yList, rune(pos.Y)) {
		return false
	}

	return true
}

func getPlayerName() (name string) {
	fmt.Println("What is your name ?")

	name, _ = reader.ReadString('\n')
	name = strings.Replace(name, "\n", "", -1)

	if len(name) == 0 {
		name = getPlayerName()
	}

	return name
}

func setShipOnGrid(shipPos Map.ShipPos, index int) {

	for _, pos := range Map.Explode(shipPos) {
		playerShipGrid[pos.X - xOffset][pos.Y - yOffset] = "\033[1;" + strconv.Itoa(31 + index) + "mo\033[0m"
	}

	showGrid(playerShipGrid)
}

func setStrikeOnGrid(pos Map.Pos, hit bool) {

	color := "1;34"
	if hit {
		color = "1;31"
	}

	strikeGrid[pos.X - xOffset][pos.Y - yOffset] = "\033[" + color + "mX\033[0m"

	showGrid(strikeGrid)
}

func whatIsAt(grid [8][8]string, pos Map.Pos) string {
	return grid[pos.X][pos.Y]
}

func showGrid(grid [8][8]string) {
	fmt.Printf("\n\n")
	fmt.Printf("   A   B   C   D   E   F   G   H\n")

	for y := 1; y <= gridSize; y++ {

		fmt.Printf("%d ", y)

		for x := 1; x <= gridSize; x++ {

			value := whatIsAt(grid, Map.Pos{x - 1, y - 1})

			fmt.Printf(" %s ", value)
			if x != gridSize {
				fmt.Printf("|")
			} else {
				fmt.Printf("\n")
			}
		}

		if y != gridSize {
			fmt.Printf("---------------------------------\n")
		}

	}

	fmt.Printf("\n\n")
}