package utils

import (
	"math/rand"
	"reflect"
)

// Full random shuffle
func Shuffle(slice interface{}) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()
	for i := length - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		swap(i, j)
	}
}

// Controlled shuffle by int64 value
func ShuffleByInt64(slice interface{}, rndInt int64) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	source := rand.NewSource(rndInt)
	random := rand.New(source)
	length := rv.Len()
	for i := length - 1; i > 0; i-- {
		j := random.Intn(i + 1)
		swap(i, j)
	}
}
