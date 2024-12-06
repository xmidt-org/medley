package medley

import (
	"fmt"
	"io"
	"unsafe"
)

// Service holds information about an endpoint. A Service could be
// as simple as a string or an arbitrary, comparable struct holding
// richer information.
type Service interface {
	comparable
}

// ServiceHasher handles writing a Service's hashable bytes.  This closure
// type is responsible for converting a Service object into a series of
// bytes to submit to a hashing function.
//
// This function will only return an error if the io.Writer
// returns an error.
type ServiceHasher[S Service] func(io.Writer, S) error

// DefaultServiceHasher uses fmt.Fprint to write a service object's
// hash bytes to dst. This function is used when no ServiceHasher is
// supplied.
//
// For most concrete Service types, one of the other ServiceHashers
// in this package or a custom ServiceHasher would be more efficient.
// In particular, the StringHasher is more efficient when a Service's
// underlying type is a string, e.g. a host name or URL.
func DefaultServiceHasher[S Service](dst io.Writer, service S) error {
	_, err := fmt.Fprint(dst, service)
	return err
}

// StringService is a service whose underlying type is a string. Hostnames,
// URLs, service locator ids, etc. are usually of this type.
type StringService interface {
	Service
	~string
}

// HashStringTo is a ServiceWriter for services whose underlying type is string.
// Typically hostnames and URLs fall into this category.
//
// The string is written is such a way as to minimize allocations.
func HashStringTo[SS StringService](dst io.Writer, service SS) error {
	serviceBytes := unsafe.Slice(unsafe.StringData(string(service)), len(service))
	_, err := dst.Write(serviceBytes)
	return err
}

// Set is a collection of services. These services are implicitly deduped
// as keys in a map.
type Set[S Service] map[S]bool

// Len returns the number of services in this set.
func (set Set[S]) Len() int {
	return len(set)
}

// Has tests if the given service is in this set.
func (set Set[S]) Has(svc S) bool {
	return set[svc]
}

// Add adds several services to this set. Duplicates are ignored.
func (set Set[S]) Add(services ...S) {
	for _, svc := range services {
		set[svc] = true
	}
}
