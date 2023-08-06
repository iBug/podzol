package portpool

import "sync"

// Pool is a thread-safe pool of port numbers.
type Pool struct {
	lower, upper int

	last int
	used map[int]bool
	mu   sync.Mutex
}

// NewPool returns a new port pool.
func NewPool(min, max int) *Pool {
	return &Pool{
		lower: min,
		upper: max,

		last: min - 1,
		used: make(map[int]bool),
	}
}

// Min returns the minimum port number.
func (p *Pool) Min() int {
	return p.lower
}

// Max returns the maximum port number.
func (p *Pool) Max() int {
	return p.upper
}

// Get returns a port number in the range [Min, Max].
// If no port is available, -1 is returned.
func (p *Pool) Get() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i := p.last + 1; i <= p.upper; i++ {
		if !p.used[i] {
			p.used[i] = true
			p.last = i
			return i
		}
	}
	for i := p.lower; i <= p.last; i++ {
		if !p.used[i] {
			p.used[i] = true
			p.last = i
			return i
		}
	}

	return -1
}

func (p *Pool) Put(port int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.used, port)
}
