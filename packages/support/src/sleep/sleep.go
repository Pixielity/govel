package sleep

import (
	"context"
	"fmt"
	"sync"
	"time"

	carbon "govel/support/src/carbon"
	"govel/support/src/traits"
)

// SleepDuration represents a duration with different time units
type SleepDuration struct {
	microseconds int64
}

// NewSleepDuration creates a new SleepDuration from microseconds
func NewSleepDuration(microseconds int64) *SleepDuration {
	return &SleepDuration{microseconds: microseconds}
}

// Microseconds returns the duration in microseconds
func (d *SleepDuration) Microseconds() int64 {
	return d.microseconds
}

// Milliseconds returns the duration in milliseconds
func (d *SleepDuration) Milliseconds() int64 {
	return d.microseconds / 1000
}

// Seconds returns the duration in seconds
func (d *SleepDuration) Seconds() float64 {
	return float64(d.microseconds) / 1000000.0
}

// Duration returns the duration as time.Duration
func (d *SleepDuration) Duration() time.Duration {
	return time.Duration(d.microseconds) * time.Microsecond
}

// String returns a human-readable representation
func (d *SleepDuration) String() string {
	if d.microseconds < 1000 {
		return fmt.Sprintf("%dμs", d.microseconds)
	} else if d.microseconds < 1000000 {
		return fmt.Sprintf("%.2fms", float64(d.microseconds)/1000.0)
	} else {
		return fmt.Sprintf("%.3fs", float64(d.microseconds)/1000000.0)
	}
}

// Copy creates a copy of the duration
func (d *SleepDuration) Copy() *SleepDuration {
	return &SleepDuration{microseconds: d.microseconds}
}

// Add adds microseconds to the duration
func (d *SleepDuration) Add(microseconds int64) *SleepDuration {
	d.microseconds += microseconds
	return d
}

// Sub subtracts microseconds from the duration
func (d *SleepDuration) Sub(microseconds int64) *SleepDuration {
	d.microseconds -= microseconds
	if d.microseconds < 0 {
		d.microseconds = 0
	}
	return d
}

// Sleep represents a Laravel-style sleep utility with testing support
type Sleep struct {
	traits.InteractsWithTime
	duration    *SleepDuration
	whileFunc   func() bool
	pending     *int64
	shouldSleep bool
	hasSlept    bool
	ctx         context.Context
}

// Global sleep configuration
var (
	sleepMutex         sync.RWMutex
	isFake             bool
	syncWithCarbon     bool
	fakeSequence       []*SleepDuration
	fakeSleepCallbacks []func(*SleepDuration)
)

// For creates a new Sleep instance for the given duration
// Duration can be:
//   - time.Duration
//   - int/int64 (treated as seconds)
//   - float64 (treated as seconds with decimal precision)
func For(duration interface{}) *Sleep {
	sleep := &Sleep{
		shouldSleep: true,
		hasSlept:    false,
		ctx:         context.Background(),
	}
	sleep.setDuration(duration)
	return sleep
}

// Until creates a new Sleep instance that sleeps until the given timestamp
func Until(timestamp interface{}) *Sleep {
	var target *carbon.Carbon

	switch t := timestamp.(type) {
	case *carbon.Carbon:
		target = t
	case int64:
		target = carbon.CreateFromTimestamp(t)
	case float64:
		target = carbon.CreateFromTimestamp(int64(t))
	case string:
		target = carbon.Parse(t)
	default:
		// Invalid timestamp, sleep for 0
		return For(0)
	}

	now := carbon.Now()
	if target.Lt(now) {
		// Target is in the past, don't sleep
		return For(0)
	}

	diff := target.DiffInSeconds(now)
	return For(time.Duration(diff) * time.Second)
}

// Usleep creates a Sleep instance for the given number of microseconds
func Usleep(microseconds int64) *Sleep {
	sleep := &Sleep{
		duration:    NewSleepDuration(microseconds),
		shouldSleep: true,
		hasSlept:    false,
		ctx:         context.Background(),
	}
	return sleep
}

