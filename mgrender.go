// Copyright (C) 2026 Murilo Gomes <profmugomes.com.br>
// SPDX-License-Identifier: MIT

// Site: https://www.profmugomes.com.br

package mgrender

import (
	"html/template"
	"maps"
	"net/http"
	"path/filepath"
	"sync"
)

type renderer struct {
	mu            sync.RWMutex
	layoutDefault string
	pageFile      []string
	data          map[string]any
}

type Context struct {
	renderer *renderer
	pageFile []string
	data     map[string]any
}

func New(layout string) *renderer {
	return &renderer{
		layoutDefault: layout,
		data:          make(map[string]any),
	}
}

func (p *renderer) AddFile(value string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.pageFile = append(
		p.pageFile,
		value,
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
				[]string{p.layoutDefault},
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
		value,
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
		filepath.Base(p.renderer.layoutDefault),
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
