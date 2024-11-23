package rocprf

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// Master key for the CPRF
// length: length of the inner product
// modulus: inner product modulus
// z0: master key
type MasterKey struct {
	length  int
	modulus *big.Int
	z0      []*big.Int
}

// Constrained key for the CPRF
// length: length of the inner product
// modulus: inner product modulus
// z1: constrained key
type ConstrainedKey struct {
	length  int
	modulus *big.Int
	z1      []*big.Int
}

// KeyGen generates a new CPRF key
// modulus: inner product modulus
// length: length of the input vector
// Outputs a CPRF master key
func KeyGen(modulus *big.Int, length int) (*MasterKey, error) {

	msk := &MasterKey{}
	msk.modulus = modulus
	msk.length = length
	msk.z0 = make([]*big.Int, length)

	var err error

	for i := 0; i < length; i++ {
		msk.z0[i], err = generateRandomBigInt(modulus)
		if err != nil {
			return nil, fmt.Errorf("failed to generate master key component %d: %w", i, err)
		}
	}

	return msk, nil
}

// Constrain outputs a constrained key for the CPRF
func (msk *MasterKey) Constrain(z []*big.Int) (*ConstrainedKey, error) {

	length := msk.length
	modulus := msk.modulus

	csk := &ConstrainedKey{}
	csk.modulus = modulus
	csk.length = length
	csk.z1 = make([]*big.Int, length)

	delta, err := generateRandomBigInt(modulus)
	if err != nil {
		return nil, fmt.Errorf("failed to generate delta for constraint: %w", err)
	}

	// the constraint key is computed as z0 - z*Delta
	// for a random Delta
	for i := 0; i < length; i++ {
		csk.z1[i] = big.NewInt(0)
		csk.z1[i].Mul(delta, z[i])          // z*Delta
		csk.z1[i].Sub(msk.z0[i], csk.z1[i]) // z0 - z*Delta
		if err != nil {
			return nil, err
		}
	}

	return csk, nil
}

func (msk *MasterKey) Eval(x []*big.Int) []byte {
	modulus := msk.modulus
	length := msk.length
	return commonEval(modulus, length, msk.z0, x)
}

func (csk *ConstrainedKey) CEval(x []*big.Int) []byte {
	modulus := csk.modulus
	length := csk.length
	return commonEval(modulus, length, csk.z1, x)
}

func commonEval(
	modulus *big.Int,
	length int,
	zb []*big.Int,
	x []*big.Int) []byte {

	tmp := big.NewInt(0)
	k := big.NewInt(0) // inner product result
	for i := 0; i < length; i++ {
		tmp.Mul(zb[i], x[i])
		k.Add(k, tmp).Mod(k, modulus)
	}

	return hashSHA256(k, x)
}

// SHÁ256 as a collision-resistant hash function.
// k: PRF key
// x: input vector
func hashSHA256(k *big.Int, x []*big.Int) []byte {

	byteInput := make([]byte, 0)

	byteInput = append(byteInput, k.Bytes()...)

	for i := 0; i < len(x); i++ {
		byteInput = append(byteInput, x[i].Bytes()...)
	}

	hasher := sha256.New()
	hasher.Write(byteInput)
	hash := hasher.Sum(nil)

	return hash
}

func generateRandomBigInt(max *big.Int) (*big.Int, error) {
	randomInt, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random number: %w", err)
	}
	return randomInt, nil
}
