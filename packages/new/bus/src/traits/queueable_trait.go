package traits

import (
	"time"
)

// Queueable provides queueing capabilities to commands/jobs
// This is implemented as a struct that can be embedded in other structs
type Queueable struct {
	// Connection specifies which queue connection to use
	Connection string `json:"connection,omitempty"`

	// Queue specifies which queue to dispatch to
	Queue string `json:"queue,omitempty"`

	// Delay specifies how long to wait before processing the job
	Delay time.Duration `json:"delay,omitempty"`

	// Tries specifies the maximum number of attempts
	Tries int `json:"tries,omitempty"`

	// Timeout specifies the maximum execution time
	Timeout time.Duration `json:"timeout,omitempty"`

	// RetryUntil specifies when to stop retrying
	RetryUntil *time.Time `json:"retry_until,omitempty"`

	// Backoff specifies the backoff strategy for retries
	Backoff []time.Duration `json:"backoff,omitempty"`

	// Middleware specifies middleware to run before the job
	Middleware []string `json:"middleware,omitempty"`

	// ChainConnection specifies the connection for chained jobs
	ChainConnection string `json:"chain_connection,omitempty"`

	// ChainQueue specifies the queue for chained jobs
	ChainQueue string `json:"chain_queue,omitempty"`

	// ChainCatchCallbacks specifies callbacks for failed chains
	ChainCatchCallbacks []func(error) error `json:"-"`
}

// NewQueueable creates a new Queueable instance with defaults
func NewQueueable() *Queueable {
	return &Queueable{
		Tries:   3,
		Timeout: 60 * time.Second,
	}
}

// OnConnection sets the queue connection
func (q *Queueable) OnConnection(connection string) *Queueable {
	q.Connection = connection
	return q
}

// OnQueue sets the queue name
func (q *Queueable) OnQueue(queue string) *Queueable {
	q.Queue = queue
	return q
}

// DelayFor sets a delay before processing
func (q *Queueable) DelayFor(delay time.Duration) *Queueable {
	q.Delay = delay
	return q
}

// DelayUntil sets an absolute time to start processing
func (q *Queueable) DelayUntil(when time.Time) *Queueable {
	q.Delay = time.Until(when)
	return q
}

// SetTries sets the maximum number of attempts
func (q *Queueable) SetTries(tries int) *Queueable {
	q.Tries = tries
	return q
}

// SetTimeout sets the maximum execution time
func (q *Queueable) SetTimeout(timeout time.Duration) *Queueable {
	q.Timeout = timeout
	return q
}

// SetRetryUntil sets when to stop retrying
func (q *Queueable) SetRetryUntil(until time.Time) *Queueable {
	q.RetryUntil = &until
	return q
}

// SetBackoff sets the backoff strategy
func (q *Queueable) SetBackoff(backoff []time.Duration) *Queueable {
	q.Backoff = backoff
	return q
}

// Through sets middleware to run before the job
func (q *Queueable) Through(middleware ...string) *Queueable {
	q.Middleware = middleware
	return q
}

// AllOnConnection sets the connection for chained jobs
func (q *Queueable) AllOnConnection(connection string) *Queueable {
	q.ChainConnection = connection
	return q
}

// AllOnQueue sets the queue for chained jobs
func (q *Queueable) AllOnQueue(queue string) *Queueable {
	q.ChainQueue = queue
	return q
}

// GetConnection returns the queue connection
func (q *Queueable) GetConnection() string {
	return q.Connection
}

// GetQueue returns the queue name
func (q *Queueable) GetQueue() string {
	return q.Queue
}

// GetDelay returns the delay duration
func (q *Queueable) GetDelay() time.Duration {
	return q.Delay
}

// GetTries returns the maximum number of attempts
func (q *Queueable) GetTries() int {
	return q.Tries
}

// GetTimeout returns the maximum execution time
func (q *Queueable) GetTimeout() time.Duration {
	return q.Timeout
}

// GetRetryUntil returns when to stop retrying
func (q *Queueable) GetRetryUntil() *time.Time {
	return q.RetryUntil
}

// GetBackoff returns the backoff strategy
func (q *Queueable) GetBackoff() []time.Duration {
	return q.Backoff
}

// GetMiddleware returns the middleware list
func (q *Queueable) GetMiddleware() []string {
	return q.Middleware
}

// ShouldQueue returns true if this job should be queued
func (q *Queueable) ShouldQueue() bool {
	return true
}

// DisplayName returns a human-readable name for the job
func (q *Queueable) DisplayName() string {
	return "Queueable Job"
}

// QueueableInterface defines the interface for queueable jobs
type QueueableInterface interface {
	OnConnection(connection string) *Queueable
	OnQueue(queue string) *Queueable
	DelayFor(delay time.Duration) *Queueable
	DelayUntil(when time.Time) *Queueable
	SetTries(tries int) *Queueable
	SetTimeout(timeout time.Duration) *Queueable
	SetRetryUntil(until time.Time) *Queueable
	SetBackoff(backoff []time.Duration) *Queueable
	Through(middleware ...string) *Queueable
	AllOnConnection(connection string) *Queueable
	AllOnQueue(queue string) *Queueable
	ShouldQueue() bool
	DisplayName() string
}
