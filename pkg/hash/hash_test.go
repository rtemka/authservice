package hash

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя")

func randString(n int) string {
	if n < 0 {
		return ""
	}
	rs := make([]rune, n)
	for i := range rs {
		rs[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(rs)
}

func TestHashPassword(t *testing.T) {
	s := randString(10)
	h, err := HashPassword(s)
	if err != nil {
		t.Errorf("HashPassword() err = %v", err)
	}

	t.Run("right_password", func(t *testing.T) {
		if err := ComparePasswords(s, h); err != nil {
			t.Error("HashPassword() = passwords don't match")
		}
	})

	t.Run("wrong_password", func(t *testing.T) {
		w := randString(10)
		if err := ComparePasswords(w, h); err == nil {
			t.Error("HashPassword() = passwords match, but they shouldn't")
		}
	})

	t.Run("hashed_twice_password", func(t *testing.T) {
		h2, err := HashPassword(s)
		if err != nil {
			t.Errorf("HashPassword() err = %v", err)
		}
		if h2 == h {
			t.Error("HashPassword() = got same passwords hashes, want different")
		}
	})
}
