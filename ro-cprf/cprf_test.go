package rocprf

import (
	"bytes"
	"crypto/rand"
	"fmt"
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

	modulus, _ := big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFF61", 16)

	length := 10
	msk, _ := KeyGen(modulus, length)

	// compute x and z such that <z,x> = 0
	z, _ := generateRandomVector(length, modulus)
	x := make([]*big.Int, length)
	for i := 0; i < length; i++ {
		x[i] = big.NewInt(0)
	}

	for i := 0; i < length; i++ {
		if generateRandomBit() == 0 {
			x[i], _ = generateRandomBigInt(modulus)
			z[i] = big.NewInt(0)
		}
	}

	csk, _ := msk.Constrain(z)

	eval := msk.Eval(x)
	ceval := csk.CEval(x)

	if bytes.Compare(eval, ceval) != 0 {
		t.Fatalf("Eval and CEval are not equal")
	}
}

func TestCPRFUnauthorized(t *testing.T) {

	modulus, _ := big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFF61", 16)

	length := 10
	msk, _ := KeyGen(modulus, length)

	z, _ := generateRandomVector(length, modulus)
	csk, _ := msk.Constrain(z)

	x, _ := generateRandomVector(length, modulus)
	eval := msk.Eval(x)
	ceval := csk.CEval(x)

	// very small probability of failure in this test case
	if bytes.Compare(eval, ceval) == 0 {
		t.Fatalf("Eval and CEval are equal")
	}

}

func BenchmarkEval(b *testing.B) {
	modulus, _ := big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFF61", 16)

	// Run the benchmark for different parameter sets
	for _, params := range []struct{ length int }{
		{10},
		{50},
		{100},
		{500},
		{1000},
	} {
		b.Run(fmt.Sprintf("length=%d", params.length), func(b *testing.B) {

			msk, _ := KeyGen(modulus, params.length)
			x, _ := generateRandomVector(params.length, modulus)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				msk.Eval(x)
			}
		})
	}

}
