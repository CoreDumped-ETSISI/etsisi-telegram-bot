package salas

import "fmt"

func indexToTime(i int) string {
	starth := 9
	h := (i / 2) + starth
	m := (i % 2) * 30

	return fmt.Sprintf("%v:%02d", h, m)
}

func timeToIndex(h, m int) int {
	sh := h - 9
	sm := m

	si := sh * 2
	if sm >= 30 {
		si++
	}

	return si
}
