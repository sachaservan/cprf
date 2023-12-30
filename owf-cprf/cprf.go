package owfcprf

import (
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github.com/sachaservan/cprf/owf-cprf/fp"
)

// Random polynomial in F_2^k for k-wise independent hashing
type PublicParameters struct {
	field *fp.Field
	poly  []*fp.FieldElement
}

// Master key for the CPRF
type MasterKey struct {
	modulus *big.Int // inner product modulus
	length  int
	n       int
	z0      [][]*big.Int
}

// Constrained key for the CPRF
type ConstrainedKey struct {
	modulus *big.Int // inner product modulus
	length  int
	n       int
	z1      [][]*big.Int
}

func KeyGen(sec int, length int, bound int) (*PublicParameters, *MasterKey, error) {

	t := length * bound // max number of inner products possible

	// parameters from Lemma 3 of the paper, computed as a function of t
	m := sec * (3*t + 5) * (t + 1)
	modbits := sec * (2*t + 6)

	// compute field modulus (fake)
	p := big.NewInt(1)
	p.Lsh(p, uint(modbits))
	p.Sub(p, big.NewInt(1))

	n := 1

	var err error

	pp := &PublicParameters{}
	pp.field = fp.NewField(p)

	pp.poly = make([]*fp.FieldElement, m)
	for i := 0; i < m; i++ {
		pp.poly[i], err = pp.field.RandomElement()
		if err != nil {
			return nil, nil, err
		}
	}

	// Step 2: Generate the master key
	msk := &MasterKey{}
	msk.length = length
	msk.modulus = pp.field.P
	msk.n = n
	msk.z0 = make([][]*big.Int, n)

	for i := 0; i < n; i++ {
		msk.z0[i] = make([]*big.Int, length)
		for j := 0; j < length; j++ {
			msk.z0[i][j], err = generateRandomBigInt(msk.modulus)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	return pp, msk, nil
}

func (msk *MasterKey) Constrain(z []*big.Int) (*ConstrainedKey, error) {
	n := msk.n
	length := msk.length

	csk := &ConstrainedKey{}
	csk.length = length
	csk.n = n
	csk.modulus = msk.modulus
	csk.z1 = make([][]*big.Int, n)

	for i := 0; i < n; i++ {

		csk.z1[i] = make([]*big.Int, length)

		deltai, err := generateRandomBigInt(msk.modulus)
		if err != nil {
			return nil, err
		}

		for j := 0; j < length; j++ {
			csk.z1[i][j] = big.NewInt(0)
			csk.z1[i][j].Mul(deltai, z[j])               // z*Delta_i
			csk.z1[i][j].Sub(msk.z0[i][j], csk.z1[i][j]) // z0 - z*Delta_i
			if err != nil {
				return nil, err
			}
		}
	}

	return csk, nil
}

func (msk *MasterKey) Eval(pp *PublicParameters, x []*big.Int) *big.Int {
	n := msk.n
	length := msk.length
	modulus := msk.modulus
	return commonEval(pp, n, length, modulus, msk.z0, x)
}

func (csk *ConstrainedKey) CEval(pp *PublicParameters, x []*big.Int) *big.Int {
	n := csk.n
	length := csk.length
	modulus := csk.modulus
	return commonEval(pp, n, length, modulus, csk.z1, x)
}

func commonEval(
	pp *PublicParameters,
	n int,
	length int,
	modulus *big.Int,
	zb [][]*big.Int,
	x []*big.Int) *big.Int {

	keys := make([]*fp.FieldElement, n)

	tmp := big.NewInt(0)
	for i := 0; i < n; i++ {
		// Step 1: compute the inner product to get the key
		key := big.NewInt(0)
		for j := 0; j < length; j++ {
			tmp.Mul(zb[i][j], x[j])
			key.Add(key, tmp)
		}

		key.Mod(key, modulus)
		keys[i] = pp.field.NewElement(key)
	}

	hashSum := pp.field.AddIdentity()
	var hash *fp.FieldElement
	for i := 0; i < n; i++ {
		// Step 2: Compute h(key) using Horner's method to evaluate the
		// random polynomial represented by pp.poly in the public parameters
		hash = pp.poly[0]
		for j := 1; j < len(pp.poly); j++ {
			pp.field.MulInplace(hash, keys[i])
			pp.field.AddInplace(hash, pp.poly[j])
		}

		// faster addition by avoiding modulo reduction
		pp.field.AddInplaceStream(hashSum, hash)
	}

	// delays the modulo reduction operation
	pp.field.AddInplaceStreamFinish(hashSum)

	// Step 3: Use h(key) as a key for any PRF
	// Here, we will just use SHA256 as our PRF and compute H(h(key) || x)
	input := make([]byte, 0)
	for i := 0; i < len(x); i++ {
		input = append(input, x[i].Bytes()...)
	}
	input = append(input, hashSum.Int.Bytes()...)

	shaHash := sha256.New()
	shaHash.Write(input)
	res := big.NewInt(0).SetBytes(shaHash.Sum(nil))

	return res
}

func generateRandomBigInt(max *big.Int) (*big.Int, error) {
	randomInt, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}

	return randomInt, nil
}
