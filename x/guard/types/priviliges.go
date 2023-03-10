package types

import (
	"math/big"
)

type Privileges interface {
	Check(query []byte) bool
	Equal(pr []byte) bool
	Мerge(pr []byte) *AccPrivileges
	МergeMore(prs [][]byte) *AccPrivileges
	SwitchOn(ids []*big.Int) *AccPrivileges
	SwitchOff(ids []*big.Int, base *big.Int) *AccPrivileges
	SetBytes(buf []byte) *AccPrivileges
	Bytes() []byte
	Empty() bool
}

var (
	_ Privileges = AccPrivileges{}
)

type AccPrivileges struct {
	num *big.Int
	raw []byte
}

func NewEmptyAccPrivileges() *AccPrivileges {
	pr := big.NewInt(0)
	return &AccPrivileges{
		raw: pr.Bytes(),
		num: pr,
	}
}
func NewAccPrivileges(base *big.Int) *AccPrivileges {
	return &AccPrivileges{
		raw: base.Bytes(),
		num: base,
	}
}

func AccPrivilegesFromBytes(bz []byte) *AccPrivileges {
	if bz == nil {
		return NewEmptyAccPrivileges()
	}
	return NewEmptyAccPrivileges().SetBytes(bz)
}

func (ap AccPrivileges) SetBytes(buf []byte) *AccPrivileges {
	ap.BigInt().SetBytes(buf)
	ap.raw = ap.BigInt().Bytes()
	return &ap
}

func (ap AccPrivileges) Empty() bool {
	return len(ap.Bytes()) == 0
}

func (ap AccPrivileges) BigInt() *big.Int {
	return ap.num
}

func (ap AccPrivileges) Bytes() []byte {
	return ap.raw
}

func (ap AccPrivileges) Check(query []byte) bool {
	queryNum := big.NewInt(0).SetBytes(query)
	return big.NewInt(0).And(ap.BigInt(), queryNum).Cmp(queryNum) == 0
}

func (ap AccPrivileges) Equal(pr []byte) bool {
	return ap.BigInt().Cmp(big.NewInt(0).SetBytes(pr)) == 0
}

func (ap AccPrivileges) Мerge(pr []byte) *AccPrivileges {
	ap.BigInt().Or(ap.BigInt(), big.NewInt(0).SetBytes(pr))
	ap.raw = ap.BigInt().Bytes()
	return &ap
}

func (ap AccPrivileges) МergeMore(prs [][]byte) *AccPrivileges {
	for _, pr := range prs {
		ap.BigInt().Or(ap.BigInt(), big.NewInt(0).SetBytes(pr))
	}
	ap.raw = ap.BigInt().Bytes()
	return &ap
}

func (ap AccPrivileges) SwitchOn(ids []*big.Int) *AccPrivileges {
	for _, id := range ids {
		ap.BigInt().Or(ap.BigInt(), big.NewInt(0).Exp(big.NewInt(2), id, nil))
	}
	ap.raw = ap.BigInt().Bytes()
	return &ap
}

func (ap AccPrivileges) SwitchOff(ids []*big.Int, base *big.Int) *AccPrivileges {
	for _, id := range ids {
		ap.BigInt().AndNot(ap.BigInt(), big.NewInt(0).Exp(big.NewInt(2), id, nil))
	}
	ap.raw = ap.BigInt().Bytes()
	return &ap
}
