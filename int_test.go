package strconv

import (
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertInt(t *testing.T, s string, ei int64) {
	i, _ := Int([]byte(s))
	assert.Equal(t, ei, i, "must match integer in "+s)
}

func TestParseInt(t *testing.T) {
	assertInt(t, "5", 5)
	assertInt(t, "99", 99)
	assertInt(t, "999", 999)
	assertInt(t, "-5", -5)
	assertInt(t, "+5", 5)
	assertInt(t, "9223372036854775807", 9223372036854775807)
	assertInt(t, "9223372036854775808", 0)
	assertInt(t, "-9223372036854775807", -9223372036854775807)
	assertInt(t, "-9223372036854775808", -9223372036854775808)
	assertInt(t, "-9223372036854775809", 0)
	assertInt(t, "18446744073709551620", 0)
	assertInt(t, "a", 0)
}

func TestLenInt(t *testing.T) {
	var lenIntTests = []struct {
		number   int64
		expected int
	}{
		{0, 1},
		{1, 1},
		{10, 2},
		{99, 2},

		// coverage
		{100, 3},
		{1000, 4},
		{10000, 5},
		{100000, 6},
		{1000000, 7},
		{10000000, 8},
		{100000000, 9},
		{1000000000, 10},
		{10000000000, 11},
		{100000000000, 12},
		{1000000000000, 13},
		{10000000000000, 14},
		{100000000000000, 15},
		{1000000000000000, 16},
		{10000000000000000, 17},
		{100000000000000000, 18},
		{1000000000000000000, 19},
	}
	for _, tt := range lenIntTests {
		assert.Equal(t, tt.expected, LenInt(tt.number), "LenInt must give expected result in "+strconv.FormatInt(tt.number, 10))
	}
}

var num []int64

func TestMain(t *testing.T) {
	for j := 0; j < 1000; j++ {
		num = append(num, rand.Int63n(1000))
	}
}

func BenchmarkLenIntLog(b *testing.B) {
	n := 0
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			n += int(math.Log10(math.Abs(float64(num[j])))) + 1
		}
	}
}

func BenchmarkLenIntSwitch(b *testing.B) {
	n := 0
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			n += LenInt(num[j])
		}
	}
}
