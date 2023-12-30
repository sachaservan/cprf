package ddhcprf

import (
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	mrand "math/rand"
	"testing"

	"github.com/sachaservan/cprf/ddh-cprf/ec"
)

func generateRandomBit() int {
	return mrand.Intn(2)
}

func generateRandomVector(length int, max *big.Int) ([]*big.Int, error) {
	res := make([]*big.Int, length)
	for i := 0; i < length; i++ {
		randomInt, err := rand.Int(rand.Reader, max)
		if err != nil {
			return nil, err
		}
		res[i] = randomInt
	}

	return res, nil
}

func TestCPRFAuthorized(t *testing.T) {
	p := elliptic.P256().Params().N
	n := 128
	length := 25

	// compute x and z such that <z,x> = 0
	z, _ := generateRandomVector(length, p)
	x := make([]*big.Int, length)
	for i := 0; i < length; i++ {
		x[i] = big.NewInt(0)
	}

	for i := 0; i < length; i++ {
		if generateRandomBit() == 0 {
			x[i], _ = generateRandomBigInt(p)
			z[i] = big.NewInt(0)
		}
	}

	pp, msk, _ := KeyGen(n, length)
	csk, _ := msk.Constrain(z)

	eval := msk.Eval(pp, x)
	ceval := csk.CEval(pp, x)

	if !ec.PointsEqual(eval, ceval) {
		t.Fatalf("Eval and CEval are not equal")
	}
}

func TestCPRFUnauthorized(t *testing.T) {
	p := elliptic.P256().Params().N
	n := 128
	length := 25
	z, _ := generateRandomVector(length, p)
	pp, msk, _ := KeyGen(n, length)
	csk, _ := msk.Constrain(z)

	x, _ := generateRandomVector(length, p)
	eval := msk.Eval(pp, x)
	ceval := csk.CEval(pp, x)

	// very small probability of failure in this test case
	if ec.PointsEqual(eval, ceval) {
		t.Fatalf("Eval and CEval are equal")
	}
}

func BenchmarkEval(b *testing.B) {
	p := elliptic.P256().Params().N
	n := 128
	length := 700
	pp, msk, _ := KeyGen(n, length)
	x, _ := generateRandomVector(length, p)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		msk.Eval(pp, x)
	}
}

func BenchmarkExp(b *testing.B) {
	_, x, y, _ := elliptic.GenerateKey(elliptic.P256(), rand.Reader)
	point, _ := ec.NewPoint(elliptic.P256(), x, y)

	_, scalar, _ := ec.RandomCurveScalar(elliptic.P256(), rand.Reader)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		point = ec.PointScalarMult(elliptic.P256(), point, scalar)
	}
}

func BenchmarkMulP(b *testing.B) {
	p := elliptic.P256().Params().N

	rand1, _ := generateRandomBigInt(p)
	rand2, _ := generateRandomBigInt(p)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rand1.Mul(rand1, rand2).Mod(rand1, p)
	}
}

func BenchmarkAddP(b *testing.B) {
	p := elliptic.P256().Params().N

	rand1, _ := generateRandomBigInt(p)
	rand2, _ := generateRandomBigInt(p)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rand1.Add(rand1, rand2).Mod(rand1, p)
	}
}
