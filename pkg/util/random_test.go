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
	r, _ := regexp.Compile("^[a-zA-Z]*$")
	randomValue := RandString(randomLength)

	if r.MatchString(randomValue) == false {
		t.Fatalf(`Expected %s to have valid characters`, randomValue)
	}

}
