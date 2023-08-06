package portpool

import "sync"

// Pool is a thread-safe pool of port numbers.
type Pool struct {
	lower, upper uint16

	last uint16
	used map[uint16]bool
	mu   sync.Mutex
}

// NewPool returns a new port pool.
func NewPool(min, max uint16) *Pool {
	return &Pool{
		lower: min,
		upper: max,

		last: min - 1,
		used: make(map[uint16]bool),
	}
}

// Min returns the minimum port number.
func (p *Pool) Min() uint16 {
	return p.lower
}

// Max returns the maximum port number.
func (p *Pool) Max() uint16 {
	return p.upper
}

// Get returns a port number in the range [Min, Max].
// If no port is available, 0 is returned.
func (p *Pool) Get() uint16 {
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

	return 0
}

func (p *Pool) Put(port uint16) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.used, port)
}
