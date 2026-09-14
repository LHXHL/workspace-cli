package cloudatlas

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/chaitin/chaitin-cli/config"
	cloudatlasspec "github.com/chaitin/chaitin-cli/products/cloudatlas/spec"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func TestNewCommandHelpShowsGeneratedCommands(t *testing.T) {
	cmd := NewCommand()
	cmd.SetArgs([]string{"asset", "seed", "enterprise", "--help"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"list", "get", "update", "batch-add", "delete", "set-confidence", "set-monitor"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("help missing %q:\n%s", want, out.String())
		}
	}
}

func TestNewCommandHelpShowsOpenAPIBaseURLExample(t *testing.T) {
	cmd := NewCommand()
	cmd.SetArgs([]string{"--help"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if want := "https://host/openapi"; !strings.Contains(out.String(), want) {
		t.Fatalf("help missing OpenAPI base URL example %q:\n%s", want, out.String())
	}
}

func TestListSendsQueryAndTokenHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/asset/ip" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.Header.Get("TOKEN"); got != "secret-token" {
			t.Fatalf("TOKEN = %q", got)
		}
		if got := r.URL.Query().Get("space"); got != "1" {
			t.Fatalf("space = %q", got)
		}
		if got := r.URL.Query().Get("status"); got != "valid" {
			t.Fatalf("status = %q", got)
		}
		writeJSON(t, w, map[string]any{"code": 200, "message": "", "data": map[string]any{"current": 1, "size": 20, "total": 1, "items": []any{map[string]any{"id": 1, "ip": "1.1.1.1"}}}})
	}))
	defer server.Close()

	ApplyRuntimeConfig(nil, rawConfig(Config{URL: server.URL, Token: "secret-token"}), false)
	cmd := NewCommand()
	cmd.SetArgs([]string{"asset", "ip", "list", "--space", "1", "--status", "valid"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), "1.1.1.1") || !strings.Contains(out.String(), "total=1") {
		t.Fatalf("unexpected output:\n%s", out.String())
	}
}

func TestListUsesConfiguredSpaceID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("space"); got != "42" {
			t.Fatalf("space = %q", got)
		}
		writeJSON(t, w, map[string]any{"code": 200, "message": "", "data": map[string]any{"current": 1, "size": 20, "total": 0, "items": []any{}}})
	}))
	defer server.Close()

	ApplyRuntimeConfig(nil, rawConfig(Config{URL: server.URL, Token: "secret-token", SpaceID: "42"}), false)
	cmd := NewCommand()
	cmd.SetArgs([]string{"asset", "ip", "list", "--status", "valid"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestSpaceIsRequiredForEveryCommand(t *testing.T) {
	ApplyRuntimeConfig(nil, rawConfig(Config{URL: "https://example.com", Token: "token"}), false)
	cmd := NewCommand()
	cmd.SetArgs([]string{"asset", "seed", "enterprise", "list", "--size", "1"})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "--space is required") {
		t.Fatalf("Execute() error = %v, want --space requirement", err)
	}
}

func TestPathParameterAndRequiredQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/task/asset/schedule" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("space"); got != "1" {
			t.Fatalf("space = %q", got)
		}
		writeJSON(t, w, map[string]any{"code": 200, "message": "", "data": []any{}})
	}))
	defer server.Close()

	ApplyRuntimeConfig(nil, rawConfig(Config{URL: server.URL, Token: "token"}), false)
	cmd := NewCommand()
	cmd.SetArgs([]string{"task", "schedule", "list", "--task-type", "asset", "--space", "1"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestBodyFileAndConfirmation(t *testing.T) {
	cmd := NewCommand()
	cmd.SetArgs([]string{"asset", "seed", "enterprise", "delete", "--body", `{"ids":[1]}`})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "--yes") {
		t.Fatalf("Execute() error = %v, want --yes requirement", err)
	}
}

func TestDryRunRedactsToken(t *testing.T) {
	ApplyRuntimeConfig(nil, rawConfig(Config{URL: "https://example.com", Token: "secret-token"}), true)
	cmd := NewCommand()
	cmd.SetArgs([]string{"asset", "seed", "enterprise", "batch-add", "--space", "1", "--body", `{"name":["demo"],"confidence":"100"}`})
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)

	oldVerbose := verbose
	verbose = true
	t.Cleanup(func() { verbose = oldVerbose })

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	combined := out.String() + errOut.String()
	if strings.Contains(combined, "secret-token") {
		t.Fatalf("dry-run leaked token:\n%s", combined)
	}
	if !strings.Contains(combined, "/v1/seed/enterprise/batch-create") {
		t.Fatalf("dry-run missing path:\n%s", combined)
	}
}

