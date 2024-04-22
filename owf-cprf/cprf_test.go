package owfcprf

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	mrand "math/rand"
	"testing"
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
	sec := 128
	bound := 5
	length := 5
	pp, msk, _ := KeyGen(sec, length, bound)

	// compute x and z such that <z,x> = 0
	z, _ := generateRandomVector(length, big.NewInt(int64(bound)))
	x := make([]*big.Int, length)
	for i := 0; i < length; i++ {
		x[i] = big.NewInt(0)
	}

	for i := 0; i < length; i++ {
		if generateRandomBit() == 0 {
			x[i], _ = generateRandomBigInt(big.NewInt(int64(bound)))
			z[i] = big.NewInt(0)
		}
	}

	csk, _ := msk.Constrain(z)

	eval := msk.Eval(pp, x)
	ceval := csk.CEval(pp, x)

	if eval.Cmp(ceval) != 0 {
		t.Fatalf("Eval and CEval are not equal")
	}
}

func TestCPRFUnauthorized(t *testing.T) {
	sec := 128
	bound := 5
	length := 5
	pp, msk, _ := KeyGen(sec, length, bound)

	z, _ := generateRandomVector(length, big.NewInt(int64(bound)))
	csk, _ := msk.Constrain(z)

	x, _ := generateRandomVector(length, big.NewInt(int64(bound)))
	eval := msk.Eval(pp, x)
	ceval := csk.CEval(pp, x)

	// very small probability of failure in this test case
	if eval.Cmp(ceval) == 0 {
		t.Fatalf("Eval and CEval are equal")
	}
}

func BenchmarkEval(b *testing.B) {
	sec := 40

	// Run the benchmark for different parameter sets
	for _, params := range []struct{ length, bound int }{
		{5, 2},
		{10, 2},
	} {
		b.Run(fmt.Sprintf("length=%d,bound=%d", params.length, params.bound), func(b *testing.B) {

			pp, msk, _ := KeyGen(sec, params.length, params.bound)
			x, _ := generateRandomVector(params.length, big.NewInt(int64(params.bound)))

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				msk.Eval(pp, x)
			}
		})
	}
}

func TestPublicParamSize(t *testing.T) {
	sec := int64(40)

	// Run the benchmark for different parameter sets
	for _, params := range []struct{ length, bound int64 }{
		{5, 2},
		{10, 2},
		{15, 2},
		{5, 4},
		{10, 4},
		{15, 4},
		{5, 8},
		{10, 8},
		{15, 8},
	} {
		t.Run(fmt.Sprintf("length=%d,bound=%d", params.length, params.bound), func(t *testing.T) {

			tt := int64(math.Pow(float64(params.bound), float64(params.length)))

			// parameters from Lemma 3 of the paper, computed as a function of t
			m := sec * (3*tt + 5) * (tt + 1)
			modbits := sec * (2*tt + 6)

			size := m * (modbits / 8) / int64(math.Pow10(6))

			t.Fatalf("Public parameter size (MB): %d", size)
		})
	}

}
