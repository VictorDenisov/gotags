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
// Updated: 2010-11-21 Tomas Heran <tomas.heran@gmail.com>

//
// usage: gotags filename [ filename... ] > tags
//

package main

import (
	"container/vector"
	"fmt"
	"go/ast"
	"go/token"
	"go/parser"
	"os"
	"sort"
	//"reflect"
)

var (
	tags vector.StringVector
)

func output_tag(fset *token.FileSet, name *ast.Ident, kind byte) {
	switch kind {
	case PKG:
		//fmt.Printf("!_DEBUG\t%#v - %c\n", name, kind)
		tags.Push(fmt.Sprintf("%s\t%s\t/\\m^\\s*package.*%s/;\"\t%c",
			name.Name, fset.Position(name.NamePos).Filename, name.Name, kind))
	case FUNC:
		tags.Push(fmt.Sprintf("%s\t%s\t/\\m^\\s*func.*%s(.*{/;\"\t%c",
			name.Name, fset.Position(name.NamePos).Filename, name.Name, kind))
	default:
		//fmt.Printf("!_DEBUG\t%#v - %c\n", name, kind)
		tags.Push(fmt.Sprintf("%s\t%s\t%d;\"\t%c",
			name.Name, fset.Position(name.NamePos).Filename, fset.Position(name.NamePos).Line, kind))
	}
}

func main() {
	parse_files()

	println("!_TAG_FILE_SORTED\t1\t")
	sort.SortStrings(tags)
	for _, s := range tags {
		println(s)
	}
}

const FUNC, TYPE, VAR, PKG = 'f', 't', 'v', 'p'

func parse_files() {
	for i, m := 1, len(os.Args); i < m; i++ {
		fset := token.NewFileSet()
		tree, ok := parser.ParseFile(fset, os.Args[i], nil, 0)
		if ok != nil {
			fmt.Fprintf(os.Stderr, "error parsing file %s - %s\n", os.Args[i], ok.String())
			os.Exit(1)
		}
		output_tag(fset, tree.Name, PKG);
		for _, node := range tree.Decls {
			switch n := node.(type) {
			case *ast.FuncDecl:
				output_tag(fset, n.Name, FUNC)
			case *ast.GenDecl:
				do_gen_decl(fset, n)
			}
		}
	}

}

func do_gen_decl(fset *token.FileSet, node *ast.GenDecl) {
	//fmt.Printf("!_DEBUG\tGenDecl\t%+v\n", node)
	for _, v := range node.Specs {
		switch n := v.(type) {
		case *ast.TypeSpec:
			//fmt.Printf("!_DEBUG\tTypeSpec(%+v)\t%+v\n", n.Type, n)
			output_tag(fset, n.Name, TYPE)

		case *ast.ValueSpec:
			//fmt.Printf("!_DEBUG\tValueSpec(%s)\t%+v\n", n.Type, n)
			for _, vv := range n.Names {
				output_tag(fset, vv, VAR)
			}
		default:
			//fmt.Printf("!_DEBUG\tSPEC:%s\t%+v\n", reflect.Typeof(n), n)
		}
	}
}
