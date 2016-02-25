package strconv

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertFloat(t *testing.T, s string, ef float64) {
	f, _ := Float([]byte(s))
	assert.Equal(t, ef, f, "must match float in "+s)
}

func assertAppendFloat(t *testing.T, f float64, es string) {
	b := AppendFloat(make([]byte, 100), f)
	assert.Equal(t, es, string(b), "must match float to "+es)
}

func TestParseFloat(t *testing.T) {
	assertFloat(t, "5", 5)
	assertFloat(t, "5.1", 5.1)
	assertFloat(t, "-5.1", -5.1)
	assertFloat(t, "5.1e-2", 5.1e-2)
	assertFloat(t, "5.1e+2", 5.1e+2)
	assertFloat(t, "0.0e1", 0.0e1)
	assertFloat(t, "18446744073709551620", 18446744073709551620.0)
	assertFloat(t, "1e23", 1e23)
	// TODO: hard to test due to float imprecision
	// assertFloat(t, "1.7976931348623e+308", 1.7976931348623e+308)
	// assertFloat(t, "4.9406564584124e-308", 4.9406564584124e-308)
}

func TestAppendFloat(t *testing.T) {
	assertAppendFloat(t, 1.2e3, "1.2e3")
}

func BenchmarkFloatToBytes1(b *testing.B) {
	n := 0
	r := []byte{}
	f := 1.2e3
	for i := 0; i < b.N; i++ {
		r = strconv.AppendFloat(r[:0], f, 'f', -1, 32)
		n += len(r)
	}
}

func BenchmarkFloatToBytes2(b *testing.B) {
	n := 0
	r := []byte{}
	f := 1.2e3
	for i := 0; i < b.N; i++ {
		r = AppendFloat(r[:0], f)
		n += len(r)
	}
}
