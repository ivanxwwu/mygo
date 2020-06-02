package myast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

/*
基础结构说明
普通Node,不是特定语法结构,属于某个语法结构的一部分.
Comment 表示一行注释 // 或者 / /
CommentGroup 表示多行注释
Field 表示结构体中的一个定义或者变量,或者函数签名当中的参数或者返回值
FieldList 表示以”{}”或者”()”包围的Filed列表
Expression & Types (都划分成Expr接口)
BadExpr 用来表示错误表达式的占位符
Ident 比如报名,函数名,变量名
Ellipsis 省略号表达式,比如参数列表的最后一个可以写成arg…
BasicLit 基本字面值,数字或者字符串
FuncLit 函数定义
CompositeLit 构造类型,比如{1,2,3,4}
ParenExpr 括号表达式,被括号包裹的表达式
SelectorExpr 选择结构,类似于a.b的结构
IndexExpr 下标结构,类似这样的结构 expr[expr]
SliceExpr 切片表达式,类似这样 expr[low:mid:high]
TypeAssertExpr 类型断言类似于 X.(type)
CallExpr 调用类型,类似于 expr()
StarExpr 指针表达式,类似于 *X
UnaryExpr 一元表达式
BinaryExpr 二元表达式
KeyValueExp 键值表达式 key:value
ArrayType 数组类型
StructType 结构体类型
FuncType 函数类型
InterfaceType 接口类型
MapType map类型
ChanType 管道类型
Statements语句
BadStmt 错误的语句
DeclStmt 在语句列表里的申明
EmptyStmt 空语句
LabeledStmt 标签语句类似于 indent:stmt
ExprStmt 包含单独的表达式语句
SendStmt chan发送语句
IncDecStmt 自增或者自减语句
AssignStmt 赋值语句
GoStmt Go语句
DeferStmt 延迟语句
ReturnStmt return 语句
BranchStmt 分支语句 例如break continue
BlockStmt 块语句 {} 包裹
IfStmt If 语句
CaseClause case 语句
SwitchStmt switch 语句
TypeSwitchStmt 类型switch 语句 switch x:=y.(type)
CommClause 发送或者接受的case语句,类似于 case x <-:
SelectStmt select 语句
ForStmt for 语句
RangeStmt range 语句
Declarations声明
Spec type
Import Spec
Value Spec
Type Spec
BadDecl 错误申明
GenDecl 一般申明(和Spec相关,比如 import “a”,var a,type a)
FuncDecl 函数申明
Files and Packages
File 代表一个源文件节点,包含了顶级元素.
Package 代表一个包,包含了很多文件.
以上内容转载自某片大神文章..具体地址忘了，知道
 */

func TestHello(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "./sample/hello.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Imports:================================================")
	for _, i := range node.Imports {
		fmt.Println(i.Path.Value)
	}

	fmt.Println("Comments:================================================")
	for _, c := range node.Comments {
		fmt.Print(c.Text())
	}

	fmt.Println("Functions:================================================")
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		fmt.Println(fn.Name.Name)
	}

	fmt.Println("Inspect:================================================")
	ast.Inspect(node, func(n ast.Node) bool {
		// Find Return Statements
		ret, ok := n.(*ast.ReturnStmt)
		if ok {
			fmt.Printf("return statement found on line %d:\n\t", fset.Position(ret.Pos()).Line)
			printer.Fprint(os.Stdout, fset, ret)
			return true
		}
		return true
	})
}

type Visitor struct {
}
func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	switch node.(type) {
	case *ast.GenDecl:
		genDecl := node.(*ast.GenDecl)
		// 查找有没有import context包
		// Notice：没有考虑没有import任何包的情况
		if genDecl.Tok == token.IMPORT {
			v.addImport(genDecl)
			// 不需要再遍历子树
			return nil
		}
	case *ast.InterfaceType:
		// 遍历所有的接口类型
		iface := node.(*ast.InterfaceType)
		addContext(iface)
		// 不需要再遍历子树
		return nil
	case *ast.Ident:
		ident := node.(*ast.Ident)
		fmt.Printf("%+v", ident)
	}

	return v
}

// addImport 引入context包
func (v *Visitor) addImport(genDecl *ast.GenDecl) {
	// 是否已经import
	hasImported := false
	for _, v := range genDecl.Specs {
		imptSpec := v.(*ast.ImportSpec)
		// 如果已经包含"context"
		if imptSpec.Path.Value == strconv.Quote("context") {
			hasImported = true
		}
	}
	// 如果没有import context，则import
	if !hasImported {
		genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote("context"),
			},
		})
	}
}

// addImport 引入context包
func AddImport(genDecl *ast.GenDecl) {
	// 是否已经import
	hasImported := false
	for _, v := range genDecl.Specs {
		imptSpec := v.(*ast.ImportSpec)
		// 如果已经包含"context"
		if imptSpec.Path.Value == strconv.Quote("context") {
			hasImported = true
		}
	}
	// 如果没有import context，则import
	if !hasImported {
		genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: strconv.Quote("context"),
			},
		})
	}
}

