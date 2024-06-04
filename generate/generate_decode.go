package generate

import (
	"log"
	"strings"

	"github.com/Nyarum/diho_bytes_generate/customtypes"

	"github.com/dave/jennifer/jen"
)

func GenerateDecodeForStruct(filename, pkg string, packetDescrs []customtypes.PacketDescr) {
	f := jen.NewFilePathName("", pkg)

	f.HeaderComment("Code generated by diho_bytes_generate " + filename + "; DO NOT EDIT.")

	for _, packetDescr := range packetDescrs {
		body := []jen.Code{
			jen.Var().Id("err").Error(),
			jen.Id("reader").Op(":=").Qual("bytes", "NewReader").Call(jen.Id("buf")),
		}

		for _, field := range packetDescr.FieldsWithTypes.Keys() {
			fieldInfo, _ := packetDescr.FieldsWithTypes.Get(field)

			if !fieldInfo.IsArray {
				switch fieldInfo.TypeName {
				case "uint16", "uint32", "uint64", "uint8", "int16", "int32", "int64", "int8", "bool":
					body = append(body, []jen.Code{
						jen.Err().Op("=").Qual("encoding/binary", "Read").Call(jen.Id("reader"), jen.Id("endian"), jen.Op("&").Id("p").Dot(field)),
						jen.If(jen.Err().Op("!=").Nil()).Block(
							jen.Return(jen.Err()),
						),
					}...)
				case "string":
					body = append(body, []jen.Code{
						jen.Id("p").Dot(field).Op(",").Id("err").Op("=").Qual("github.com/Nyarum/diho_bytes_generate/utils", "ReadStringNull").Call(jen.Id("reader")),
						jen.If(jen.Err().Op("!=").Nil()).Block(
							jen.Return(jen.Err()),
						),
					}...)
				default:
					endianSwitch := jen.Id("endian")
					if fieldInfo.IsLittle {
						endianSwitch = jen.Qual("encoding/binary", "LittleEndian")
					}

					body = append(body, []jen.Code{
						jen.If(jen.Err().Op("=").Parens(jen.Id("&").Id("p").Dot(field)).Dot("Decode").Call(jen.Id("ctx"), jen.Id("buf"), endianSwitch),
							jen.Err().Op("!=").Nil()).Block(
							jen.Return(jen.Err()),
						),
					}...)
				}
			} else {
				switch fieldInfo.TypeName {
				case "uint16", "uint32", "uint64", "uint8", "int16", "int32", "int64", "int8", "bool":
					body = append(body, []jen.Code{
						jen.For(jen.Id("k").Op(":=").Range().Id("p").Dot(field)).Block(
							jen.Var().Id("tempValue").Id(fieldInfo.TypeName),
							jen.If(jen.Err().Op("=").Qual("encoding/binary", "Read").Call(jen.Id("reader"), jen.Id("endian"), jen.Op("&").Id("tempValue")),
								jen.Err().Op("!=").Nil()).Block(
								jen.Return(jen.Err()),
							),
							jen.Id("p").Dot(field).Index(jen.Id("k")).Op("=").Id("tempValue"),
						),
					}...)
				case "byte":
					body = append(body, []jen.Code{
						jen.Id("p").Dot(field).Op(",").Id("err").Op("=").Qual("github.com/Nyarum/diho_bytes_generate/utils", "ReadBytes").Call(jen.Id("reader")),
						jen.If(jen.Err().Op("!=").Nil()).Block(
							jen.Return(jen.Err()),
						),
					}...)
				default:
					body = append(body, []jen.Code{
						jen.For(jen.Id("k").Op(":=").Range().Id("p").Dot(field)).Block(
							jen.If(jen.Err().Op("=").Parens(jen.Op("&").Id("p").Dot(field).Index(jen.Id("k"))).Dot("Decode").Call(jen.Id("ctx"), jen.Id("buf"), jen.Id("endian")),
								jen.Err().Op("!=").Nil()).Block(
								jen.Return(jen.Err()),
							),
						),
					}...)
				}
			}

			if packetDescr.IsFilterMethod {
				body = append(body, []jen.Code{
					jen.If(jen.Id("p").Dot("Filter").Call(jen.Id("ctx"))).Op("==").Id("true").Block(
						jen.Return(
							jen.Err(),
						),
					),
				}...)
			}
		}

		body = append(body, jen.Return(
			jen.Nil(),
		))

		f.Func().Params(jen.Id("p").Op("*").Id(packetDescr.StructName)).Id("Decode").Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("buf").Index().Byte(), jen.Id("endian").Qual("encoding/binary", "ByteOrder"),
		).Params(
			jen.Error(),
		).Block(body...)
	}

	outputFilename := strings.TrimSuffix(filename, ".go") + "_decode.gen.go"
	if err := f.Save(outputFilename); err != nil {
		log.Fatalf("Failed to save file: %s", err)
	}
}
