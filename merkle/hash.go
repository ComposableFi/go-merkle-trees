package merkle

// ConcatAndHash appends two bytes and the uses hasher to hash the appended bytes
func ConcatAndHash(hasher Hasher, left []byte, right []byte) ([]byte, error) {
	if right == nil {
		return left, nil
	}
	return hasher.Hash(append(left[:], right[:]...))
}
