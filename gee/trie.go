package gee

import (
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由，eg: /p/:lang
	part     string  // 路由中的一部分，eg: :lang
	children []*node // 子节点, eg: [doc, tutorial, intro]
	isWild   bool    // 是否精准匹配, part含有 : 或 * 时为true
}

// matchChild 第一个匹配成功的节点, 用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	// 递归回溯条件
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	// 如果child不存在，就执行插入操作
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 存在之后指针移到child节点，递归在trie树上插入
	child.insert(pattern, parts, height+1)
}

// search 递归查询结果节点 yes return result node, or no return nil
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 表示匹配失败
		if n.pattern == "" {
			return nil
		}
		// 匹配成功
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	// 可能有多个匹配成功的路由节点，都要dfs搜索一遍
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
