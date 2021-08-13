// Copyright 2012 The Go Asm Authors
// SPDX-License-Identifier: BSD-3-Clause

// Command genasmdb auto-generate an assembly database from asmjit/asmdb.
package main

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-json-experiment/json"
)

func init() {
	spew.Config = spew.ConfigState{
		Indent:                  "  ",
		SortKeys:                true, // maps should be spewed in a deterministic order
		DisablePointerAddresses: true, // don't spew the addresses of pointers
		DisableCapacities:       true, // don't spew capacities of collections
		ContinueOnMethod:        true, // recursion should continue once a custom error or Stringer interface is invoked
		SpewKeys:                true, // if unable to sort map keys then spew keys to strings and sort those
		MaxDepth:                4,    // maximum number of levels to descend into nested data structures.
	}
}

const (
	// asmdbX86DataJS filepath of x86data.js.
	asmdbX86DataJS = "asmdb/x86data.js"

	// asmdbArmDataJS filepath of armdata.js.
	asmdbArmDataJS = "asmdb/armdata.js"
)

var (
	//go:embed asmdb/x86data.js
	asmdbX86 embed.FS

	//go:embed asmdb/armdata.js
	asmdbArm embed.FS
)

func main() {
	if err := gen(); err != nil {
		log.Fatal(err)
	}
}

func gen() error {
	fsX86, err := asmdbX86.Open(asmdbX86DataJS)
	if err != nil {
		return fmt.Errorf("read %s embeded file: %w", asmdbX86DataJS, err)
	}
	defer fsX86.Close()

	x86AsmData, err := parse(fsX86)
	if err != nil {
		return fmt.Errorf("parse asmdb data: %w", err)
	}

	var x86Asm X86
	if err := json.Unmarshal(x86AsmData, &x86Asm); err != nil {
		return fmt.Errorf("unmarshal X86: %w", err)
	}
	instructions := x86Asm.Instructions // copy
	x86Asm.Instructions = nil

	fmt.Printf("x86asm: %s\n", spew.Sdump(x86Asm))

	insts := make([]X86Instruction, len(instructions))
	for i, inst := range instructions {
		// _ = inst[4] // BCE hint // TODO(zchee): still needs?
		insts[i].Name = inst[0]
		insts[i].Operands = inst[1]
		insts[i].Encoding = inst[2]
		insts[i].OpCode = inst[3]
		insts[i].Metadata = inst[4]
	}
	fmt.Printf("Instructions: %s\n", spew.Sdump(insts))

	return nil
}

const (
	// markJSONBegin is a magic comment that marks the beginning of the JSON data in the asmjit/asmdb JavaScript file.
	markJSONBegin = "// ${JSON:BEGIN}"

	// markJSONEnd is the magic comment that marks the end of the JSON data in the asmjit/asmdb JavaScript file.
	markJSONEnd = "// ${JSON:END}"
)

func parse(r io.Reader) (data []byte, err error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read r reader: %w", err)
	}

	// split buf by markJSONBegin magic comment
	splitted := bytes.SplitN(buf, []byte(markJSONBegin), 2)
	if len(splitted) <= 1 {
		return nil, fmt.Errorf("could not split asmdb data by %q magic comment: splitted length: %d", markJSONBegin, len(splitted))
	}

	data = splitted[1]
	if len(data) == 0 {
		return nil, fmt.Errorf("incorrect splitted asmdb data: data length: %d", len(data))
	}
	data = data[1:] // 1 means trim first newline

	// trim after the markJSONEnd magic comment
	idx := bytes.Index(data, []byte(markJSONEnd))
	if idx <= 0 {
		return nil, fmt.Errorf("could not find %q magic comment from asmdb data: %s", markJSONEnd, string(data))
	}
	data = data[:idx-1] // -1 means also trim end of newline

	return
}
