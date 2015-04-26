package toolkit

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash/adler32"
	"hash/crc32"
	"os"

	"bytes"
	"math/rand"
	"time"
)

func RandomString(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		seed := RandInt(0, 10)
		mod := seed % 3
		if mod == 0 {
			temp = string(RandInt(65, 90))
		} else if mod == 1 {
			temp = string(RandInt(48, 57))
		} else {
			temp = string(RandInt(97, 122))
		}
		result.WriteString(temp)
		i++
	}
	return result.String()
}

func RandInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func HashCrc32(data []byte) uint32 {
	return uint32(crc32.ChecksumIEEE(data))
}

func Hash32(data []byte) uint32 {
	h := adler32.New()
	h.Write(data)
	return h.Sum32()
}

func MakeMd5(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	return h.Sum(nil)
}

func MakeMd5Str(data []byte) string {
	return hex.EncodeToString(MakeMd5(data))
}

func MakeMd5StrWithMask(data, mask []byte) string {
	buf := make([]byte, len(data)+len(mask))
	copy(buf, data)
	copy(buf[len(data):], mask)
	return MakeMd5Str(buf)
}

func MakeMd5FromFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 2048)
	h := md5.New()
	for {
		n, _ := file.Read(buf)
		if n <= 0 {
			break
		}
		h.Write(buf)
	}
	return h.Sum(nil), nil
}

func MakeMd5StrFromFile(filePath string) (string, error) {
	bytes, err := MakeMd5FromFile(filePath)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func MakeSha1(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

func MakeSha1Str(data []byte) string {
	return hex.EncodeToString(MakeSha1(data))
}

func MakeSha1FromFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 2048)
	h := sha1.New()
	for {
		n, _ := file.Read(buf)
		if n <= 0 {
			break
		}
		h.Write(buf)
	}
	return h.Sum(nil), nil
}

func MakeSha1StrFromFile(filePath string) (string, error) {
	bytes, err := MakeSha1FromFile(filePath)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
