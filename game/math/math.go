package math

import "math/big"

const (
	NCR52_4 = 270725 // 270K - number of entries in scores/four summaries; each with 48 possible cuts

	NCR52_6 = 20358520 // 20M - number of entries in scores/six summaries; each with 15 possible 4s and 46 possible cuts

	NCR6_4 = 15 // number of four hand possibilities from each six hand; same as 6 choose 2
)

func NChooseR(n int, r int) int {
	z := &big.Int{}
	z.Binomial(int64(n), int64(r))
	return int(z.Int64())
}
