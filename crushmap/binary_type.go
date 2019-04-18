package crushmap

import (
	"encoding/binary"
)

func (p *binaryParser) handleType() ([]*Type, error) {
	var err error
	var n uint32
	var itypes []*Type

	err = binary.Read(p.r, binary.LittleEndian, &n)
	if err != nil {
		return nil, err
	}

	for i := n; i > 0; i-- {
		var key int32
		err = binary.Read(p.r, binary.LittleEndian, &key)
		if err != nil {
			return nil, err
		}

		var l uint32
		err = binary.Read(p.r, binary.LittleEndian, &l)
		if err != nil {
			return nil, err
		}
		if l == 0 {
			err = binary.Read(p.r, binary.LittleEndian, &l)
			if err != nil {
				return nil, err
			}
		}

		val := make([]byte, l)
		err = binary.Read(p.r, binary.LittleEndian, &val)
		if err != nil {
			return nil, err
		}
		itypes = append(itypes, &Type{ID: key, Name: string(val)})
	}

	return itypes, nil
}
