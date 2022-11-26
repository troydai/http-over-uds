package payload

import "crypto/rand"

func MustGenRandomBytes(size int) []byte {
	blk := make([]byte, size)
	_, err := rand.Read(blk)
	if err != nil {
		panic(err)
	}

	return blk
}
