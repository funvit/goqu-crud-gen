package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/scanner"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/fatih/structtag"
)

const ( //known dialects
	dMysql     = "mysql"
	dMysql8    = "mysql8"
	dPostgres  = "postgres"
	dSqlite3   = "sqlite3"
	dSqlserver = "sqlserver"
)

var knownDialects = []string{
	dMysql,
	dMysql8,
	dPostgres,
	dSqlite3,
	dSqlserver,
}

var ( // loggers
	logInfo, logError, logDebug *log.Logger
)

func init() {
	logInfo = log.New(os.Stdout, "INFO ", log.LstdFlags)
	logError = log.New(os.Stderr, "ERRO ", log.LstdFlags)
	logDebug = log.New(ioutil.Discard, "DEBU ", log.LstdFlags|log.Lshortfile)
}

func main() {
	fmt.Println("GOQU CRUD generator.")

	fModel := flag.String("model", "", "model name")
	fTable := flag.String("table", "", "table name")
	fDialect := flag.String("dialect", "", "database dialect")
	fPath := flag.String("path", ".", "path to folder with model")
	fRepo := flag.String("repo", "", "custom repository struct name")
	fPrivateCrud := flag.Bool("private-crud-methods", false, "create CRUD methods as private")
	fNoGen := flag.Bool("g", false, "don't put //go:generate instruction to the generated code")
	fWithTranName := flag.String("rename-with-tran", "WithTran", "rename WithTran helper method")

	fDebug := flag.Bool("d", false, "debug")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()

		fmt.Fprintf(
			flag.CommandLine.Output(),
			"\nKnown dialects: %s.\n",
			strings.Join(knownDialects, ", "),
		)

		fmt.Fprintf(
			flag.CommandLine.Output(),
			"\nExample:",
		)

		fmt.Fprintf(
			flag.CommandLine.Output(),
			"\n\n    %s -model User -table user -dialect mysql \n\n",
			os.Args[0],
		)
	}

	flag.Parse()

	// region validate some flags

	if *fPath == "" {
		logError.Fatalln("path flag must be specified")
	}
	if *fModel == "" {
		logError.Fatalln("model flag must be specified")
	}
	if *fTable == "" {
		logError.Fatalln("table flag must be specified")
	}
	if *fWithTranName == "" {
		logError.Fatalln("WithTran new name flag must be specified")
	}

	// endregion

	// region handle flags

	if *fDebug {
		logDebug.SetOutput(os.Stdout)
	}

	// endregion

	// remove current repo file to avoid parser errors
	repoFileName := camelToSake(*fModel) + "_repo.go"
	outputFilepath := filepath.Join(*fPath, repoFileName)
	_ = os.Remove(outputFilepath)

	// region parse path

	fset := token.NewFileSet()
	m, err := parser.ParseDir(
		fset,
		filepath.Dir(*fPath),
		func(info os.FileInfo) bool {
			return !info.IsDir()
		},
		parser.ParseComments,
	)
	if err != nil {
		logError.Fatalln(err)
	}
	if len(m) == 0 {
		logError.Fatalln("no packages found")
	}
	if len(m) > 1 {
		logError.Fatalln("more than one package found")
	}

	// endregion

	// region fill tpl dto

	var p *ast.Package
	for _, v := range m {
		p = v
		break
	}

	d, err := modelToTplDTO(p, *fModel)
	if err != nil {
		logError.Fatalln("model spec load error:", err)
	}

	d.WithTranName = *fWithTranName

	d.Repo = Repo{
		Name:    fmt.Sprintf("%sRepo", *fModel),
		Table:   *fTable,
		Dialect: *fDialect,
	}

	if *fRepo != "" {
		d.Repo.Name = *fRepo
	}
	if *fPrivateCrud {
		d.PrivateCRUD = true
	}
	if !*fNoGen {
		var sb strings.Builder
		sb.WriteString(os.Args[0])

		sb.WriteString(" ")
		sb.WriteString("-model")
		sb.WriteString(" ")
		sb.WriteString(*fModel)

		sb.WriteString(" ")
		sb.WriteString("-table")
		sb.WriteString(" ")
		sb.WriteString(*fTable)

		sb.WriteString(" ")
		sb.WriteString("-dialect")
		sb.WriteString(" ")
		sb.WriteString(*fDialect)

		sb.WriteString(" ")
		sb.WriteString("-path")
		sb.WriteString(" ")
		sb.WriteString(".")

		if *fRepo != "" {
			sb.WriteString(" ")
			sb.WriteString("-repo")
			sb.WriteString(" ")
			sb.WriteString(*fRepo)
		}

		if *fPrivateCrud {
			sb.WriteString(" ")
			sb.WriteString("-private-crud-methods")
		}

		if *fWithTranName != "" {
			sb.WriteString(" ")
			sb.WriteString("-rename-with-tran")
			sb.WriteString(" ")
			sb.WriteString(*fWithTranName)
		}

		if *fDebug {
			sb.WriteString(" ")
			sb.WriteString("-d")
		}

		d.GenerateCmd = sb.String()
	}

	// endregion

	// region validate model

	if !d.Model.HasPrimaryKeyField() {
		logError.Fatalln("model have no field with \"primary\" mark (db tag option)")
	}
	logDebug.Printf(
		"model primary key field: %s with type: %s",
		d.Model.GetPrimaryKeyField().Name,
		d.Model.GetPrimaryKeyField().Type,
	)

	// endregion

	// region generate file

	buf := bytes.NewBuffer(nil)
	err = generateRepoFile(buf, *d)
	if err != nil {
		_ = os.Remove(repoFileName)
		logError.Fatalln("file generate error:", err)
	}
	logInfo.Println("file written:", outputFilepath)

	// endregion

	// region format generated file

	f, err := os.Create(outputFilepath)
	if err != nil {
		log.Fatalln("output file open error:", err)
	}
	defer func() {
		_ = f.Close()
	}()

	logInfo.Println("formatting generated code")
	b, err := format.Source(buf.Bytes())
	if err != nil {
		logError.Printf("generated code format error: %s", err.Error())

		switch sel := err.(type) {
		case scanner.ErrorList:
			if len(sel) > 0 {
				err = sel[0]
			}
		}

		var scannerError *scanner.Error
		if errors.As(err, &scannerError) {
			const paddingLines = 5

			lineScanner := bufio.NewScanner(bytes.NewBuffer(buf.Bytes()))
			startLine := scannerError.Pos.Line

			if startLine > paddingLines {
				startLine -= paddingLines
			}

			for i := 0; i < startLine; i++ {
				lineScanner.Scan()
			}

			// calc max line number len as str
			var maxLineNumberLen int
			{
				ln := scannerError.Pos.Line + paddingLines
				lnStr := strconv.Itoa(ln)
				maxLineNumberLen = len(lnStr)
			}

			fmt.Println("")
			i := 0
			for lineScanner.Scan() {
				if startLine != scannerError.Pos.Line {
					if i >= paddingLines*2+1 {
						break
					}
				} else {
					if i >= paddingLines+1 {
						break
					}
				}
				if startLine+i == scannerError.Pos.Line {
					fmt.Print("--> ")
				} else {
					fmt.Print("    ")
				}
				fmt.Printf(
					"line %"+strconv.Itoa(maxLineNumberLen)+"d: %s\n",
					startLine+i,
					string(lineScanner.Bytes()),
				)
				i++
			}
			fmt.Println("")
		}

		os.Exit(1)
	}

	// endregion

	// region write generated file

	n, err := f.Write(b)
	if err != nil {
		logError.Fatalln("file write error: %w", err)
	}
	logDebug.Printf("written %d byte(s)", n)

	// endregion

	fmt.Println("Done.")
}

