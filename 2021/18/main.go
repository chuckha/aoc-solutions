package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	mag := 0
	for i := 0; i < len(lines); i++ {
		for j := i; j < len(lines); j++ {
			s := newSnapfish(lines[i]).add(newSnapfish(lines[j]))
			s.reduce()
			m := s.magnitude()
			if m > mag {
				mag = m
			}
		}
	}
	for i := 0; i < len(lines); i++ {
		for j := i; j < len(lines); j++ {
			s := newSnapfish(lines[j]).add(newSnapfish(lines[i]))
			s.reduce()
			m := s.magnitude()
			if m > mag {
				mag = m
			}
		}
	}
	fmt.Println(mag)
	//fmt.Println(newSnapfish(lines[0]))
	//	s := newSnapfish(lines[0])
	// [[1,2],3]

	//	fmt.Println("[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]")
	//	s := newSnapfish("[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]")
	//	fmt.Println(s.needsExplode(0).right.rightNumber())
	///	s.explode()

	// example
	// s := newSnapfish("[[[[4,3],4],4],[7,[[8,4],9]]]")
	// a := newSnapfish("[1,1]")
	// s = s.add(a)
	// for {
	// 	//		fmt.Println(s)
	// 	if s.needsExplode(0) == nil && s.needsSplit() == nil {
	// 		break
	// 	}
	// 	if s.needsExplode(0) != nil {
	// 		s.explode()
	// 		continue
	// 	}
	// 	if s.needsSplit() != nil {
	// 		s.split()
	// 	}
	// }
	// fmt.Println(s)
	//	fmt.Println("[[[[0,7],4],[[7,8],[0,[6,7]]]],[1,1]]")
	//	fmt.Println("[[[[0,7],4],[[7,8],[6,0]]],[8,1]]")
}

type snapfish struct {
	data   int
	parent *snapfish
	left   *snapfish
	right  *snapfish
}

func newSnapfish(input string) *snapfish {
	//	fmt.Printf("sn(%q)\n", input)
	if input == "" {
		return nil
	}
	if len(input) == 1 {
		d, _ := strconv.Atoi(string(input[0]))
		return &snapfish{data: d}
	}
	if input[0] != '[' {
		parts := strings.SplitN(input, ",", 2)
		l, _ := strconv.Atoi(parts[0])
		if len(parts) == 1 {
			return &snapfish{data: l}
		}
		parent := &snapfish{data: -1}
		parent.setLeft(&snapfish{data: l})
		parent.setRight(newSnapfish(parts[1]))
		return parent
	}
	idx := matchingParenIndex(input)
	if idx == len(input)-1 {
		return newSnapfish(input[1:idx])
	}
	left := newSnapfish(input[1:idx])
	var right *snapfish
	right = newSnapfish(input[idx+2:])
	parent := &snapfish{data: -1}
	parent.setLeft(left)
	parent.setRight(right)
	return parent
}

func (s *snapfish) reduce() {
	for {
		//		fmt.Println("needs explode", s.needsExplode(1), "needs split", s.needsSplit())
		if s.needsExplode(1) == nil && s.needsSplit() == nil {
			break
		}
		if s.needsExplode(1) != nil {
			s.explode()
			continue
		}
		if s.needsSplit() != nil {
			s.split()
		}
	}
}

func (s *snapfish) setLeft(sn *snapfish) {
	if sn == nil {
		return
	}
	s.left = sn
	sn.parent = s
}

func (s *snapfish) setRight(sn *snapfish) {
	if sn == nil {
		return
	}
	s.right = sn
	sn.parent = s
}

func (s *snapfish) String() string {
	if s == nil {
		return ""
	}
	if s.isSnapfish() {
		return "[" + s.left.String() + "," + s.right.String() + "]"
	}
	return fmt.Sprintf("%d", s.data)
	// prefix notation
	// pass a snapfish node on the left: add [
	// pass a value node on the left: add #
	// pass a snapfish value on the bottom: add ,
	// pass a snapfish node on the right: add ]

}

func (s *snapfish) isSnapfish() bool {
	return s.data == -1
}

