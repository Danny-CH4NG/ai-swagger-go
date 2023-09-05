package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type APIInfo struct {
	Method    string
	URL       string
	Params    []string
	Responses []string
}

var (
	infoMap = make(map[string]*APIInfo)
)

func processFuncDecl(funcDecl *ast.FuncDecl, apiInfo *APIInfo, outFile *os.File) {
	// Traverse the function body
	for _, stmt := range funcDecl.Body.List {
		// Look for expression statements (like function calls)
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}

		// Check for function call
		callExpr, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}

		// Check for "ctx.Param"
		if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if selectorExpr.Sel.Name == "Param" {
				// Extract parameter name
				if lit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
					param := lit.Value // This is the parameter name in string form, including quotes
					apiInfo.Params = append(apiInfo.Params, param)
				}
			}

			// Check for "ctx.JSON"
			if selectorExpr.Sel.Name == "JSON" {
				// TODO: Extract JSON content, you'll probably need to get complex types, not just basic literals
				apiInfo.Responses = append(apiInfo.Responses, "ExampleJSONResponse")
			}
		}
	}
	// Write to output file
	if len(apiInfo.Params) > 0 || len(apiInfo.Responses) > 0 {
		infoMap[funcDecl.Name.Name] = apiInfo
		outFile.WriteString("Param:" + strings.Join(apiInfo.Params, ", ") + "\n")
		outFile.WriteString("Responses:" + strings.Join(apiInfo.Responses, ", ") + "\n")
		if apiInfo.Method != "" && apiInfo.URL != "" {
			outFile.WriteString("Method:" + apiInfo.Method + "\n")
			outFile.WriteString("URL:" + apiInfo.URL + "\n")
		}
		outFile.WriteString("\n\n")
	}
}

func processFile(filePath string, fset *token.FileSet, outFile *os.File) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	node, err := parser.ParseFile(fset, filePath, data, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.FuncDecl:
			apiInfo := APIInfo{}
			processFuncDecl(v, &apiInfo, outFile)
			outFile.WriteString("Param:" + strings.Join(apiInfo.Params, ", ") + "\n")
			outFile.WriteString("Responses:" + strings.Join(apiInfo.Responses, ", ") + "\n")
		case *ast.CallExpr:
			if selectorExpr, ok := v.Fun.(*ast.SelectorExpr); ok {
				if _, ok := selectorExpr.X.(*ast.Ident); ok {
					methodName := selectorExpr.Sel.Name
					if methodName == "GET" || methodName == "POST" || methodName == "PUT" || methodName == "DELETE" {
						// TODO: 從 callExpr.Args 中提取路徑和其他參數
						// 提取 URL
						if len(v.Args) > 0 {
							funcExpr, ok := v.Args[1].(*ast.SelectorExpr)
							if ok {
								calledFuncName := funcExpr.Sel.Name
								urlExpr, ok := v.Args[0].(*ast.BasicLit)
								if ok {
									if apiInfo, exists := infoMap[calledFuncName]; exists {
										apiInfo.Method = methodName
										apiInfo.URL = urlExpr.Value
									} else {
										apiInfo := APIInfo{
											Method: methodName,
											URL:    urlExpr.Value,
										}
										infoMap[calledFuncName] = &apiInfo
									}
									_, err = outFile.WriteString("Method:" + methodName + "\n")
									_, err = outFile.WriteString("URL:" + urlExpr.Value + "\n")
								}
							}
						}
						_, err = outFile.WriteString("\n\n")
					}
				}
			}
		}
		return true
	})
}

func walkDirectory(dir string, fset *token.FileSet, outFile *os.File) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		fullPath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			walkDirectory(fullPath, fset, outFile) // 遞迴進入子目錄
		} else if filepath.Ext(file.Name()) == ".go" {
			processFile(fullPath, fset, outFile)
		}
	}
}

func main() {
	fset := token.NewFileSet()
	rootDir := "../api-service/controllers" // 更換為你的專案根目錄
	rootDir2 := "../api-service/routers"    // 更換為你的專案根目錄

	// 創建 txt 文件
	outFile, err := os.Create("result.txt")
	if err != nil {
		fmt.Println("Error creating txt file:", err)
		return
	}
	defer outFile.Close()

	walkDirectory(rootDir, fset, outFile)
	walkDirectory(rootDir2, fset, outFile)
}
