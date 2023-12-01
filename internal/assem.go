package internal

import (
	"fmt"
	"strconv"
)

type Assem struct {
	Registers    map[string]int
	Instructions []Instruction
	InstIdx      int
	ExtraMemory  map[string]int
}

func NewAssem(instructions []Instruction) *Assem {
	return &Assem{
		Registers:    make(map[string]int),
		Instructions: instructions,
		InstIdx:      0,
		ExtraMemory:  make(map[string]int),
	}
}

func (a *Assem) AddInstruction(i Instruction) {
	a.Instructions = append(a.Instructions, i)
}

func (a *Assem) Run() bool {
	instr := a.Instructions[a.InstIdx]
	instr.Run(a)
	return a.InstIdx >= len(a.Instructions)
}

func (a *Assem) Exec(inst Instruction) {
	inst.Run(a)
}

func (a *Assem) GetValOrNum(v string) int {
	if v[0] >= 'a' && v[0] <= 'z' {
		return a.GetVal(v)
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		panic("weird register name: " + v)
	}
	return n
}

func (a *Assem) SetReg(reg string, val int) {
	a.Registers[reg] = val
}

func (a *Assem) GetVal(reg string) int {
	return a.Registers[reg]
}

type Instruction interface {
	Run(a *Assem)
}

type JumpNotZero struct {
	X, Y string
}

func (j JumpNotZero) Run(a *Assem) {
	v := a.GetValOrNum(j.X)
	if v == 0 {
		a.InstIdx++
		return
	}
	h := a.GetValOrNum(j.Y)
	a.InstIdx += h
}

func (j JumpNotZero) String() string {
	return fmt.Sprintf("jnz %v %v", j.X, j.Y)
}

type Increase struct {
	X string
}

func (i Increase) Run(a *Assem) {
	a.SetReg(i.X, a.GetVal(i.X)+1)
	a.InstIdx++
}
func (i Increase) String() string {
	return fmt.Sprintf("inc %v ", i.X)
}

type Decrease struct {
	X string
}

func (d Decrease) Run(a *Assem) {
	a.SetReg(d.X, a.GetVal(d.X)-1)
	a.InstIdx++
}
func (d Decrease) String() string {
	return fmt.Sprintf("dec %v ", d.X)
}

type Copy struct {
	X, Y string
}

func (c Copy) Run(a *Assem) {
	a.SetReg(c.Y, a.GetValOrNum(c.X))
	a.InstIdx++
}

func (c Copy) String() string {
	return fmt.Sprintf("cpy %v %v", c.X, c.Y)
}

type Noop struct{}

func (n Noop) Run(a *Assem) {}
func (n Noop) String() string {
	return "noop"
}

type AddX struct {
	X   string
	Val int
}

func (a AddX) Run(assem *Assem) {
	assem.SetReg(a.X, assem.GetVal(a.X)+a.Val)
}
func (a AddX) String() string {
	return fmt.Sprintf("addx %d", a.Val)
}
