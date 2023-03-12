package main

import (
	"fmt"
	"os"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const deprecationComment = "// Deprecated: Do not use."

const (
	methodGet     = "Get"
	methodHead    = "Head"
	methodPost    = "Post"
	methodPut     = "Put"
	methodPatch   = "Patch"
	MethodDelete  = "Delete"
	methodConnect = "Connect"
	methodOptions = "Options"
	methodTrace   = "Trace"
)

var (
	contextPackage = protogen.GoImportPath("context")
	netHttpPackage = protogen.GoImportPath("github.com/iobrother/zoo/core/transport/http")
)

var methodSets = make(map[string]int)

func run(gen *protogen.Plugin) error {
	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	for _, f := range gen.Files {
		if !f.Generate {
			continue
		}
		generateFile(gen, f, *omitempty)
	}
	return nil
}

func generateFile(gen *protogen.Plugin, file *protogen.File, omitempty bool) *protogen.GeneratedFile {
	if len(file.Services) == 0 || (omitempty && !hasHTTPRule(file.Services)) {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_http.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	generateHeader(gen, file, g)
	generateImports(g)
	generateFileContent(file, g, omitempty)
	return g
}

func generateHeader(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	g.P("// Code generated by protoc-gen-zoo-http. DO NOT EDIT.")
	g.P("// versions:")
	g.P("// - protoc-gen-zoo-http v", version)
	g.P("// - protoc                ", protocVersion(gen))
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("// ", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
}

func generateImports(g *protogen.GeneratedFile) {
	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var _ = ", contextPackage.Ident("TODO"))
	g.P("var _ = ", netHttpPackage.Ident("NewServer"))
	g.P()
}

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}

// generateFileContent generates the errors definitions, excluding the package statement.
func generateFileContent(file *protogen.File, g *protogen.GeneratedFile, omitempty bool) {
	if len(file.Services) == 0 {
		return
	}
	for _, service := range file.Services {
		genService(file, g, service, omitempty)
	}
}

func genService(file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, omitempty bool) {
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}
	// HTTP Server.
	sd := &serviceDesc{
		ServiceType: service.GoName,
		ServiceName: string(service.Desc.FullName()),
		Metadata:    file.Desc.Path(),
	}
	for _, method := range service.Methods {
		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			continue
		}
		rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
		if rule != nil && ok {
			for _, bind := range rule.AdditionalBindings {
				sd.Methods = append(sd.Methods, buildHTTPRule(g, method, bind))
			}
			sd.Methods = append(sd.Methods, buildHTTPRule(g, method, rule))
		} else if !omitempty {
			path := fmt.Sprintf("/%s/%s", service.Desc.FullName(), method.Desc.Name())
			sd.Methods = append(sd.Methods, buildMethodDesc(g, method, methodPost, path))
		}
	}
	if len(sd.Methods) == 0 {
		return
	}
	err := sd.execute(g)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr,
			"\u001B[31mWARN\u001B[m: execute template failed.\n")
	}
}

func hasHTTPRule(services []*protogen.Service) bool {
	for _, service := range services {
		for _, method := range service.Methods {
			if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
				continue
			}
			rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
			if rule != nil && ok {
				return true
			}
		}
	}
	return false
}

