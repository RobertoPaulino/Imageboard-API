package main

func roundToThreeDecimals(value float64) float64 {
	return float64(int(value*1000)) / 1000
}