// addContext 添加context参数
func addContext(iface *ast.InterfaceType) {
	// 接口方法不为空时，遍历接口方法
	if iface.Methods != nil || iface.Methods.List != nil {
		for _, v := range iface.Methods.List {
			ft := v.Type.(*ast.FuncType)
			hasContext := false
			// 判断参数中是否包含context.Context类型
			for _, v := range ft.Params.List {
				if expr, ok := v.Type.(*ast.SelectorExpr); ok {
					if ident, ok := expr.X.(*ast.Ident); ok {
						if ident.Name == "context" {
							hasContext = true
						}
					}
				}
			}
			// 为没有context参数的方法添加context参数
			if !hasContext {
				ctxField := &ast.Field{
					Names: []*ast.Ident{
						ast.NewIdent("ctx"),
					},
					// Notice: 没有考虑import别名的情况
					Type: &ast.SelectorExpr{
						X:   ast.NewIdent("context"),
						Sel: ast.NewIdent("Context"),
					},
				}
				list := []*ast.Field{
					ctxField,
				}
				ft.Params.List = append(list, ft.Params.List...)
			}
		}
	}
}


func TestDemo(t *testing.T) {
	fset := token.NewFileSet()
	// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
	path, _ := filepath.Abs("./sample/demo.go")
	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		log.Println(err)
	}
	// 打印语法树
	ast.Print(fset, f)


}

func TestDemo2(t *testing.T) {
	fset := token.NewFileSet()
	// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
	path, _ := filepath.Abs("./sample/demo.go")
	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		log.Println(err)
	}
	ast.Print(fset, f)
	v := Visitor{}
	ast.Walk(&v, f)

	var output []byte
	buffer := bytes.NewBuffer(output)
	err = format.Node(buffer, fset, f)
	if err != nil {
		log.Fatal(err)
	}

	// 输出Go代码
	fmt.Println(buffer.String())
}


func TestDemo3(t *testing.T) {
	fset := token.NewFileSet()
	// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
	path, _ := filepath.Abs("/root/goworkspace/nfa_autotest/nfa_gotests/example/example_test.go")
	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		log.Println(err)
	}
	ast.Print(fset, f)
}



func TestImports(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "./sample/ast_traversal.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
		return
	}

	for _, i := range f.Imports {
		t.Logf("import:	%s", i.Path.Value)
	}

	return
}

func TestComments(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "./sample/ast_traversal.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
		return
	}

	for _, i := range f.Comments {
		t.Logf("comment:	%s", i.Text())
	}

	return
}

func TestFunctions(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "./sample/ast_traversal.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
		return
	}

	for _, i := range f.Decls {
		fn, ok := i.(*ast.FuncDecl)
		if !ok {
			continue
		}
		t.Logf("function:	%s", fn.Name.Name)
	}

	return
}

func TestInspect(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "./sample/demo.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
		return
	}

	//ast.Print(fset, f)

	ast.Inspect(f, func(n ast.Node) bool {
		// Find Return Statements
		gen, ok := n.(*ast.GenDecl)
		if ok {
			t.Logf("return GenDecl found on line %d:\n\t", fset.Position(gen.Pos()).Line)
			if gen.Tok == token.IMPORT {
				AddImport(gen)
			}
		}
		ret, ok := n.(*ast.ReturnStmt)
		if ok {
			t.Logf("return statement found on line %d:\n\t", fset.Position(ret.Pos()).Line)
			printer.Fprint(os.Stdout, fset, ret)
			return true
		}

		// Find Functions
		fn, ok := n.(*ast.FuncDecl)
		if ok {
			var exported string
			if fn.Name.IsExported() {
				exported = "exported "
			}
			t.Logf("%sfunction declaration found on line %d: %s", exported, fset.Position(fn.Pos()).Line, fn.Name.Name)
			return true
		}

		call, ok := n.(*ast.CallExpr)
		if ok {
			t.Logf("call:%+v", call)
		}

		return true
	})

	var output []byte
	buffer := bytes.NewBuffer(output)
	err = format.Node(buffer, fset, f)
	if err != nil {
		log.Fatal(err)
	}

	// 输出Go代码
	fmt.Println(buffer.String())

	return
}


func TestGomonkey(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "./sample/demo.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
		return
	}

	ast.Inspect(f, func(n ast.Node) bool {
		// Find Return Statements
		gen, ok := n.(*ast.GenDecl)
		if ok {
			if gen.Tok != token.VAR {
				return false
			}
			if len(gen.Specs) < 1 {
				return false
			}
			spec := gen.Specs[0]
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				return false
			}
			if len(vs.Names) < 1 || vs.Names[0].Name != "mocks" {
				return false
			}

			if len(vs.Values) < 1 {
				return false
			}
			fnDecl, ok := vs.Values[0].(*ast.FuncLit)
			if !ok {
				return false
			}
			if len(fnDecl.Type.Results.List) < 1 {
				return false
			}
			retType, ok := fnDecl.Type.Results.List[0].Type.(*ast.ArrayType)
			if !ok {
				return false
			}
			sexpr, ok := retType.Elt.(*ast.StarExpr)
			if !ok {
				return false
			}
			selexpr, ok := sexpr.X.(*ast.SelectorExpr)
			if !ok {
				return false
			}
			pack, ok := selexpr.X.(*ast.Ident)
			if !ok {
				return false
			}
			if pack.Name != "gomonkey" || selexpr.Sel.Name != "Patches" {
				return false
			}
			t.Logf("return GenDecl found on line %d %+v:\n\t", fset.Position(gen.Pos()).Line, gen)
			var output []byte
			buf := bytes.NewBuffer(output)
			err = format.Node(buf, fset, gen)
			t.Logf(buf.String())
			//ast.Inspect(gen, func(n ast.Node) bool {
			//
			//	return true
			//})
			//if gen.Tok == token.IMPORT {
			//	AddImport(gen)
			//}
		}
		return true
	})

	var output []byte
	buffer := bytes.NewBuffer(output)
	err = format.Node(buffer, fset, f)
	if err != nil {
		log.Fatal(err)
	}

	// 输出Go代码
	fmt.Println(buffer.String())

	return
}
