package main

//import (
//	"bytes"
//	"encoding/binary"
//	"github.com/Binject/debug/elf"
//	"io/fs"
//	"io/ioutil"
//	"log"
//	"os"
//)
//
//func main() {
//	f, err := ioutil.ReadFile(os.Args[1])
//	if err != nil {
//		panic(err)
//	}
//
//	elfF, err := elf.NewFile(bytes.NewReader(f))
//	if err != nil {
//		panic(err)
//	}
//
//	data, shoff, err := elfF.Bytes()
//	if err != nil {
//		panic(err)
//	}
//
//	log.Printf("bytes written: %d, shoff: %x\n", len(data), shoff)
//	log.Printf("The 40 byte is: %x\n", data[40:48])
//	b := make([]byte, 8)
//	binary.LittleEndian.PutUint64(b, shoff)
//	j := 0
//	for i := 40; i < 48; i++ {
//		data[i] = b[j]
//		j++
//	}
//
//	log.Printf("The 40 byte is after change: %x\n", data[40:48])
//	err = ioutil.WriteFile("./eden123", data, fs.ModePerm)
//	if err != nil {
//		panic(err)
//	}
//}

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
	//if err := initLogger(); err != nil {
	//	log.Printf("error init logger: %s\n", err)
	//	return
	//}

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

	//output, shoff, err := bj.ElfBinject(input, payloadBytes, &bj.BinjectConfig{
	//	InjectionMethod: bj.SilvioInject,
	//})

	elfFile, err := elf.NewFile(bytes.NewReader(input))
	if err != nil {
		log.Printf("Error while opening input binary: %s\n", err)
		return
	}
	output, shoff, err := bj.StaticSilvioMethod(elfFile, payloadBytes)
	//output, shoff, err := elfFile.Bytes()
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

//func initLogger() error {
//	curDir := path.Dir(os.Args[0])
//	fileName := "launch.log"
//	logPath := path.Join(curDir, fileName)
//	f, err := os.Create(logPath)
//	if err != nil {
//		return err
//	}
//
//	log.Default().SetOutput(f)
//	return nil
//}
