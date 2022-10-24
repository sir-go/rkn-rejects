package main

import (
	"fmt"
)

type TLSPayload struct {
	len int
	raw []byte
	pos int
}

func (tp *TLSPayload) GetLenW() (int, error) {
	if tp.len < tp.pos+2 {
		return 0, fmt.Errorf("too small payload len (%d)", tp.len)
	}
	b1 := int(tp.raw[tp.pos])
	b2 := int(tp.raw[tp.pos+1])
	tp.pos += 2
	return (b1 << 8) + b2, nil
}

func (tp *TLSPayload) GetLenB() (int, error) {
	if tp.len < tp.pos+1 {
		return 0, fmt.Errorf("too small payload len (%d)", tp.len)
	}
	tp.pos++
	return int(tp.raw[tp.pos-1]), nil
}

func (tp *TLSPayload) Skip(n int) error {
	if tp.len < tp.pos+n {
		return fmt.Errorf("too small payload len (%d)", tp.len)
	}
	tp.pos += n
	return nil
}

func (tp *TLSPayload) GetString(n int) (res string, err error) {
	if tp.len < tp.pos+n {
		return "", fmt.Errorf("too small payload len (%d)", tp.len)
	}
	res = string(tp.raw[tp.pos : tp.pos+n])
	tp.pos += n
	return
}

func GetSNIForced(d []byte) (sni string, err error) {
	// SessionIdLength offset = 43
	pl := TLSPayload{len: len(d), raw: d, pos: 43}

	var _t int
	// sesIDLen
	if _t, err = pl.GetLenB(); err != nil {
		return
	}
	if err = pl.Skip(_t); err != nil {
		return
	}

	// cipherSuitsLen
	if _t, err = pl.GetLenW(); err != nil {
		return
	}
	if err = pl.Skip(_t); err != nil {
		return
	}

	// compressionMethodLen
	if _t, err = pl.GetLenB(); err != nil {
		return
	}
	if err = pl.Skip(_t); err != nil {
		return
	}

	// extensionsLen
	if _t, err = pl.GetLenW(); err != nil {
		return
	}

	var (
		extType int
		extLen  int
		SNILen  int
	)

	for pl.pos+1 < pl.len {
		if extType, err = pl.GetLenW(); err != nil {
			return
		}
		if extLen, err = pl.GetLenW(); err != nil {
			return
		}

		// SN
		if extType != 0x00 {
			if err = pl.Skip(extLen); err != nil {
				return
			}
			continue
		}
		if err = pl.Skip(3); err != nil {
			return
		}
		if SNILen, err = pl.GetLenW(); err != nil {
			return
		}

		return pl.GetString(SNILen)
	}
	return
}
