package main

const maxCount = 5

func Sum(numbers [maxCount]int) int {
	sum := 0
	//first type of for loops
	// for i := 0; i < maxCount; i++ {
	// 	sum += numbers[i]
	// }

	//second type of for loops as foreach in c#
	for _, number := range numbers {
		sum += number
	}
	return sum
}

func SumAll(numbersToSum ...[]int) []int {

	var sums []int
	for _, numbers := range numbersToSum {
		sums = append(sums, SumSlice(numbers))
	}

	return sums
}

func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			tail := numbers[1:]
			sums = append(sums, SumSlice(tail))
		}
	}

	return sums
}

func SumSlice(numbers []int) int {
	sum := 0

	for _, number := range numbers {
		sum += number
	}
	return sum
}
