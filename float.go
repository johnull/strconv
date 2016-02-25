package strconv

import (
	"math"
	"strconv"
)

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

const prec64 = 1e18                    // 2^63 = 10^18.96...
const minLen = 1 + 18 + 1 + 18 + 1 + 2 // minus + whole + point + fractional + e + exponent

func FormatFloat(b []byte, f float64) []byte {
	b = b[:0]
	// take slow path for really small or large numbers, NaN and Inf
	if !(f > -prec64 && f < prec64) {
		return strconv.AppendFloat(b, f, 'g', -1, 64)
	}

	// neg := false
	// if f < 0 {
	// 	f = -f
	// 	neg = true
	// }

	whole := int64(f)
	frac := uint64((f - float64(whole)) * prec64)

	if whole == 0 && frac == 0 {
		b = append(b, '0')
		return b[:1]
	}

	wholeLen := LenInt(whole)
	maxLen := 1 + wholeLen + 1 + 18 + 1 + 2

	i := 0
	if cap(b) < maxLen {
		b = make([]byte, maxLen)
		//fmt.Println(cap(b))
	} else {
		b = b[:maxLen]
	}

	j := i + wholeLen
	for whole > 0 {
		j--
		b[j] = '0' + byte(whole%10)
		whole /= 10
	}
	i += wholeLen

	if frac > 0 {
		b[i] = '.'
		i += 19
		j = i
		foundNonZero := false
		for frac > 0 {
			digit := frac % 10
			if !foundNonZero && digit > 0 {
				i = j
				foundNonZero = true
			}
			j--
			b[j] = '0' + byte(digit)
			frac /= 10
		}
	}

	return b[:i]
}
