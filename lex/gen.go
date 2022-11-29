package lex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	EncodingCBOR = "application/cbor"
	EncodingJSON = "application/json"
)

type Schema struct {
	prefix string

	Lexicon     int                   `json:"lexicon"`
	ID          string                `json:"id"`
	Type        string                `json:"type"`
	Key         string                `json:"key"`
	Description string                `json:"description"`
	Parameters  *TypeSchema           `json:"parameters"`
	Input       *InputType            `json:"input"`
	Output      *OutputType           `json:"output"`
	Defs        map[string]TypeSchema `json:"defs"`
	Record      *TypeSchema           `json:"record"`
}

type Param struct {
	Type     string `json:"type"`
	Maximum  int    `json:"maximum"`
	Required bool   `json:"required"`
}

type OutputType struct {
	Encoding string     `json:"encoding"`
	Schema   TypeSchema `json:"schema"`
}

type InputType struct {
	Encoding string     `json:"encoding"`
	Schema   TypeSchema `json:"schema"`
}

type TypeSchema struct {
	Type       string                `json:"type"`
	Ref        string                `json:"$ref"`
	Required   []string              `json:"required"`
	Properties map[string]TypeSchema `json:"properties"`
	MaxLength  int                   `json:"maxLength"`
	Items      *TypeSchema           `json:"items"`
	OneOf      []TypeSchema          `json:"oneOf"`
	Const      *string               `json:"const"`
	Enum       []string              `json:"enum"`
	Not        *TypeSchema           `json:"not"`
}

func (s *Schema) Name() string {
	p := strings.Split(s.ID, ".")
	return p[len(p)-2] + p[len(p)-1]
}

type outputType struct {
	Name   string
	Type   TypeSchema
	Record bool
}

func (s *Schema) AllTypes(prefix string) []outputType {
	var out []outputType

	var walk func(name string, ts TypeSchema, record bool)
	walk = func(name string, ts TypeSchema, record bool) {
		if ts.Type == "object" ||
			(ts.Type == "" && len(ts.OneOf) > 0) {
			out = append(out, outputType{
				Name:   name,
				Type:   ts,
				Record: record,
			})
		}

		for childname, val := range ts.Properties {
			walk(name+"_"+strings.Title(childname), val, false)
		}

		if ts.Items != nil {
			walk(name+"_Elem", *ts.Items, false)
		}
	}

	tname := s.nameFromID(s.ID, prefix)

	for name, def := range s.Defs {
		walk(tname+"_"+strings.Title(name), def, false)
	}

	if s.Input != nil {
		walk(tname+"_Input", s.Input.Schema, false)
	}
	if s.Output != nil {
		walk(tname+"_Output", s.Output.Schema, false)
	}

	if s.Type == "record" {
		walk(tname, *s.Record, false)
	}

	return out
}

func ReadSchema(f string) (*Schema, error) {
	fi, err := os.Open(f)
	if err != nil {
		return nil, err
	}

	var s Schema
	if err := json.NewDecoder(fi).Decode(&s); err != nil {
		return nil, err
	}

	return &s, nil
}

func GenCodeForSchema(pkg string, prefix string, fname string, reqcode bool, s *Schema) error {
	buf := new(bytes.Buffer)

	s.prefix = prefix

	fmt.Fprintf(buf, "package %s\n\n", pkg)
	fmt.Fprintf(buf, "import (\n")
	fmt.Fprintf(buf, "\t\"context\"\n")
	fmt.Fprintf(buf, "\t\"fmt\"\n")
	fmt.Fprintf(buf, "\t\"encoding/json\"\n")
	fmt.Fprintf(buf, "\t\"github.com/whyrusleeping/gosky/xrpc\"\n")
	fmt.Fprintf(buf, "\t\"github.com/whyrusleeping/gosky/lex/util\"\n")
	fmt.Fprintf(buf, ")\n\n")
	fmt.Fprintf(buf, "// schema: %s\n\n", s.ID)

	tps := s.AllTypes(prefix)

	for _, ot := range tps {
		if err := s.WriteType(ot.Name, ot.Type, buf); err != nil {
			return err
		}
	}

	if reqcode {
		if err := writeMethods(prefix, s, buf); err != nil {
			return err
		}
	}

	if err := writeCodeFile(buf.Bytes(), fname); err != nil {
		return err
	}

	return nil
}