// SleepFor creates a Sleep instance for the given number of seconds
func SleepFor(seconds float64) *Sleep {
	return For(time.Duration(seconds * float64(time.Second)))
}

// setDuration sets the duration from various input types
func (s *Sleep) setDuration(duration interface{}) {
	switch d := duration.(type) {
	case time.Duration:
		s.duration = NewSleepDuration(int64(d / time.Microsecond))
		s.pending = nil
	case int:
		pending := int64(d)
		s.pending = &pending
		s.duration = NewSleepDuration(0)
	case int64:
		pending := d
		s.pending = &pending
		s.duration = NewSleepDuration(0)
	case float64:
		pending := int64(d)
		s.pending = &pending
		s.duration = NewSleepDuration(0)
	default:
		// Default to 0 microseconds
		s.duration = NewSleepDuration(0)
		s.pending = nil
	}
}

// Duration modification methods

// Minutes treats the pending duration as minutes
func (s *Sleep) Minutes() *Sleep {
	if s.pending != nil {
		s.duration.Add(*s.pending * 60 * 1000000) // Convert minutes to microseconds
		s.pending = nil
	} else {
		panic("No pending duration specified")
	}
	return s
}

// Minute treats the pending duration as 1 minute (alias for Minutes)
func (s *Sleep) Minute() *Sleep {
	return s.Minutes()
}

// Seconds treats the pending duration as seconds
func (s *Sleep) Seconds() *Sleep {
	if s.pending != nil {
		s.duration.Add(*s.pending * 1000000) // Convert seconds to microseconds
		s.pending = nil
	} else {
		panic("No pending duration specified")
	}
	return s
}

// Second treats the pending duration as 1 second (alias for Seconds)
func (s *Sleep) Second() *Sleep {
	return s.Seconds()
}

// Milliseconds treats the pending duration as milliseconds
func (s *Sleep) Milliseconds() *Sleep {
	if s.pending != nil {
		s.duration.Add(*s.pending * 1000) // Convert milliseconds to microseconds
		s.pending = nil
	} else {
		panic("No pending duration specified")
	}
	return s
}

// Millisecond treats the pending duration as 1 millisecond (alias for Milliseconds)
func (s *Sleep) Millisecond() *Sleep {
	return s.Milliseconds()
}

// Microseconds treats the pending duration as microseconds
func (s *Sleep) Microseconds() *Sleep {
	if s.pending != nil {
		s.duration.Add(*s.pending) // Already in microseconds
		s.pending = nil
	} else {
		panic("No pending duration specified")
	}
	return s
}

// Microsecond treats the pending duration as 1 microsecond (alias for Microseconds)
func (s *Sleep) Microsecond() *Sleep {
	return s.Microseconds()
}

// And adds additional time to sleep for
func (s *Sleep) And(duration interface{}) *Sleep {
	switch d := duration.(type) {
	case int:
		pending := int64(d)
		s.pending = &pending
	case int64:
		pending := d
		s.pending = &pending
	case float64:
		pending := int64(d)
		s.pending = &pending
	}
	return s
}

// Conditional methods

// While sets a condition that must be true for sleep to continue
func (s *Sleep) While(condition func() bool) *Sleep {
	s.whileFunc = condition
	return s
}

// When only sleeps when the given condition is true
func (s *Sleep) When(condition interface{}) *Sleep {
	switch c := condition.(type) {
	case bool:
		s.shouldSleep = c
	case func() bool:
		s.shouldSleep = c()
	case func(*Sleep) bool:
		s.shouldSleep = c(s)
	}
	return s
}

// Unless only sleeps when the given condition is false
func (s *Sleep) Unless(condition interface{}) *Sleep {
	switch c := condition.(type) {
	case bool:
		s.shouldSleep = !c
	case func() bool:
		s.shouldSleep = !c()
	case func(*Sleep) bool:
		s.shouldSleep = !c(s)
	}
	return s
}

// Execution methods

// Then executes a callback after sleeping and returns the result
func (s *Sleep) Then(callback func() interface{}) interface{} {
	s.execute()
	s.hasSlept = true
	return callback()
}

