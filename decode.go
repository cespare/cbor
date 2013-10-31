package cbor

import (
	"runtime"
)

func Unmarshal(data []byte, v interface{}) error {
	// TODO: encoding/json checks for well-formedness before actually unmarshaling. Should we do that here?
	d := newDecodeState(data)
	return d.unmarshal(v)
}

type decodeState struct {
	data   []byte
	offset int // into data
}

func (d *decodeState) unmarshal(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	// TODO: WIP
}
