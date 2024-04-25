package core

import (
	"fmt"
	"io"
	"log"
)

// key memory locations
const (
	zeroPageEnd = 0x00FF
	stackStart  = 0x0100
	nmiVector   = 0xFFFA
	rstVector   = 0xFFFC
	irqVector   = 0xFFFE
)

// interrupt types
const (
	nmi = iota
	rst
	irq
)

// extra cycle costs
const (
	branchSuccCycleCost = 1
	interruptCycleCost  = 7
)

// Status holds data for each status flag
// in the 6502 status register.
type status struct {
	c bool // Carry
	z bool // Zero result
	i bool // Interrupt disable
	d bool // Decimal mode
	b bool // Break command
	u bool // Unused (here for better logging)
	v bool // Overflow
	n bool // Zero result
}

// GenerateInterrupt causes cpu c to generate the interrupt
// specified by interruptType.
func (c *cpu) GenerateInterrupt(interruptType int) {
	c.mustHandleInterrupt = true
	c.interruptType = interruptType
}

// handleInterrupt causes cpu c to handle the interrupt specified by c.interruptType.
// Briefly, this consists of:
// 1. Push the program counter and status register on to the stack.
// 2. Set the interrupt disable flag to prevent further interrupts.
// 3. Load the address of the interrupt handling routine from the vector table into the program
// counter.
// See http://www.nesdev.com/NESDoc.pdf.
func (c *cpu) handleInterrupt() {
	c.mustHandleInterrupt = false

	// 1. Push PC and SP onto stack
	c.push16(c.pc)
	c.pushStack(c.status.asByte())

	// 2. Set interrupt disable flag
	c.status.i = true

	// 3. Load address of interrupt handling routine into PC from vector table
	switch c.interruptType {
	case nmi:
		c.pc = c.read16(nmiVector)
	case irq:
		c.pc = c.read16(irqVector)
	case rst:
		c.pc = c.read16(rstVector)
	default:
	}
}

// convert converts a boolean bit into its byte form.
func convert(bit bool) (num byte) {
	if bit {
		return 1
	}
	return 0
}

// String implements Stringer.
func (sr *status) String() (repr string) {
	return fmt.Sprintf("%X", sr.asByte())
}

// asByte returns the status register in byte format.
func (sr *status) asByte() (data byte) {
	var status byte = 0x00
	status |= convert(sr.c)
	status |= convert(sr.z) << 1
	status |= convert(sr.i) << 2
	status |= convert(sr.d) << 3
	status |= convert(sr.b) << 4
	status |= convert(sr.u) << 5
	status |= convert(sr.v) << 6
	status |= convert(sr.n) << 7
	return status
}

// fromByte sets the status register according to the contents of a byte of data.
func (sr *status) fromByte(data byte) {
	sr.c = data&mask0 != 0
	sr.z = data&mask1 != 0
	sr.i = data&mask2 != 0
	sr.d = data&mask3 != 0
	sr.b = data&mask4 != 0
	sr.u = data&mask5 != 0
	sr.v = data&mask6 != 0
	sr.n = data&mask7 != 0
}

// setZN sets both the zero flag and negative flag of sr
// according to the contents of reg.
func (sr *status) setZN(reg byte) {
	sr.z = reg == 0
	sr.n = reg&mask7 != 0
}

// Registers holds data for each register
// used by the 6502.
type registers struct {
	// Special purpose registers
	status *status // Status register
	pc     uint16  // Program counter
	sp     byte    // Stack pointer

	// General purpose registers
	a byte // Accumulator register
	x byte // Index register X
	y byte // Index register Y
}

// String implements Stringer.
func (r *registers) String() (repr string) {
	return fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02s SP:%02X", r.a, r.x, r.y, r.status, r.sp)
}

// cpu represents to 6502 and its associated registers and memory map.
// This should be declared and used as a singleton during emulator execution.
type cpu struct {
	*memory    // Pointer to main memory
	*registers // Set of registers

	instructions        map[byte]instruction // Instructions available to cpu
	cycles              int                  // Number of cpu cycles
	pageCrossed         bool                 // Whether or not the most recently executed instruction crossed a page
	branchSucceeded     bool                 // Whether or not the most recently executed branch instruction succeeded
	mustHandleInterrupt bool                 // Whether or not the cpu must handle an interrupt on its next step
	interruptType       int                  // The type of interrupt that must be handled (assuming cpu is interrupted)

	// For logging only
	debug  bool      // Whether or not to output logs
	logger io.Writer // Writer through which to output logs
}