func (s *snapfish) add(sn *snapfish) *snapfish {
	root := &snapfish{data: -1}
	root.setLeft(s)
	root.setRight(sn)
	return root
}

func (s *snapfish) needsExplode(i int) *snapfish {
	//	fmt.Println(s, i)
	if s.isSnapfish() && i > 4 {
		return s
	}
	var leftE *snapfish
	if s.left.isSnapfish() {
		leftE = s.left.needsExplode(i + 1)
	}
	if leftE != nil {
		return leftE
	}
	var rightE *snapfish
	if s.right.isSnapfish() {
		rightE = s.right.needsExplode(i + 1)
	}
	if rightE != nil {
		return rightE
	}
	return nil
}

/*

x  a,[6,7],z   y
x+6 0 y+7

[6, 7]

*/

func (s *snapfish) explode() {
	nodeToExplode := s.needsExplode(1)
	if nodeToExplode == nil {
		return
	}
	// fmt.Println("exploding this node", nodeToExplode)
	// fmt.Println("parent of exploding node", nodeToExplode.parent)
	leftNum := nodeToExplode.left.leftNumber()
	// fmt.Println("left number", leftNum)
	if leftNum != nil {
		// 	p.setLeft(&snapfish{data: 0})
		// } else {
		leftNum.data += nodeToExplode.left.data
		//		p.setLeft(&snapfish{data: nodeToExplode.left.data + leftNum.data})
	}
	rightNum := nodeToExplode.right.rightNumber()
	//	fmt.Println("right number", rightNum)
	if rightNum != nil {
		// 	p.setRight(&snapfish{data: 0})
		// } else {
		rightNum.data += nodeToExplode.right.data
		//		p.setRight(&snapfish{data: nodeToExplode.right.data + rightNum.data})
	}
	zero := &snapfish{
		data: 0,
	}
	if nodeToExplode.parent.left == nodeToExplode {
		nodeToExplode.parent.setLeft(zero)
		return
	}
	nodeToExplode.parent.setRight(zero)
}

func (s *snapfish) needsSplit() *snapfish {
	if s == nil {
		return nil
	}
	if s.data > 9 {
		return s
	}
	l := s.left.needsSplit()
	if l != nil {
		return l
	}
	r := s.right.needsSplit()
	if r != nil {
		return r
	}
	return nil
}

func (s *snapfish) split() {
	nodeToSplit := s.needsSplit()
	if nodeToSplit == nil {
		return
	}

	left := nodeToSplit.data / 2
	right := nodeToSplit.data - left
	nodeToSplit.setLeft(&snapfish{data: left})
	nodeToSplit.setRight(&snapfish{data: right})
	nodeToSplit.data = -1
}

func (s *snapfish) leftNumber() *snapfish {
	// go parent until you can make a left or you get to root node
	// if you make it to root node and you can't go left (becasue left == whre you just came from)
	// then return nil
	// if you can at some point go left, then go as far right as you can
	if s.parent == nil {
		return nil
	}
	p := s.parent
	cur := s
	for p.left == cur {
		if p.parent == nil {
			return nil
		}
		p = p.parent
		cur = cur.parent
	}
	p = p.left
	for p.data == -1 {
		if p.right.data != -1 {
			return p.right
		}
		p = p.right
	}
	return p
}
func (s *snapfish) rightNumber() *snapfish {
	if s.parent == nil {
		return nil
	}
	p := s.parent
	cur := s
	for p.right == cur {
		if p.parent == nil {
			return nil
		}
		p = p.parent
		cur = cur.parent
	}
	p = p.right
	for p.data == -1 {
		if p.left.data != -1 {
			return p.left
		}
		p = p.left
	}
	return p
}

func matchingParenIndex(input string) int {
	s := internal.NewStack[byte]()
	if input[0] != '[' {
		return -1 // no brakets
	}
	count := 0
	s.Push(input[0])
	for !s.Empty() {
		switch input[count+1] {
		case '[':
			count++
			s.Push('[')
		case ']':
			count++
			s.Pop()
		default:
			count++
		}
	}
	return count
}

func (s *snapfish) magnitude() int {
	if !s.isSnapfish() {
		return s.data
	}
	return s.left.magnitude()*3 + s.right.magnitude()*2
}
