package inspected

import (
	"github.com/iamgoroot/dbietool/template"
	"github.com/iamgoroot/merge"
	"path/filepath"
	"sort"
	"strings"
)

type Result struct {
	Pkg          string
	Renderers    merge.Slice[template.RendererResult]
	ImportLookup merge.Map[string, string]
	Imports      merge.Map[string, bool]
}

func (r *Result) Merge(rr *Result) *Result {
	r.Renderers.Merge(rr.Renderers)
	r.ImportLookup.Merge(rr.ImportLookup)
	r.Imports = r.Imports.Merge(rr.Imports)
	return r
}

func (r *Result) Add(res ...template.RendererResult) {
	r.Renderers = append(r.Renderers, res...)
}

func (r *Result) ImportByName(name string) *Result {
	if r.Imports == nil {
		r.Imports = merge.Map[string, bool]{
			name: true,
		}
		return r
	}
	if r.ImportLookup == nil {
		r.ImportLookup = merge.Map[string, string]{}
	}
	if _, ok := r.ImportLookup[name]; ok {
		r.Imports[name] = true
	} else {
		base := filepath.Base(strings.Trim(name, `"`))
		r.ImportLookup[base] = name
		r.Imports.Merge(map[string]bool{name: true})
	}
	return r
}

func (r Result) GetRenderers() []template.RendererResult {
	sort.Slice(r.Renderers, func(i, j int) bool {
		return r.Renderers[i].Weight() < r.Renderers[j].Weight() &&
			r.Renderers[i].ID() < r.Renderers[j].ID()
	})

	//?TODO: handle dup result id ?
	return r.Renderers
}
func (r Result) GetImports() map[string]string {
	imports := map[string]string{}
	for k, v := range r.Imports {
		if v {
			imports[k] = r.ImportLookup[k]
		}
	}
	return imports
}
