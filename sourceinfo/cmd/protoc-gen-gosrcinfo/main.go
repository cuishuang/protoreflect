// Command protoc-gen-gosrcinfo is a protoc plugin. It emits Go code, into files
// named "<file>.pb.srcinfo.go". These source files include source code info for
// processed proto files and register that info with the srcinfo package.
package main

import (
	"bytes"
	"compress/gzip"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	protogen.Options{}.Run(genSourceInfo)
}

func genSourceInfo(plugin *protogen.Plugin) error {
	plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL |
		pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
	plugin.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_2023
	plugin.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_2023
	for _, f := range plugin.Files {
		if f.Generate {
			if err := generateSourceInfo(f, plugin); err != nil {
				return fmt.Errorf("%s: %v", f.Desc.Path(), err)
			}
		}
	}
	return nil
}

func generateSourceInfo(f *protogen.File, plugin *protogen.Plugin) error {
	si := f.Proto.GetSourceCodeInfo()
	if len(si.GetLocation()) == 0 {
		return nil
	}
	out := plugin.NewGeneratedFile(f.GeneratedFilenamePrefix+".pb.srcinfo.go", f.GoImportPath)
	out.P("// Code generated by protoc-gen-gosrcinfo. DO NOT EDIT.")
	out.P("// source: ", f.Desc.Path())
	out.P()
	out.P("package ", f.GoPackageName)
	out.P()

	siBytes, err := proto.Marshal(si)
	if err != nil {
		return fmt.Errorf("failed to serialize source code info: %w", err)
	}
	// Source code info has lots of repeated sequences of bytes in the 'path' field
	// of locations, and the comments tend to be text that is reasonably well
	// compressed. So we can make binaries a little smaller by gzipping first.
	var encodedBuf bytes.Buffer
	zipWriter := gzip.NewWriter(&encodedBuf)
	if _, err := zipWriter.Write(siBytes); err != nil {
		return fmt.Errorf("failed to compress source code info: %w", err)
	}
	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("failed to compress source code info: %w", err)
	}

	srcInfoPkg := protogen.GoImportPath("github.com/jhump/protoreflect/v2/sourceinfo")

	out.P("func init() {")
	out.P("  srcInfo := []byte{")
	encodedBytes := encodedBuf.Bytes()
	var buf bytes.Buffer
	for len(encodedBytes) > 0 {
		var chunk []byte
		if len(encodedBytes) < 16 {
			chunk = encodedBytes
			encodedBytes = nil
		} else {
			chunk = encodedBytes[:16]
			encodedBytes = encodedBytes[16:]
		}
		buf.Reset()
		for _, b := range chunk {
			_, _ = fmt.Fprintf(&buf, " 0x%02x,", b)
		}
		out.P(buf.String())
	}
	out.P("  }")
	out.P("  ", srcInfoPkg.Ident("Register"), "(", fmt.Sprintf("%q", f.Desc.Path()), ", srcInfo)")
	out.P("}")
	return nil
}
