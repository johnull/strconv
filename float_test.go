package strconv // import "github.com/tdewolff/strconv"

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFloat(t *testing.T) {
	var floatTests = []struct {
		f        string
		expected float64
	}{
		{"5", 5},
		{"5.1", 5.1},
		{"-5.1", -5.1},
		{"5.1e-2", 5.1e-2},
		{"5.1e+2", 5.1e+2},
		{"0.0e1", 0.0e1},
		{"18446744073709551620", 18446744073709551620.0},
		{"1e23", 1e23},
		// TODO: hard to test due to float imprecision
		// {"1.7976931348623e+308", 1.7976931348623e+308)
		// {"4.9406564584124e-308", 4.9406564584124e-308)
	}
	for _, tt := range floatTests {
		f, _ := ParseFloat([]byte(tt.f))
		assert.Equal(t, tt.expected, f, "ParseFloat must give expected result in "+tt.f)
	}
}

func TestAppendFloat(t *testing.T) {
	var floatTests = []struct {
		f        float64
		prec     int
		expected string
	}{
		{0, 6, "0"},
		{1, 6, "1"},
		{123, 6, "123"},
		{0.123456, 6, ".123456"},
		{12e2, 6, "1200"},
		{12e3, 6, "12e3"},
		{0.1, 6, ".1"},
		{0.001, 6, ".001"},
		{0.0001, 6, "1e-4"},
		{-1, 6, "-1"},
		{-123, 6, "-123"},
		{-123.456, 6, "-123.456"},
		{-12e3, 6, "-12e3"},
		{-0.1, 6, "-.1"},
		{-0.0001, 6, "-1e-4"},
		{0.000100009, 10, "100009e-9"},
		{0.0001000009, 10, "1.000009e-4"},
		{1e18, 0, "1e18"},
		{1e19, 0, ""}, // overflow
	}
	for _, tt := range floatTests {
		f, _ := AppendFloat([]byte{}, tt.f, tt.prec)
		assert.Equal(t, tt.expected, string(f), "AppendFloat must give expected result with "+strconv.FormatFloat(tt.f, 'f', -1, 64))
	}
}

////////////////////////////////////////////////////////////////

func TestAppendFloatStress(t *testing.T) {
	r := rand.New(rand.NewSource(99))
	prec := 10
	for i := 0; i < 0; i++ {
		f := r.ExpFloat64()
		f = math.Floor(f*float64(prec)) / float64(prec)

		b, _ := AppendFloat([]byte{}, f, prec)
		f2, _ := strconv.ParseFloat(string(b), 64)
		if f != f2 {
			fmt.Println("Bad:", f, "!=", f2, "in", string(b))
		}
	}
}

func BenchmarkFloatToBytes1(b *testing.B) {
	r := []byte{} //make([]byte, 10)
	f := 123.456
	for i := 0; i < b.N; i++ {
		r = strconv.AppendFloat(r[:0], f, 'g', 6, 64)
	}
}

func BenchmarkFloatToBytes2(b *testing.B) {
	r := make([]byte, 10)
	f := 123.456
	for i := 0; i < b.N; i++ {
		r, _ = AppendFloat(r[:0], f, 6)
	}
}

func BenchmarkModf1(b *testing.B) {
	f := 123.456
	x := 0.0
	for i := 0; i < b.N; i++ {
		a, b := math.Modf(f)
		x += a + b
	}
}

func BenchmarkModf2(b *testing.B) {
	f := 123.456
	x := 0.0
	for i := 0; i < b.N; i++ {
		a := float64(int64(f))
		b := f - a
		x += a + b
	}
}

func BenchmarkPrintInt1(b *testing.B) {
	X := int64(123456789)
	n := LenInt(X)
	r := make([]byte, n)
	for i := 0; i < b.N; i++ {
		x := X
		j := n
		for x > 0 {
			j--
			r[j] = '0' + byte(x%10)
			x /= 10
		}
	}
}

func BenchmarkPrintInt2(b *testing.B) {
	X := int64(123456789)
	n := LenInt(X)
	r := make([]byte, n)
	for i := 0; i < b.N; i++ {
		x := X
		j := n
		for x > 0 {
			j--
			newX := x / 10
			r[j] = '0' + byte(x-10*newX)
			x = newX
		}
	}
}

var int64pow10 = []int64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18,
}

func BenchmarkPrintInt3(b *testing.B) {
	X := int64(123456789)
	n := LenInt(X)
	r := make([]byte, n)
	for i := 0; i < b.N; i++ {
		x := X
		j := 0
		for j < n {
			pow := int64pow10[n-j-1]
			tmp := x / pow
			r[j] = '0' + byte(tmp)
			j++
			x -= tmp * pow
		}
	}
}