// New initializes a new 6502 cpu with all status bits, register, and memory
// initialized to zero. Memory is the shared memory that the cpu will access.
func newCpu() (c *cpu) {
	cpu := &cpu{
		registers: &registers{
			status: &status{},
		},
	}
	return cpu
}

// OutputTo sets the cpu to log its execution to io.Writer w.
func (c *cpu) OutputTo(w io.Writer) {
	c.debug = true
	c.logger = w
}

// UseMemory associates the cpu c with main memory m.
func (c *cpu) UseMemory(m *memory) {
	c.memory = m
}

// Init implements nes.Component.
// See https://wiki.nesdev.com/w/index.php/cpu_power_up_state#cite_note-1.
// TODO: Make robust
func (c *cpu) Init() {
	c.initInstructions()
	c.status.i = true
	c.status.u = true
	c.sp = 0xFD
	// TODO: APU start-up state
}

// Clear sets every register in c (including PC, SP, and status) to 0x00.
// Retains memory linked through UseMemory.
// TODO: Make robust
func (c *cpu) Clear() {
	*c.registers = registers{
		status: &status{},
	}
}

// setPageCrossed sets the cpu to whether or not a page has been crossed
// according to the current PC and the provided address.
func (c *cpu) setPageCrossed(address uint16) {
	// TODO: implement page cross logic
	c.pageCrossed = false
}

// branchTo branches the cpu program counter to address.
// This function counts cycles correctly.
func (c *cpu) branchTo(address uint16) {
	c.setPageCrossed(address)
	c.branchSucceeded = true
	c.pc = address
}

// pushStack pushes a byte of data onto the stack.
// The stack pointer always points to the next free location on the stack.
func (c *cpu) pushStack(data byte) {
	c.write(stackStart+uint16(c.sp), data)
	c.sp--
}

// push16 pushes a 16 byte word onto the stack, low byte and then high byte.
func (c *cpu) push16(word uint16) {
	lo := byte(word & 0x00FF)
	hi := byte(word & 0xFF00)
	c.pushStack(lo)
	c.pushStack(hi)
}

// pullStack pulls a byte of data from the stack.
// The stack pointer always points to the next free location on the stack.
func (c *cpu) pullStack() (data byte) {
	c.sp++
	return c.Read(stackStart + uint16(c.sp))
}

// pull16 pulls a 16 byte word from the stack, high byte and then low byte.
func (c *cpu) pull16() (word uint16) {
	hi := uint16(c.pullStack())
	lo := uint16(c.pullStack())
	return hi<<8 | lo
}

// decode decodes opcode opcode and returns relevant information.
func (c *cpu) decode(opcode byte) (name string, addressingMode, byteCost, cycleCost, pageCrossCost int, execute func(uint16), err error) {
	if instruction, ok := c.instructions[opcode]; ok {
		return instruction.name,
			instruction.addressingMode,
			instruction.byteCost,
			instruction.cycleCost,
			instruction.pageCrossCycleCost,
			instruction.execute,
			nil
	}
	return "", 0, 0, 0, 0, nil, ErrInvalidOpcode(opcode)
}

