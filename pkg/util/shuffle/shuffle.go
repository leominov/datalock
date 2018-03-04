package shuffle

import (
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

func isShuffleEnabledByCookie(r *http.Request) (int64, bool) {
	cookie, err := r.Cookie("shuffle")
	if err != nil {
		return 0, false
	}
	if cookie == nil {
		return 0, false
	}
	if len(cookie.Value) == 0 {
		return 0, false
	}
	i, err := strconv.ParseInt(cookie.Value, 10, 64)
	if err != nil {
		return 0, false
	}
	if i == 0 {
		return time.Now().UnixNano(), true
	}
	return i, true
}

func isShuffleEnabledByQuery(r *http.Request) (int64, bool) {
	value := r.URL.Query().Get("_shuffle")
	if len(value) == 0 {
		return 0, false
	}
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}
	if i == 0 {
		return time.Now().UnixNano(), true
	}
	return i, true
}

func IsShuffleEnabled(r *http.Request) (int64, bool) {
	if i, ok := isShuffleEnabledByQuery(r); ok {
		return i, true
	}
	if i, ok := isShuffleEnabledByCookie(r); ok {
		return i, true
	}
	return 0, false
}

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
