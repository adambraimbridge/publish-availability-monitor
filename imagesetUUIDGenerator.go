package main

import (
	"strconv"
	"strings"

	"github.com/willf/bitset"
)

var magic = toBitSet(NewNameUUIDFromBytes([]byte("imageset")).lsb)

func GenerateImageSetUUID(imageUUID UUID) (UUID, error) {
	uuidBits := toBitSet(imageUUID.lsb)
	uuidBits.InPlaceSymmetricDifference(&magic) //XOR

	lsb, err := strconv.ParseInt(strings.TrimSuffix(uuidBits.DumpAsBits(), "."), 2, 64)
	return UUID{imageUUID.msb, uint64(lsb)}, err
}

func toBitSet(number uint64) bitset.BitSet {
	lsbString := strconv.FormatUint(number, 2)

	var b bitset.BitSet
	for i, character := range lsbString {
		if character == '1' {
			b.Set(uint(len(lsbString)-1) - uint(i))
		}
	}
	return b
}
