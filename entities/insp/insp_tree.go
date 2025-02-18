package insp

import (
	"PgInspector/entities/config"
	"fmt"
	"sort"
	"strings"
)

/**
 * @description: insp的树
 * @author Wg
 * @date 2025/1/19
 */

// Node 分为Insp节点和索引节点，Insp节点也是叶子节点
type Node struct {
	Name     string
	SQL      string
	Children Map

	AlertID   config.ID
	AlertFunc func(Result) error //包括检查是否符合报警条件，并且发送报警
}

type Tree struct {
	Roots   Map
	Num     int     //Insp节点数量
	AllInsp []*Node //所有的Insp节点
}

type Map map[string]*Node

func NewTree() *Tree {
	return &Tree{
		Roots:   make(Map),
		AllInsp: []*Node{},
	}
}

// GetNode 获取单个子节点
func (t *Tree) GetNode(path string) *Node {
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

// Arr 获取Map的第一层子节点，会按照字典序排序
func (m Map) Arr() []*Node {
	keys := make([]string, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	res := make([]*Node, 0, len(m))
	for _, k := range keys {
		res = append(res, m[k])
	}
	return res
}

func (n *Node) IsInsp() bool {
	return n.SQL != ""
}

// GetAllInsp 获取该节点的所有insp节点，会寻找至叶子节点
func (n *Node) GetAllInsp() []*Node {
	var res, idx []*Node
	idx = n.Children.Arr()
	for len(idx) != 0 {
		var nextIdx []*Node
		for _, v := range idx {
			if v.IsInsp() {
				res = append(res, v)
			} else {
				nextIdx = append(nextIdx, v)
			}
		}
		idx = nextIdx
	}
	return res
}

func (n *Node) AddChild(node *Node) error {
	if n == nil {
		return fmt.Errorf("insp parent is nil")
	}
	if node == nil {
		return fmt.Errorf("new insp child is nil")
	}
	if n.Children == nil {
		n.Children = make(Map)
	}
	n.Children[node.Name] = node
	return nil
}

func (t *Tree) AddChild(path string, node *Node) error {
	if node == nil {
		return fmt.Errorf("Insp tree new child node is nil, path: %s\n", path)
	}
	if node.IsInsp() {
		t.AllInsp = append(t.AllInsp, node)
		t.Num++
	}
	//根目录层
	if path == "" {
		t.Roots[node.Name] = node
	}

	//子目录
	n := t.GetNode(path)
	if n == nil {
		return fmt.Errorf("path is not exist: %s\n", path)
	}
	if _, ok := n.Children[node.Name]; ok {
		return fmt.Errorf("node is already exist, path: %s, name: %s \n", path, node.Name)
	}
	return n.AddChild(node)
}
