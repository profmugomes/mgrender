// Copyright (c) 2026 Murilo Gomes <profmugomes.com.br>. All Rights Reserved.

// Licensed under the PolyForm Perimeter License 1.0.1.
// See LICENSE.md for details.

package mgrender

import (
	"html/template"
	"maps"
	"net/http"
	"path/filepath"
	"sync"
)

type renderer struct {
	mu       sync.RWMutex
	layout   string
	pageFile []string
	data     map[string]any
}

type Context struct {
	renderer *renderer
	pageFile []string
	data     map[string]any
}

func New() *renderer {
	return &renderer{
		layout: "layout/layout.html",
		data:   make(map[string]any),
	}
}

func (p *renderer) AddFile(value string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.pageFile = append(
		p.pageFile,
		filepath.Join("layout", value),
	)
}

func (p *renderer) SetData(key string, value any) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.data[key] = value
}

func (p *renderer) Handle(
	handler func(*Context, http.ResponseWriter, *http.Request),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.mu.RLock()

		ctx := &Context{
			renderer: p,
			pageFile: append(
				[]string{p.layout},
				p.pageFile...,
			),
			data: maps.Clone(p.data),
		}

		p.mu.RUnlock()

		handler(ctx, w, r)

		ctx.Show(w, http.StatusOK)
	}
}

func (p *Context) AddFile(value string) {
	p.pageFile = append(
		p.pageFile,
		filepath.Join("layout", value),
	)
}

func (p *Context) SetData(key string, value any) {
	p.data[key] = value
}

func (p *Context) Show(w http.ResponseWriter, status int) {
	tpl, err := template.ParseFiles(p.pageFile...)
	if err != nil {
		http.Error(
			w,
			"Erro ao carregar templates: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	w.WriteHeader(status)

	err = tpl.ExecuteTemplate(
		w,
		"layout.html",
		p.data,
	)

	if err != nil {
		http.Error(
			w,
			"Erro ao renderizar template: "+err.Error(),
			http.StatusInternalServerError,
		)
	}
}
