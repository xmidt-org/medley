package medley

import (
	"hash"
	"hash/fnv"
	"strings"

	"github.com/spaolacci/murmur3"
)

const (
	// AlgorithmFNV is the configuration value for an fnv.New64a algorithm
	AlgorithmFNV = "fnv"

	// AlgorithmMurmur3 is the configuration value for a murmur3.New64 algorithm
	AlgorithmMurmur3 = "murmur3"
)

// UnknownAlgorithmError indicates that no algorithm could be created using
// the given Name
type UnknownAlgorithmError struct {
	// Name is the name which is unrecognized.  This will never be blank,
	// as a blank name is interpreted as AlgorithmDefault.
	Name string
}

// Error fulfills the error interface
func (e *UnknownAlgorithmError) Error() string {
	var o strings.Builder
	o.WriteString(" algorithm with name: ")
	o.WriteString(e.Name)
	return o.String()
}

// Algorithm is a constructor for a 64-bit hash object.  The various NewXXX function
// in the stdlib hash subpackages are of this type.
type Algorithm func() hash.Hash64

// DefaultAlgorithm returns the Algorithm to be used when none is supplied or configured.
// Currently, this package defaults to github.com/spaolacci/murmur3.
func DefaultAlgorithm() Algorithm {
	return murmur3.New64
}

var builtinAlgorithms = map[string]Algorithm{
	"fnv":     fnv.New64a,
	"murmur3": murmur3.New64,
}

// findAlgorithm first consults builtinAlgorithms, then the extensions, for the named algorithm.
// If name is empty, DefaultAlgorithm() is returned.
func findAlgorithm(name string, extensions map[string]Algorithm) (alg Algorithm, err error) {
	if len(name) > 0 {
		var found bool
		alg, found = builtinAlgorithms[name]
		if !found {
			alg, found = extensions[name]
		}

		if !found {
			err = &UnknownAlgorithmError{Name: name}
		}
	} else {
		alg = DefaultAlgorithm()
	}

	return
}

// GetAlgorithm accepts a configured name and attempts to locate the built-in algorithm
// associated with that name.  If name is empty, DefaultAlgorithm() is returned.
//
// This function will return an error of type *UnknownAlgorithmError if no such algorithm
// was found.
func GetAlgorithm(name string) (Algorithm, error) {
	return findAlgorithm(name, nil)
}

// FindAlgorithm accepts a configured name and attempts to locate the appropriate Algorithm.
// The set of built-in algorithms is consulted first, followed by the extensions (if supplied).
// The set of extensions can be nil.  If name is empty, DefaultAlgorithm() is returned.
//
// This function will return an error of type *UnknownAlgorithmError if no such algorithm
// was found.
func FindAlgorithm(name string, extensions map[string]Algorithm) (Algorithm, error) {
	return findAlgorithm(name, extensions)
}
