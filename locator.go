package medley

import (
	"errors"
	"sync"
	"sync/atomic"
	"unsafe"
)

var (
	// ErrNoServices is returned by a Locator to indicate that the Locator contains
	// no service entries.
	ErrNoServices = errors.New("no services defined in this locator")
)

// Locator is a service locator based on hashing input objects.
type Locator[S Service] interface {
	// Find locates a service for a particular key.
	Find([]byte) (S, error)
}

// FindString locates a service for a string key.
func FindString[S Service](l Locator[S], v string) (S, error) {
	return l.Find(
		unsafe.Slice(unsafe.StringData(v), len(v)),
	)
}

// MultiLocator represents an aggregate set of locators, each of which is
// consulted for services. Methods on this type are safe for concurrent usage.
// The zero value for this type is usable, but will return ErrNoServices.
// To initialize a MultiLocator with some locators, use NewMultiLocator.
//
// A MultiLocator must not be copied after creation.
type MultiLocator[S Service] struct {
	lock     sync.RWMutex
	locators []Locator[S]
}

// NewMultiLocator returns a MultiLocator initialized with the give set of Locators.
func NewMultiLocator[S Service](ls ...Locator[S]) *MultiLocator[S] {
	return &MultiLocator[S]{
		locators: append(
			[]Locator[S]{},
			ls...,
		),
	}
}

// Add adds another locator to this MultiLocator. This method does not protect
// against adding a Locator more than once.
func (ml *MultiLocator[S]) Add(l Locator[S]) {
	ml.lock.Lock()
	ml.locators = append(ml.locators, l)
	ml.lock.Unlock()
}

// Remove removes a locator from this MultiLocator. If the same Locator
// was added multiple times, this method only removes the first one.
func (ml *MultiLocator[S]) Remove(l Locator[S]) {
	defer ml.lock.Unlock()
	ml.lock.Lock()

	for i, candidate := range ml.locators {
		if candidate == l {
			last := len(ml.locators) - 1
			ml.locators[i], ml.locators[last] = ml.locators[last], nil
			ml.locators = ml.locators[:last]
			return
		}
	}
}

// Find returns the services from each locator in this aggregate. This method
// will halt early on error if any Locator returned an error other than ErrNoServices.
//
// This method only returns ErrNoServices if and only if every locator returned
// no services.
func (ml *MultiLocator[S]) Find(object []byte) ([]S, error) {
	defer ml.lock.RUnlock()
	ml.lock.RLock()

	services := make([]S, 0, len(ml.locators))
	for _, l := range ml.locators {
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

// FindString locates services based on a string key.
func (ml *MultiLocator[S]) FindString(object string) ([]S, error) {
	return ml.Find(
		unsafe.Slice(unsafe.StringData(object), len(object)),
	)
}

// UpdatableLocator is a Locator whose actual implementation can be swapped
// out atomically. Useful for dynamic Locators such as would be driven
// by service discovery or DNS.
//
// The zero value of this type is usable, but will return ErrNoServices. Use
// NewUpdatableLocator to return an initialized UpdatableLocator.
type UpdatableLocator[S Service] struct {
	impl atomic.Pointer[Locator[S]]
}

// NewUpdatableLocator returns an UpdatableLocator initialized with the given
// implementation.
func NewUpdatableLocator[S Service](impl Locator[S]) *UpdatableLocator[S] {
	ul := new(UpdatableLocator[S])
	ul.Set(impl)
	return ul
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
func (ul *UpdatableLocator[S]) Find(object []byte) (svc S, err error) {
	if l := ul.impl.Load(); l != nil {
		svc, err = (*l).Find(object)
	} else {
		err = ErrNoServices
	}

	return
}
