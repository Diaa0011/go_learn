package generics

import (
	"strings"
	"testing"
)

// func TestReduce(t *testing.T) {
// 	t.Run("multiplication of all elements", func(t *testing.T) {
// 		multiply := func(x, y int) int {
// 			return x * y
// 		}

// 		AssertEqual(t, Reduce([]int{1, 2, 3}, multiply, 1), 6)
// 	})

// 	t.Run("concatenate strings", func(t *testing.T) {
// 		concatenate := func(x, y string) string {
// 			return x + y
// 		}

// 		AssertEqual(t, Reduce([]string{"a", "b", "c"}, concatenate, ""), "abc")
// 	})
// }

func TestBadBank(t *testing.T) {
	var (
		riya  = Account{Name: "Riya", Balance: 100}
		chris = Account{Name: "Chris", Balance: 75}
		adil  = Account{Name: "Adil", Balance: 200}

		transactions = []Transaction{
			NewTransaction(chris, riya, 100),
			NewTransaction(adil, chris, 25),
		}
	)

	newBalanceFor := func(account Account) float64 {
		return NewBalanceFor(account, transactions).Balance
	}

	AssertEqual(t, newBalanceFor(riya), 200)
	AssertEqual(t, newBalanceFor(chris), 0)
	AssertEqual(t, newBalanceFor(adil), 175)
}

func TestFind(t *testing.T) {
	t.Run("find first even number", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		firstEvenNumber, found := Find(numbers, func(x int) bool {
			return x%2 == 0
		})
		AssertTrue(t, found)
		AssertEqual(t, firstEvenNumber, 2)
	})

	t.Run("Find the best programmer", func(t *testing.T) {
		people := []Person{
			Person{Name: "Kent Beck"},
			Person{Name: "Martin Fowler"},
			Person{Name: "Chris James"},
		}

		king, found := Find(people, func(p Person) bool {
			return strings.Contains(p.Name, "Chris")
		})

		AssertTrue(t, found)
		AssertEqual(t, king, Person{Name: "Chris James"})
	})
}

func BalanceFor(transactions []Transaction, name string) float64 {

	adjustBalance := func(currentBalance float64, t Transaction) float64 {
		if t.From == name {
			return currentBalance - t.Sum
		}

		if t.To == name {
			return currentBalance + t.Sum
		}

		return currentBalance
	}

	return Reduce(transactions, adjustBalance, 0.0)
}

func Sum(numbers []int) int {
	add := func(acc, x int) int { return acc + x }

	return Reduce(numbers, add, 0)
}

func SumAllTails(numbers ...[]int) []int {
	sumTail := func(acc, x []int) []int {
		if len(x) == 0 {
			return append(acc, 0)
		} else {
			tail := x[1:]
			return append(acc, Sum(tail))
		}
	}
	return Reduce(numbers, sumTail, []int{})
}

// func Reduce[A any](collection []A, f func(A, A) A, initialValue A) A {
// 	var result = initialValue
// 	for _, x := range collection {
// 		result = f(result, x)
// 	}

// 	return result
// }

func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func AssertTrue(t *testing.T, got bool) {
	t.Helper()
	if !got {
		t.Errorf("got %v, want true", got)
	}
}