func writeCodeFile(b []byte, fname string) error {
	formatted, err := format.Source(b)
	if err != nil {
		fmt.Println(string(b))
		return fmt.Errorf("failed to format generated file: %w", err)
	}

	fixed, err := fixImports(formatted)
	if err != nil {
		return err
	}

	if err := os.WriteFile(fname, fixed, 0664); err != nil {
		return err
	}

	return nil
}

func fixImports(b []byte) ([]byte, error) {
	cmd := exec.Command("goimports")

	cmd.Stdin = bytes.NewReader(b)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func writeMethods(prefix string, s *Schema, w io.Writer) error {
	switch s.Type {
	case "token":
		n := s.nameFromID(s.ID, prefix)
		fmt.Fprintf(w, "const %s = %q\n", n, s.ID)
		return nil
	case "record":
		return nil
	case "query":
		return s.WriteRPC(w, prefix)
	case "procedure":
		return s.WriteRPC(w, prefix)
	default:
		return fmt.Errorf("unrecognized lexicon type %q", s.Type)
	}
}

func (s *Schema) nameFromID(id, prefix string) string {
	parts := strings.Split(strings.TrimPrefix(id, prefix), ".")
	var tname string
	for _, s := range parts {
		tname += strings.Title(s)
	}

	return tname

}

func orderedMapIter[T any](m map[string]T, cb func(string, T) error) error {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		if err := cb(k, m[k]); err != nil {
			return err
		}
	}
	return nil
}

func (s *Schema) WriteRPC(w io.Writer, prefix string) error {
	fname := s.nameFromID(s.ID, s.prefix)

	params := "ctx context.Context, c *xrpc.Client"
	inpvar := "nil"
	inpenc := ""
	if s.Input != nil {
		inpvar = "input"
		inpenc = s.Input.Encoding
		switch s.Input.Encoding {
		case EncodingCBOR:
			params = fmt.Sprintf("%s, input io.Reader", params)
		case EncodingJSON:
			params = fmt.Sprintf("%s, input %s_Input", params, fname)
		default:
			return fmt.Errorf("unsupported input encoding: %q", s.Input.Encoding)
		}
	}

	if s.Parameters != nil {
		if err := orderedMapIter[TypeSchema](s.Parameters.Properties, func(name string, t TypeSchema) error {
			tn, err := s.typeNameForField(name, "", t)
			if err != nil {
				return err
			}

			// TODO: deal with optional params
			params = params + fmt.Sprintf(", %s %s", name, tn)
			return nil
		}); err != nil {
			return err
		}
	}

	out := "error"
	if s.Output != nil {
		switch s.Output.Encoding {
		case EncodingCBOR:
			out = "([]byte, error)"
		case EncodingJSON:
			out = fmt.Sprintf("(*%s_Output, error)", fname)
		default:
			return fmt.Errorf("unrecognized encoding scheme: %q", s.Output.Encoding)
		}
	}

	fmt.Fprintf(w, "func %s(%s) %s {\n", fname, params, out)

	outvar := "nil"
	errRet := "err"
	outRet := "nil"
	if s.Output != nil {
		switch s.Output.Encoding {
		case EncodingCBOR:
			fmt.Fprintf(w, "buf := new(bytes.Buffer)\n")
			outvar = "buf"
			errRet = "nil, err"
			outRet = "buf.Bytes(), nil"

		case EncodingJSON:
			fmt.Fprintf(w, "\tvar out %s_Output\n", fname)
			outvar = "&out"
			errRet = "nil, err"
			outRet = "&out, nil"
		default:
			return fmt.Errorf("unrecognized output encoding: %q", s.Output.Encoding)
		}
	}

	queryparams := "nil"
	if s.Parameters != nil {
		queryparams = "params"
		fmt.Fprintf(w, `
	params := map[string]interface{}{
`)
		if err := orderedMapIter[TypeSchema](s.Parameters.Properties, func(name string, t TypeSchema) error {
			fmt.Fprintf(w, `"%s": %s,
`, name, name)
			return nil
		}); err != nil {
			return err
		}
		fmt.Fprintf(w, "}\n")
	}

	var reqtype string
	switch s.Type {
	case "procedure":
		reqtype = "xrpc.Procedure"
	case "query":
		reqtype = "xrpc.Query"
	default:
		return fmt.Errorf("can only generate RPC for Query or Procedure (got %s)", s.Type)
	}

	fmt.Fprintf(w, "\tif err := c.Do(ctx, %s, %q, \"%s\", %s, %s, %s); err != nil {\n", reqtype, inpenc, s.ID, queryparams, inpvar, outvar)
	fmt.Fprintf(w, "\t\treturn %s\n", errRet)
	fmt.Fprintf(w, "\t}\n\n")
	fmt.Fprintf(w, "\treturn %s\n", outRet)
	fmt.Fprintf(w, "}\n\n")

	return nil
}

