// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protogen

import (
	"flag"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/internal/scalar"
	"google.golang.org/protobuf/reflect/protoreflect"

	descriptorpb "google.golang.org/protobuf/types/descriptor"
	pluginpb "google.golang.org/protobuf/types/plugin"
)

func TestPluginParameters(t *testing.T) {
	var flags flag.FlagSet
	value := flags.Int("integer", 0, "")
	opts := &Options{
		ParamFunc: flags.Set,
	}
	const params = "integer=2"
	_, err := New(&pluginpb.CodeGeneratorRequest{
		Parameter: scalar.String(params),
	}, opts)
	if err != nil {
		t.Errorf("New(generator parameters %q): %v", params, err)
	}
	if *value != 2 {
		t.Errorf("New(generator parameters %q): integer=%v, want 2", params, *value)
	}
}

func TestPluginParameterErrors(t *testing.T) {
	for _, parameter := range []string{
		"unknown=1",
		"boolean=error",
	} {
		var flags flag.FlagSet
		flags.Bool("boolean", false, "")
		opts := &Options{
			ParamFunc: flags.Set,
		}
		_, err := New(&pluginpb.CodeGeneratorRequest{
			Parameter: scalar.String(parameter),
		}, opts)
		if err == nil {
			t.Errorf("New(generator parameters %q): want error, got nil", parameter)
		}
	}
}

