package ddhcprf

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/sachaservan/cprf/ddh-cprf/ec"
)

// Public parameters consists of k random group elements
// and is used to compute a variant of the Damgard hash function
// based on the hardness of the discrete logarithm problem.
type PublicParameters struct {
	hashElements []*ec.Point // hashing group elements
}

// Master key for the CPRF
// length: length of the inner product
// n: number of elements in the Naor-Reingold PRF key
// z0: master key
type MasterKey struct {
	length int
	n      int
	z0     [][]*big.Int
}

// Constrained key for the CPRF
// length: length of the inner product
// n: number of elements in the Naor-Reingold PRF key
// z1: constrained key
type ConstrainedKey struct {
	length int
	n      int
	z1     [][]*big.Int
}

// KeyGen generates a new CPRF key
// n: number of elements in the Naor-Reingold PRF key
// length: length of the inner product
// Outputs public parameters and a master key
func KeyGen(n int, length int) (*PublicParameters, *MasterKey, error) {

	// p is the order of the eliptic curve
	p := elliptic.P256().Params().N

	msk := &MasterKey{}
	msk.n = n
	msk.length = length
	msk.z0 = make([][]*big.Int, n)

	var err error

	for i := 0; i < n; i++ {
		msk.z0[i] = make([]*big.Int, length)
		for j := 0; j < length; j++ {
			msk.z0[i][j], err = generateRandomBigInt(p)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to generate master key component (%d,%d): %w", i, j, err)
			}
		}
	}

	// hash elements for the public parameters
	// bound ensures there are enough elements
	bound := 2 * (length + n)
	hashElements := make([]*ec.Point, bound)
	for i := 0; i < bound; i++ {
		_, hashElements[i], _ = ec.NewRandomPoint()
	}

	pp := &PublicParameters{}
	pp.hashElements = hashElements

	return pp, msk, nil
}

// Constrain outputs a constrained key for the CPRF
func (msk *MasterKey) Constrain(z []*big.Int) (*ConstrainedKey, error) {

	// p is the order of the eliptic curve
	p := elliptic.P256().Params().N
	length := msk.length
	n := msk.n

	csk := &ConstrainedKey{}
	csk.n = n
	csk.length = length
	csk.z1 = make([][]*big.Int, n)

	// the constraint key is computed as z0 - z*Delta_i
	// for a random Delta_i with i = 1 ... n
	for i := 0; i < n; i++ {
		csk.z1[i] = make([]*big.Int, length)

		deltai, err := generateRandomBigInt(p)
		if err != nil {
			return nil, fmt.Errorf("failed to generate delta_%d for constraint: %w", i, err)
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

func (msk *MasterKey) Eval(pp *PublicParameters, x []*big.Int) *ec.Point {
	n := msk.n
	length := msk.length
	return commonEval(pp, n, length, msk.z0, x)
}

func (csk *ConstrainedKey) CEval(pp *PublicParameters, x []*big.Int) *ec.Point {
	n := csk.n
	length := csk.length
	return commonEval(pp, n, length, csk.z1, x)
}

func commonEval(
	pp *PublicParameters,
	n int,
	length int,
	zb [][]*big.Int,
	x []*big.Int) *ec.Point {

	curve := elliptic.P256()
	p := elliptic.P256().Params().N

	keys := make([]*big.Int, n)
	keyFPs := make([]*ec.Point, n) // key fingerprint curve points

	tmp := big.NewInt(0)
	for i := 0; i < n; i++ {
		acc := big.NewInt(0)

		for j := 0; j < length; j++ {
			tmp.Mul(zb[i][j], x[j])
			acc.Add(acc, tmp).Mod(acc, p)
		}

		keys[i] = acc
		keyFPs[i] = ec.BaseScalarMult(curve, acc)
	}

	bits := hashDL(pp, x, keyFPs)[:n] // hashes to n points

	// Alternative: use SHA256
	// bits := hashSHA256(x, keyFPs)[:n]

	prod := big.NewInt(1)

	// Recall: the input is always prefixed by 11
	prod.Mul(prod, keys[0]).Mod(prod, p)
	prod.Mul(prod, keys[1]).Mod(prod, p)

	// Compute a_i^{x_i}
	for i := 2; i < n; i++ {
		if bits[i] {
			prod.Mul(prod, keys[i]).Mod(prod, p)
		}
	}

	res := ec.BaseScalarMult(curve, prod)
	return res
}

// Variant of the Damgard group-based hash function.
// pp: public parameters of the DL hash
// x: input vector to the PRF
// keyFPs: curve points of the key fingerprint
func hashDL(
	pp *PublicParameters,
	x []*big.Int,
	keyFPs []*ec.Point) []bool {

	// convert everything into a byte array
	byteInput := make([]byte, 0)
	for i := 0; i < len(x); i++ {
		byteInput = append(byteInput, x[i].Bytes()...)
	}

	for i := 0; i < len(keyFPs); i++ {
		byteInput = append(byteInput, keyFPs[i].MarshalCompressed()...)
	}

	// chunk everything up into 256 bit chunks
	curve := elliptic.P256()
	blocklen := 256 / 8

	// pad out to the block length if needed
	paddingRequired := len(byteInput) % blocklen
	for i := 0; i < paddingRequired; i++ {
		byteInput = append(byteInput, byte(0))
	}

	numBlocks := len(byteInput) / blocklen

	// compute hash of the message as PROD h_i^b_i
	// where h_i is the i-th element in the public parameters
	// and b_i is the i-th block of bytes
	res := ec.BaseScalarMult(curve, big.NewInt(1))
	for i := 0; i < numBlocks; i++ {
		start := i * blocklen
		end := start + blocklen
		blockNext := big.NewInt(0).SetBytes(byteInput[start:end])
		a := ec.PointScalarMult(curve, pp.hashElements[i], blockNext)
		res = ec.PointAdd(curve, res, a)
	}

	hash := big.NewInt(0).SetBytes(res.MarshalCompressed()).Bytes()

	// Apply randomness extractor to the output
	// bits of the group representation to ensure uniform distribution.
	// Doesn't need to be sha256 but convenient and doesn't add much overhead.
	hasher := sha256.New()
	hasher.Write(hash)
	hash = hasher.Sum(nil)

	// Convert the hash to a bit-wise representation
	hashBits := make([]bool, 0, len(hash)*8)
	for _, b := range hash {
		for i := 7; i >= 0; i-- {
			bit := (b>>uint(i))&1 == 1
			hashBits = append(hashBits, bit)
		}
	}

	return hashBits
}

// SHÁ256 as a collision-resistant hash function.
// pp: public parameters of the DL hash
// x: input vector to the PRF
// keyFPs: curve points of the key fingerprint
func hashSHA256(x []*big.Int, keyFPs []*ec.Point) []bool {

	byteInput := make([]byte, 0)
	for i := 0; i < len(x); i++ {
		byteInput = append(byteInput, x[i].Bytes()...)
	}

	for i := 0; i < len(keyFPs); i++ {
		byteInput = append(byteInput, keyFPs[i].MarshalCompressed()...)
	}

	hasher := sha256.New()
	hasher.Write(byteInput)
	hash := hasher.Sum(nil)

	// Convert the hash to a bit-wise representation
	hashBits := make([]bool, 0, len(hash)*8)
	for _, b := range hash {
		for i := 7; i >= 0; i-- {
			bit := (b>>uint(i))&1 == 1
			hashBits = append(hashBits, bit)
		}
	}

	return hashBits
}

func generateRandomBigInt(max *big.Int) (*big.Int, error) {
	randomInt, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random number: %w", err)
	}
	return randomInt, nil
}
