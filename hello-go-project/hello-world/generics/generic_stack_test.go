package generics

// import "testing"

// func TestStack(t *testing.T) {
// 	// t.Run("integer stack", func(t *testing.T) {
// 	// 	myStackOfInts := new(StackOfInts)

// 	// 	AssertTrue(t, myStackOfInts.IsEmpty())

// 	// 	myStackOfInts.Push(123)
// 	// 	AssertFalse(t, myStackOfInts.IsEmpty())

// 	// 	myStackOfInts.Push(456)
// 	// 	value, _ := myStackOfInts.Pop()
// 	// 	AssertEqual(t, value, 456)
// 	// 	value, _ = myStackOfInts.Pop()
// 	// 	AssertEqual(t, value, 123)
// 	// 	AssertTrue(t, myStackOfInts.IsEmpty())

// 	// })

// 	// t.Run("string stack", func(t *testing.T) {
// 	// 	myStackOfStrings := new(StackOfStrings)

// 	// 	// check stack is empty
// 	// 	AssertTrue(t, myStackOfStrings.IsEmpty())

// 	// 	// add a thing, then check it's not empty
// 	// 	myStackOfStrings.Push("123")
// 	// 	AssertFalse(t, myStackOfStrings.IsEmpty())

// 	// 	// add another thing, pop it back again
// 	// 	myStackOfStrings.Push("456")
// 	// 	value, _ := myStackOfStrings.Pop()
// 	// 	AssertEqual(t, value, "456")
// 	// 	value, _ = myStackOfStrings.Pop()
// 	// 	AssertEqual(t, value, "123")
// 	// 	AssertTrue(t, myStackOfStrings.IsEmpty())
// 	// })

// 	// t.Run("interface stack DX is horrid", func(t *testing.T) {
// 	// 	myStackOfInts := new(StackOfInts)

// 	// 	myStackOfInts.Push(1)
// 	// 	myStackOfInts.Push(2)
// 	// 	firstNum, _ := myStackOfInts.Pop()
// 	// 	secondNum, _ := myStackOfInts.Pop()
// 	// 	reallyFirstNum, ok := firstNum.(int)
// 	// 	AssertTrue(t, ok)
// 	// 	reallySecondNum, ok := secondNum.(int)
// 	// 	AssertTrue(t, ok)

// 	// 	AssertEqual(t, reallyFirstNum+reallySecondNum, 3)
// 	// })

// 	t.Run("integer stack", func(t *testing.T) {
// 		myStackOfInts := new(Stack[int])

// 		AssertTrue(t, myStackOfInts.IsEmpty())

// 		myStackOfInts.Push(123)
// 		AssertFalse(t, myStackOfInts.IsEmpty())

// 		myStackOfInts.Push(456)
// 		value, _ := myStackOfInts.Pop()
// 		AssertEqual(t, value, 456)
// 		value, _ = myStackOfInts.Pop()
// 		AssertEqual(t, value, 123)
// 		AssertTrue(t, myStackOfInts.IsEmpty())

// 		myStackOfInts.Push(1)
// 		myStackOfInts.Push(2)
// 		firstNum, _ := myStackOfInts.Pop()
// 		secondNum, _ := myStackOfInts.Pop()
// 		AssertEqual(t, firstNum+secondNum, 3)

// 	})
// }

// func TestNeoStack(t *testing.T) {
// 	t.Run("integer stack", func(t *testing.T) {
// 		myStackOfInts := NewStack[int]()

// 		// check stack is empty
// 		AssertTrue(t, myStackOfInts.IsEmpty())

// 		// add a thing, then check it's not empty
// 		myStackOfInts.Push(123)
// 		AssertFalse(t, myStackOfInts.IsEmpty())

// 		// add another thing, pop it back again
// 		myStackOfInts.Push(456)
// 		value, _ := myStackOfInts.Pop()
// 		AssertEqual(t, value, 456)
// 		value, _ = myStackOfInts.Pop()
// 		AssertEqual(t, value, 123)
// 		AssertTrue(t, myStackOfInts.IsEmpty())

// 		// can get the numbers we put in as numbers, not untyped interface{}
// 		myStackOfInts.Push(1)
// 		myStackOfInts.Push(2)
// 		firstNum, _ := myStackOfInts.Pop()
// 		secondNum, _ := myStackOfInts.Pop()
// 		AssertEqual(t, firstNum+secondNum, 3)
// 	})
// }

// func AssertTrue(t *testing.T, got bool) {
// 	t.Helper()
// 	if !got {
// 		t.Errorf("got %v, want true", got)
// 	}
// }

// func AssertFalse(t *testing.T, got bool) {
// 	t.Helper()
// 	if got {
// 		t.Errorf("got %v, want false", got)
// 	}
// }

// // func AssertFalse(t *testing.T, got bool) {
// // 	t.Helper()
// // 	if got {
// // 		t.Errorf("got %v, want false", got)
// // 	}
// // }

// // func AssertEqual[T comparable](t *testing.T, got, want T) {
// // 	t.Helper()
// // 	if got != want {
// // 		t.Errorf("got %v, want %v", got, want)
// // 	}
// // }

// // type StackOfInts struct {
// // 	values []int
// // }

// // func (s *StackOfInts) Push(value int) {
// // 	s.values = append(s.values, value)
// // }

// // func (s *StackOfInts) IsEmpty() bool {
// // 	return len(s.values) == 0
// // }

// // func (s *StackOfInts) Pop() (int, bool) {
// // 	if s.IsEmpty() {
// // 		return 0, false
// // 	}

// // 	index := len(s.values) - 1
// // 	el := s.values[index]
// // 	s.values = s.values[:index]
// // 	return el, true
// // }

// // type StackOfStrings struct {
// // 	values []string
// // }

// // func (s *StackOfStrings) Push(value string) {
// // 	s.values = append(s.values, value)
// // }

// // func (s *StackOfStrings) IsEmpty() bool {
// // 	return len(s.values) == 0
// // }

// // func (s *StackOfStrings) Pop() (string, bool) {
// // 	if s.IsEmpty() {
// // 		return "", false
// // 	}
// // 	index := len(s.values) - 1
// // 	el := s.values[index]
// // 	s.values = s.values[:index]
// // 	return el, true
// // }

// type Stack[T any] struct {
// 	values []T
// }

// func NewStack[T any]() *Stack[T] {
// 	return new(Stack[T])
// }

// func (s *Stack[T]) Push(value T) {
// 	s.values = append(s.values, value)
// }

// func (s *Stack[T]) IsEmpty() bool {
// 	return len(s.values) == 0
// }

// func (s *Stack[T]) Pop() (T, bool) {
// 	if s.IsEmpty() {
// 		var zero T
// 		return zero, false
// 	}

// 	index := len(s.values) - 1
// 	el := s.values[index]
// 	s.values = s.values[:index]
// 	return el, true
// }
// func AssertEqual(t *testing.T, got, want interface{}) {
// 	t.Helper()

// 	if got != want {
// 		t.Errorf("got %v, want %v", got, want)
// 	}
// }
