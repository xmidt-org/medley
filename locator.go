package medley

import (
	"errors"
	"sync/atomic"
)

var (
	// ErrNoServices is returned by a Locator to indicate that the Locator contains
	// no service entries.
	ErrNoServices = errors.New("no services defined in this locator")
)

// Locator is a service locator based on hashing input objects.
type Locator[S Service] interface {
	// Find locates a service for a particular string.
	Find(string) (S, error)
}

// MultiLocator represents an aggregate set of locators, each of which is
// consulted for services.
type MultiLocator[S Service] []Locator[S]

// Find returns the services from each locator in this aggregate. This method
// will halt early on error if any Locator returned an error other than ErrNoServices.
//
// This method only returns ErrNoServices if and only if every locator returned
// no services.
func (ml MultiLocator[S]) Find(object string) ([]S, error) {
	services := make([]S, 0, len(ml))
	for _, l := range ml {
		if svc, findErr := l.Find(object); findErr == nil {
			services = append(services, svc)
		} else if !errors.Is(findErr, ErrNoServices) {
			return nil, findErr
		}
	}

	if len(services) == 0 {
		return nil, ErrNoServices
	}

	return services, nil
}

// UpdatableLocator is a Locator whose actual implementation can be swapped
// out atomically. Useful for dynamic Locators such as would be driven
// by service discovery or DNS.
type UpdatableLocator[S Service] struct {
	impl atomic.Pointer[Locator[S]]
}

var _ Locator[string] = &UpdatableLocator[string]{}

// Set atomically changes this locator's implementation. If the implementation
// is nil, methods of this UpdatableLocator will generally return ErrNoServices.
// Setting an implementation to nil effectively "turns off" this locator.
func (ul *UpdatableLocator[S]) Set(impl Locator[S]) {
	if impl != nil {
		ul.impl.Store(&impl)
	} else {
		ul.impl.Store(nil)
	}
}

// Find consults the current Locator implementation for the given object.
// This method returns ErrNoServices is no implementation has been set
// yet.
func (ul *UpdatableLocator[S]) Find(object string) (svc S, err error) {
	if l := ul.impl.Load(); l != nil {
		svc, err = (*l).Find(object)
	} else {
		err = ErrNoServices
	}

	return
}
