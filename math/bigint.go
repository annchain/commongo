// Copyright © 2019 Annchain Authors <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package math

import (
	"encoding/json"
	"fmt"
	"math/big"
)

// DO NOT USE MSGP FOR AUTO-GENERATING HERE.
// THE bigint_gen.go has been modified intentionally to adapt big.Int

// FORBID: go:generate msgp
// FORBID: msgp:tuple BigInt

// A BigInt represents an unsigned multi-precision integer.
// The BigInt should always be unsigned since the sign of the number is not serialized.
type BigInt struct {
	Value *big.Int
}

// NewBigInt allocates and returns a new BigInt set to x.
func NewBigInt(x int64) *BigInt {
	return &BigInt{big.NewInt(x)}
}

// NewBigInt allocates and returns a new BigInt set to x.
func NewBigIntFromString(x string, base int) (*BigInt, bool) {
	v, success := big.NewInt(0).SetString(x, base)
	return &BigInt{v}, success
}

// NewBigInt allocates and returns a new BigInt set to x.
func NewBigIntFromBigInt(x *big.Int) *BigInt {
	return &BigInt{big.NewInt(0).SetBytes(x.Bytes())}
}

// GetBytes returns the absolute value of x as a big-endian byte slice.
func (bi *BigInt) GetBytes() []byte {
	return bi.Value.Bytes()
}

// String returns the value of x as a formatted decimal string.
func (bi *BigInt) String() string {
	return bi.Value.String()
}

// GetInt64 returns the int64 representation of x. If x cannot be represented in
// an int64, the result is undefined.
func (bi *BigInt) GetInt64() int64 {
	return bi.Value.Int64()
}

// SetInt64 sets the big int to x.
func (bi *BigInt) SetInt64(x int64) {
	bi.Value.SetInt64(x)
}

// GetSigBytes returns the bytes of bigint. This bytes are for signatures only.
func (bi *BigInt) GetSigBytes() []byte {
	if bi.GetInt64() == 0 {
		return []byte{0}
	}
	return bi.GetBytes()
}

// Sign returns:
//
//	-1 if x <  0
//	 0 if x == 0
//	+1 if x >  0
//
func (bi *BigInt) Sign() int {
	return bi.Value.Sign()
}

// Set sets bi to x and returns bi.
func (bi *BigInt) Set(x *BigInt) *BigInt {
	bi.Value.Set(x.Value)
	return bi
}

// Add sets bi to the sum (bi + increment) and returns bi.
func (bi *BigInt) Add(increment *BigInt) *BigInt {
	return NewBigIntFromBigInt(new(big.Int).Add(bi.Value, increment.Value))
}

// Sub sets bi to the difference (bi - decrement) and returns bi.
func (bi *BigInt) Sub(decrement *BigInt) *BigInt {
	return NewBigIntFromBigInt(new(big.Int).Sub(bi.Value, decrement.Value))
}

// SetString sets the big int to x.
//
// The string prefix determines the actual conversion base. A prefix of "0x" or
// "0X" selects base 16; the "0" prefix selects base 8, and a "0b" or "0B" prefix
// selects base 2. Otherwise the selected base is 10.
func (bi *BigInt) SetString(x string, base int) {
	if bi.Value == nil {
		bi.Value = big.NewInt(0)
	}
	bi.Value.SetString(x, base)
}

// GetString returns the value of x as a formatted string in some number base.
func (bi *BigInt) GetString(base int) string {
	return bi.Value.Text(base)
}

func (bi *BigInt) MarshalJSON() ([]byte, error) {
	res := fmt.Sprintf("%d", bi.Value)
	//fmt.Println("Marshaling into ", res)
	return json.Marshal(res)
}

func (bi *BigInt) UnmarshalJSON(b []byte) error {
	var val string
	err := json.Unmarshal(b, &val)
	if err != nil {
		panic(err)
	}

	bi.SetString(val, 10)
	return nil
}