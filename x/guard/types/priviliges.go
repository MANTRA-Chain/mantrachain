package types

import (
	"math/big"
)

type IPrivileges interface {
	Check(query []byte) bool
	Equal(pr []byte) bool
	Merge(pr []byte) *Privileges
	MergeMore(prs [][]byte) *Privileges
	SwitchOn(ids []*big.Int) *Privileges
	SwitchOff(ids []*big.Int) *Privileges
	SetBytes(buf []byte) *Privileges
	Bytes() []byte
	Empty() bool
}

var (
	_ IPrivileges = Privileges{}
)

type Privileges struct {
	num *big.Int
	raw []byte
}

func NewEmptyPrivileges() *Privileges {
	pr := big.NewInt(0)
	return &Privileges{
		raw: pr.Bytes(),
		num: pr,
	}
}
func NewPrivileges(base *big.Int) *Privileges {
	return &Privileges{
		raw: base.Bytes(),
		num: base,
	}
}

func PrivilegesFromBytes(bz []byte) *Privileges {
	if bz == nil {
		return NewEmptyPrivileges()
	}
	return NewEmptyPrivileges().SetBytes(bz)
}

func (p Privileges) SetBytes(buf []byte) *Privileges {
	p.BigInt().SetBytes(buf)
	p.raw = p.BigInt().Bytes()
	return &p
}

func (p Privileges) Empty() bool {
	return len(p.Bytes()) == 0
}

func (p Privileges) BigInt() *big.Int {
	return p.num
}

func (p Privileges) Bytes() []byte {
	return p.raw
}

func (p Privileges) CheckPrivileges(privileges *Privileges, defPriv []byte) bool {
	defPrNum := big.NewInt(0).SetBytes(defPriv)

	defReqPr := big.NewInt(0).And(defPrNum, privileges.BigInt())
	defAccPr := big.NewInt(0).And(defPrNum, p.BigInt())

	if big.NewInt(0).And(defAccPr, defReqPr).Cmp(defReqPr) == 0 {
		invDefPrNum := big.NewInt(0).Not(defPrNum)

		reqPr := big.NewInt(0).And(invDefPrNum, privileges.BigInt())
		accPr := big.NewInt(0).And(invDefPrNum, p.BigInt())

		return big.NewInt(0).And(accPr, reqPr).Cmp(big.NewInt(0)) == 1
	}

	return false
}

func (p Privileges) Check(query []byte) bool {
	queryNum := big.NewInt(0).SetBytes(query)
	return big.NewInt(0).And(p.BigInt(), queryNum).Cmp(queryNum) == 0
}

func (p Privileges) Equal(pr []byte) bool {
	return p.BigInt().Cmp(big.NewInt(0).SetBytes(pr)) == 0
}

func (p Privileges) Merge(pr []byte) *Privileges {
	p.BigInt().Or(p.BigInt(), big.NewInt(0).SetBytes(pr))
	p.raw = p.BigInt().Bytes()
	return &p
}

func (p Privileges) MergeMore(prs [][]byte) *Privileges {
	for _, pr := range prs {
		p.BigInt().Or(p.BigInt(), big.NewInt(0).SetBytes(pr))
	}
	p.raw = p.BigInt().Bytes()
	return &p
}

func (p Privileges) SwitchOn(ids []*big.Int) *Privileges {
	for _, id := range ids {
		p.BigInt().Or(p.BigInt(), big.NewInt(0).Exp(big.NewInt(2), id, nil))
	}
	p.raw = p.BigInt().Bytes()
	return &p
}

func (p Privileges) SwitchOff(ids []*big.Int) *Privileges {
	for _, id := range ids {
		p.BigInt().AndNot(p.BigInt(), big.NewInt(0).Exp(big.NewInt(2), id, nil))
	}
	p.raw = p.BigInt().Bytes()
	return &p
}
