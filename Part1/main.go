package main

import (
	"fmt"
	"math"
)

func main() {

	fmt.Println(calculateFinalPrice(100, "TH", true))
}

func calculateFinalPrice(basePrice int, countryCode string, isFirstOrder bool) float64 {

	var calculatedPrice float64

	calculatedPrice = float64(basePrice)

	switch countryCode {
	case "TH":
		calculatedPrice += float64(basePrice) * 0.07
	case "FR":
		calculatedPrice += float64(basePrice) * 0.2
	}

	if isFirstOrder {
		calculatedPrice -= 100
	}

	if calculatedPrice <= 0 {
		return 0
	}

	finalPrice := roundPrice(calculatedPrice)

	return finalPrice

}

func roundPrice(price float64) float64 {
	return math.Round(price*100) / 100
}