func camelToSake(s string) string {
	var sb strings.Builder

	for i, r := range s {
		if string(r) == strings.ToUpper(string(r)) {
			if i > 0 {
				sb.WriteRune('_')
			}
			sb.WriteString(strings.ToLower(string(r)))
			continue
		}
		sb.WriteRune(r)
	}

	return sb.String()
}

func modelToTplDTO(p *ast.Package, modelName string) (*tplDTO, error) {
	logInfo.Println("found package:", p.Name)

	// find model
	var model *ast.TypeSpec
	var modelFileImports []*ast.ImportSpec
	for _, v := range p.Files {
		ast.Inspect(v, func(node ast.Node) bool {
			ts, ok := node.(*ast.TypeSpec)
			if !ok {
				return true
			}

			if ts.Name.Name == modelName {
				model = ts
				modelFileImports = v.Imports
			}
			return false
		})
	}

	if model == nil {
		return nil, fmt.Errorf("model not found")
	}

	r := tplDTO{
		Package: p.Name,
		Model: &Model{
			MustImport: map[string]struct{}{},
			Name:       modelName,
		},
	}

	logInfo.Println("found struct:", model.Name)
	for _, field := range model.Type.(*ast.StructType).Fields.List {
		logDebug.Println(fmt.Sprintf(
			"field:%s, tag: %s",
			field.Names[0].Name,
			field.Tag.Value,
		))

		tags, err := structtag.Parse(strings.Trim(field.Tag.Value, "`"))
		if err != nil {
			return nil, fmt.Errorf("field %s tag parse error: %w", field.Names[0].Name, err)
		}

		t, err := tags.Get("db")
		if err != nil {
			return nil, fmt.Errorf("tag \"db\" value read error: %w", err)
		}

		mf := ModelField{
			Name:    field.Names[0].Name,
			ColName: t.Name,
			Options: t.Options,
		}

		if v, ok := field.Type.(*ast.Ident); ok {
			mf.Type = v.Name
		} else if v, ok := field.Type.(*ast.SelectorExpr); ok {
			if x, ok := v.X.(*ast.Ident); ok {
				mf.Type = x.String() + "." + v.Sel.String()
				r.Model.MustImport[x.String()] = struct{}{}
			} else {
				mf.Type = v.Sel.String()
			}
		} else if v, ok := field.Type.(*ast.ArrayType); ok {
			if v.Len == nil {
				mf.Type = fmt.Sprintf(
					"[]%s",
					v.Elt.(*ast.Ident).String(),
				)
			} else {
				mf.Type = fmt.Sprintf(
					"[%s]%s",
					v.Len.(*ast.BasicLit).Value,
					v.Elt.(*ast.Ident).String(),
				)
			}
		} else {
			return nil, errors.New(fmt.Sprintln(
				"cannot determine field:", field.Names,
				"type:", field.Type,
			))
		}

		r.Model.Fields = append(r.Model.Fields, mf)
	}

	// collect required imports
	for k := range r.Model.MustImport {
		for _, imp := range modelFileImports {
			if imp.Name != nil && imp.Name.Name == k {
				r.Imports = append(r.Imports, fmt.Sprintf(`%s %s`, imp.Name.Name, imp.Path.Value))
				break
			}

			last := filepath.Base(strings.Trim(imp.Path.Value, "\""))
			if last == k {
				r.Imports = append(r.Imports, imp.Path.Value)
				break
			}
		}
	}

	return &r, nil
}

