package crush

func lnScore(child Node, weight float32, input uint32, round uint32) int64 {
	var draw int64

	if weight > 0 {
		hash := hash3(input, Btoi(digestString(child.GetID())), round)
		hash = hash & 0xFFFF
		ln := int64(crushLn(hash) - 0x1000000000000)
		draw = int64(float32(ln) / weight)
	} else {
		draw = S64_MIN
	}
	return draw
}

const (
	U8_MAX  uint8  = 255
	S8_MAX  int8   = 127
	S8_MIN  int8   = (-S8_MAX - 1)
	U16_MAX uint16 = 65535
	S16_MAX int16  = 32767
	S16_MIN int16  = (-S16_MAX - 1)
	U32_MAX uint32 = 4294967295
	S32_MAX int32  = 2147483647
	S32_MIN int32  = (-S32_MAX - 1)
	U64_MAX uint64 = 18446744073709551615
	S64_MAX int64  = 9223372036854775807
	S64_MIN int64  = (-S64_MAX - 1)
)
