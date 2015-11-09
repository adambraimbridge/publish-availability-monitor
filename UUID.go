package main

import (
	"crypto/md5"
	"errors"
	"strconv"
	"strings"
)

type UUID struct {
	msb uint64
	lsb uint64
}

func NewNameUUIDFromBytes(bytes []byte) *UUID {
	md5Hash := md5.Sum(bytes)
	md5Hash[6] &= 0x0f /* clear version        */
	md5Hash[6] |= 0x30 /* set to version 3     */
	md5Hash[8] &= 0x3f /* clear variant        */
	md5Hash[8] |= 0x80 /* set to IETF variant  */

	var msb uint64
	var lsb uint64

	for i := 0; i < 8; i++ {
		msb = (msb << 8) | (uint64(md5Hash[i]) & 0xff)
	}
	for i := 8; i < 16; i++ {
		lsb = (lsb << 8) | (uint64(md5Hash[i]) & 0xff)
	}

	return &UUID{msb, lsb}
}

func NewUUIDFromString(uuidString string) (*UUID, error) {
	components := strings.Split(uuidString, "-")
	if len(components) != 5 {
		return &UUID{0, 0}, errors.New("Invalid UUID string")
	}

	msb := hexToInt(components[0])
	msb <<= 16
	msb |= hexToInt(components[1])
	msb <<= 16
	msb |= hexToInt(components[2])

	lsb := hexToInt(components[3])
	lsb <<= 48
	lsb |= hexToInt(components[4])

	return &UUID{msb, lsb}, nil
}

func (uuid *UUID) String() string {
	parts := make([]string, 5)
	parts[0] = digits(uuid.msb>>32, 8)
	parts[1] = digits(uuid.msb>>16, 4)
	parts[2] = digits(uuid.msb, 4)
	parts[3] = digits(uuid.lsb>>48, 4)
	parts[4] = digits(uuid.lsb, 12)

	return strings.Join(parts, "-")
}

func digits(val uint64, digits uint) string {
	hi := 1 << (digits * 4)
	result := uint64(hi) | (val & (uint64(hi) - uint64(1)))
	return strconv.FormatInt(int64(result), 16)[1:]
}

func hexToInt(hexString string) uint64 {
	result, err := strconv.ParseUint(hexString, 16, 0)
	if err != nil {
		//TODO how and where to handle this...
	}
	return result
}