func doTemplate(w io.Writer, info interface{}, templ string) error {
	t := template.Must(template.New("").
		Funcs(template.FuncMap{
			"TODO": func(thing string) string {
				return "//TODO: " + thing
			},
		}).Parse(templ))

	return t.Execute(w, info)
}

func CreateHandlerStub(pkg string, impmap map[string]string, dir string, schemas []*Schema) error {
	buf := new(bytes.Buffer)

	if err := WriteXrpcServer(buf, schemas, pkg, impmap); err != nil {
		return err
	}

	fname := filepath.Join(dir, "stubs.go")
	if err := writeCodeFile(buf.Bytes(), fname); err != nil {
		return err
	}

	return nil
}

func importNameForPrefix(prefix string) string {
	return strings.Join(strings.Split(prefix, "."), "") + "types"
}

func WriteXrpcServer(w io.Writer, schemas []*Schema, pkg string, impmap map[string]string) error {
	fmt.Fprintf(w, "package %s\n\n", pkg)
	fmt.Fprintf(w, "import (\n")
	fmt.Fprintf(w, "\t\"context\"\n")
	fmt.Fprintf(w, "\t\"fmt\"\n")
	fmt.Fprintf(w, "\t\"encoding/json\"\n")
	fmt.Fprintf(w, "\t\"github.com/whyrusleeping/gosky/xrpc\"\n")
	fmt.Fprintf(w, "\t\"github.com/labstack/echo/v4\"\n")
	for k, v := range impmap {
		fmt.Fprintf(w, "\t%s\"%s\"\n", importNameForPrefix(k), v)
	}
	fmt.Fprintf(w, ")\n\n")

	fmt.Fprintf(w, "func (s *Server) RegisterHandlers(e echo.Echo) error {\n")
	for _, s := range schemas {
		var verb string
		switch s.Type {
		case "query":
			verb = "GET"
		case "procedure":
			verb = "POST"
		default:
			continue
		}

		fmt.Fprintf(w, "e.%s(\"/xrpc/%s\", s.Handle%s)\n", verb, s.ID, idToTitle(s.ID))
	}

	fmt.Fprintf(w, "return nil\n}\n\n")

	for _, s := range schemas {
		var prefix string
		for k := range impmap {
			if strings.HasPrefix(s.ID, k) {
				prefix = k
				break
			}
		}

		if s.Type == "procedure" || s.Type == "query" {
			if err := s.WriteRPCHandler(w, prefix); err != nil {
				return err
			}
		}
	}

	return nil
}

func idToTitle(id string) string {
	var fname string
	for _, p := range strings.Split(id, ".") {
		fname += strings.Title(p)
	}
	return fname
}

