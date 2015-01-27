package main

import (
	"math"
	"os/user"
)

// roundTheTimestamp round the timestamp based on precision
func roundTheTimestamp(timestamp, precision int64) int64 {
	floor := math.Floor(float64(timestamp) / float64(precision))
	return int64(floor) * precision
}

// roundFloat round the number
func roundFloat(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	x = .5
	if frac < 0.0 {
		x = -.5
	}
	if frac >= x {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}

// usePercent convert metrict to percent (integer)
func usePercent(total, free uint64) uint64 {
	used := (total - free)

	if total != 0 {
		u100 := used * 100
		pct := float64(u100) / float64(total)

		return uint64(roundFloat(pct, 0))
	}

	return 0
}

func checkPrivileges() bool {
	if user, err := user.Current(); err == nil && user.Gid == "0" {
		return true
	}

	return false
}