// GetAddressWithMode uses addressing mode addressingMode to get
// an address on which any instruction can execute.
// Must be used when c.PC is on an opcode address, otherwise
// the following addresses will be interpreted incorrectly.
func (c *cpu) getAddressWithMode(addressingMode int) (addr uint16) {
	switch addressingMode {
	case modeImplied:
		// Address will be unused for following two addressing modes; return 0
		fallthrough

	case modeAccumulator:
		return 0

	case modeRelative:
		// Instructions with modeRelative take 2 bytes:
		// 1. opcode
		// 2. 8 bit constant value
		// The address will only be jumped to if the branch succeeeds.
		// Note: relative addressing uses twos complement to branch both
		// forwards and backwards.
		offset := uint16(c.Read(c.pc + 1))
		if offset >= 0x80 {
			// interpret as negative number
			return c.pc + offset - 0x100
		}
		return c.pc + offset

	case modeImmediate:
		// Instructions with modeImmediate take 2 bytes:
		// 1. opcode
		// 2. 8 bit constant value
		return c.pc + 1

	case modeZeroPage:
		// Instructions with modeZeroPage take 2 bytes:
		// 1. opcode
		// 2. zero-page address
		return uint16(c.Read(c.pc + 1))

	case modeZeroPageX:
		// Same as modeZeroPage, but with zero page address being added to X register with wraparound
		return uint16(c.Read(c.pc+1)+c.x) & zeroPageEnd

	case modeZeroPageY:
		// Same as modeZeroPage, but with zero page address being added to Y register with wraparound
		return uint16(c.Read(c.pc+1)+c.y) & zeroPageEnd

	case modeAbsolute:
		// Instructions with modeAbsolute take 3 bytes:
		// 1. opcode
		// 2. least significant byte of address
		// 3. most significant byte of address
		return c.read16(c.pc + 1)

	case modeAbsoluteX:
		// Same as modeAbsolute, with address being added to contents of X register
		return c.read16(c.pc+1) + uint16(c.x)

	case modeAbsoluteY:
		// Same as modeAbsolute, with address being added to contents of Y register
		return c.read16(c.pc+1) + uint16(c.y)

	case modeIndirect:
		// Instructions with modeIndirect take 3 bytes:
		// 1. opcode
		// 2. least significant byte of address
		// 3. most significant byte of address
		// The formulated address, along with the next,
		// are then accessed again to get the final address.
		return c.read16(c.read16(c.pc + 1))

	case modeIndirectX:
		// Instructions with modeIndirectX take 2 bytes:
		// 1. opcode
		// 2. single byte
		// The byte is then added to the X register, which then
		// gives the least significant byte of the target address.
		return c.read16(uint16(c.Read(c.pc+1) + c.x))

	case modeIndirectY:
		// Instructions with modeIndirectY take 2 bytes:
		// 1. opcode
		// 2. least significant byte of zero page address
		// The zero page address is then accessed, and the data
		// is added to the Y register. The resulting data is the
		// target address.
		return c.read16(uint16(c.Read(c.pc+1))) + uint16(c.y)

	default:
		// shouldn't happen, but handle gracefully
		return 0
	}
}

// Step performs a single step of the cpu.
// Briefly, this consists of:
// 0. Handling any interrupts (i.e. loading PC with interrupt handling routine if needed)
// 1. Retrieving the opcode at current PC.
// 2. Decoding the opcode.
// 3. Incrementing the program counter by the correct amount.
// 4. Performing the instruction. This is done after (3) because jump instructions may directly change the PC.
// 5. Add cpu cycles based on instruction execution.
func (c *cpu) step() {
	// Reset instruction-wise flags
	c.pageCrossed = false
	c.branchSucceeded = false

	// 0. Handle any interrupts
	if c.mustHandleInterrupt {
		c.handleInterrupt()
		c.cycles += interruptCycleCost
	}

	// 1. Retrieve opcode at current PC
	opcode := c.Read(c.pc)

	// 2. Decode opcode
	name, addressingMode, byteCost, cycleCost, pageCrossCycleCost, execute, err := c.decode(opcode)
	if IsInvalidOpcodeErr(err) {
		// If the opcode is invalid, shut down everything for now.
		log.Fatalln(err)
		return
	}
	instructionAddress := c.getAddressWithMode(addressingMode)

	// 2.5. Log cpu execution if necessary
	// Format according to ideal nestest log
	if c.debug {
		// Retrieve the raw next bytes used for this instruction.  Used purely for logging.
		var nextBytes []byte
		for i := 0; i < byteCost; i++ {
			nextBytes = append(nextBytes, c.Read(c.pc+uint16(i)))
		}

		// Form a trace (a line of logs for this single instruction including status state, cycles, all registers, etc.)
		trace := fmt.Sprintf("%-6X% -10X%-7s%s CYC:%d\n", c.pc, nextBytes, name, c.registers, c.cycles)
		io.WriteString(c.logger, trace)
	}

	// 3. Increment program counter
	c.pc += uint16(byteCost)

	// 4. Perform instruction
	// Done after (3) because some instructions (relative addressing) will directly change the PC.
	execute(instructionAddress)

	// 5. Add cpu cycles based on instruction execution.
	// cpu cycles are used to keep the CPU in sync with other modules (like the PPU).
	c.cycles += cycleCost
	if c.branchSucceeded {
		c.cycles += branchSuccCycleCost
	}
	if c.pageCrossed {
		c.cycles += pageCrossCycleCost
	}
}
