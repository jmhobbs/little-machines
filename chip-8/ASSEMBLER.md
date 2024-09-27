# chip-8 assembler

There is a very basic assembler for chip-8 in `asm/`, with a frontend in `cmd/asm/`

# Syntax

The general syntax is `<instruction> <operand>(, <operand>(, <operand>))`.

Operands can be a register, byte value, address or flag.

Capitalization does not matter, `load v0, 0x3f` is equivalent to `LOAD V0, 0X3F`.

Comments may come anywhere on a line, and begin with `;`

## Register

The main registers are represented by `VX`, where `X` is `0-F`.

The `I` register is just `I`.  Delay timer is `DT`, sound timer is `ST`.

## Byte

Byte values _must_ be hex encoded, and _may_ be prefixed by `0x`.

## Address

Address values _must_ be hex encoded, _must_ be three characters wide, and _may_ be prefixed by `0x`.

Left pad with `0` if required.

e.g. `0x001`, **not** `0x1`

## Flags

There are special case flags for keypress (`KEY`), BCD encoding (`BCD`), and sprite fonts (`FONT`).

These are limited use to a few instructions.

# Opcodes

Based on http://devernay.free.fr/hacks/chip8/C8TECH10.HTM, with changes to suit my aesthetic needs.

| Instruction | Operands                 | Description |
|-------------|--------------------------|-------------|
| ADD         | I, Register              | Add register value to I register. |
| ADD         | Register, Byte           | Add byte to register. |
| ADD         | Register, Register       | Add second register to first register, sets VF if carry. |
| AND         | Register, Register       | Bitwise AND, stored into first register. |
| CALL        | Address                  | Call subroutine at address. |
| CLEAR       |                          | Clear screen. |
| DRAW        | Register, Register, Byte | Draw n-byte sprite at (VX, VY) reading from memory address I. |
| JUMP        | Address                  | Jump to address. |
| JUMP        | V0, Address              | Jump to address + V0. |
| LOAD        | BCD, Register            | Store BCD representation of register in memory locations I, I+1, I+2. |
| LOAD        | Delay Timer, Register    | Load register value into delay timer. |
| LOAD        | FONT, Register           | Load font sprite address for digit in register into I register. |
| LOAD        | I, Address               | Load address into I register. |
| LOAD        | I, Register              | Store registers V0 through VX in memory starting at I. |
| LOAD        | Register, Byte           | Load byte into register. |
| LOAD        | Register, DT             | Load delay timer into register. |
| LOAD        | Register, I              | Read registers V0 through VX from memory starting at I. |
| LOAD        | Register, Keypress       | Wait for keypress, store in register. |
| LOAD        | Register, Register       | Load value of second register into first register. |
| LOAD        | Sound Timer, Register    | Load register value into sound timer. |
| OR          | Register, Register       | Bitwise OR, stored into first register. |
| RAND        | Register, Byte           | Load random byte AND byte into register. |
| RETURN      |                          | Return from subroutine. |
| SHIFTL      | Register                 | Shift left, sets VF to most significant bit. |
| SHIFTR      | Register                 | Shift right, sets VF to least significant bit. |
| SKIP        | Register, Byte           | Skip next instruction if register equals byte. |
| SKIP        | Register, Register       | Skip next instruction if registers are equal. |
| SKIPN       | Register, Byte           | Skip next instruction if register does not equal byte. |
| SKIPN       | Register, Register       | Skip next instruction if registers are not equal. |
| SKIPNP      | Register                 | Skip next instruction if key with value of register is not pressed. |
| SKIPP       | Register                 | Skip next instruction if key with value of register is pressed. |
| SUB         | Register, Register       | Subtract second register from first register, sets VF if no borrow. |
| SUBN        | Register, Register       | Subtract first register from second register, sets VF if no borrow. |
| SYS         | Address                  | Jump to machine code location, unused. |
| XOR         | Register, Register       | Bitwise XOR, stored into first register. |

# Example

```
; Clear the screen
CLEAR

; Set our x position in V0
LOAD V0, 0x2

; Set out y position in V1
LOAD V1, 0x2

; Draw the 5 font sprite at (2, 2)
DRAW V0, V1, 0x05

; Increment x position
ADD V0, 0x08

; Increment y position
ADD V1, V0

; Draw the sprite again, now at (10, 10)
DRAW V0, V1, 0x05
```