func (s *Schema) WriteRPCHandler(w io.Writer, prefix string) error {
	fname := idToTitle(s.ID)

	tname := s.nameFromID(s.ID, prefix)

	fmt.Fprintf(w, "func (s *Server) Handle%s(c echo.Context) error {\n", fname)

	fmt.Fprintf(w, "ctx, span := otel.Tracer(\"server\").Start(c.Request().Context(), %q)\n", "Handle"+fname)
	fmt.Fprintf(w, "defer span.End()\n")

	impname := importNameForPrefix(prefix)

	paramtypes := []string{"ctx context.Context"}
	params := []string{"ctx"}
	if s.Type == "query" {
		if s.Parameters != nil {
			orderedMapIter[TypeSchema](s.Parameters.Properties, func(k string, t TypeSchema) error {
				switch t.Type {
				case "string":
					params = append(params, k)
					paramtypes = append(paramtypes, k+" string")
					fmt.Fprintf(w, "%s := c.QueryParam(\"%s\")\n", k, k)
				case "number":
					params = append(params, k)
					paramtypes = append(paramtypes, k+" int")
					fmt.Fprintf(w, `
%s, err := strconv.Atoi(c.QueryParam("%s"))
if err != nil {
	return err
}
`, k, k)
				default:
					return fmt.Errorf("unsupported handler parameter type: %s", t.Type)
				}
				return nil
			})
		}
	} else if s.Type == "procedure" {
		fmt.Fprintf(w, `
var body %s.%s
if err := c.Bind(&body); err != nil {
	return err
}
`, impname, tname+"_Input")
	} else {
		return fmt.Errorf("can only generate handlers for queries or procedures")
	}

	assign := "handleErr"
	returndef := "error"
	if s.Output != nil {
		assign = "out, handleErr"
		fmt.Fprintf(w, "var out *%s.%s\n", impname, tname+"_Output")
		returndef = fmt.Sprintf("(*%s.%s_Output, error)", impname, tname)
	}
	fmt.Fprintf(w, "var handleErr error\n")
	fmt.Fprintf(w, "// func (s *Server) handle%s(%s) %s\n", fname, strings.Join(paramtypes, ","), returndef)
	fmt.Fprintf(w, "%s = s.handle%s(%s)\n", assign, fname, strings.Join(params, ","))
	fmt.Fprintf(w, "if handleErr != nil {\nreturn handleErr\n}\n")

	fmt.Fprintf(w, "return c.JSON(200, out)\n}\n\n")

	return nil
}

func (s *Schema) typeNameFromRef(r string) string {
	sname := s.nameFromID(s.ID, s.prefix)
	p := strings.Split(r, "/")
	return sname + "_" + strings.Title(p[len(p)-1])
}

func (s *Schema) typeNameForField(name, k string, v TypeSchema) (string, error) {
	switch v.Type {
	case "string":
		return "string", nil
	case "number":
		return "int64", nil
	case "boolean":
		return "bool", nil
	case "object":
		return "*" + name + "_" + strings.Title(k), nil
	case "":
		if v.Ref != "" {
			return "*" + s.typeNameFromRef(v.Ref), nil
		}

		if len(v.OneOf) > 0 {
			return "*" + name + "_" + strings.Title(k), nil
		}

		if v.Const != nil {
			return "string", nil
		}

		return "", fmt.Errorf("field %q in %s does not have discernable type name", k, name)
	case "array":
		subt, err := s.typeNameForField(name+"_"+strings.Title(k), "Elem", *v.Items)
		if err != nil {
			return "", err
		}

		return "[]" + subt, nil
	default:
		return "", fmt.Errorf("field %q in %s has unsupported type name", k, name)
	}
}

func (s *Schema) lookupRef(ref string) (*TypeSchema, error) {
	parts := strings.Split(ref, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid ref: %q", ref)
	}

	if parts[1] != "defs" {
		return nil, fmt.Errorf("ref lookups outside of defs not supported")
	}
	t, ok := s.Defs[parts[2]]
	if !ok {
		return nil, fmt.Errorf("no such def: %q", ref)
	}

	return &t, nil
}

func (s *Schema) WriteType(name string, t TypeSchema, w io.Writer) error {
	name = strings.Title(name)
	if err := s.writeTypeDefinition(name, t, w); err != nil {
		return err
	}

	if err := s.writeTypeMethods(name, t, w); err != nil {
		return err
	}

	return nil
}

