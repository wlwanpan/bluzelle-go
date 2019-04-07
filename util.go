package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
	"unicode/utf8"
)

func stringToASCII(s string) []byte {
	t := make([]byte, utf8.RuneCountInString(s))
	i := 0
	for _, r := range s {
		t[i] = byte(r)
		i++
	}
	return t
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