func TestGeneratedHelpExplainsRequiredParameterValues(t *testing.T) {
	cmd := NewCommand()
	cmd.SetArgs([]string{"asset", "ip", "list", "--help"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	help := out.String()
	for _, want := range []string{
		"--status",
		"可选值: valid=已确认, await=待确认, ignored=已排除",
		"--page",
		"类型: integer",
		"示例: 1",
		"可由 --space-id 或 cloudAtlas.space_id 提供默认值",
	} {
		if !strings.Contains(help, want) {
			t.Fatalf("help missing %q:\n%s", want, help)
		}
	}
}

func TestParameterUsageSupportsApifoxEnumFormats(t *testing.T) {
	tests := []struct {
		name   string
		schema *Schema
		want   string
	}{
		{
			name: "enumDescriptions",
			schema: &Schema{
				Type: "string",
				Enum: []any{"valid", "await"},
				XApifox: ApifoxSchema{EnumDescriptions: map[string]string{
					"valid": "已确认",
					"await": "待确认",
				}},
			},
			want: "valid=已确认, await=待确认",
		},
		{
			name: "x-apifox-enum",
			schema: &Schema{
				Type: "string",
				Enum: []any{"ignored", "invalid"},
				XApifoxEnum: []ApifoxEnumItem{
					{Value: "ignored", Description: "已排除"},
					{Value: "invalid", Description: "历史记录"},
				},
			},
			want: "ignored=已排除, invalid=历史记录",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage := formatParameterUsage(Parameter{Name: "status", Description: "状态", Schema: tt.schema}, true)
			if !strings.Contains(usage, tt.want) {
				t.Fatalf("usage = %q, want %q", usage, tt.want)
			}
		})
	}
}

func TestRequestBodyHelpExplainsFields(t *testing.T) {
	body := &RequestBody{Content: map[string]MediaType{
		"application/json": {Schema: &Schema{
			Type:     "object",
			Required: []string{"status"},
			Properties: map[string]Schema{
				"status": {
					Type:        "string",
					Description: "状态",
					Enum:        []any{"enabled", "disabled"},
					XApifoxEnum: []ApifoxEnumItem{{Value: "enabled", Description: "启用"}, {Value: "disabled", Description: "停用"}},
				},
				"name": {
					Type:        "string",
					Description: "名称",
					Example:     "demo",
				},
			},
		}},
	}}

	help := formatRequestBodyHelp(body)
	for _, want := range []string{"status: 状态；必填；类型: string；可选值: enabled=启用, disabled=停用", "name: 名称；类型: string；示例: demo"} {
		if !strings.Contains(help, want) {
			t.Fatalf("body help missing %q:\n%s", want, help)
		}
	}
}

func TestVulnerabilityListSerializesRepeatedQueryArrays(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assertQueryValues(t, query["status"], []string{"open", "following"})
		assertQueryValues(t, query["severity"], []string{"1", "2"})
		assertQueryValues(t, query["subject_id"], []string{"id1", "id2"})
		assertQueryValues(t, query["bu_id"], []string{"1", "2"})
		assertQueryValues(t, query["tag_id"], []string{"1", "2"})

		for key, want := range map[string]string{
			"space":           "1",
			"page":            "1",
			"size":            "20",
			"query":           "keyword",
			"created_at_days": "30",
		} {
			if got := query.Get(key); got != want {
				t.Errorf("query[%q] = %q, want %q", key, got, want)
			}
		}
		writeJSON(t, w, map[string]any{"code": 200, "message": "", "data": map[string]any{"current": 1, "size": 20, "total": 0, "items": []any{}}})
	}))
	defer server.Close()

	ApplyRuntimeConfig(nil, rawConfig(Config{URL: server.URL, Token: "token"}), false)
	cmd := NewCommand()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{
		"risk", "vulnerability", "list",
		"--space", "1",
		"--status", "open", "--status", "following",
		"--severity", "1", "--severity", "2",
		"--subject-id", "id1", "--subject-id", "id2",
		"--bu-id", "1", "--bu-id", "2",
		"--tag-id", "1", "--tag-id", "2",
		"--page", "1", "--size", "20", "--query", "keyword", "--created-at-days", "30",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestQueryArrayRejectsLegacyJSONArraySyntax(t *testing.T) {
	ApplyRuntimeConfig(nil, rawConfig(Config{URL: "https://example.com", Token: "token"}), false)
	cmd := NewCommand()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{
		"risk", "vulnerability", "list",
		"--space", "1",
		"--status", `["open","following"]`,
	})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want legacy JSON array syntax error")
	}
	for _, want := range []string{"invalid value for --status", "--status open --status following"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Execute() error = %q, want it to contain %q", err, want)
		}
	}
}

