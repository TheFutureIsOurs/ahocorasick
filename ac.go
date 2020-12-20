/*
 * @Author: Daiming Liu (xingrufeng)
 * Copyright (C) Daiming Liu (xingrufeng)
 */

package ahocorasick

import (
	"bufio"
	"errors"
	"os"
	"sort"
	"strings"
)

const (
	initSize  int = 64
	rootIndex int = 0
	rootBase  int = 1
	failState int = -1
)

// Ac result shape of AhoCorasick
type Ac struct {
	doubleArrayTrie
	fail   []int
	output []int // maxLength of suffix
}

// doubleArrayTrie the AhoCorasick's base implication
type doubleArrayTrie struct {
	base  []int
	check []int
}

// readLine read file line by line and drop one line's length >4096
func readLine(r *bufio.Reader) (string, error) {
	line, isprefix, err := r.ReadLine()
	// drop
	for isprefix && err == nil {
		_, isprefix, err = r.ReadLine()
	}
	return string(line), err
}

// BuildFromFile build ac from file
func BuildFromFile(inputfile string) (*Ac, error) {
	file, err := os.Open(inputfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bufReader := bufio.NewReader(file)
	keywords := [][]rune{}
	for {
		line, err := readLine(bufReader)
		if err != nil {
			break
		}
		keyword := strings.TrimSpace(line)
		if keyword == "" {
			continue
		}
		keywords = append(keywords, []rune(keyword))
	}
	if len(keywords) == 0 {
		return nil, errors.New("Empty keywords to build")
	}
	ac := &Ac{}
	ac.buildTrie(keywords)
	return ac, nil
}

// Build a ahocorasick based on double array trie
func Build(keywords []string) (*Ac, error) {
	if len(keywords) == 0 {
		return nil, errors.New("Empty keywords to build")
	}
	kws := make([][]rune, len(keywords))
	for k, v := range keywords {
		kws[k] = []rune(v)
	}
	ac := &Ac{}
	ac.buildTrie(kws)
	return ac, nil
}

// node a node of tree
type node struct {
	code                            rune
	depth, base, index, left, right int
	term                            bool // check is the leaf node
	children                        []*node
}

// tree.
type tree struct {
	root *node
}

type dartsKey []rune
type dartsKeySlice []dartsKey

func (k dartsKeySlice) Len() int {
	return len(k)
}

func (k dartsKeySlice) Less(i, j int) bool {
	var l int
	if len(k[i]) < len(k[j]) {
		l = len(k[i])
	} else {
		l = len(k[j])
	}

	for m := 0; m < l; m++ {
		if k[i][m] < k[j][m] {
			return true
		} else if k[i][m] == k[j][m] {
			continue
		} else {
			return false
		}
	}
	if len(k[i]) < len(k[j]) {
		return true
	}
	return false
}

func (k dartsKeySlice) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

type dartsBuild struct {
	dat          *doubleArrayTrie
	tree         *tree
	fail         []int
	output       []int
	keys         dartsKeySlice
	nextCheckPos int
	used         []bool
}

func (darts *dartsBuild) resize(newSize int) {
	darts.dat.base = append(darts.dat.base, make([]int, newSize-len(darts.dat.base))...)
	darts.dat.check = append(darts.dat.check, make([]int, newSize-len(darts.dat.check))...)
	darts.fail = append(darts.fail, make([]int, newSize-len(darts.fail))...)
	darts.output = append(darts.output, make([]int, newSize-len(darts.output))...)
}

// getChildren get the parent's children from dartsBuild
func (darts *dartsBuild) getChildren(parent *node) []*node {
	children := []*node{}
	var prev rune = 0
	for i := parent.left; i < parent.right; i++ {
		var cur rune = 0
		if len(darts.keys[i]) != parent.depth {
			cur = darts.keys[i][parent.depth]
		} else {
			parent.term = true
		}

		if cur != prev {
			tmpNode := &node{
				depth: parent.depth + 1,
				code:  cur,
				left:  i,
				base:  parent.base,
				term:  false,
			}
			if len(children) != 0 {
				children[len(children)-1].right = i
				parent.children[len(parent.children)-1].right = i
			}
			children = append(children, tmpNode)
			parent.children = append(parent.children, tmpNode)
		}
		prev = cur
	}

	if len(children) != 0 {
		children[len(children)-1].right = parent.right
		parent.children[len(children)-1].right = parent.right
	}

	return children
}

func max(m, n int) int {
	if m > n {
		return m
	}
	return n
}

// getBegin return the no use begin to fill parent's children
func (darts *dartsBuild) getBegin(parent *node) int {
	begin := 0
	pos := max(int(parent.children[0].code), darts.nextCheckPos)
	first := false
	for {
	next:
		pos++
		if len(darts.dat.base) <= pos {
			darts.resize(pos + initSize)
		}
		if 0 != darts.dat.base[pos] {
			continue
		} else if !first {
			darts.nextCheckPos = pos
			first = true
		}
		begin = pos - int(parent.children[0].code)
		if len(darts.dat.base) <= (begin + int(parent.children[len(parent.children)-1].code)) {
			darts.resize(begin + int(parent.children[len(parent.children)-1].code) + initSize)
		}

		for _, v := range parent.children {
			index := begin + int(v.code)
			if 0 != darts.dat.base[index] {
				goto next
			}
		}
		break
	}
	return begin
}

// setBC set base and check
func (darts *dartsBuild) setBC(parent *node) {
	if len(parent.children) == 0 {
		begin := parent.base
		darts.dat.base[parent.index] = -begin
	} else {
		begin := 0
		if parent.depth == 0 {
			begin = parent.base
		} else {
			begin = darts.getBegin(parent)
			parent.base = begin
		}

		if parent.term {
			darts.dat.base[parent.index] = -begin
		} else {
			darts.dat.base[parent.index] = begin
		}
		for _, v := range parent.children {
			pos := begin + int(v.code)
			v.index = pos
			v.base = begin
			if len(darts.dat.base) <= pos {
				darts.resize(pos + initSize)
			}
			darts.dat.base[pos] = begin
			darts.dat.check[pos] = parent.index
		}
	}
}

func getAbs(n int) int {
	if n >= 0 {
		return n
	}
	return -n
}

// getState give a inState, output index
func (dat *doubleArrayTrie) getState(inState int, code rune) int {
	b := getAbs(dat.base[inState])
	p := b + int(code)
	if p >= len(dat.base) {
		if inState == rootIndex {
			return rootIndex
		}
		return failState
	}

	if dat.base[p] != 0 && inState == dat.check[p] {
		return p
	}
	if inState == rootIndex {
		return rootIndex
	}
	return failState
}

// buildTrie build trie what we need
func (ac *Ac) buildTrie(keywords [][]rune) {
	darts := &dartsBuild{}
	// the length we know is equal to keywords
	darts.keys = make(dartsKeySlice, len(keywords))

	for k, v := range keywords {
		darts.keys[k] = v
	}
	sort.Sort(darts.keys)

	darts.dat = &doubleArrayTrie{}
	darts.resize(initSize)

	darts.tree = &tree{}
	darts.tree.root = &node{
		depth: 0,
		left:  0,
		right: len(darts.keys),
		base:  rootBase,
		index: rootIndex,
		term:  false,
	}

	queue := []*node{darts.tree.root}

	for len(queue) != 0 {
		node := queue[0]
		queue = queue[1:]

		children := darts.getChildren(node)
		if len(children) != 0 {
			queue = append(queue, children...)
		}
		darts.setBC(node)

		if node.term {
			darts.output[node.index] = len(darts.keys[node.left])
		}

		if node.depth == 0 || node.depth == 1 {
			darts.fail[node.index] = rootIndex
			continue
		}
		pIndex := darts.dat.check[node.index]
		inState := darts.fail[pIndex]
	set_state:
		outState := darts.dat.getState(inState, node.code)
		if outState == failState {
			inState = darts.fail[inState]
			goto set_state
		}
		if value := darts.output[outState]; value != 0 && value > darts.output[node.index] {
			darts.output[node.index] = value
		}
		darts.fail[node.index] = outState
	}
	ac.base = darts.dat.base
	ac.check = darts.dat.check
	ac.fail = darts.fail
	ac.output = darts.output
}

// Hit the result hit
type Hit struct {
	Begin int
	End   int
	Value []rune
}

// MultiPatternSearch return all find begin end and value
func (ac *Ac) MultiPatternSearch(content []rune) []Hit {
	hits := []Hit{}
	state := rootIndex
	for k, v := range content {
	start:
		if ac.getState(state, v) == failState {
			state = ac.fail[state]
			goto start
		} else {
			state = ac.getState(state, v)
			if val := ac.output[state]; val != 0 {
				hit := Hit{
					Begin: k - val + 1,
					End:   k,
					Value: content[k-val+1 : k+1],
				}
				hits = append(hits, hit)
			}
		}
	}
	return hits
}

// MultiPatternIndexes return the all find indexes of the content
func (ac *Ac) MultiPatternIndexes(content []rune) []int {
	hits := []int{}
	state := rootIndex
	for k, v := range content {
	start:
		if ac.getState(state, v) == failState {
			state = ac.fail[state]
			goto start
		} else {
			state = ac.getState(state, v)
			if val := ac.output[state]; val != 0 {
				hits = append(hits, k-val+1)
			}
		}
	}
	return hits
}

// MultiPatternHit return is the content is hit dictionary,
// it will return when first find
func (ac *Ac) MultiPatternHit(content []rune) bool {
	state := rootIndex
	for _, v := range content {
	start:
		if ac.getState(state, v) == failState {
			state = ac.fail[state]
			goto start
		} else {
			state = ac.getState(state, v)
			if ac.output[state] != 0 {
				return true
			}
		}
	}
	return false
}