func (s *Schema) writeTypeDefinition(name string, t TypeSchema, w io.Writer) error {
	switch t.Type {
	case "string":
		// TODO: deal with max length
		fmt.Fprintf(w, "type %s string\n", name)
	case "number":
		fmt.Fprintf(w, "type %s int64\n", name)
	case "boolean":
		fmt.Fprintf(w, "type %s bool\n", name)
	case "object":
		if len(t.Properties) == 0 {
			fmt.Fprintf(w, "type %s interface{}\n", name)
			return nil
		}

		fmt.Fprintf(w, "type %s struct {\n", name)

		for k, v := range t.Properties {
			goname := strings.Title(k)

			tname, err := s.typeNameForField(name, k, v)
			if err != nil {
				return err
			}

			fmt.Fprintf(w, "\t%s %s `json:\"%s\"`\n", goname, tname, k)
		}
		fmt.Fprintf(w, "}\n\n")

	case "array":
		tname, err := s.typeNameForField(name, "elem", *t.Items)
		if err != nil {
			return err
		}

		fmt.Fprintf(w, "type %s []%s\n", name, tname)

	case "":
		if len(t.OneOf) > 0 {
			// check if this is actually just a string enum
			first, err := s.lookupRef(t.OneOf[0].Ref)
			if err != nil {
				return fmt.Errorf("oneOf pre-check failed: %w", err)
			}

			if first.Type == "string" {
				// okay, this is just a string enum, do something different
				fmt.Fprintf(w, "type %s string\n", name)
			} else {

				fmt.Fprintf(w, "type %s struct {\n", name)
				for _, e := range t.OneOf {
					// TODO: for now, asserting that all enum options are refs
					if e.Ref == "" {
						return fmt.Errorf("Enums must only contain refs")
					}

					tname := s.typeNameFromRef(e.Ref)
					fmt.Fprintf(w, "\t%s *%s\n", tname, tname)
				}
				fmt.Fprintf(w, "}\n\n")
			}

		}
	default:
		return fmt.Errorf("%s has unrecognized type type %s", name, t.Type)
	}

	return nil
}

func (s *Schema) writeTypeMethods(name string, t TypeSchema, w io.Writer) error {
	switch t.Type {
	case "string", "number", "array", "boolean":
		return nil
	case "object":
		if err := s.writeJsonMarshalerObject(name, t, w); err != nil {
			return err
		}

		if err := s.writeJsonUnmarshalerObject(name, t, w); err != nil {
			return err
		}

		return nil
	case "":
		if len(t.OneOf) > 0 {
			reft, err := s.lookupRef(t.OneOf[0].Ref)
			if err != nil {
				return err
			}

			if reft.Type == "string" {
				return nil
			}

			if err := s.writeJsonMarshalerEnum(name, t, w); err != nil {
				return err
			}

			if err := s.writeJsonUnmarshalerEnum(name, t, w); err != nil {
				return err
			}

			return nil
		}

		return fmt.Errorf("%q unsupported for marshaling", name)
	default:
		return fmt.Errorf("%q has unrecognized type type %s", name, t.Type)
	}
}

func forEachProp(t TypeSchema, cb func(k string, ts TypeSchema) error) error {
	var keys []string
	for k := range t.Properties {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		subv := t.Properties[k]

		if err := cb(k, subv); err != nil {
			return err
		}
	}
	return nil
}

func (s *Schema) writeJsonMarshalerObject(name string, t TypeSchema, w io.Writer) error {
	if len(t.Properties) == 0 {
		// TODO: this is a hacky special casing of record types...
		return nil
	}

	fmt.Fprintf(w, "func (t *%s) MarshalJSON() ([]byte, error) {\n", name)

	if err := forEachProp(t, func(k string, ts TypeSchema) error {
		if ts.Const != nil {
			// TODO: maybe check for mutations before overwriting? mutations would mean bad code
			fmt.Fprintf(w, "\tt.%s = %q\n", strings.Title(k), *ts.Const)
		}

		return nil
	}); err != nil {
		return err
	}

	// TODO: this is ugly since i can't just pass things through to json.Marshal without causing an infinite recursion...
	fmt.Fprintf(w, "\tout := make(map[string]interface{})\n")
	if err := forEachProp(t, func(k string, ts TypeSchema) error {
		fmt.Fprintf(w, "\tout[%q] = t.%s\n", k, strings.Title(k))
		return nil
	}); err != nil {
		return err
	}

	fmt.Fprintf(w, "\treturn json.Marshal(out)\n}\n\n")
	return nil
}

