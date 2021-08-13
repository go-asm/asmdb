// Copyright 2012 The Go Asm Authors
// SPDX-License-Identifier: BSD-3-Clause

package main

// x86data.js
//
// X86/X64 instruction-set data.
//
// License
//
// Public Domain.
//
//
// INSTRUCTIONS
//
// Each instruction definition consists of 5 strings:
//
//   [0] - Instruction name.
//   [1] - Instruction operands.
//   [2] - Instruction encoding.
//   [3] - Instruction opcode.
//   [4] - Instruction metadata - CPU features, FLAGS (read/write), and other metadata.
//
// The definition tries to match Intel and AMD instruction set manuals, but there
// are small differences to make the definition more informative and compact.
//
//
// OPERANDS
//
//   * "op"    - Explicit operand, must always be part of the instruction. If a fixed
//               register (like "cl") is used, it means that the instruction uses this
//               register implicitly, but it must be specified anyway.
//
//   * "<op>"  - Implicit operand - some assemblers allow implicit operands the be passed
//               explicitly for documenting purposes. And some assemblers like AsmJit's
//               Compiler infrastructure requires implicit operands to be passed explicitly
//               for register allocation purposes.
//
//   * "{op}"  - Optional operand. Mostly used by AVX_512:
//
//               - {k} mask selector.
//               - {z} zeroing.
//               - {1tox} broadcast.
//               - {er} embedded-rounding.
//               - {sae} suppress-all-exceptions.
//
//   * "?:Op"  - Each operand can provide metadata that can be used to describe which
//               operands are used as a destination, and which operands are source-only.
//               Each instruction in general assumes that the first operand is always
//               read/write and all following operands are read-only. However, this is
//               not correct for all instructions, thus, instructions that don't match
//               this assumption must provide additional information:
//
//               - "R:Op" - The operand is read-only.
//               - "w:Op" - The operand is write-only (does not zero-extend).
//               - "W:Op" - The operand is write-only (implicit zero-extend).
//               - "x:Op" - The operand is read/write (does not zero-extend).
//               - "X:Op" - The operand is read/write (implicit zero-extend).
//
//   * Op[A:B] - Optional bit-range that describes which bits are read and written.
//
//   * "~Op"   - Operand is commutative with other operands prefixed by "~". Commutativity
//               means that all operands marked by '~' can be swapped and the result of the
//               instruction would be the same.
//
// WHAT IS MISSING
//
// Here is a list of missing instructions to keep track of it:
//
// [ ] xlat/xlatb

// X86 represents a x86_x64 instruction set data.
type X86 struct {
	Architectures []string         `json:"architectures"`
	Extensions    []*X86Extension  `json:"extensions"`
	Attributes    []*X86Attribute  `json:"attributes"`
	SpecialRegs   []*X86SpecialReg `json:"specialRegs"`
	Shortcuts     []*X86Shortcut   `json:"shortcuts"`
	Register      *X86Register     `json:"registers"`
	Instructions  [][5]string      `json:"instructions,omitempty"`
}

// X86Extension represents a available extension, instruction can specify extension in metadata
type X86Extension struct {
	Name string `json:"name"`
}

// X86Attribute represents a available attribute, instruction can specify attribute in metadata.
type X86Attribute struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Doc  string `json:"doc"`
}

// X86SpecialReg represents a special registers (and their parts) that instructions can read/write to/from.
type X86SpecialReg struct {
	Name  string `json:"name"`
	Group string `json:"group"`
	Doc   string `json:"doc"`
}

// X86Shortcut represents a shortcuts that can be used inside instruction's metadata, these shortcuts then expand to the expand key.
type X86Shortcut struct {
	Name   string `json:"name"`
	Expand string `json:"expand"`
}

// X86Register represents a x86 processors general purpose registers.
type X86Register struct {
	Bnd  *X86RegisterData `json:"bnd"`
	Creg *X86RegisterData `json:"creg"`
	Dreg *X86RegisterData `json:"dreg"`
	K    *X86RegisterData `json:"k"`
	Mm   *X86RegisterData `json:"mm"`
	R16  *X86RegisterData `json:"r16"`
	R32  *X86RegisterData `json:"r32"`
	R64  *X86RegisterData `json:"r64"`
	R8   *X86RegisterData `json:"r8"`
	R8hi *X86RegisterData `json:"r8hi"`
	Rxx  *X86RegisterData `json:"rxx"`
	Sreg *X86RegisterData `json:"sreg"`
	St   *X86RegisterData `json:"st"`
	Tmm  *X86RegisterData `json:"tmm"`
	Xmm  *X86RegisterData `json:"xmm"`
	Ymm  *X86RegisterData `json:"ymm"`
	Zmm  *X86RegisterData `json:"zmm"`
}

// X86RegisterData represents ax86 processors general purpose registers data.
type X86RegisterData struct {
	Names []string `json:"names"`
	Kind  string   `json:"kind"`
	Any   string   `json:"any,omitempty"`
}

// X86Instruction represents a x86_x64 instruction set.
type X86Instruction struct {
	Name     string `json:"name"`
	Operands string `json:"operands,omitempty"`
	Encoding string `json:"encoding"`
	OpCode   string `json:"opcode"`
	Metadata string `json:"metadata"`
}
