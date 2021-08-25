package main

import (
	"fmt"
	"goLocalMemory/setGroup"
)

func main() {
	setGroup := setGroup.New()
	setGroup.Add("programming", "php")
	setGroup.Add("programming", "go")
	setGroup.Add("programming", "python")

	setGroup.Add("book", "php")
	setGroup.Add("book", "go")
	setGroup.Add("book", "java")
	setGroup.Add("book", "python")

	setGroup.Add("love", "java")
	setGroup.Add("love", "python")
	setGroup.Add("love", "php")

	_, _ = setGroup.Remove("book", "java")
	setGroup.Add("programming", "c++")
	keyExists, memberExists := setGroup.Remove("book", "c++")
	setGroup.Add("programming", "java")
	fmt.Println("key exists:", keyExists, "member exists:", memberExists)
	setGroup.Add("love", "go")

	intersect := setGroup.Intersect("programming", "book")
	for _, intersectKey := range intersect {
		fmt.Println(intersectKey)
	}

	setGroup.FPrint()
}
