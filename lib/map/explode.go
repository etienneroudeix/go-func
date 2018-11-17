package Map

type Pos struct {
	X int
	Y int
}

type ShipPos struct {
	Start Pos
	End Pos
}

func (sp ShipPos) IsHorizontal() bool {
	return sp.Start.Y == sp.End.Y
}

func (sp ShipPos) IsVertical() bool {
	return sp.Start.X == sp.End.X
}


func Explode(shipPos ShipPos) (explosion []Pos) {

	if shipPos.IsHorizontal() {
		if shipPos.Start.X < shipPos.End.X {
			for i := shipPos.Start.X; i <= shipPos.End.X; i++ {
				explosion = append(explosion, Pos{i, shipPos.Start.Y})
			}
			return
		}

		for i := shipPos.End.X; i <= shipPos.Start.X; i++ {
			explosion = append(explosion, Pos{i, shipPos.Start.Y})
		}
		return
	}

	if shipPos.Start.Y < shipPos.End.Y {
		for i := shipPos.Start.Y; i <= shipPos.End.Y; i++ {
			explosion = append(explosion, Pos{shipPos.Start.X, i})
		}
		return
	}

	for i := shipPos.End.Y; i <= shipPos.Start.Y; i++ {
		explosion = append(explosion, Pos{shipPos.Start.X, i})
	}
	return
}