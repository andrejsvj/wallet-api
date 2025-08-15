package utils

import (
	"testing"
)

func TestNewMoneyFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{"positive decimal", "123.45", 12345, false},
		{"positive integer", "100", 10000, false},
		{"zero", "0", 0, false},
		{"negative decimal", "-50.25", -5025, false},
		{"small decimal", "0.01", 1, false},
		{"large number", "999999.99", 99999999, false},
		{"invalid format", "abc", 0, true},
		{"empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMoneyFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMoneyFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Raw != tt.want {
				t.Errorf("NewMoneyFromString() = %v, want %v", got.Raw, tt.want)
			}
		})
	}
}

func TestMoney_String(t *testing.T) {
	tests := []struct {
		name string
		m    Money
		want string
	}{
		{"positive decimal", Money{Raw: 12345}, "123.45"},
		{"positive integer", Money{Raw: 10000}, "100.00"},
		{"zero", Money{Raw: 0}, "0.00"},
		{"negative decimal", Money{Raw: -5025}, "-50.25"},
		{"small amount", Money{Raw: 1}, "0.01"},
		{"large amount", Money{Raw: 99999999}, "999999.99"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("Money.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	tests := []struct {
		name  string
		m     Money
		other Money
		want  Money
	}{
		{"positive addition", Money{Raw: 1000}, Money{Raw: 500}, Money{Raw: 1500}},
		{"zero addition", Money{Raw: 1000}, Money{Raw: 0}, Money{Raw: 1000}},
		{"negative addition", Money{Raw: 1000}, Money{Raw: -300}, Money{Raw: 700}},
		{"both negative", Money{Raw: -1000}, Money{Raw: -500}, Money{Raw: -1500}},
		{"large numbers", Money{Raw: 99999999}, Money{Raw: 1}, Money{Raw: 100000000}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Add(tt.other); got.Raw != tt.want.Raw {
				t.Errorf("Money.Add() = %v, want %v", got.Raw, tt.want.Raw)
			}
		})
	}
}

func TestMoney_Sub(t *testing.T) {
	tests := []struct {
		name  string
		m     Money
		other Money
		want  Money
	}{
		{"positive subtraction", Money{Raw: 1000}, Money{Raw: 300}, Money{Raw: 700}},
		{"zero subtraction", Money{Raw: 1000}, Money{Raw: 0}, Money{Raw: 1000}},
		{"negative result", Money{Raw: 500}, Money{Raw: 1000}, Money{Raw: -500}},
		{"negative from negative", Money{Raw: -1000}, Money{Raw: -300}, Money{Raw: -700}},
		{"large numbers", Money{Raw: 100000000}, Money{Raw: 1}, Money{Raw: 99999999}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Sub(tt.other); got.Raw != tt.want.Raw {
				t.Errorf("Money.Sub() = %v, want %v", got.Raw, tt.want.Raw)
			}
		})
	}
}

func TestMoney_IsNegative(t *testing.T) {
	tests := []struct {
		name string
		m    Money
		want bool
	}{
		{"positive", Money{Raw: 1000}, false},
		{"zero", Money{Raw: 0}, false},
		{"negative", Money{Raw: -1000}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsNegative(); got != tt.want {
				t.Errorf("Money.IsNegative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_IsZero(t *testing.T) {
	tests := []struct {
		name string
		m    Money
		want bool
	}{
		{"positive", Money{Raw: 1000}, false},
		{"zero", Money{Raw: 0}, true},
		{"negative", Money{Raw: -1000}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsZero(); got != tt.want {
				t.Errorf("Money.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Scan(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    int64
		wantErr bool
	}{
		{"int64", int64(12345), 12345, false},
		{"bytes", []byte("12345"), 12345, false},
		{"invalid bytes", []byte("abc"), 0, true},
		{"unsupported type", "string", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m Money
			err := m.Scan(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Money.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && m.Raw != tt.want {
				t.Errorf("Money.Scan() = %v, want %v", m.Raw, tt.want)
			}
		})
	}
}

func TestMoney_Value(t *testing.T) {
	m := Money{Raw: 12345}
	value, err := m.Value()
	if err != nil {
		t.Errorf("Money.Value() unexpected error = %v", err)
	}
	if value != int64(12345) {
		t.Errorf("Money.Value() = %v, want %v", value, int64(12345))
	}
}
