package main

import "fmt"

type A struct {
	name string
}

func (a *A) printName() {
	fmt.Println(a.name)
}

type B struct {
	val int
}

func (b *B) printName() {
	fmt.Println(b.val)
}

type C struct {
	A
	B
	count int
}

func main() {
	//遍历gridX集合中每个格子的gid
	gridX := make([]int, 5)
	grids := map[string]int{"2": 1, "3": 4}
	for _, v := range grids {
		gridX = append(gridX, v)
	}
	fmt.Println(gridX)
}
