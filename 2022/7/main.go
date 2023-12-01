package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

const avail = 70000000
const need = 30000000

func main() {
	input := internal.ReadInput()
	root := newDir("/")
	fs := root
	for _, line := range input {
		if strings.HasPrefix(line, "$") {
			cmd := strings.TrimPrefix(line, "$ ")
			if strings.HasPrefix(cmd, "cd") {
				parts := strings.Split(cmd, " ")
				if parts[1] == ".." {
					fs = fs.parent
					continue
				}
				if parts[1] == "/" {
					fs = root
					continue
				}
				fs = fs.entries[parts[1]]
			}
			continue
		}
		if strings.HasPrefix(line, "dir") {
			parts := strings.Split(line, " ")
			fs.add(newDir(parts[1]))
			continue
		}
		// it's a file
		parts := strings.Split(line, " ")
		size, _ := strconv.Atoi(parts[0])
		fs.add(newFile(parts[1], size))
	}

	dirsToDelete := []dirdata{}
	root.walk(func(f *file) {
		if !f.mode.IsDir() {
			return
		}
		ts := f.totalSize()
		dirsToDelete = append(dirsToDelete, dirdata{name: f.name, size: ts})
	})
	sort.Sort(dd(dirsToDelete))
	// for _, d := range dirsToDelete {
	// 	fmt.Println(d)
	// }
	used := root.totalSize()
	free := avail - used
	fmt.Println(free)
	for _, x := range dirsToDelete {
		fmt.Println(free + x.size)
		if free+x.size >= need {
			fmt.Println(x.name, x.size)
			break
		}
	}
}

type dirdata struct {
	name string
	size int
}
type dd []dirdata

func (d dd) Len() int           { return len(d) }
func (d dd) Less(i, j int) bool { return d[i].size < d[j].size }
func (d dd) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func newDir(name string) *file {
	return &file{
		name:    name,
		mode:    os.ModeDir,
		entries: make(map[string]*file),
	}
}

func newFile(name string, size int) *file {
	return &file{
		name: name,
		size: size,
	}
}

type file struct {
	name    string
	parent  *file
	depth   int
	mode    os.FileMode
	size    int
	entries map[string]*file
}

func (f *file) walk(fn func(*file)) {
	for _, e := range f.entries {
		if e.mode.IsDir() {
			fn(e)
			e.walk(fn)
			continue
		}
		fn(e)
	}
}

func (f *file) String() string {
	out := f.name
	for _, e := range f.entries {
		out += "\n"
		tabs := strings.Repeat("\t", e.depth)
		if e.mode.IsDir() {
			out += tabs + e.String()
			continue
		}
		out += fmt.Sprintf("%s%s - %d\n", tabs, e.name, e.size)
	}
	return out
}

func (f *file) add(f2 *file) {
	if !f.mode.IsDir() {
		panic("cannot add a file to a non dir file")
	}
	f2.parent = f
	f2.depth = f.depth + 1
	f.entries[f2.name] = f2
}

func (f *file) totalSize() int {
	if !f.mode.IsDir() {
		return f.size
	}
	sum := 0
	for _, entry := range f.entries {
		if entry.mode.IsDir() {
			sum += entry.totalSize()
			continue
		}
		sum += entry.size
	}
	return sum
}
