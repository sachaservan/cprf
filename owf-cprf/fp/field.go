package fp

import (
	"crypto/rand"
	"errors"
	"math/big"
)

type Field struct {
	P *big.Int // field modulus
}

type FieldElement struct {
	Int *big.Int
}

// new field of order p ~ 2^bits
// such that p is as close as possible to 2^bits
func NewFieldWithBitlen(bits int) (*Field, error) {

	if bits < 2 {
		return nil, errors.New("bit length must be at least 2")
	}

	one := big.NewInt(1)
	two := big.NewInt(2)
	p := big.NewInt(1)
	p.Lsh(p, uint(bits))
	p.Sub(p, one)
	for !p.ProbablyPrime(4) { // TODO: should be higher confidence IRL
		p.Sub(p, two)
	}

	if p.BitLen() < int(bits) {
		return nil, errors.New("could not find a prime of suitable bit length")
	}

	return &Field{P: p}, nil
}

// new field of order p
func NewField(p *big.Int) *Field {

	return &Field{P: p}
}

// add modulo P
func (f *Field) Add(a, b *FieldElement) *FieldElement {
	newValue := new(big.Int).Mod(new(big.Int).Add(a.Int, b.Int), f.P)
	return f.NewElement(newValue)
}

func (f *Field) AddInplace(a *FieldElement, b *FieldElement) {
	a.Int.Add(a.Int, b.Int).Mod(a.Int, f.P)
}

func (f *Field) AddInplaceStream(a *FieldElement, b *FieldElement) {
	a.Int.Add(a.Int, b.Int)
}

func (f *Field) AddInplaceStreamFinish(a *FieldElement) {
	a.Int.Mod(a.Int, f.P)
}

// sub modulo P
func (f *Field) Sub(a, b *FieldElement) *FieldElement {
	newValue := new(big.Int).Mod(new(big.Int).Sub(a.Int, b.Int), f.P)
	return f.NewElement(newValue)
}

func (f *Field) SubInplace(a *FieldElement, b *FieldElement) {
	a.Int.Sub(a.Int, b.Int).Mod(a.Int, f.P)
}

func (f *Field) Negate(a *FieldElement) *FieldElement {
	newValue := new(big.Int).Mod(new(big.Int).Sub(f.P, a.Int), f.P)
	return f.NewElement(newValue)
}

// return multiplicative inverse with mod P
func (f *Field) MulInv(a *FieldElement) *FieldElement {
	newValue := new(big.Int).ModInverse(a.Int, f.P)
	return f.NewElement(newValue)
}

// multiply mod P
func (f *Field) Mul(a, b *FieldElement) *FieldElement {
	newValue := new(big.Int).Mul(a.Int, b.Int)
	return f.NewElement(newValue)
}

func (f *Field) MulInplace(a *FieldElement, b *FieldElement) {
	a.Int.Mul(a.Int, b.Int).Mod(a.Int, f.P)
}

// exponentiation mod P
func (f *Field) Exp(a *FieldElement, c *big.Int) *FieldElement {
	newValue := exp(a.Int, c, f.P)
	return f.NewElement(newValue)
}

func (f *Field) ExpInplace(a *FieldElement, c *big.Int) {
	expInplace(a.Int, c, f.P)
}

// new element (a mod P)
func (f *Field) NewElement(a *big.Int) *FieldElement {
	newValue := new(big.Int).Mod(a, f.P)
	return &FieldElement{newValue}
}

// returns a random element in the field
func (f *Field) RandomElement() (*FieldElement, error) {
	a, err := randomInt(f.P)
	if err != nil {
		return nil, err
	}
	return f.NewElement(a), nil
}

func (f *Field) AddIdentity() *FieldElement {
	return &FieldElement{big.NewInt(0)}
}

func (f *Field) MulIdentity() *FieldElement {
	return &FieldElement{big.NewInt(1)}
}

func (f *Field) IsAddIdentity(e *FieldElement) bool {
	return f.AddIdentity().Cmp(e) == 0
}

func (f *Field) IsMulIdentity(e *FieldElement) bool {
	return f.MulIdentity().Cmp(e) == 0
}

func (f *Field) IsZero(e *FieldElement) bool {
	return f.AddIdentity().Cmp(e) == 0
}

func (elem *FieldElement) Cmp(b *FieldElement) int {
	return elem.Int.Cmp(b.Int)
}

func randomInt(max *big.Int) (*big.Int, error) {
	randomBig, err := rand.Int(rand.Reader, new(big.Int).SetBytes(max.Bytes()))
	return new(big.Int).SetBytes(randomBig.Bytes()), err
}

func exp(a, b, n *big.Int) *big.Int {
	return new(big.Int).Exp(a, b, n)
}

func expInplace(a, b, n *big.Int) {
	new(big.Int).Exp(a, b, n)
}
