package insp

import (
	"strings"
)

/**
 * @description: TODO
 * @author Wg
 * @date 2025/1/19
 */

type Node struct {
	Name     string
	SQL      string
	Children Map
}

type Tree struct {
	Roots Map
	Len   int
	All   []Node
}

type Map map[string]*Node

func (t *Tree) Get(path string) *Node {
	if len(path) == 0 {
		return nil
	}
	paths := strings.Split(path, ".")

	var cur *Node
	if val, ok := t.Roots[paths[0]]; ok {
		cur = val
	} else {
		return nil
	}

	for i := 1; i < len(paths); i++ {
		if val, ok := cur.Children[paths[i]]; ok {
			cur = val
		} else {
			return nil
		}
	}
	return cur
}
