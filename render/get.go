package render

import (
	"github.com/iamgoroot/dbietool/tr"
	"sort"
)

func (r Result) GetRenderers() []tr.RendererResult {
	sort.Slice(r.Renderers, func(i, j int) bool {
		return r.Renderers[i].Weight() < r.Renderers[j].Weight() &&
			r.Renderers[i].ID() < r.Renderers[j].ID()
	})

	//?TODO: handle dup render id ?
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
