package Map

func Intersect(l1 []Pos, l2 []Pos) (result []Pos) {

	for _, pos1 := range l1 {
		for _, pos2 := range l2 {
			if pos1.X == pos2.X && pos1.Y == pos2.Y {
				result = append(result, pos1)
			}
		}
	}

	return
}