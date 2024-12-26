// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package medley

import (
	"fmt"
	"io"
	"iter"
	"unsafe"
)

// Service represents some sort of endpoint that objects can be hashed to.
// The only requirement is that a Service satisfy comparable so that in can
// be used as map keys.
//
// A service's underlying type may be a string, such as a host name or URL.
// Or a service could be a struct that gives a richer description of an endpoint.
type Service interface {
	comparable
}

// ServiceHasher handles writing a Service's hashable bytes.  This closure
// type is responsible for converting a Service object into a series of
// bytes to submit to a hashing function.
type ServiceHasher[S Service] func(io.Writer, S) error

// DefaultServiceHasher uses fmt.Fprint to write a service object's
// hash bytes to dst. This function can be used when no ServiceHasher is
// supplied.
//
// For most concrete Service types, a custom ServiceHasher is preferable
// as it will be more efficient than the fmt package's reflection.
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

// HashStringTo is a ServiceHasher for services whose underlying type is string.
// Typically hostnames and URLs fall into this category.
//
// The string is written is such a way as to minimize allocations.
func HashStringTo[SS StringService](dst io.Writer, service SS) error {
	serviceBytes := unsafe.Slice(unsafe.StringData(string(service)), len(service))
	_, err := dst.Write(serviceBytes)
	return err
}

// Map is a mapping of services onto arbitrary values. A Map can be used to
// dedupe services or to provide fast access to some associated value object.
type Map[S Service, V any] map[S]V

// Len returns the count of services in this map.
func (m Map[S, V]) Len() int { return len(m) }

// Update indicates the disposition of a service object that is (possibly) an update.
type Update[S Service, V any] struct {
	// Service is the service object. One result per service in the update list
	// will be generated.
	Service S

	// Value is the value object associated with Service. If the Service didn't
	// exist in the map, this field will just be the zero value for its type.
	Value V

	// Exists indicates whether the Service field existed in the Map.
	Exists bool
}

// Update compares an updated slice of services to this Map. One Update is
// generated for each service object in the update slice, and calling code can
// iterate over the update to examine its disposition.
//
// This method is primarily useful when determining what to do when a set
// of services needs to be updated and rehashed.
func (m Map[S, V]) Update(services ...S) iter.Seq[Update[S, V]] {
	return func(f func(Update[S, V]) bool) {
		for _, svc := range services {
			v, exists := m[svc]
			if !f(Update[S, V]{Service: svc, Value: v, Exists: exists}) {
				return
			}
		}
	}
}

// BasicService is a URI-based service object that represents the typical things an application
// needs when describing a service. A string can hold the information in this struct,
// but sometimes an application needs easy access to the individual parts of a URI.
type BasicService struct {
	// Scheme is the URI scheme for the service, e.g. http.
	Scheme string

	// Host is the domain name or IP address of this service.
	Host string

	// Port is the IP port the service listens on.
	Port int

	// Path is the URI path for this service.
	Path string
}

// HashBasicServiceTo is a ServiceHasher for BasicService instances.
func HashBasicServiceTo(dst io.Writer, s BasicService) error {
	return NewHashBuilder(dst).
		WriteString(s.Scheme).
		WriteString(s.Host).
		WriteUint16(uint16(s.Port)).
		WriteString(s.Path).
		Err()
}

var _ ServiceHasher[BasicService] = HashBasicServiceTo