// Context sets a context for cancellation
func (s *Sleep) WithContext(ctx context.Context) *Sleep {
	s.ctx = ctx
	return s
}

// execute performs the actual sleep operation
func (s *Sleep) execute() {
	if s.hasSlept || !s.shouldSleep {
		return
	}

	if s.pending != nil {
		panic("Unknown duration unit - call a unit method like .Seconds() or .Milliseconds()")
	}

	sleepMutex.RLock()
	fake := isFake
	sync := syncWithCarbon
	callbacks := make([]func(*SleepDuration), len(fakeSleepCallbacks))
	copy(callbacks, fakeSleepCallbacks)
	sleepMutex.RUnlock()

	if fake {
		// Record the sleep for testing
		sleepMutex.Lock()
		fakeSequence = append(fakeSequence, s.duration.Copy())
		sleepMutex.Unlock()

		// Sync with Carbon if requested
		if sync {
			carbon.SetTestNow(carbon.Now().AddSeconds(int(s.duration.Seconds())))
		}

		// Execute fake sleep callbacks
		for _, callback := range callbacks {
			callback(s.duration)
		}
		return
	}

	// Perform actual sleep
	remaining := s.duration.Copy()

	whileCondition := s.whileFunc
	if whileCondition == nil {
		// Default: sleep once
		executed := false
		whileCondition = func() bool {
			if executed {
				return false
			}
			executed = true
			return true
		}
	}

	// Execute the sleep with while condition support
	for whileCondition() {
		// Check context cancellation
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		duration := remaining.Duration()
		if duration <= 0 {
			break
		}

		// Use Go's time.Sleep for the actual sleeping
		time.Sleep(duration)

		// For simple sleep operations, we break after first iteration
		// The while condition controls whether we continue
		if s.whileFunc == nil {
			// No custom while function, sleep once and exit
			break
		}
		// If there's a custom while function, it controls the loop continuation
		// The loop will continue or break based on whileCondition() result
	}
}

// Note: In Go, we can't rely on destructors like PHP's __destruct
// Users should call Sleep() explicitly, use Then(), or rely on the
// automatic execution when the sleep duration is fully configured.

// Testing methods

// Fake enables fake sleeping for testing
func Fake(sync ...bool) {
	sleepMutex.Lock()
	defer sleepMutex.Unlock()

	isFake = true
	syncWithCarbon = false
	if len(sync) > 0 {
		syncWithCarbon = sync[0]
	}

	fakeSequence = []*SleepDuration{}
	fakeSleepCallbacks = []func(*SleepDuration){}
}

// StopFaking disables fake sleeping
func StopFaking() {
	sleepMutex.Lock()
	defer sleepMutex.Unlock()

	isFake = false
	syncWithCarbon = false
	fakeSequence = []*SleepDuration{}
	fakeSleepCallbacks = []func(*SleepDuration){}
}

// AssertSlept asserts that sleep occurred with the expected duration a specific number of times
func AssertSlept(expected func(*SleepDuration) bool, times int) error {
	sleepMutex.RLock()
	sequence := make([]*SleepDuration, len(fakeSequence))
	copy(sequence, fakeSequence)
	sleepMutex.RUnlock()

	count := 0
	for _, duration := range sequence {
		if expected(duration) {
			count++
		}
	}

	if count != times {
		return fmt.Errorf("expected sleep occurred %d times, but found %d times", times, count)
	}
	return nil
}

// AssertSleptTimes asserts the total number of sleep operations
func AssertSleptTimes(expected int) error {
	sleepMutex.RLock()
	actual := len(fakeSequence)
	sleepMutex.RUnlock()

	if actual != expected {
		return fmt.Errorf("expected %d sleeps but found %d", expected, actual)
	}
	return nil
}

// AssertNeverSlept asserts that no sleeping occurred
func AssertNeverSlept() error {
	return AssertSleptTimes(0)
}

// AssertInsomniac asserts that all sleep durations were 0 (no actual sleeping)
func AssertInsomniac() error {
	sleepMutex.RLock()
	sequence := make([]*SleepDuration, len(fakeSequence))
	copy(sequence, fakeSequence)
	sleepMutex.RUnlock()

	for _, duration := range sequence {
		if duration.Microseconds() > 0 {
			return fmt.Errorf("unexpected sleep duration of %s found", duration.String())
		}
	}
	return nil
}

