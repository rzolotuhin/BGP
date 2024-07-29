package tools

import (
	"errors"
	"fmt"
	"math/big"
	"net/netip"
)

func IpRangeToCIDR(start, end string) ([]string, error) {
	ips, err := netip.ParseAddr(start)
	if err != nil {
		return nil, err
	}
	ipe, err := netip.ParseAddr(end)
	if err != nil {
		return nil, err
	}

	isV4 := ips.Is4()
	if isV4 != ipe.Is4() {
		return nil, errors.New("start and end types are different")
	}
	if ips.Compare(ipe) > 0 {
		return nil, errors.New("start > end")
	}

	var (
		ipsInt = new(big.Int).SetBytes(ips.AsSlice())
		ipeInt = new(big.Int).SetBytes(ipe.AsSlice())
		nextIp = new(big.Int)
		maxBit = new(big.Int)
		cmpSh  = new(big.Int)
		bits   = new(big.Int)
		mask   = new(big.Int)
		one    = big.NewInt(1)
		buf    []byte
		cidr   []string
		bitSh  uint
	)
	if isV4 {
		maxBit.SetUint64(32)
		buf = make([]byte, 4)
	} else {
		maxBit.SetUint64(128)
		buf = make([]byte, 16)
	}

	for {
		bits.SetUint64(1)
		mask.SetUint64(1)
		for bits.Cmp(maxBit) < 0 {
			nextIp.Or(ipsInt, mask)

			bitSh = uint(bits.Uint64())
			cmpSh.Lsh(cmpSh.Rsh(ipsInt, bitSh), bitSh)
			if (nextIp.Cmp(ipeInt) > 0) || (cmpSh.Cmp(ipsInt) != 0) {
				bits.Sub(bits, one)
				mask.Rsh(mask, 1)
				break
			}
			bits.Add(bits, one)
			mask.Add(mask.Lsh(mask, 1), one)
		}

		addr, _ := netip.AddrFromSlice(ipsInt.FillBytes(buf))
		cidr = append(cidr, fmt.Sprintf("%s/%s",
			addr.String(),
			bits.Sub(maxBit, bits).String(),
		))

		if nextIp.Or(ipsInt, mask); nextIp.Cmp(ipeInt) >= 0 {
			break
		}
		ipsInt.Add(nextIp, one)
	}
	return cidr, nil
}
