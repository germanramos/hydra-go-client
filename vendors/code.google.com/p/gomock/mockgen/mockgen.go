// Copyright 2010 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// MockGen generates mock implementations of Go interfaces.
package main

// TODO: This does not support recursive embedded interfaces.
// TODO: This does not support embedding package-local interfaces in a separate file.

import (
	"flag"
	"fmt"
	"go/token"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"unicode"

	"github.com/innotech/hydra-go-client/vendors/code.google.com/p/gomock/mockgen/model"
)

const (
	gomockImportPath = "code.google.com/p/gomock/gomock"
)

var (
	source		= flag.String("source", "", "(source mode) Input Go source file; enables source mode.")
	destination	= flag.String("destination", "", "Output file; defaults to stdout.")
	packageOut	= flag.String("package", "", "Package of the generated code; defaults to the package of the input with a 'mock_' prefix.")
	selfPackage	= flag.String("self_package", "", "If set, the package this mock will be part of.")

	debugParser	= flag.Bool("debug_parser", false, "Print out parser results only.")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	var pkg *model.Package
	var err error
	if *source != "" {
		pkg, err = ParseFile(*source)
	} else {
		if flag.NArg() != 2 {
			log.Fatal("Expected exactly two arguments")
		}
		pkg, err = Reflect(flag.Arg(0), strings.Split(flag.Arg(1), ","))
	}
	if err != nil {
		log.Fatalf("Loading input failed: %v", err)
	}

	if *debugParser {
		pkg.Print(os.Stdout)
		return
	}

	dst := os.Stdout
	if len(*destination) > 0 {
		f, err := os.Create(*destination)
		if err != nil {
			log.Fatalf("Failed opening destination file: %v", err)
		}
		defer f.Close()
		dst = f
	}

	packageName := *packageOut
	if packageName == "" {
		// pkg.Name in reflect mode is the base name of the import path,
		// which might have characters that are illegal to have in package names.
		packageName = "mock_" + sanitize(pkg.Name)
	}

	g := generator{
		w: dst,
	}
	if *source != "" {
		g.filename = *source
	} else {
		g.srcPackage = flag.Arg(0)
		g.srcInterfaces = flag.Arg(1)
	}
	if err := g.Generate(pkg, packageName); err != nil {
		log.Fatalf("Failed generating mock: %v", err)
	}
}

func usage() {
	io.WriteString(os.Stderr, usageText)
	flag.PrintDefaults()
}

const usageText = `mockgen has two modes of operation: source and reflect.

Source mode generates mock interfaces from a source file.
It is enabled by using the -source flag. Other flags that
may be useful in this mode are -imports and -aux_files.
Example:
	mockgen -source=foo.go [other options]

Reflect mode generates mock interfaces by building a program
that uses reflection to understand interfaces. It is enabled
by passing two non-flag arguments: an import path, and a
comma-separated list of symbols.
Example:
	mockgen database/sql/driver Conn,Driver

`

type generator struct {
	w	io.Writer
	indent	string

	filename			string	// may be empty
	srcPackage, srcInterfaces	string	// may be empty

	packageMap	map[string]string	// map from import path to package name
}

func (g *generator) p(format string, args ...interface{}) {
	fmt.Fprintf(g.w, g.indent+format+"\n", args...)
}

func (g *generator) in() {
	g.indent += "\t"
}

func (g *generator) out() {
	if len(g.indent) > 0 {
		g.indent = g.indent[0 : len(g.indent)-1]
	}
}

func removeDot(s string) string {
	if len(s) > 0 && s[len(s)-1] == '.' {
		return s[0 : len(s)-1]
	}
	return s
}

// sanitize cleans up a string to make a suitable package name.
func sanitize(s string) string {
	t := ""
	for _, r := range s {
		if t == "" {
			if unicode.IsLetter(r) || r == '_' {
				t += string(r)
				continue
			}
		} else {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				t += string(r)
				continue
			}
		}
		t += "_"
	}
	if t == "_" {
		t = "x"
	}
	return t
}

