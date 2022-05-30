package render

import (
	"github.com/iamgoroot/dbietool/tr"
	"github.com/iamgoroot/merge"
	"path/filepath"
	"strings"
)

type Result struct {
	Pkg          string
	Renderers    merge.Slice[tr.RendererResult]
	ImportLookup merge.Map[string, string]
	Imports      merge.Map[string, bool]
}

func (r *Result) Merge(rr *Result) *Result {
	r.Renderers.Merge(rr.Renderers)
	r.ImportLookup.Merge(rr.ImportLookup)
	r.Imports = r.Imports.Merge(rr.Imports)
	return r
}

func (r *Result) Add(res ...tr.RendererResult) {
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
		//r.ImportLookup[name] = pkg
		r.Imports[name] = true
	} else {
		base := filepath.Base(strings.Trim(name, `"`))
		r.ImportLookup[base] = name
		r.Imports.Merge(map[string]bool{name: true})
	}
	return r
}
