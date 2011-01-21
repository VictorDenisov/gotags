//  gotags generates a tags file for the Go Programming Language in the format used by exuberant-ctags.
//  Copyright (C) 2009  Michael R. Elkins <me@sigpipe.org>
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Date: 2009-11-12
// Updated: 2010-08-21

//
// usage: gotags filename [ filename... ] > tags
//

package main

import (
	"container/vector"
	"fmt"
	"go/ast"
	"go/parser"
	"os"
	"sort"
)

var (
	tags vector.StringVector
)

func output_tag(name *ast.Ident, kind byte) {
	tags.Push(fmt.Sprintf("%s\t%s\t%d;\"\t%c",
		name.Name, name.Position.Filename, name.Position.Line, kind))
}

func main() {
	parse_files()

	println("!_TAG_FILE_SORTED\t1\t")
	sort.SortStrings(tags)
	for _, s := range tags {
		println(s)
	}
}

const FUNC, TYPE, VAR = 'f', 't', 'v'

func parse_files() {
	for i, m := 1, len(os.Args); i < m; i++ {
		tree, ok := parser.ParseFile(os.Args[i], nil, 0)
		if ok != nil {
			println("error parsing file", os.Args[i], ok.String())
			panic(nil)
		}

		for _, node := range tree.Decls {
			switch n := node.(type) {
			case *ast.FuncDecl:
				output_tag(n.Name, FUNC)
			case *ast.GenDecl:
				do_gen_decl(n)
			}
		}
	}

}

func do_gen_decl(node *ast.GenDecl) {
	for _, v := range node.Specs {
		switch n := v.(type) {
		case *ast.TypeSpec:
			output_tag(n.Name, TYPE)

		case *ast.ValueSpec:
			for _, vv := range n.Names {
				output_tag(vv, VAR)
			}
		}
	}
}
