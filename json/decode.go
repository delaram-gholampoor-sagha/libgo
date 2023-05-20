/* For license and copyright information please see the LEGAL file in the code repository */

package json

import (
	"bytes"
	"encoding/base64"
	"strconv"

	"libgo/convert"
	"libgo/protocol"
)

// Decoder store data to decode data by each method!
type Decoder struct {
	buf      []byte
	token    byte
	lastItem []byte

	Options struct {
		ErrorOnDuplicateKey bool
		ErrorOnNotFoundKey  bool
	}
}

//libgo:impl libgo/protocol.ObjectLifeCycle
func (d *Decoder) Init(buf []byte) (err protocol.Error) {
	d.buf = buf
	return
}
func (d *Decoder) Reinit() (err protocol.Error) { return }
func (d *Decoder) Deinit() (err protocol.Error) { return }

// Offset make d.buf to start of given offset
func (d *Decoder) Offset(o int) {
	d.buf = d.buf[o:]
}

// FindEndToken find next end json token
func (d *Decoder) FindEndToken() {
	for i, c := range d.buf {
		switch c {
		case ',':
			d.token = ','
			d.lastItem = d.buf[:i]
			d.buf = d.buf[i:]
			return
		case ']':
			d.token = ']'
			d.lastItem = d.buf[:i]
			d.buf = d.buf[i:]
			return
		case '}':
			d.token = '}'
			d.lastItem = d.buf[:i]
			d.buf = d.buf[i:]
			return
		}
	}
}

// TrimToDigit remove any data and set d.buf to first character as number
func (d *Decoder) TrimToDigit() {
	for i, c := range d.buf {
		if '0' <= c && c <= '9' {
			d.buf = d.buf[i:]
			return
		}
	}
}

// TrimSpaces remove any spaces from d.buf
func (d *Decoder) TrimSpaces() {
	for i, c := range d.buf {
		if c != ' ' {
			d.buf = d.buf[i:]
			return
		}
	}
}

// TrimToStringStart remove any data from d.buf to first character after quotation mark as '"'
func (d *Decoder) TrimToStringStart() {
	for i, c := range d.buf {
		if c == '"' {
			d.buf = d.buf[i+1:] // remove any byte before " due to don't need them
			return
		}
	}
}

// CheckNullValue check if null exist as value. pass d.buf start from after : and receive from from after , if null exist
func (d *Decoder) CheckNullValue() (null bool) {
	for i, c := range d.buf {
		switch c {
		case 'n':
			if bytes.Equal(d.buf[i:i+4], []byte("null")) {
				null = true
			}
		case '"':
			return false
		case ',':
			return
		}
	}
	return
}

// ResetToken set d.token to nil
func (d *Decoder) ResetToken() {
	d.token = 0
}

// CheckToken set d.token to nil
func (d *Decoder) CheckToken(t byte) bool {
	if d.token == t {
		d.ResetToken()
		return true
	}
	return false
}

// DecodeKey return key very safe for each decode iteration. pass d.buf start from any where and receive from after :
func (d *Decoder) DecodeKey() string {
	d.TrimToStringStart()
	var loc = bytes.IndexByte(d.buf, '"')
	if loc < 0 {
		return ""
	}

	var key []byte = d.buf[:loc]

	d.buf = d.buf[loc+1:] // remove any byte before last " due to don't need them
	loc = bytes.IndexByte(d.buf, ':')
	d.buf = d.buf[loc+1:]
	return convert.UnsafeByteSliceToString(key)
}

// NotFoundKey call in default switch of each decode iteration
func (d *Decoder) NotFoundKey() (err protocol.Error) {
	d.FindEndToken()
	return
}

// NotFoundKeyStrict call in default switch of each decode iteration in strict mode.
func (d *Decoder) NotFoundKeyStrict() protocol.Error {
	return &ErrEncodedIncludeNotDefinedKey
}

func (d *Decoder) End() bool {
	if len(d.buf) < 3 || d.token == '}' {
		return true
	}
	return false
}

// DecodeBool convert string base boolean to bool. pass d.buf start from after : and receive from after ,
func (d *Decoder) DecodeBool() (b bool, err protocol.Error) {
	d.TrimSpaces()
	if d.buf[0] == 't' {
		b = true
		d.Offset(5) // true,
	} else {
		// b = false
		d.Offset(6) // false,
	}
	return
}

