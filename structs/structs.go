package main

import "fmt"

type position struct {
	x float32
	y float32
}

type badGuy struct {
	name   string
	health int
	pos    position
}

func whereIsBadGuy(guy badGuy) {
	x := guy.pos.x
	y := guy.pos.y
	fmt.Println("(", x, ",", y, ")")
}

func main() {

	p := position{4, 2}
	fmt.Println(p.x)

	badguy := badGuy{"Jabba The Hut", 100, p}
	fmt.Println(badguy)
	whereIsBadGuy(badguy)
}
