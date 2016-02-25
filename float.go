package strconv

import "math"

var float64pow10 = []float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9,
	1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19,
	1e20, 1e21, 1e22,
}

// Float parses a byte-slice and returns the float it represents.
// If an invalid character is encountered, it will stop there.
func ParseFloat(b []byte) (float64, bool) {
	i := 0
	neg := false
	if i < len(b) && (b[i] == '+' || b[i] == '-') {
		neg = b[i] == '-'
		i++
	}
	dot := -1
	trunk := -1
	n := uint64(0)
	for ; i < len(b); i++ {
		c := b[i]
		if c >= '0' && c <= '9' {
			if trunk == -1 {
				if n > math.MaxUint64/10 {
					trunk = i
				} else {
					n *= 10
					n += uint64(c - '0')
				}
			}
		} else if dot == -1 && c == '.' {
			dot = i
		} else {
			break
		}
	}
	f := float64(n)
	if neg {
		f = -f
	}
	mantExp := int64(0)
	if dot != -1 {
		if trunk == -1 {
			trunk = i
		}
		mantExp = int64(trunk - dot - 1)
	} else if trunk != -1 {
		mantExp = int64(trunk - i)
	}
	expExp := int64(0)
	if i < len(b) && (b[i] == 'e' || b[i] == 'E') {
		i++
		if e, ok := Int(b[i:]); ok {
			expExp = e
		}
	}
	exp := expExp - mantExp
	// copied from strconv/atof.go
	if exp == 0 {
		return f, true
	} else if exp > 0 && exp <= 15+22 { // int * 10^k
		// If exponent is big but number of digits is not,
		// can move a few zeros into the integer part.
		if exp > 22 {
			f *= float64pow10[exp-22]
			exp = 22
		}
		if f <= 1e15 && f >= -1e15 {
			return f * float64pow10[exp], true
		}
	} else if exp < 0 && exp >= -22 { // int / 10^k
		return f / float64pow10[-exp], true
	}
	f *= math.Pow10(int(-mantExp))
	return f * math.Pow10(int(expExp)), true
}

func AppendFloat(b []byte, f float64, prec int) ([]byte, bool) {
	neg := false
	if f < 0.0 {
		f = -f
		neg = true
	}
	if prec >= len(float64pow10) {
		return b, false
	}
	f *= float64pow10[prec]
	if f >= float64(math.MaxInt64) {
		return b, false
	}

	// calculate mantissa and exponent
	mant := int64(f)
	mantLen := LenInt(mant)
	mantExp := mantLen - prec - 1
	if mant == 0 {
		return append(b, '0'), true
	}

	// expLen is zero for positive exponents, because positive exponents are determined later on in the big conversion loop
	exp := 0
	expLen := 0
	if mantExp < -3 {
		exp = mantExp
		mantExp = 0
		expLen = 3      // e + minus + digit
		if exp <= -10 { // exp is never lower than -18
			expLen++
		}
	} else if mantExp < -1 {
		mantLen += -mantExp - 1 // extra zero between dot and first digit
	}

	// reserve space in b
	i := len(b)
	maxLen := 1 + mantLen + expLen // dot + mantissa digits + exponent
	if neg {
		maxLen++
	}
	if i+maxLen > cap(b) {
		b = append(b, make([]byte, maxLen)...)
	} else {
		b = b[:i+maxLen]
	}

	// write to string representation
	if neg {
		b[i] = '-'
		i++
	}

	// big conversion loop, start at the end and move to the front
	// initially print trailing zeros and remove them later on
	// for example if the first non-zero digit is three positions in front of the dot, it will overwrite the zeros with a positive exponent
	zero := true
	last := i + mantLen      // right-most position of digit that is non-zero
	dot := last - prec - exp // position of dot
	j := last
	for mant > 0 {
		if j == dot {
			b[j] = '.'
			j--
		}
		digit := mant % 10
		if zero && digit > 0 {
			// first non-zero digit, if we are still behind the dot we can trim the end to this position
			// otherwise trim to the dot (including the dot)
			if j > dot {
				i = j + 1
				// decrease negative exponent further to get rid of dot
				if exp < 0 {
					relExp := j - dot
					// getting rid of the dot shouldn't lower exponent to two digits, unless it's already two digits
					if exp-relExp > -10 || exp <= -10 { // exp is never lower than -18
						exp -= relExp
						dot = j
						j--
						i--
					}
				}
			} else {
				i = dot
			}
			last = j
			zero = false
		}
		b[j] = '0' + byte(digit)
		j--
		mant /= 10
	}

	// handle 0.1
	if j == dot {
		b[j] = '.'
		j--
	}

	// extra zeros between dot and first digit
	if j > dot+1 {
		for j > dot {
			b[j] = '0'
			j--
		}
		b[j] = '.'
	}

	// add positive exponent because we have 3 or more zeros before the dot
	if last+3 < dot {
		i = last + 1
		exp = dot - last - 1
	}

	// exponent
	if exp != 0 {
		b[i] = 'e'
		i++
		if exp < 0 {
			b[i] = '-'
			i++
			exp = -exp
		}
		if exp >= 10 {
			b[i+1] = '0' + byte(exp%10)
			b[i] = '0' + byte(exp/10)
			i += 2
		} else {
			b[i] = '0' + byte(exp%10)
			i++
		}
	}
	return b[:i], true
}
