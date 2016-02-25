package strconv

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