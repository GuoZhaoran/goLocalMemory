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
	setGroup.Add("love", "c++")
	setGroup.Add("love", "javascript")

	_, _ = setGroup.Remove("book", "java")
	_, _ = setGroup.Remove("book", "python")
	setGroup.Add("programming", "c++")
	keyExists, memberExists := setGroup.Remove("book", "c++")
	setGroup.Add("programming", "java")
	fmt.Println("key exists:", keyExists, "member exists:", memberExists)
	setGroup.Add("love", "go")
	setGroup.Add("book", "c#")
	setGroup.Remove("book", "java")

	intersect := setGroup.Intersect( "book", "programming")
	fmt.Println("----------- intersect -----------")
	for _, intersectKey := range intersect {
		fmt.Println(intersectKey)
	}
	fmt.Println("----------- different -----------")
	different := setGroup.Different("love", "programming", "book")
	for _, different := range different {
		fmt.Println(different)
	}
	fmt.Println("----------- union -----------")
	unions := setGroup.Union("love", "programming", "book")
	for _, unions := range unions {
		fmt.Println(unions)
	}

	setGroup.FPrint()
}