func (g *generator) Generate(pkg *model.Package, pkgName string) error {
	g.p("// Automatically generated by MockGen. DO NOT EDIT!")
	if g.filename != "" {
		g.p("// Source: %v", g.filename)
	} else {
		g.p("// Source: %v (interfaces: %v)", g.srcPackage, g.srcInterfaces)
	}
	g.p("")

	// Get all required imports, and generate unique names for them all.
	im := pkg.Imports()
	im[gomockImportPath] = true
	g.packageMap = make(map[string]string, len(im))
	localNames := make(map[string]bool, len(im))
	for pth := range im {
		base := sanitize(path.Base(pth))

		// Local names for an imported package can usually be the basename of the import path.
		// A couple of situations don't permit that, such as duplicate local names
		// (e.g. importing "html/template" and "text/template"), or where the basename is
		// a keyword (e.g. "foo/case").
		// try base0, base1, ...
		pkgName := base
		i := 0
		for localNames[pkgName] || token.Lookup(pkgName).IsKeyword() {
			pkgName = base + strconv.Itoa(i)
			i++
		}

		g.packageMap[pth] = pkgName
		localNames[pkgName] = true
	}

	g.p("package %v", pkgName)
	g.p("")
	g.p("import (")
	g.in()
	for path, pkg := range g.packageMap {
		if path == *selfPackage {
			continue
		}
		g.p("%v %q", pkg, path)
	}
	for _, path := range pkg.DotImports {
		g.p(". %q", path)
	}
	g.out()
	g.p(")")

	for _, intf := range pkg.Interfaces {
		if err := g.GenerateMockInterface(intf); err != nil {
			return err
		}
	}

	return nil
}

// The name of the mock type to use for the given interface identifier.
func mockName(typeName string) string {
	return "Mock" + typeName
}

func (g *generator) GenerateMockInterface(intf *model.Interface) error {
	mockType := mockName(intf.Name)

	g.p("")
	g.p("// Mock of %v interface", intf.Name)
	g.p("type %v struct {", mockType)
	g.in()
	g.p("ctrl     *gomock.Controller")
	g.p("recorder *_%vRecorder", mockType)
	g.out()
	g.p("}")
	g.p("")

	g.p("// Recorder for %v (not exported)", mockType)
	g.p("type _%vRecorder struct {", mockType)
	g.in()
	g.p("mock *%v", mockType)
	g.out()
	g.p("}")
	g.p("")

	// TODO: Re-enable this if we can import the interface reliably.
	//g.p("// Verify that the mock satisfies the interface at compile time.")
	//g.p("var _ %v = (*%v)(nil)", typeName, mockType)
	//g.p("")

	g.p("func New%v(ctrl *gomock.Controller) *%v {", mockType, mockType)
	g.in()
	g.p("mock := &%v{ctrl: ctrl}", mockType)
	g.p("mock.recorder = &_%vRecorder{mock}", mockType)
	g.p("return mock")
	g.out()
	g.p("}")
	g.p("")

	// XXX: possible name collision here if someone has EXPECT in their interface.
	g.p("func (_m *%v) EXPECT() *_%vRecorder {", mockType, mockType)
	g.in()
	g.p("return _m.recorder")
	g.out()
	g.p("}")

	g.GenerateMockMethods(mockType, intf, *selfPackage)

	return nil
}

func (g *generator) GenerateMockMethods(mockType string, intf *model.Interface, pkgOverride string) {
	for _, m := range intf.Methods {
		g.p("")
		g.GenerateMockMethod(mockType, m, pkgOverride)
		g.p("")
		g.GenerateMockRecorderMethod(mockType, m)
	}
}

