package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"astswagger/models"

	"gopkg.in/yaml.v3"
)

var (
	apiMap = make(map[string]*models.Draft)
)

// func extractResponseInfo(compLit *ast.CompositeLit) map[string]interface{} {
// 	responseInfo := make(map[string]interface{})
// 	for _, elt := range compLit.Elts {
// 		// 斷言該元素為KeyValue表達式
// 		kvExpr, ok := elt.(*ast.KeyValueExpr)
// 		if ok {
// 			// 斷言key為標識符並獲得其名稱
// 			keyIdent, ok := kvExpr.Key.(*ast.Ident)
// 			if ok {
// 				// 在這裡我們將僅處理基本的字面量值，例如字符串或數字
// 				// 你可能想擴展它以處理更複雜的情況
// 				switch val := kvExpr.Value.(type) {
// 				case *ast.BasicLit:
// 					responseInfo[keyIdent.Name] = val.Value
// 				}
// 			}
// 		}
// 	}
// 	return responseInfo
// }

func extractResponseInfo(compLit *ast.CompositeLit) string {
	var buf bytes.Buffer
	err := format.Node(&buf, token.NewFileSet(), compLit)
	if err != nil {
		return fmt.Sprintf("error formatting node: %v", err)
	}
	responseStr := buf.String()
	responseStr = strings.Replace(responseStr, "gin.H", "", 1)
	return responseStr
}

func processFuncDecl(funcDecl *ast.FuncDecl, outFile *os.File) {
	// printer.Fprint(os.Stdout, token.NewFileSet(), funcDecl)
	currentFuncName := funcDecl.Name.Name // 獲得當前函數名稱
	ast.Inspect(funcDecl, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if ok {
			selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if ok {
				// 檢查這是否是一個 "ctx.JSON" 或 "ctx.Param" 調用
				switch selectorExpr.Sel.Name {
				case "JSON":
					// 如果是，你可以從callExpr.Args獲得參數來獲得你需要的信息
					// ... 你的代碼來處理這個情況
					if len(callExpr.Args) > 1 {
						// 提取 HTTP 狀態碼
						statusCodeExpr := callExpr.Args[0]
						var statusCode string
						if basicLit, ok := statusCodeExpr.(*ast.BasicLit); ok {
							statusCode = basicLit.Value
						} else {
							statusCode = "UnknownStatusCode"
						}

						// 提取 JSON 響應體
						jsonResponseExpr := callExpr.Args[1]
						var jsonResponse string
						if compositeLit, ok := jsonResponseExpr.(*ast.CompositeLit); ok {
							// TODO: Inspect compositeLit to build a proper JSON string representation
							// Here, we just print the Go syntax representation as a placeholder
							jsonResponse = fmt.Sprintf("%+v\n", extractResponseInfo(compositeLit))
						} else if mapExpr, ok := jsonResponseExpr.(*ast.CallExpr); ok {
							// TODO: Handle the map expression properly to build the JSON string
							jsonResponse = fmt.Sprintf("%#v", mapExpr)
						} else {
							jsonResponse = "UnknownJSONResponse"
						}

						if apiInfo, exists := apiMap[currentFuncName]; exists {
							apiInfo.Responses = append(apiInfo.Responses, fmt.Sprintf("Status code: %s, Response: %s", statusCode, jsonResponse))
						} else {
							apiInfo := models.Draft{
								Responses: []string{fmt.Sprintf("Status code: %s, Response: %s", statusCode, jsonResponse)},
							}
							apiMap[currentFuncName] = &apiInfo
						}
					}
				case "Param":
					if lit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
						param := lit.Value // This is the parameter name in string form, including quotes
						if apiInfo, exists := apiMap[currentFuncName]; exists {
							apiInfo.Params = append(apiInfo.Params, param)
						} else {
							apiInfo := models.Draft{
								Params: []string{param},
							}
							apiMap[currentFuncName] = &apiInfo
						}
					}
				}
			}
		}
		return true
	})
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
			processFuncDecl(v, outFile)
		case *ast.CallExpr:
			if selectorExpr, ok := v.Fun.(*ast.SelectorExpr); ok {
				if _, ok := selectorExpr.X.(*ast.Ident); ok {
					methodName := selectorExpr.Sel.Name
					// 從 controller 提取 method & url
					if methodName == "GET" || methodName == "POST" || methodName == "PUT" || methodName == "DELETE" {
						if len(v.Args) > 0 {
							funcExpr, ok := v.Args[1].(*ast.SelectorExpr)
							if ok {
								calledFuncName := funcExpr.Sel.Name
								urlExpr, ok := v.Args[0].(*ast.BasicLit)
								if ok {
									if apiInfo, exists := apiMap[calledFuncName]; exists {
										apiInfo.Method = methodName
										apiInfo.URL = urlExpr.Value
									} else {
										apiInfo := models.Draft{
											Method: methodName,
											URL:    urlExpr.Value,
										}
										apiMap[calledFuncName] = &apiInfo
									}
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

func ConvertStructToYaml(target models.Draft) {
	outFile, err := os.Create("result.yaml")
	if err != nil {
		fmt.Println("Error creating yaml file:", err)
		return
	}
	defer outFile.Close()

	b, err := yaml.Marshal(target)
	if err != nil {
		fmt.Println("Marshal yaml fail:", err)
		return
	}
	err = ioutil.WriteFile("result.yaml", b, 0644)
	if err != nil {
		fmt.Println("Error creating yaml file:", err)
		return
	}
}

func ConvertStructToJson() {
	jsonData, err := json.MarshalIndent(apiMap, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling json:", err)
		return
	}

	err = ioutil.WriteFile("drafts.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
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

	ConvertStructToJson()

	for k, v := range apiMap {
		fmt.Println(k, v)
		// ConvertStructToYaml(*v)
	}

}
