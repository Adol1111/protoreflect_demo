package main

import (
	"errors"
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"google.golang.org/genproto/googleapis/api/annotations"
)

func main() {
	pb, err := loadProto("a.proto", "./example", "./third_parts")
	if err != nil {
		fmt.Println(err)
	}
	pbs := getAllFile(nil, pb)
	opts := protogen.Options{}

	names := make([]string, 0, len(pbs))
	for _, fdp := range pbs {
		names = append(names, *fdp.Name)
	}
	plugin, err := opts.New(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: names,
		ProtoFile:      pbs,
	})
	for _, file := range plugin.Files {
		for _, service := range file.Services {
			for _, method := range service.Methods {
				proto.GetExtension(method.Desc.Options(), annotations.E_Http)
			}
		}
	}
	fmt.Println("ok")
}

func loadProto(entryFile string, dirs ...string) (*desc.FileDescriptor, error) {
	parser := protoparse.Parser{
		IncludeSourceCodeInfo: true,
		ImportPaths:           dirs,
	}
	ds, err := parser.ParseFiles(entryFile)
	if err != nil {
		return nil, err
	}
	if len(ds) <= 0 {
		return nil, errors.New("invalid proto file")
	}
	return ds[0], nil
}

func getAllFile(files []*descriptorpb.FileDescriptorProto, proto *desc.FileDescriptor) []*descriptorpb.FileDescriptorProto {
	if proto.GetDependencies() != nil {
		for _, fd := range proto.GetDependencies() {
			files = getAllFile(files, fd)
		}
	}
	for _, fdp := range files {
		if fdp.GetName() == proto.GetName() {
			return files
		}
	}
	files = append(files, proto.AsFileDescriptorProto())
	return files
}