func buildHTTPRule(g *protogen.GeneratedFile, m *protogen.Method, rule *annotations.HttpRule) *methodDesc {
	var (
		path         string
		method       string
		body         string
		responseBody string
	)
	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		path = pattern.Get
		method = methodGet
	case *annotations.HttpRule_Put:
		path = pattern.Put
		method = methodPut
	case *annotations.HttpRule_Post:
		path = pattern.Post
		method = methodPost
	case *annotations.HttpRule_Delete:
		path = pattern.Delete
		method = MethodDelete
	case *annotations.HttpRule_Patch:
		path = pattern.Patch
		method = methodPatch
	case *annotations.HttpRule_Custom:
		path = pattern.Custom.Path
		method = pattern.Custom.Kind
	}
	body = rule.Body
	responseBody = rule.ResponseBody
	md := buildMethodDesc(g, m, method, path)
	switch {
	case method == methodGet:
		if body != "" {
			_, _ = fmt.Fprintf(os.Stderr,
				"\u001B[31mWARN\u001B[m: %s %s body should not be declared.\n", method, path)
		}
		md.HasBody = false
	case method == MethodDelete:
		if body != "" {
			md.HasBody = true
			if !*allowDeleteBody {
				md.HasBody = false
				_, _ = fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: %s %s body should not be declared.\n", method, path)
			}
		} else {
			md.HasBody = false
		}
	case method == methodPatch:
		if body != "" {
			md.HasBody = true
		} else {
			md.HasBody = false
			if !*allowEmptyPatchBody {
				_, _ = fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: %s %s is does not declare a body.\n", method, path)
			}
		}
	case body == "*":
		md.HasBody = true
		md.Body = ""
	case body != "":
		md.HasBody = true
		md.Body = "." + camelCaseVars(body)
	default:
		md.HasBody = false
		_, _ = fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: %s %s is does not declare a body.\n", method, path)
	}
	if responseBody == "*" {
		md.ResponseBody = ""
	} else if responseBody != "" {
		md.ResponseBody = "." + camelCaseVars(responseBody)
	}
	return md
}

func buildMethodDesc(g *protogen.GeneratedFile, m *protogen.Method, method, path string) *methodDesc {
	defer func() { methodSets[m.GoName]++ }()
	params := extractParams(m, path)
	fields := m.Input.Desc.Fields()
	for _, v := range params {
		for _, field := range strings.Split(v, ".") {
			if strings.TrimSpace(field) == "" {
				continue
			}
			if strings.Contains(field, ":") {
				field = strings.Split(field, ":")[0]
			}
			fd := fields.ByName(protoreflect.Name(field))
			if fd == nil {
				fmt.Fprintf(os.Stderr, "\u001B[31mERROR\u001B[m: The corresponding field '%s' declaration in message could not be found in '%s'\n", v, path)
				os.Exit(2)
			}
			switch {
			case fd.IsMap():
				fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: The field in path:'%s' shouldn't be a map.\n", v)
			case fd.IsList():
				fmt.Fprintf(os.Stderr, "\u001B[31mWARN\u001B[m: The field in path:'%s' shouldn't be a list.\n", v)
			case fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind:
				fields = fd.Message().Fields()
			}
		}
	}
	comment := m.Comments.Leading.String() + m.Comments.Trailing.String()
	if comment != "" {
		comment = strings.TrimPrefix(strings.TrimSuffix(comment, "\n"), "//")
		comment = "// " + m.GoName + comment
	} else {
		comment = "// " + m.GoName + " ..."
	}
	return &methodDesc{
		Name:      m.GoName,
		Num:       methodSets[m.GoName],
		Request:   g.QualifiedGoIdent(m.Input.GoIdent),
		Reply:     g.QualifiedGoIdent(m.Output.GoIdent),
		Path:      transformPath(path),
		Method:    method,
		Comment:   comment,
		HasParams: len(params) > 0,
	}
}

// transformPath {xxx} --> :xxx
func transformPath(path string) string {
	paths := strings.Split(path, "/")
	for i, p := range paths {
		if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") || strings.HasPrefix(p, ":") {
			paths[i] = ":" + p[1:len(p)-1]
		}
	}

	return strings.Join(paths, "/")
}

func extractParams(_ *protogen.Method, path string) (params []string) {
	for _, v := range strings.Split(path, "/") {
		if strings.HasPrefix(v, "{") && strings.HasSuffix(v, "}") {
			params = append(params, strings.TrimSuffix(strings.TrimPrefix(v, "{"), "}"))
		}
	}
	return
}

func camelCaseVars(s string) string {
	vars := make([]string, 0)
	subs := strings.Split(s, ".")
	for _, sub := range subs {
		vars = append(vars, camelCase(sub))
	}
	return strings.Join(vars, ".")
}

// camelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func camelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
