package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"

	"github.com/Nyarum/diho_bytes_generate/customtypes"
	"github.com/fatih/structtag"

	"github.com/elliotchance/orderedmap/v2"
)

func isNormalType(t string) bool {
	switch t {
	case "uint16", "uint32", "uint64", "uint8", "int16", "int32", "int64", "int8", "string", "byte":
		return true
	default:
		return false
	}
}

func ParseBinaryFile(filename string) (pkgName string, packetsDescrs []customtypes.PacketDescr) {
	// Parse the Go source file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse file: %s", err)
	}

	pkgName = node.Name.Name

	packetsDescrs = make([]customtypes.PacketDescr, 0)

	filterMethodsByPacket := make(map[string]bool)

	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				// Check if the receiver is *Packet
				if starExpr, ok := fn.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok {
						// Check if the function name is Filter
						if fn.Name.Name == "Filter" {
							filterMethodsByPacket[ident.Name] = true
						}
					}
				}
			}
		}

		return true
	})

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			packetDescr := customtypes.PacketDescr{
				FieldsWithTypes: orderedmap.NewOrderedMap[string, customtypes.Field](),
				PackageName:     node.Name.Name,
			}

			packetDescr.StructName = typeSpec.Name.Name

			if v, ok := filterMethodsByPacket[packetDescr.StructName]; ok && v {
				packetDescr.IsFilterMethod = true
			}

			if v := typeSpec.Type.(*ast.StructType); v != nil {
				packetDescr.StructName = typeSpec.Name.Name

			outerFor:
				for _, field := range v.Fields.List {
					var isLittle bool
					if field.Tag != nil {
						fmt.Println(field.Tag.Value)

						tags, err := structtag.Parse(field.Tag.Value)
						if err != nil {
							fmt.Println("can't parse field tag", err)
						}

						dbgTag, err := tags.Get("dbg")
						if err != nil {
							fmt.Println("can't get dbg tag", err)
						}

						for _, option := range dbgTag.Options {
							if option == "ignore" {
								continue outerFor
							}

							if option == "little" {
								isLittle = true
							}
						}
					}

					if field.Tag != nil && strings.Contains(field.Tag.Value, "ignore") {
						continue
					}

					if v, ok := field.Type.(*ast.Ident); ok {
						packetDescr.FieldsWithTypes.Set(field.Names[0].Name, customtypes.Field{
							TypeName: v.Name,
							IsLittle: isLittle,
						})
					}

					if v, ok := field.Type.(*ast.ArrayType); ok {
						if v, ok := v.Elt.(*ast.Ident); ok {
							packetDescr.FieldsWithTypes.Set(field.Names[0].Name, customtypes.Field{
								IsArray:  true,
								TypeName: v.Name,
								IsLittle: isLittle,
							})
						}
					}
				}
			}

			packetsDescrs = append(packetsDescrs, packetDescr)
		}
	}

	return
}
