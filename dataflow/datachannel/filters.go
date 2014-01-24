/* filters.go: Contains some Filters that can be reused for stream programming.
Author: James Fairbanks
Date: 2014-01-24
License : BSD
*/
package datachannel

/*
 * Specific filters and streams
 */

//mkStep: Returns a function that acts as a threshold Filter, returning true if x > t.
func MkStep(t Data) Filter {
	return func(x Data) Data {
		if x > t {
			return 1
		} else {
			return 0
		}
	}
}

//MkShift: Returns a function that shifts each Data point by a constant.
func MkShift(k Data) Filter {
	return func(x Data) Data {
		return x + k
	}
}

//MkScale: Returns a function that scales each Data point by a constant.
func MkScale(k Data) Filter {
	return func(x Data) Data {
		return x * k
	}
}

//MkSquareFilter: Returns a Filter returning 1 if x in [lowerbound, upperbound] else zero
func MkSquareFilter(lowerbound Data, upperbound Data) Filter {
	return func(x Data) Data {
		var bigenough, smallenough Data
		bigenough, smallenough = 0, 0
		if x >= lowerbound {
			bigenough = 1
		}
		if x <= upperbound {
			smallenough = 1
		}
		return bigenough * smallenough
	}
}