func TestQueryParamPreservesRepeatedKeys(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertQueryValues(t, r.URL.Query()["status"], []string{"open", "following"})
		writeJSON(t, w, map[string]any{"code": 200, "message": "", "data": map[string]any{"current": 1, "size": 20, "total": 0, "items": []any{}}})
	}))
	defer server.Close()

	ApplyRuntimeConfig(nil, rawConfig(Config{URL: server.URL, Token: "token"}), false)
	cmd := NewCommand()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{
		"risk", "vulnerability", "list",
		"--space", "1",
		"--query-param", "status=open",
		"--query-param", "status=following",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestVulnerabilityListDryRunSerializesArrayQuery(t *testing.T) {
	oldDryRun, oldOutputFormat := dryRun, outputFormat
	t.Cleanup(func() {
		dryRun = oldDryRun
		outputFormat = oldOutputFormat
	})

	tests := []struct {
		name   string
		values []string
	}{
		{name: "single value", values: []string{"open"}},
		{name: "multiple values", values: []string{"open", "following"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApplyRuntimeConfig(nil, rawConfig(Config{}), true)
			cmd := NewCommand()
			args := []string{"risk", "vulnerability", "list", "--output", "json", "--space", "1"}
			for _, value := range tt.values {
				args = append(args, "--status", value)
			}
			cmd.SetArgs(args)
			var out bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&bytes.Buffer{})

			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			var result struct {
				URL string `json:"url"`
			}
			if err := json.Unmarshal(bytes.TrimSpace(out.Bytes()), &result); err != nil {
				t.Fatalf("parse dry-run output %q: %v", out.String(), err)
			}
			parsedURL, err := url.Parse(result.URL)
			if err != nil {
				t.Fatalf("parse dry-run URL %q: %v", result.URL, err)
			}
			assertQueryValues(t, parsedURL.Query()["status"], tt.values)
		})
	}
}

func TestBuildQuerySerializesArrayStyles(t *testing.T) {
	falseValue := false
	tests := []struct {
		name      string
		parameter Parameter
		want      []string
	}{
		{
			name:      "default form explode true",
			parameter: Parameter{Name: "status", In: "query", Schema: &Schema{Type: "array", Items: &Schema{Type: "string"}}},
			want:      []string{"open", "following"},
		},
		{
			name:      "form explode false",
			parameter: Parameter{Name: "status", In: "query", Style: "form", Explode: &falseValue, Schema: &Schema{Type: "array", Items: &Schema{Type: "string"}}},
			want:      []string{"open,following"},
		},
		{
			name:      "space delimited",
			parameter: Parameter{Name: "status", In: "query", Style: "spaceDelimited", Schema: &Schema{Type: "array", Items: &Schema{Type: "string"}}},
			want:      []string{"open following"},
		},
		{
			name:      "pipe delimited",
			parameter: Parameter{Name: "status", In: "query", Style: "pipeDelimited", Schema: &Schema{Type: "array", Items: &Schema{Type: "string"}}},
			want:      []string{"open|following"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{Use: "test"}
			cmd.Flags().StringArray("status", nil, "")
			cmd.Flags().String("space", "", "")
			cmd.Flags().StringArray("query", nil, "")
			if err := cmd.Flags().Set("status", "open"); err != nil {
				t.Fatal(err)
			}
			if err := cmd.Flags().Set("status", "following"); err != nil {
				t.Fatal(err)
			}
			if err := cmd.Flags().Set("space", "1"); err != nil {
				t.Fatal(err)
			}

			query, err := buildQuery(cmd, []Parameter{
				{Name: "space", In: "query", Schema: &Schema{Type: "integer"}},
				tt.parameter,
			})
			if err != nil {
				t.Fatalf("buildQuery() error = %v", err)
			}
			assertQueryValues(t, query["status"], tt.want)
		})
	}
}