// AssertSequence asserts the given sleep sequence was encountered
func AssertSequence(expected []*Sleep) error {
	if err := AssertSleptTimes(len(expected)); err != nil {
		return err
	}

	sleepMutex.RLock()
	sequence := make([]*SleepDuration, len(fakeSequence))
	copy(sequence, fakeSequence)
	sleepMutex.RUnlock()

	for i, expectedSleep := range expected {
		if expectedSleep == nil {
			continue
		}

		expectedSleep.shouldSleep = false // Mark as shouldn't actually sleep

		if i >= len(sequence) {
			return fmt.Errorf("expected sleep at position %d not found", i)
		}

		actual := sequence[i]
		expected := expectedSleep.duration

		if actual.Microseconds() != expected.Microseconds() {
			return fmt.Errorf("expected sleep duration of %s but actually slept for %s",
				expected.String(), actual.String())
		}
	}

	return nil
}

// WhenFakingSleep registers a callback to be called when faking sleep
func WhenFakingSleep(callback func(*SleepDuration)) {
	sleepMutex.Lock()
	defer sleepMutex.Unlock()

	fakeSleepCallbacks = append(fakeSleepCallbacks, callback)
}

// SyncWithCarbon enables syncing Carbon's "now" time when sleeping
func SyncWithCarbon(sync ...bool) {
	sleepMutex.Lock()
	defer sleepMutex.Unlock()

	syncWithCarbon = true
	if len(sync) > 0 {
		syncWithCarbon = sync[0]
	}
}

// GetFakeSequence returns the sequence of fake sleeps (for testing)
func GetFakeSequence() []*SleepDuration {
	sleepMutex.RLock()
	defer sleepMutex.RUnlock()

	sequence := make([]*SleepDuration, len(fakeSequence))
	copy(sequence, fakeSequence)
	return sequence
}

// Helper functions for creating common sleep durations

// Minutes creates a sleep for the given number of minutes
func Minutes(minutes int64) *Sleep {
	return For(minutes).Minutes()
}

// Seconds creates a sleep for the given number of seconds
func Seconds(seconds int64) *Sleep {
	return For(seconds).Seconds()
}

// Milliseconds creates a sleep for the given number of milliseconds
func Milliseconds(milliseconds int64) *Sleep {
	return For(milliseconds).Milliseconds()
}

// Microseconds creates a sleep for the given number of microseconds
func Microseconds(microseconds int64) *Sleep {
	return For(microseconds).Microseconds()
}

// Sleep executes the sleep immediately (convenience method)
func (s *Sleep) Sleep() {
	s.execute()
	s.hasSlept = true
}

// Duration returns the total sleep duration
func (s *Sleep) Duration() *SleepDuration {
	return s.duration.Copy()
}

// String returns a string representation of the sleep
func (s *Sleep) String() string {
	return fmt.Sprintf("Sleep{duration: %s, shouldSleep: %t, hasSlept: %t}",
		s.duration.String(), s.shouldSleep, s.hasSlept)
}

// Example usage and patterns:
//
// Basic usage:
//   For(5).Seconds().Sleep()                    // Sleep for 5 seconds
//   For(100).Milliseconds().Sleep()             // Sleep for 100ms
//   Usleep(1000)                               // Sleep for 1000 microseconds
//
// Conditional sleeping:
//   For(1).Second().When(condition).Sleep()    // Only sleep if condition is true
//   For(1).Second().Unless(busy).Sleep()       // Sleep unless busy
//
// Complex durations:
//   For(2).Minutes().And(30).Seconds().Sleep() // Sleep for 2.5 minutes
//
// With callbacks:
//   For(1).Second().Then(func() interface{} {
//       fmt.Println("Done sleeping")
//       return "result"
//   })
//
// Testing:
//   Fake()
//   For(5).Seconds().Sleep()
//   AssertSleptTimes(1)
//   StopFaking()
