package main

import (
	"fmt"
	"goLocalMemory/setGroup"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var once sync.Once
var randomSet string

func init() {
	rand.Seed(time.Now().UnixNano())
	randomSet = getRandomString(12)
}

func getRandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

func BenchmarkAddOneSetGroup(b *testing.B) {
	setGroup := setGroup.New()
	for n := 0; n < b.N; n++ {
		setGroup.Add(randomSet, getRandomString(20))
	}
}

func BenchmarkAddMultiSetGroup(b *testing.B) {
	setGroup := setGroup.New()
	for n := 0; n < b.N; n++ {
		setGroup.Add(getRandomString(20), getRandomString(20))
	}
}

func BenchmarkRemoveMultiSetGroup(b *testing.B) {
	setGroup := setGroup.New()
	for n := 0; n < b.N; n++ {
		setGroup.Remove(getRandomString(20), getRandomString(20))
	}
}

func BenchmarkMultiIntersectSetGroup(b *testing.B) {
	setGroup := setGroup.New()
	for n := 0; n < b.N; n++ {
		setGroup.Intersect(getRandomString(20), getRandomString(20), getRandomString(20), getRandomString(20))
	}
}




