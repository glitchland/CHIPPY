package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	oc "pkg/opcodes"
)

func main() {
	var opcodes []uint16

	romFile := flag.String("rom", "", "the name of the rom to load and emulate")
	flag.Parse()

	fmt.Println(romFile)
	if *romFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fp, err := os.Open(*romFile)

	if err != nil {
		fmt.Printf("Could not find the '%s' rom file\n", *romFile)
		os.Exit(1)
	}

	defer fp.Close()

	fi, err := fp.Stat()
	if err != nil {
		fmt.Printf("Problem stating rom file.\n")
		os.Exit(1)
	}

	romLength := fi.Size()
	romBuffer := make([]byte, romLength)

	// read the ROM contents
	_, err = fp.Read(romBuffer)

	if err == io.EOF {
		fmt.Printf("Problem readiing rom file.\n")
		os.Exit(1)
	}

	for i := 0; i < len(romBuffer)-1; i +=2 {
		wVal := uint16(binary.BigEndian.Uint16(romBuffer[i:]))
		opcodes = append(opcodes, wVal)
	}

	addr := 0
	for _, opcode := range opcodes {
		oc.DisassembleAndPrint(0x200+addr, opcode)
		addr += 2
	}
}