// DecodeUInt8 convert 8bit integer number string to number. pass d.buf start from number and receive from after ,
func (d *Decoder) DecodeUInt8() (ui uint8, err protocol.Error) {
	d.TrimToDigit()
	d.FindEndToken()
	ui, err = convert.StringToUint8Base10(convert.UnsafeByteSliceToString(d.lastItem))
	if err != nil {
		err = &ErrEncodedIntegerCorrupted
		return
	}
	return
}

// DecodeUInt16 convert 16bit integer number string to number. pass d.buf start from number and receive from after ,
func (d *Decoder) DecodeUInt16() (ui uint16, err protocol.Error) {
	d.TrimToDigit()
	d.FindEndToken()
	ui, err = convert.StringToUint16Base10(convert.UnsafeByteSliceToString(d.lastItem))
	if err != nil {
		err = &ErrEncodedIntegerCorrupted
		return
	}
	return
}

// DecodeUInt32 convert 32bit integer number string to number. pass d.buf start from number and receive from after ,
func (d *Decoder) DecodeUInt32() (ui uint32, err protocol.Error) {
	d.TrimToDigit()
	d.FindEndToken()
	ui, err = convert.StringToUint32Base10(convert.UnsafeByteSliceToString(d.lastItem))
	if err != nil {
		err = &ErrEncodedIntegerCorrupted
		return
	}
	return
}

// DecodeUInt64 convert 64bit integer number string to number. pass d.buf start from after : and receive from after ,
func (d *Decoder) DecodeUInt64() (ui uint64, err protocol.Error) {
	d.TrimToDigit()
	d.FindEndToken()
	ui, err = convert.StringToUint64Base10(convert.UnsafeByteSliceToString(d.lastItem))
	if err != nil {
		err = &ErrEncodedIntegerCorrupted
		return
	}
	return
}

// DecodeInt64 convert 64bit number string to number. pass d.buf start from number and receive from after ,
func (d *Decoder) DecodeInt64() (i int64, err protocol.Error) {
	d.TrimToDigit()
	d.FindEndToken()
	var goErr error
	i, goErr = strconv.ParseInt(convert.UnsafeByteSliceToString(d.lastItem), 10, 64)
	if goErr != nil {
		return 0, &ErrEncodedIntegerCorrupted
	}
	return
}

// DecodeFloat64AsNumber convert float64 number string to float64 number. pass d.buf start from after : and receive from ,
func (d *Decoder) DecodeFloat64AsNumber() (f float64, err protocol.Error) {
	d.TrimToDigit()
	d.FindEndToken()
	var goErr error
	f, goErr = strconv.ParseFloat(convert.UnsafeByteSliceToString(d.lastItem), 64)
	if goErr != nil {
		return 0, &ErrEncodedIntegerCorrupted
	}
	return
}

// DecodeString return string. pass d.buf start from after : and receive from from after "
func (d *Decoder) DecodeString() (s string, err protocol.Error) {
	if d.CheckNullValue() {
		return
	}

	d.TrimToStringStart()
	var loc = bytes.IndexByte(d.buf, '"')
	if loc < 0 {
		err = &ErrEncodedStringCorrupted
		return
	}

	var slice []byte = d.buf[:loc]

	d.Offset(loc + 1)
	s = string(slice)
	return
}

/*
	Array part
*/

// DecodeByteArrayAsBase64 convert base64 string to [n]byte
func (d *Decoder) DecodeByteArrayAsBase64(array []byte) (err protocol.Error) {
	if d.CheckNullValue() {
		return
	}

	d.TrimToStringStart()
	var loc = bytes.IndexByte(d.buf, '"')
	if loc < 0 {
		err = &ErrEncodedArrayCorrupted
		return
	}

	var goErr error
	_, goErr = base64.RawStdEncoding.Decode(array, d.buf[:loc])
	if goErr != nil {
		return &ErrEncodedArrayCorrupted
	}

	d.FindEndToken()
	return
}

// DecodeByteArrayAsNumber convert number array to [n]byte
func (d *Decoder) DecodeByteArrayAsNumber(array []byte) (err protocol.Error) {
	if d.CheckNullValue() {
		return
	}

	var value uint8
	for i := 0; i < len(array); i++ {
		value, err = d.DecodeUInt8()
		if err != nil {
			err = &ErrEncodedArrayCorrupted
			return
		}
		array[i] = value
		d.FindEndToken()
	}
	if d.token != ']' {
		err = &ErrEncodedArrayCorrupted
	}
	return
}