func TestNoGoPackage(t *testing.T) {
	gen, err := New(&pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{
			{
				Name:    scalar.String("testdata/go_package/no_go_package.proto"),
				Syntax:  scalar.String(protoreflect.Proto3.String()),
				Package: scalar.String("goproto.testdata"),
			},
			{
				Name:       scalar.String("testdata/go_package/no_go_package_import.proto"),
				Syntax:     scalar.String(protoreflect.Proto3.String()),
				Package:    scalar.String("goproto.testdata"),
				Dependency: []string{"go_package/no_go_package.proto"},
			},
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	for i, f := range gen.Files {
		if got, want := string(f.GoPackageName), "goproto_testdata"; got != want {
			t.Errorf("gen.Files[%d].GoPackageName = %v, want %v", i, got, want)
		}
		if got, want := string(f.GoImportPath), "testdata/go_package"; got != want {
			t.Errorf("gen.Files[%d].GoImportPath = %v, want %v", i, got, want)
		}
	}
}

func TestPackageNamesAndPaths(t *testing.T) {
	const (
		filename         = "dir/filename.proto"
		protoPackageName = "proto.package"
	)
	for _, test := range []struct {
		desc               string
		parameter          string
		goPackageOption    string
		generate           bool
		wantPackageName    GoPackageName
		wantImportPath     GoImportPath
		wantFilenamePrefix string
	}{
		{
			desc:               "no parameters, no go_package option",
			generate:           true,
			wantPackageName:    "proto_package",
			wantImportPath:     "dir",
			wantFilenamePrefix: "dir/filename",
		},
		{
			desc:               "go_package option sets import path",
			goPackageOption:    "golang.org/x/foo",
			generate:           true,
			wantPackageName:    "foo",
			wantImportPath:     "golang.org/x/foo",
			wantFilenamePrefix: "golang.org/x/foo/filename",
		},
		{
			desc:               "go_package option sets import path and package",
			goPackageOption:    "golang.org/x/foo;bar",
			generate:           true,
			wantPackageName:    "bar",
			wantImportPath:     "golang.org/x/foo",
			wantFilenamePrefix: "golang.org/x/foo/filename",
		},
		{
			desc:               "go_package option sets package",
			goPackageOption:    "foo",
			generate:           true,
			wantPackageName:    "foo",
			wantImportPath:     "dir",
			wantFilenamePrefix: "dir/filename",
		},
		{
			desc:               "command line sets import path for a file",
			parameter:          "Mdir/filename.proto=golang.org/x/bar",
			goPackageOption:    "golang.org/x/foo",
			generate:           true,
			wantPackageName:    "foo",
			wantImportPath:     "golang.org/x/bar",
			wantFilenamePrefix: "golang.org/x/foo/filename",
		},
		{
			desc:               "import_path parameter sets import path of generated files",
			parameter:          "import_path=golang.org/x/bar",
			goPackageOption:    "golang.org/x/foo",
			generate:           true,
			wantPackageName:    "foo",
			wantImportPath:     "golang.org/x/bar",
			wantFilenamePrefix: "golang.org/x/foo/filename",
		},
		{
			desc:               "import_path parameter does not set import path of dependencies",
			parameter:          "import_path=golang.org/x/bar",
			goPackageOption:    "golang.org/x/foo",
			generate:           false,
			wantPackageName:    "foo",
			wantImportPath:     "golang.org/x/foo",
			wantFilenamePrefix: "golang.org/x/foo/filename",
		},
	} {
		context := fmt.Sprintf(`
TEST: %v
  --go_out=%v:.
  file %q: generate=%v
  option go_package = %q;

  `,
			test.desc, test.parameter, filename, test.generate, test.goPackageOption)

		req := &pluginpb.CodeGeneratorRequest{
			Parameter: scalar.String(test.parameter),
			ProtoFile: []*descriptorpb.FileDescriptorProto{
				{
					Name:    scalar.String(filename),
					Package: scalar.String(protoPackageName),
					Options: &descriptorpb.FileOptions{
						GoPackage: scalar.String(test.goPackageOption),
					},
				},
			},
		}
		if test.generate {
			req.FileToGenerate = []string{filename}
		}
		gen, err := New(req, nil)
		if err != nil {
			t.Errorf("%vNew(req) = %v", context, err)
			continue
		}
		gotFile, ok := gen.FileByName(filename)
		if !ok {
			t.Errorf("%v%v: missing file info", context, filename)
			continue
		}
		if got, want := gotFile.GoPackageName, test.wantPackageName; got != want {
			t.Errorf("%vGoPackageName=%v, want %v", context, got, want)
		}
		if got, want := gotFile.GoImportPath, test.wantImportPath; got != want {
			t.Errorf("%vGoImportPath=%v, want %v", context, got, want)
		}
		if got, want := gotFile.GeneratedFilenamePrefix, test.wantFilenamePrefix; got != want {
			t.Errorf("%vGeneratedFilenamePrefix=%v, want %v", context, got, want)
		}
	}
}

func TestPackageNameInference(t *testing.T) {
	gen, err := New(&pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{
			{
				Name:    scalar.String("dir/file1.proto"),
				Package: scalar.String("proto.package"),
			},
			{
				Name:    scalar.String("dir/file2.proto"),
				Package: scalar.String("proto.package"),
				Options: &descriptorpb.FileOptions{
					GoPackage: scalar.String("foo"),
				},
			},
		},
		FileToGenerate: []string{"dir/file1.proto", "dir/file2.proto"},
	}, nil)
	if err != nil {
		t.Fatalf("New(req) = %v", err)
	}
	if f1, ok := gen.FileByName("dir/file1.proto"); !ok {
		t.Errorf("missing file info for dir/file1.proto")
	} else if f1.GoPackageName != "foo" {
		t.Errorf("dir/file1.proto: GoPackageName=%v, want foo; package name should be derived from dir/file2.proto", f1.GoPackageName)
	}
}

func TestInconsistentPackageNames(t *testing.T) {
	_, err := New(&pluginpb.CodeGeneratorRequest{
		ProtoFile: []*descriptorpb.FileDescriptorProto{
			{
				Name:    scalar.String("dir/file1.proto"),
				Package: scalar.String("proto.package"),
				Options: &descriptorpb.FileOptions{
					GoPackage: scalar.String("golang.org/x/foo"),
				},
			},
			{
				Name:    scalar.String("dir/file2.proto"),
				Package: scalar.String("proto.package"),
				Options: &descriptorpb.FileOptions{
					GoPackage: scalar.String("golang.org/x/foo;bar"),
				},
			},
		},
		FileToGenerate: []string{"dir/file1.proto", "dir/file2.proto"},
	}, nil)
	if err == nil {
		t.Fatalf("inconsistent package names for the same import path: New(req) = nil, want error")
	}
}

func TestImports(t *testing.T) {
	gen, err := New(&pluginpb.CodeGeneratorRequest{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	g := gen.NewGeneratedFile("foo.go", "golang.org/x/foo")
	g.P("package foo")
	g.P()
	for _, importPath := range []GoImportPath{
		"golang.org/x/foo",
		// Multiple references to the same package.
		"golang.org/x/bar",
		"golang.org/x/bar",
		// Reference to a different package with the same basename.
		"golang.org/y/bar",
		"golang.org/x/baz",
		// Reference to a package conflicting with a predeclared identifier.
		"golang.org/z/string",
	} {
		g.P("var _ = ", GoIdent{GoName: "X", GoImportPath: importPath}, " // ", importPath)
	}
	want := `package foo

import (
	bar "golang.org/x/bar"
	baz "golang.org/x/baz"
	bar1 "golang.org/y/bar"
	string1 "golang.org/z/string"
)

var _ = X         // "golang.org/x/foo"
var _ = bar.X     // "golang.org/x/bar"
var _ = bar.X     // "golang.org/x/bar"
var _ = bar1.X    // "golang.org/y/bar"
var _ = baz.X     // "golang.org/x/baz"
var _ = string1.X // "golang.org/z/string"
`
	got, err := g.Content()
	if err != nil {
		t.Fatalf("g.Content() = %v", err)
	}
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Fatalf("content mismatch (-want +got):\n%s", diff)
	}
}

func TestImportRewrites(t *testing.T) {
	gen, err := New(&pluginpb.CodeGeneratorRequest{}, &Options{
		ImportRewriteFunc: func(i GoImportPath) GoImportPath {
			return "prefix/" + i
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	g := gen.NewGeneratedFile("foo.go", "golang.org/x/foo")
	g.P("package foo")
	g.P("var _ = ", GoIdent{GoName: "X", GoImportPath: "golang.org/x/bar"})
	want := `package foo

import (
	bar "prefix/golang.org/x/bar"
)

var _ = bar.X
`
	got, err := g.Content()
	if err != nil {
		t.Fatalf("g.Content() = %v", err)
	}
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Fatalf("content mismatch (-want +got):\n%s", diff)
	}
}
