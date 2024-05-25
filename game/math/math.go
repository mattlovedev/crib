package math

import "math/big"

const NCR52_4 = 270725

func NChooseR(n int, r int) int {
	z := &big.Int{}
	z.Binomial(int64(n), int64(r))
	return int(z.Int64())
}