func (s *Schema) writeJsonMarshalerEnum(name string, t TypeSchema, w io.Writer) error {
	fmt.Fprintf(w, "func (t *%s) MarshalJSON() ([]byte, error) {\n", name)

	for _, e := range t.OneOf {
		tname := s.typeNameFromRef(e.Ref)
		fmt.Fprintf(w, "\tif t.%s != nil {\n", tname)
		fmt.Fprintf(w, "\t\treturn json.Marshal(t.%s)\n\t}\n", tname)
	}

	fmt.Fprintf(w, "\treturn nil, fmt.Errorf(\"cannot marshal empty enum\")\n}\n")
	return nil
}

func (s *Schema) writeJsonUnmarshalerObject(name string, t TypeSchema, w io.Writer) error {
	// TODO: would be nice to add some validation...
	return nil
	//fmt.Fprintf(w, "func (t *%s) UnmarshalJSON(b []byte) (error) {\n", name)
}

func (s *Schema) getTypeConstValueForType(t TypeSchema) (string, []string, error) {
	parts := strings.Split(t.Ref, "/")
	if len(parts) == 3 && parts[0] == "#" && parts[1] == "defs" {
		defs, ok := s.Defs[parts[2]]
		if !ok {
			return "", nil, fmt.Errorf("bad reference %q", parts[2])
		}

		typ, ok := defs.Properties["type"]
		if !ok {
			return "", nil, fmt.Errorf("referenced enum value %q does not have type property", parts[2])
		}

		if typ.Const == nil && typ.Not == nil {
			return "", nil, fmt.Errorf("referenced enum value %q has non-const type property and no not", parts[2])
		}

		if typ.Const != nil {
			return *typ.Const, nil, nil
		}

		if len(typ.Not.Enum) == 0 {
			return "", nil, fmt.Errorf("final clause 'not' enum must not be empty")
		}

		return "", typ.Not.Enum, nil
	}

	return "", nil, fmt.Errorf("type had bad Ref value: %q", t.Ref)
}

func (s *Schema) writeJsonUnmarshalerEnum(name string, t TypeSchema, w io.Writer) error {
	fmt.Fprintf(w, "func (t *%s) UnmarshalJSON(b []byte) (error) {\n", name)
	fmt.Fprintf(w, "\ttyp, err := util.EnumTypeExtract(b)\n")
	fmt.Fprintf(w, "\tif err != nil {\n\t\treturn err\n\t}\n\n")
	fmt.Fprintf(w, "\tswitch typ {\n")
	for i, e := range t.OneOf {
		tc, nots, err := s.getTypeConstValueForType(e)
		if err != nil {
			return err
		}

		if len(nots) > 0 {
			if i == len(t.OneOf)-1 {
				tnref := s.typeNameFromRef(e.Ref)
				fmt.Fprintf(w, `
	default:
		var out %s
		if err := json.Unmarshal(b, &out); err != nil {
			return err
		}
		t.%s = &out
		return nil
`, tnref, tnref)

			} else {
				return fmt.Errorf("enum member with a not clause must be the last in a oneOf")
			}
			break
		}

		goname := s.typeNameFromRef(e.Ref)

		fmt.Fprintf(w, "\t\tcase \"%s\":\n", tc)
		fmt.Fprintf(w, "\t\t\tt.%s = new(%s)\n", goname, goname)
		fmt.Fprintf(w, "\t\t\treturn json.Unmarshal(b, t.%s)\n", goname)
	}
	fmt.Fprintf(w, "\t}\n")
	fmt.Fprintf(w, "}\n\n")

	return nil
}
