package types

// Hasher is an interface used to provide a hashing algorithm for the library.
type Hasher interface {
	Hash(data []byte) ([]byte, error)
}
