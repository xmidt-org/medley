package consistent

import (
	"github.com/xmidt-org/medley"
)

const (
	// DefaultVNodes is the default number of nodes used when none is supplied.
	// This value is consistent with the default in github.com/billhathaway/consistentHash.
	DefaultVNodes = 200
)

// Builder is a fluent builder for hash Rings. This type can be used
// through normal instantiation or by starting a build chain with
// the Build function.
type Builder[S medley.Service] struct {
	hasher   hasher[S]
	services medley.ServiceSet[S]
}

// Build starts a fluent chain to initialize a Ring.
// More services can be added via the builder's Services method.
func Build[S medley.Service](services ...S) *Builder[S] {
	b := new(Builder[S])
	return b.Services(services...)
}

// VNodes sets the number of hash nodes used per service. By default,
// DefaultVNodes is used.
func (b *Builder[S]) VNodes(v int) *Builder[S] {
	b.hasher.vnodes = v
	return b
}

// Algorithm sets the medley hash algorithm to use. By default,
// medley.Murmur3 is used.
func (b *Builder[S]) Algorithm(a medley.Algorithm) *Builder[S] {
	b.hasher.alg = a
	return b
}

// ServiceHasher establishes the sequence of bytes used to hash a
// service object. By default, medley.DefaultServiceHasher is used.
//
// It's usually a good idea to set this, as you can generally get better
// performance with custom hash bytes.
func (b *Builder[S]) ServiceHasher(sh medley.ServiceHasher[S]) *Builder[S] {
	b.hasher.serviceHasher = sh
	return b
}

// Services adds services to the Ring that is built by this Builder. Multiple
// uses of this method are cumulative. Duplicate services are ignored.
//
// When Build is called, the set of services known to this builder is reset.
func (b *Builder[S]) Services(services ...S) *Builder[S] {
	if b.services == nil {
		b.services = make(medley.ServiceSet[S], len(services))
	}

	for _, svc := range services {
		b.services[svc] = true
	}

	return b
}

// newHasher creates a token hasher using this builder's configuration.
// This method enforces defaults, so the returned hasher is ready to use.
func (b *Builder[S]) newHasher() (h hasher[S]) {
	h = b.hasher
	if h.vnodes < 1 {
		h.vnodes = DefaultVNodes
	}

	if h.alg == nil {
		h.alg = medley.Murmur3{}
	}

	if h.serviceHasher == nil {
		h.serviceHasher = medley.DefaultServiceHasher[S]
	}

	return
}

// Build creates a brand new Ring instance. The set of services known to this
// builder is reset, and a distinct new Ring is returned.
//
// This Builder can be reused to create multiple Rings, although Services will
// need to be added between calls to Build. However, the Update function more
// efficiently handles creating a new Ring with an updated set of services.
func (b *Builder[S]) Build() *Ring[S] {
	hasher := b.newHasher()
	r := &Ring[S]{
		hasher:   hasher,
		services: make(serviceNodes[S], b.services.Len()),
		nodes:    make(nodes[S], 0, hasher.totalCount(b.services.Len())),
	}

	for svc := range b.services {
		snodes := hasher.serviceNodes(svc)
		r.services[svc] = snodes
		r.nodes = append(r.nodes, snodes...)
	}

	r.nodes.sort()
	b.services = nil
	return r
}