func mustWriteLn(w io.Writer, format string, arg ...interface{}) {
	_, err := fmt.Fprintln(w, fmt.Sprintf(format, arg...))
	if err != nil {
		panic(err)
	}
}

func generateRepoFile(w io.Writer, p tplDTO) error {
	funcMap := template.FuncMap{
		"CRUD": func(s string) string {
			if p.PrivateCRUD {
				return fmt.Sprintf("_%s", s)
			}
			return s
		},
		"Private": func(s string) string {
			if s == "" {
				return ""
			}
			return strings.ToLower(string(s[0])) + s[1:]
		},
		"StrConcat": func(a, b string) string { return a + b },
	}

	//
	// header
	//
	logDebug.Println("writing header")

	mustWriteLn(w, "// Code generated by generator; DO NOT EDIT.")
	mustWriteLn(w, "")
	if p.GenerateCmd != "" {
		mustWriteLn(w, fmt.Sprintf("//go:generate %s", p.GenerateCmd))
	}
	mustWriteLn(w, "package %s", p.Package)

	importPackages := []string{
		"context",
		"fmt",
		"github.com/doug-martin/goqu/v9",
		"github.com/doug-martin/goqu/v9/exp",
		"github.com/jmoiron/sqlx",
		". github.com/funvit/goqu-crud-gen",
		"time",
	}
	mustWriteLn(w, "import (")
	for _, v := range importPackages {
		strs := strings.SplitN(v, " ", 2)
		if len(strs) == 2 {
			mustWriteLn(w, fmt.Sprintf(`%s "%s"`, strs[0], strs[1]))
			continue
		}
		mustWriteLn(w, fmt.Sprintf(`"%s"`, v))
	}

	mustWriteLn(w, dialectImport(p.Repo.Dialect))

	// import packages if required for PK
	for _, v := range p.Imports {
		mustWriteLn(w, v)
	}

	mustWriteLn(w, ")")
	mustWriteLn(w, "")

	// name => tpl
	templates := [][]string{
		{"repo", repoTpl},
		{"create", createTpl},
		{"get", getTpl},
		{"update", updateTpl},
		{"delete", deleteTpl},
	}

	for _, v := range templates {
		logDebug.Printf("writing %s tpl execute result", v[0])

		t, err := template.New("").Funcs(funcMap).Parse(v[1])
		if err != nil {
			return err
		}

		err = t.Execute(w, p)
		if err != nil {
			return fmt.Errorf("%s tpl processing error: %w", v[0], err)
		}
	}

	return nil
}

func dialectImport(dialect string) string {
	const comment = "import is need for proper dialect selection"

	var v string

	switch dialect {
	case dMysql, dMysql8:
		v = "github.com/doug-martin/goqu/v9/dialect/mysql"
	case dPostgres:
		v = "github.com/doug-martin/goqu/v9/dialect/postgres"
	case dSqlite3:
		v = "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	case dSqlserver:
		v = "github.com/doug-martin/goqu/v9/dialect/sqlserver"
	default:
		logError.Fatalf("unknown dialect: %s", dialect)
	}

	return fmt.Sprintf(`_ "%s" // %s`, v, comment)
}
