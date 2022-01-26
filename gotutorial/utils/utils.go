package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToJSON(i interface{}) []byte {
	b, err := json.Marshal(i)
	HandleErr(err)
	return b
}

func ToBytes(i interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	HandleErr(enc.Encode(i))
	return buf.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	HandleErr(gob.NewDecoder(bytes.NewReader(data)).Decode(i))
}

func Hash(anything interface{}) string {
	s := fmt.Sprintf("%v", anything)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func Splitter(s string, sep string, i int) string {
	arr := strings.Split(s, sep)
	if len(arr) > i {
		return arr[i]
	}
	return ""
}