// GenerateMockMethod generates a mock method implementation.
// If non-empty, pkgOverride is the package in which unqualified types reside.
func (g *generator) GenerateMockMethod(mockType string, m *model.Method, pkgOverride string) error {
	args := make([]string, len(m.In))
	argNames := make([]string, len(m.In))
	for i, p := range m.In {
		name := p.Name
		if name == "" {
			name = fmt.Sprintf("_param%d", i)
		}
		ts := p.Type.String(g.packageMap, pkgOverride)
		args[i] = name + " " + ts
		argNames[i] = name
	}
	if m.Variadic != nil {
		name := m.Variadic.Name
		if name == "" {
			name = fmt.Sprintf("_param%d", len(m.In))
		}
		ts := m.Variadic.Type.String(g.packageMap, pkgOverride)
		args = append(args, name+" ..."+ts)
		argNames = append(argNames, name)
	}
	argString := strings.Join(args, ", ")

	rets := make([]string, len(m.Out))
	for i, p := range m.Out {
		rets[i] = p.Type.String(g.packageMap, pkgOverride)
	}
	retString := strings.Join(rets, ", ")
	if len(rets) > 1 {
		retString = "(" + retString + ")"
	}
	if retString != "" {
		retString = " " + retString
	}

	g.p("func (_m *%v) %v(%v)%v {", mockType, m.Name, argString, retString)
	g.in()

	callArgs := strings.Join(argNames, ", ")
	if callArgs != "" {
		callArgs = ", " + callArgs
	}
	if m.Variadic != nil {
		// Non-trivial. The generated code must build a []interface{},
		// but the variadic argument may be any type.
		g.p("_s := []interface{}{%s}", strings.Join(argNames[:len(argNames)-1], ", "))
		g.p("for _, _x := range %s {", argNames[len(argNames)-1])
		g.in()
		g.p("_s = append(_s, _x)")
		g.out()
		g.p("}")
		callArgs = ", _s..."
	}
	if len(m.Out) == 0 {
		g.p(`_m.ctrl.Call(_m, "%v"%v)`, m.Name, callArgs)
	} else {
		g.p(`ret := _m.ctrl.Call(_m, "%v"%v)`, m.Name, callArgs)

		// Go does not allow "naked" type assertions on nil values, so we use the two-value form here.
		// The value of that is either (x.(T), true) or (Z, false), where Z is the zero value for T.
		// Happily, this coincides with the semantics we want here.
		retNames := make([]string, len(rets))
		for i, t := range rets {
			retNames[i] = fmt.Sprintf("ret%d", i)
			g.p("%s, _ := ret[%d].(%s)", retNames[i], i, t)
		}
		g.p("return " + strings.Join(retNames, ", "))
	}

	g.out()
	g.p("}")
	return nil
}

func (g *generator) GenerateMockRecorderMethod(mockType string, m *model.Method) error {
	nargs := len(m.In)
	args := make([]string, nargs)
	for i := 0; i < nargs; i++ {
		args[i] = "arg" + strconv.Itoa(i)
	}
	argString := strings.Join(args, ", ")
	if nargs > 0 {
		argString += " interface{}"
	}
	if m.Variadic != nil {
		if nargs > 0 {
			argString += ", "
		}
		argString += fmt.Sprintf("arg%d ...interface{}", nargs)
	}

	g.p("func (_mr *_%vRecorder) %v(%v) *gomock.Call {", mockType, m.Name, argString)
	g.in()

	callArgs := strings.Join(args, ", ")
	if nargs > 0 {
		callArgs = ", " + callArgs
	}
	if m.Variadic != nil {
		if nargs == 0 {
			// Easy: just use ... to push the arguments through.
			callArgs = ", arg0..."
		} else {
			// Hard: create a temporary slice.
			g.p("_s := append([]interface{}{%s}, arg%d...)", strings.Join(args, ", "), nargs)
			callArgs = ", _s..."
		}
	}
	g.p(`return _mr.mock.ctrl.RecordCall(_mr.mock, "%v"%v)`, m.Name, callArgs)

	g.out()
	g.p("}")
	return nil
}