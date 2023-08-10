// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package attributes

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	pb "github.com/DataDog/datadog-agent/pkg/proto/pbgo/trace"
	"gopkg.in/yaml.v3"
)

type SpanChangeAttrs struct {
	AttributeMap map[string]string `yaml:"attribute_map"`
	ApplyToSpans []string          `yaml:"apply_to_spans"`
}

type SpanChanges struct {
	RenameAttributes      *SpanChangeAttrs `yaml:"rename_attributes"`
	ChangeAttributeValues *SpanChangeAttrs `yaml:"change_attribute_values"`
}

type SpanMigration struct {
	Changes []*SpanChanges `yaml:"changes"`
}

type Migration struct {
	Spans *SpanMigration `yaml:"spans"`
}

type SpanMigrationConfig struct {
	Migrations []*Migration `yaml:"migrations"`
}

func GetConfig(file_path string) (*SpanMigrationConfig, error) {
	cfg := SpanMigrationConfig{}
	yamlFile, err := os.ReadFile(file_path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *SpanMigrationConfig) MigrateSpans(tp *pb.TracerPayload) {
	for _, m := range cfg.Migrations {
		for _, cg := range m.Spans.Changes {
			switch {
			case cg.RenameAttributes != nil:
				applyRenaming(tp, cg.RenameAttributes)
			case cg.ChangeAttributeValues != nil:
				applyValueChanges(tp, cg.ChangeAttributeValues)
			}
		}
	}
}

func applyRenaming(tp *pb.TracerPayload, attrs *SpanChangeAttrs) {
	for _, ch := range tp.GetChunks() {
		for _, s := range ch.GetSpans() {
			if !appliesToSpan(s, attrs.ApplyToSpans) {
				continue
			}
			for from, to := range attrs.AttributeMap {
				remapSpanMeta(s, from, to)
				remapSpanMetrics(s, from, to)
			}
		}
	}
}

func remapSpanMeta(s *pb.Span, from, to string) {
	val, ok := s.Meta[from]
	if !ok {
		return
	}
	s.Meta[to] = val
	delete(s.Meta, from)
}

func remapSpanMetrics(s *pb.Span, from, to string) {
	val, ok := s.Metrics[from]
	if !ok {
		return
	}
	s.Metrics[to] = val
	delete(s.Metrics, from)
}

func applyValueChanges(tp *pb.TracerPayload, attrs *SpanChangeAttrs) {
	for _, ch := range tp.GetChunks() {
		for _, s := range ch.GetSpans() {
			if !appliesToSpan(s, attrs.ApplyToSpans) {
				continue
			}
			for tag, tplVal := range attrs.AttributeMap {
				newVal, ok := resolveTemplatedValue(s, tplVal)
				if !ok {
					return
				}
				if tag == "name" {
					s.Name = newVal
					continue
				}
				if tag == "service.name" {
					s.Service = newVal
				}
				applyValChangeSpanMeta(s, tag, newVal)
				applyValChangeSpanMetrics(s, tag, newVal)
			}
		}
	}
}

func applyValChangeSpanMeta(s *pb.Span, tag, val string) {
	s.Meta[tag] = val
}

func applyValChangeSpanMetrics(s *pb.Span, tag, val string) {
	n, err := strconv.ParseFloat(val, 64)
	if err == nil {
		s.Metrics[tag] = n
	}
}

func resolveTemplatedValue(s *pb.Span, tplStr string) (string, bool) {
	tplCtx := templateContext(s)
	tpl, err := template.New("expression.tpl").Funcs(template.FuncMap{
		"hasTag": func(tag string) bool {
			_, ok := tplCtx.Tags[tag]
			return ok
		},
		"eqTag": func(tag string, cmpVal string) bool {
			val, ok := tplCtx.Tags[tag]
			if !ok {
				return false
			}
			return val == cmpVal
		},
		"eqTagAny": func(tag string, cmpVals ...string) bool {
			val, ok := tplCtx.Tags[tag]
			if !ok {
				return false
			}
			for _, cmpVal := range cmpVals {
				if val == cmpVal {
					return true
				}
			}
			return false
		},
		"asMetric": func(val any) int64 {
			n, _ := strconv.ParseInt(val.(string), 10, 64)
			return n
		},
		"getTag": func(tag string) any {
			return tplCtx.Tags[tag]
		},
		"contains": func(val any, substr string) bool {
			strVal := val.(string)
			return strings.Contains(strVal, substr)
		},
		"split": func(val any, sep string) []string {
			strVal := val.(string)
			return strings.Split(strVal, sep)
		},
	}).Parse(tplStr)
	if err != nil {
		return "", false
	}
	var out bytes.Buffer
	if err := tpl.Execute(&out, tplCtx); err != nil {
		return "", false
	}
	return out.String(), true
}

type tplContext struct {
	Name string
	Tags map[string]any
}

func templateContext(s *pb.Span) *tplContext {
	tplCtx := &tplContext{}
	allTags := map[string]any{}
	for tag, val := range s.Meta {
		allTags[tag] = val
	}
	for tag, val := range s.Metrics {
		allTags[tag] = val
	}
	tplCtx.Tags = allTags
	tplCtx.Name = s.Name
	return tplCtx
}

func appliesToSpan(s *pb.Span, filters []string) (b bool) {
	if len(filters) == 0 {
		return true
	}
	joinedFilter := ""
	if len(filters) == 1 {
		joinedFilter = fmt.Sprintf("(%s)", filters[0])
	} else {
		joinedFilter = "and "
		for i := 0; i < len(filters); i++ {
			f := filters[i]
			joinedFilter = joinedFilter + fmt.Sprintf("(%s)", f)
			if i < len(filters)-1 {
				joinedFilter = joinedFilter + " "
			}
		}
	}

	tf := fmt.Sprintf("{{if %s}}true{{else}}false{{end}}", joinedFilter)
	if rf, ok := resolveTemplatedValue(s, tf); ok {
		b, _ := strconv.ParseBool(rf)
		return b
	}
	return false
}
