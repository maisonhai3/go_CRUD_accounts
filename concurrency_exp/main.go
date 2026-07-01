package main

import (
	"concurrency_exp/cache_stampede"
)

func main() {
	// DataRace()
	// MultipleWithdraw()
	// MultipleWithdrawError()
	// cache_stampede.TrySingleFlight()
	cache_stampede.DirtyHandSingleFlight()
}