func TestParseSpecPreservesParameterStyleAndExplode(t *testing.T) {
	api, err := ParseSpec([]byte(`
openapi: 3.0.1
paths:
  /items:
    get:
      parameters:
        - name: status
          in: query
          style: form
          explode: false
          schema:
            type: array
            items:
              type: string
`))
	if err != nil {
		t.Fatalf("ParseSpec() error = %v", err)
	}
	parameter := api.Paths["/items"].Get.Parameters[0]
	if parameter.Style != "form" {
		t.Fatalf("Style = %q, want form", parameter.Style)
	}
	if parameter.Explode == nil || *parameter.Explode {
		t.Fatalf("Explode = %v, want pointer to false", parameter.Explode)
	}
}

func TestVulnerabilityListHelpUsesRepeatableArrayFlags(t *testing.T) {
	cmd := NewCommand()
	cmd.SetArgs([]string{"risk", "vulnerability", "list", "--help"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	help := out.String()
	for _, want := range []string{
		"--status stringArray",
		"--severity stringArray",
		"--subject-id stringArray",
		"--bu-id stringArray",
		"--tag-id stringArray",
		"可重复指定",
		"示例: --status open",
	} {
		if !strings.Contains(help, want) {
			t.Fatalf("help missing %q:\n%s", want, help)
		}
	}
	for _, unwanted := range []string{`示例: ["open"]`, `示例: ["1"]`} {
		if strings.Contains(help, unwanted) {
			t.Fatalf("help contains misleading JSON array example %q:\n%s", unwanted, help)
		}
	}
}

func TestAllOpenAPIQueryArraysUseRepeatableFlags(t *testing.T) {
	api, err := ParseSpec(cloudatlasspec.OpenAPIYAML)
	if err != nil {
		t.Fatalf("ParseSpec() error = %v", err)
	}
	parser := NewParser(api)
	arrayCount := 0
	for path, pathItem := range api.Paths {
		for _, operation := range operationsForPath(pathItem) {
			if operation.op == nil {
				continue
			}
			cmd := parser.createOperationCommand("test", operation.method, path, operation.op)
			for _, parameter := range operation.op.Parameters {
				if !isQueryArrayParameter(parameter) {
					continue
				}
				arrayCount++
				flagName := normalizeFlagName(parameter.Name)
				flag := cmd.Flags().Lookup(flagName)
				if flag == nil {
					t.Errorf("%s %s query array %q has no flag", operation.method, path, parameter.Name)
					continue
				}
				if got := flag.Value.Type(); got != "stringArray" {
					t.Errorf("%s %s --%s type = %q, want stringArray", operation.method, path, flagName, got)
				}
				if strings.Contains(flag.Usage, "示例: [") {
					t.Errorf("%s %s --%s contains JSON array example: %q", operation.method, path, flagName, flag.Usage)
				}
			}
		}
	}
	if arrayCount == 0 {
		t.Fatal("embedded OpenAPI spec contains no query array parameters")
	}
}

func assertQueryValues(t *testing.T, got, want []string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("query values = %#v, want %#v", got, want)
	}
}

func rawConfig(value Config) config.Raw {
	var node yaml.Node
	if err := node.Encode(value); err != nil {
		panic(err)
	}
	return config.Raw{"cloudAtlas": node}
}

func writeJSON(t *testing.T, w http.ResponseWriter, value any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
}
