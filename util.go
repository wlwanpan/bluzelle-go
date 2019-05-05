package bluzelle

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
	"unicode/utf8"
)

func encodeToUTF8(str string) string {
	if utf8.ValidString(str) {
		return str
	}
	buffer := make([]rune, 0, len(str))
	for i, r := range str {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(str[i:])
			if size == 1 {
				continue
			}
		}
		buffer = append(buffer, r)
	}
	return string(buffer)
}

func randUint64() uint64 {
	now := time.Now().UTC().Unix()
	r := rand.New(rand.NewSource(now))
	return r.Uint64()
}

// ReadPemFile reads a private pem file.
func ReadPemFile(path string) ([]byte, error) {
	privKeyFile, err := os.Open(path)
	if err != nil {
		return []byte{}, err
	}
	defer privKeyFile.Close()
	pemStats, err := privKeyFile.Stat()
	if err != nil {
		return []byte{}, err
	}
	log.Println("Loaded pem file", pemStats.Name())

	return ioutil.ReadAll(privKeyFile)
}
