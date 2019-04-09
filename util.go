package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
	"unicode/utf8"
)

func stringToASCII(str string) []byte {
	acsii := make([]byte, utf8.RuneCountInString(str))
	c := 0
	for _, r := range str {
		acsii[c] = byte(runeToASCII(r))
		c++
	}
	return acsii
}

func runeToASCII(r rune) rune {
	switch {
	case 97 <= r && r <= 122:
		return r - 32
	case 65 <= r && r <= 90:
		return r + 32
	default:
		return r
	}
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
