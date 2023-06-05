package types

import (
	"math/big"
)

type IPrivileges interface {
	Check(query []byte) bool
	Equal(pr []byte) bool
	Мerge(pr []byte) *Privileges
	МergeMore(prs [][]byte) *Privileges
	SwitchOn(ids []*big.Int) *Privileges
	SwitchOff(ids []*big.Int, base *big.Int) *Privileges
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

func (ap Privileges) SetBytes(buf []byte) *Privileges {
	ap.BigInt().SetBytes(buf)
	ap.raw = ap.BigInt().Bytes()
	return &ap
}

func (ap Privileges) Empty() bool {
	return len(ap.Bytes()) == 0
}

func (ap Privileges) BigInt() *big.Int {
	return ap.num
}

func (ap Privileges) Bytes() []byte {
	return ap.raw
}

func (ap Privileges) CheckPrivileges(privileges *Privileges, defPriv []byte) bool {
	defPrNum := big.NewInt(0).SetBytes(defPriv)

	defReqPr := big.NewInt(0).And(defPrNum, privileges.BigInt())
	defAccPr := big.NewInt(0).And(defPrNum, ap.BigInt())

	if big.NewInt(0).And(defAccPr, defReqPr).Cmp(defReqPr) == 0 {
		invDefPrNum := big.NewInt(0).Not(defPrNum)

		reqPr := big.NewInt(0).And(invDefPrNum, privileges.BigInt())
		accPr := big.NewInt(0).And(invDefPrNum, ap.BigInt())

		return big.NewInt(0).And(accPr, reqPr).Cmp(big.NewInt(0)) == 1
	}

	return false
}

func (ap Privileges) Check(query []byte) bool {
	queryNum := big.NewInt(0).SetBytes(query)
	return big.NewInt(0).And(ap.BigInt(), queryNum).Cmp(queryNum) == 0
}

func (ap Privileges) Equal(pr []byte) bool {
	return ap.BigInt().Cmp(big.NewInt(0).SetBytes(pr)) == 0
}

func (ap Privileges) Мerge(pr []byte) *Privileges {
	ap.BigInt().Or(ap.BigInt(), big.NewInt(0).SetBytes(pr))
	ap.raw = ap.BigInt().Bytes()
	return &ap
}

func (ap Privileges) МergeMore(prs [][]byte) *Privileges {
	for _, pr := range prs {
		ap.BigInt().Or(ap.BigInt(), big.NewInt(0).SetBytes(pr))
	}
	ap.raw = ap.BigInt().Bytes()
	return &ap
}

func (ap Privileges) SwitchOn(ids []*big.Int) *Privileges {
	for _, id := range ids {
		ap.BigInt().Or(ap.BigInt(), big.NewInt(0).Exp(big.NewInt(2), id, nil))
	}
	ap.raw = ap.BigInt().Bytes()
	return &ap
}

func (ap Privileges) SwitchOff(ids []*big.Int, base *big.Int) *Privileges {
	for _, id := range ids {
		ap.BigInt().AndNot(ap.BigInt(), big.NewInt(0).Exp(big.NewInt(2), id, nil))
	}
	ap.raw = ap.BigInt().Bytes()
	return &ap
}
