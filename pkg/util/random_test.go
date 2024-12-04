package util

import (
	"regexp"
	"testing"
)

func TestRandomZeroLength(t *testing.T) {
	randomLength := 0
	randomValue := RandString(randomLength)

	if len(randomValue) != randomLength {
		t.Fatalf(`Expected %s to have %d length`, randomValue, randomLength)
	}

}

func TestRandomLength(t *testing.T) {
	randomLength := 10
	randomValue := RandString(randomLength)

	if len(randomValue) != randomLength {
		t.Fatalf(`Expected %s to have %d length`, randomValue, randomLength)
	}

}

func TestBigLength(t *testing.T) {
	randomLength := 10000
	randomValue := RandString(randomLength)

	if len(randomValue) != randomLength {
		t.Fatalf(`Expected %s to have %d length`, randomValue, randomLength)
	}

}

func TestValidCharacters(t *testing.T) {
	randomLength := 100
	r, _ := regexp.Compile("^[0-9A-F]*$")
	randomValue := RandString(randomLength)

	if r.MatchString(randomValue) == false {
		t.Fatalf(`Expected %s to have valid characters`, randomValue)
	}

}

func BenchmarkSmallRandom(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandString(5)
	}
}

func BenchmarkBigRandom(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandString(500)
	}
}