/*
	Slice as Number
*/

// DecodeByteSliceAsNumber convert number string slice to []byte. pass buf start from after [ and receive from after ]
func (d *Decoder) DecodeByteSliceAsNumber() (slice []byte, err protocol.Error) {
	slice = make([]byte, 0, 8) // TODO::: Is cap efficient enough?

	var num uint8
	for !d.CheckToken(']') {
		num, err = d.DecodeUInt8()
		if err != nil {
			err = &ErrEncodedSliceCorrupted
			return
		}
		slice = append(slice, num)

		d.FindEndToken()
	}
	return
}

// DecodeUInt16SliceAsNumber convert uint16 number string slice to []byte. pass buf start from after [ and receive from after ]
func (d *Decoder) DecodeUInt16SliceAsNumber() (slice []uint16, err protocol.Error) {
	slice = make([]uint16, 0, 8) // TODO::: Is cap efficient enough?

	var num uint16
	for !d.CheckToken(']') {
		num, err = d.DecodeUInt16()
		if err != nil {
			err = &ErrEncodedSliceCorrupted
			return
		}
		slice = append(slice, num)

		d.FindEndToken()
	}
	return
}

// DecodeUInt32SliceAsNumber convert uint32 number string slice to []byte. pass buf start from after [ and receive from after ]
func (d *Decoder) DecodeUInt32SliceAsNumber() (slice []uint32, err protocol.Error) {
	slice = make([]uint32, 0, 8) // TODO::: Is cap efficient enough?

	var num uint32
	for !d.CheckToken(']') {
		num, err = d.DecodeUInt32()
		if err != nil {
			err = &ErrEncodedSliceCorrupted
			return
		}
		slice = append(slice, num)

		d.FindEndToken()
	}
	return
}

// DecodeUInt64SliceAsNumber convert uint64 number string slice to []byte. pass buf start from after [ and receive from after ]
func (d *Decoder) DecodeUInt64SliceAsNumber() (slice []uint64, err protocol.Error) {
	slice = make([]uint64, 0, 8) // TODO::: Is cap efficient enough?

	var num uint64
	for !d.CheckToken(']') {
		num, err = d.DecodeUInt64()
		if err != nil {
			err = &ErrEncodedSliceCorrupted
			return
		}
		slice = append(slice, num)

		d.FindEndToken()
	}
	return
}

/*
	Slice as Base64
*/

// DecodeByteSliceAsBase64 convert base64 string to []byte
func (d *Decoder) DecodeByteSliceAsBase64() (slice []byte, err protocol.Error) {
	d.TrimToStringStart()
	var loc = bytes.IndexByte(d.buf, '"')
	if loc < 0 {
		err = &ErrEncodedSliceCorrupted
		return
	}

	slice = make([]byte, base64.RawStdEncoding.DecodedLen(len(d.buf[:loc])))
	var n int
	var goErr error
	n, goErr = base64.RawStdEncoding.Decode(slice, d.buf[:loc])
	if goErr != nil {
		return slice, &ErrEncodedSliceCorrupted
	}
	slice = slice[:n]

	d.FindEndToken()
	return
}

// Decode32ByteArraySliceAsBase64 decode [32]byte base64 string slice. pass buf start from after [ and receive from after ]
func (d *Decoder) Decode32ByteArraySliceAsBase64() (slice [][32]byte, err protocol.Error) {
	const base64Len = 43 // base64.RawStdEncoding.EncodedLen(len(32))	>>	(32*8 + 5) / 6
	slice = make([][32]byte, 0, 8)

	var openBracketLoc = bytes.IndexByte(d.buf, '[')
	if openBracketLoc < 0 {
		err = &ErrEncodedSliceCorrupted
		return
	}
	d.buf = d.buf[openBracketLoc+1:]
	var goErr error
	var array [32]byte
	for !d.CheckToken(']') {
		d.TrimToStringStart()
		_, goErr = base64.RawStdEncoding.Decode(array[:], d.buf[:base64Len])
		if goErr != nil {
			err = &ErrEncodedSliceCorrupted
			return
		}
		slice = append(slice, array)
		d.buf = d.buf[base64Len:]
		d.FindEndToken()
	}
	return
}
