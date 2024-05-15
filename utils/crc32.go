package utils

import "hash/crc32"

func CRC32(strKey string) uint32 {
	table := crc32.MakeTable(crc32.IEEE)
	ret := crc32.Checksum([]byte(strKey), table)
	return ret
}

func CRC32Slice(strKey []string) uint32 {
	table := crc32.MakeTable(crc32.IEEE)
	ret := crc32.Checksum([]byte(Strcatslice(strKey)), table)
	return ret
}
