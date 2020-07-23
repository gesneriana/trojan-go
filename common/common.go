package common

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Runnable interface {
	Run() error
	Close() error
}

func SHA224String(password string) string {
	hash := sha256.New224()
	hash.Write([]byte(password))
	val := hash.Sum(nil)
	str := ""
	for _, v := range val {
		str += fmt.Sprintf("%02x", v)
	}
	return str
}

var ProgramDir string

func GetProgramDir() string {
	if len(ProgramDir) > 0 {
		return ProgramDir
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
