package utils

import (
	"database/sql/driver"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const (
	roundHalfUpDivisor = 2
)

type Money struct {
	Raw int64
}

func NewMoneyFromString(s string) (Money, error) {
	r := new(big.Rat)
	if _, ok := r.SetString(s); !ok {
		return Money{}, fmt.Errorf("invalid money string")
	}

	scale := new(big.Rat).SetFloat64(math.Pow10(2))
	r.Mul(r, scale)

	val := new(big.Int)
	r.Add(r, big.NewRat(1, roundHalfUpDivisor))
	r.FloatString(0)
	val.Div(r.Num(), r.Denom())

	return Money{Raw: val.Int64()}, nil
}

func (m Money) String() string {
	r := new(big.Rat).SetInt64(m.Raw)
	scale := new(big.Rat).SetInt64(100)
	r.Quo(r, scale)
	return r.FloatString(2)
}

func (m Money) Add(other Money) Money {
	return Money{Raw: m.Raw + other.Raw}
}

func (m Money) Sub(other Money) Money {
	return Money{Raw: m.Raw - other.Raw}
}

func (m Money) IsNegative() bool {
	return m.Raw < 0
}

func (m Money) IsZero() bool {
	return m.Raw == 0
}

func (m *Money) Scan(value interface{}) error {
	switch v := value.(type) {
	case int64:
		m.Raw = v
		return nil
	case []byte:
		i, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse money from bytes: %w", err)
		}
		m.Raw = i
		return nil
	default:
		return fmt.Errorf("unsupported scan type for Money: %T", v)
	}
}

func (m Money) Value() (driver.Value, error) {
	return m.Raw, nil
}
