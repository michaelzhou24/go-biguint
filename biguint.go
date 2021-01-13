package biguint

import (
	"errors"
	"fmt"
)

// BigUInt type definition, containing a slice of unsigned bytes
// unsigned ints should be split up into 2 digit base 16 chunks,
// indexed from least to most significant, e.g:
// []uint8{ 0x00, 0xff } <=> 0xff00
//
// this is also an example of slice syntax, which are
// discussed in more detail here https://blog.golang.org/slices-intro
type BigUInt struct {
	data []uint8
}

// ErrUnderflow is the underflow error for subtraction. See https://blog.golang.org/go1.13-errors
// for an up-to-date discussion of how to define and/or deal with errors.
//
// in this assignment, you just have to return this error in the correct situation
var ErrUnderflow = errors.New("arithmetic underflow")

// bytesFromUInt64 takes an unsigned 64-bit integer and converts it into an array of bytes,
// following the established scheme for this assignment (least to most significant bytes).
//
// notice that the resulting slice does not include any leading zeroes, stopping at the
// most significant non-zero byte
func bytesFromUInt64(src uint64) []uint8 {
	res := make([]uint8, 0, 8) // allocates a slice with capacity 8 but size 0, which
	// will "grow" as needed, up to the size of a uint64
	acc := src
	for acc != 0 {
		res = append(res, uint8(acc&0xFF)) // casts, like from 64 to 8 bit ints, are almost always explicit in golang
		acc >>= 8
	}
	return res
}

// NewBigUInt is the constructor for a BigUInt, based on a uint64.
// This function relies on bytesFromUInt64.
func NewBigUInt(i uint64) *BigUInt {
	return &BigUInt{data: bytesFromUInt64(i)}
}

// Add method for BigUInt.
//
// Increases x by the number represented by y, returning x.
// Note that x's slice's size may increase as a result of this operation.
// TODO
func (x *BigUInt) Add(y *BigUInt) *BigUInt {
	var carry uint16
	var result uint16
	if len(x.data) < len(y.data) {
		tmp := x
		x = y.Copy()
		y = tmp
	}
	i := 0
	for i < len(x.data) {
		if i < len(y.data)  {
			result = uint16(x.data[i]) + uint16(y.data[i]) + carry
			//fmt.Printf("Adding bytes x=%d + y=%d + carry=%d\n", uint16(x.data[i]), uint16(y.data[i]), carry)
			carry = result / 256
			result = result % 256
			//fmt.Printf("Got result %d, carry %d\n", result, carry)
			x.data[i] = uint8(result)
		} else {
			if carry > 0 {
				result = uint16(x.data[i]) + carry
				carry = result / 256
				result = result % 256
				x.data[i] = uint8(result)
			}
		}
		result = 0
		i++
	}
	if carry > 0 {
		//fmt.Printf("Adding carry at the front -> %d\n", carry)
		x.data = append(x.data, uint8(carry))
	}

	return x
}

// Subtract method for BigUInt.
//
// Decreases x by the number represented by y, returning x.
// Note that x's slice's size may decrease as a result of this operation.
//
// If y > x, then (nil, ErrUnderflow) should be returned, and
// x should be unchanged.
// TODO
func (x *BigUInt) Subtract(y *BigUInt) (*BigUInt, error) {
	return nil, errors.New("not implemented")
}

// Bytes provides access to the raw bytes underlying a given BigUInt
func (x *BigUInt) Bytes() []uint8 {
	return x.data
}

// String generates a string representing x, under the following scheme:
// - digits should be printed in base 16, with lowercase letters
// - groups of 8 digits should be separated by underscores
// - no leading zeroes should be printed
// - the string should be prefixed with "0x"
//
// see https://golang.org/pkg/fmt/#Formatter for reference material on
// golang's printf-style string formatting
func (x *BigUInt) String() string {
	if len(x.data) == 0 {
		return "0x0"
	}
	str := "0x"
	for i := len(x.data) - 1; i >= 0; i-- {
		if x.data[i] > 0xF || i == len(x.data)-1 {
			str += fmt.Sprintf("%x", x.data[i])
		} else {
			str += fmt.Sprintf("0%x", x.data[i])
		}
		if i != 0 && i%4 == 0 {
			str += "_"
		}
	}
	return str
}

// Copy generates a fully independent (deep) copy of a given BigUInt
func (x *BigUInt) Copy() *BigUInt {
	len := len(x.data)
	ret := make([]uint8, len)
	for i, v := range x.data {
		ret[i] = v
	}
	bigUInt := BigUInt{data: ret}
	return &bigUInt
}
