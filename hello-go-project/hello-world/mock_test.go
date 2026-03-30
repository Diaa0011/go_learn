package main

// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"reflect"
// 	"testing"
// 	"time"
// )

// const finalWord = "Go!"
// const countdownStart = 3

// func TestCountdown(t *testing.T) {

// 	t.Run("prints 3 to Go!", func(t *testing.T) {
// 		buffer := &bytes.Buffer{}
// 		Countdown(buffer, &SpyCountdownOperations{})

// 		got := buffer.String()
// 		want := `3
// 2
// 1
// Go!`

// 		if got != want {
// 			t.Errorf("got %q want %q", got, want)
// 		}
// 	})

// 	t.Run("sleep before every print", func(t *testing.T) {
// 		spySleepPrinter := &SpyCountdownOperations{}
// 		Countdown(spySleepPrinter, spySleepPrinter)

// 		want := []string{
// 			write,
// 			sleep,
// 			write,
// 			sleep,
// 			write,
// 			sleep,
// 			write,
// 		}

// 		if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
// 			t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
// 		}
// 	})
// }

// func Countdown(out io.Writer, sleeper Sleeper) {
// 	for i := countdownStart; i > 0; i-- {
// 		fmt.Fprintln(out, i)
// 		sleeper.Sleep()
// 	}
// 	fmt.Fprint(out, finalWord)
// }

// type Sleeper interface {
// 	Sleep()
// }

// type SpySleeper struct {
// 	Calls int
// }

// func (s *SpySleeper) Sleep() {
// 	s.Calls++
// }

// type SpyCountdownOperations struct {
// 	Calls []string
// }

// func (s *SpyCountdownOperations) Sleep() {
// 	s.Calls = append(s.Calls, sleep)
// }

// func (s *SpyCountdownOperations) Write(p []byte) (n int, err error) {
// 	s.Calls = append(s.Calls, write)
// 	return
// }

// type ConfigurableSleeper struct {
// 	duration time.Duration
// 	sleep    func(time.Duration)
// }

// type SpyTime struct {
// 	durationSlept time.Duration
// }

// func (s *SpyTime) SetDurationSlept(duration time.Duration) {
// 	s.durationSlept = duration
// }

// func TestConfigurableSleeper(t *testing.T) {
// 	sleepTime := 5 * time.Second

// 	spyTime := &SpyTime{}
// 	sleeper := ConfigurableSleeper{sleepTime, spyTime.SetDurationSlept}
// 	sleeper.Sleep()

// 	if spyTime.durationSlept != sleepTime {
// 		t.Errorf("should have slept for %v but slept for %v", sleepTime, spyTime.durationSlept)
// 	}
// }

// func (c *ConfigurableSleeper) Sleep() {
// 	c.sleep(c.duration)
// }

// const write = "write"
// const sleep = "sleep"
