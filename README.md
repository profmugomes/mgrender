# MGRender

[![License](https://img.shields.io/badge/license-PolyForm%20Perimeter%201.0.1-5351FB)](LICENSE.md)

Biblioteca para renderização de templates HTML em aplicações Go.

O MGRender permite configurar arquivos e dados globais no `main` e criar um `Context` independente para cada requisição HTTP.

## Características

- Renderização de templates HTML
- Dados globais compartilhados entre páginas
- Arquivos globais compartilhados entre páginas
- Contexto independente por requisição
- Compatível com `net/http`
- Compatível com roteadores como Chi
- Seguro para uso concorrente
- Sem dependências externas

## Instalação

```bash
go get github.com/profmugomes/mgrender
````

## Estrutura do projeto

Um projeto pode ser organizado desta forma:

```text
mgcadastro/
├── main.go
├── cadastro/
│   ├── lista.go
│   └── editar.go
└── layout/
    ├── layout.html
    ├── menu.html
    ├── lista.html
    └── editar.html
```

## Uso básico

### main.go

O `renderer` deve ser criado uma única vez e configurado no `main`.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/profmugomes/mgrender"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	tpl := mgrender.New()

	tpl.AddFile("menu.html")

	tpl.SetData("SiteName", "Sistema de Cadastro")
	tpl.SetData("MenuData", []string{
		"Início",
		"Clientes",
		"Configurações",
	})

	rt := chi.NewRouter()

	rt.Use(middleware.RedirectSlashes)

	rt.Get("/", tpl.Handle(Lista))

	fmt.Println("Servidor iniciado em 127.0.0.1:8000")

	err := http.ListenAndServe(":8000", rt)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
```

### Handler no próprio main

Também é possível criar o handler diretamente no `main`:

```go
func Lista(ctx *mgrender.Context,w http.ResponseWriter, r *http.Request) {
	ctx.AddFile("lista.html")

	ctx.SetData("Title", "Lista de Cadastro")
	ctx.SetData("Data", []map[string]string{
		{
			"Nome": "João",
			"Idade": "29",
		},
	})
}
```

A rota utiliza o renderer configurado anteriormente:

```go
rt.Get("/", tpl.Handle(Lista))
```

## Usando handlers em subpastas

Em projetos maiores, os handlers podem ficar separados por funcionalidade.

Exemplo:

```text
mgcadastro/
├── main.go
├── cadastro/
│   ├── lista.go
│   └── editar.go
└── layout/
    ├── layout.html
    ├── menu.html
    ├── lista.html
    └── editar.html
```

### cadastro/lista.go

```go
package cadastro

import (
	"net/http"

	"github.com/profmugomes/mgrender"
)

func Lista(ctx *mgrender.Context, w http.ResponseWriter, r *http.Request) {
	ctx.AddFile("lista.html")

	ctx.SetData("Title", "Lista de Cadastro")

	ctx.SetData("Data", []map[string]string{
		{
			"Nome":  "João",
			"Idade": "29",
		},
		{
			"Nome":  "Maria",
			"Idade": "32",
		},
	})
}
```

### cadastro/editar.go

```go
package cadastro

import (
	"net/http"

	"github.com/profmugomes/mgrender"
)

func Editar(ctx *mgrender.Context, w http.ResponseWriter, r *http.Request) {
	ctx.AddFile("editar.html")

	ctx.SetData("Title", "Editar Cadastro")
}
```

### main.go

O `main` continua sendo responsável pela configuração do renderer e pelas rotas:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/profmugomes/mgcadastro/cadastro"
	"github.com/profmugomes/mgcadastro/menu"
	"github.com/profmugomes/mgrender"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	tpl := mgrender.New()

	tpl.AddFile("menu.html")

	tpl.SetData("SiteName", "MG Cadastro")
	tpl.SetData("MenuData", menu.Menu())

	rt := chi.NewRouter()

	rt.Use(middleware.RedirectSlashes)

	rt.Get("/", tpl.Handle(cadastro.Lista))
	rt.Get("/editar", tpl.Handle(cadastro.Editar))

	fmt.Println("Servidor iniciado em 127.0.0.1:8000")

	err := http.ListenAndServe(":8000", rt)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}
```

## Dados globais e dados da página

Os dados configurados através do `renderer` são globais:

```go
tpl.SetData("SiteName", "Cadastro")
tpl.SetData("MenuData", menu.Menu())
```

Eles ficam disponíveis em todas as páginas.

O handler pode adicionar ou substituir dados para aquela requisição:

```go
ctx.SetData("Title", "Lista de Cadastro")
```

Se existir um dado global com a mesma chave, o valor definido pelo `Context` terá prioridade naquela requisição.

Por exemplo:

```go
tpl.SetData("Title", "MG Cadastro")
```

e:

```go
ctx.SetData("Title", "Lista de Cadastro")
```

Na página:

```gotemplate
<title>{{.Title}}</title>
```

o resultado será:

```text
Lista de Cadastro
```

O valor global continua sendo:

```text
MG Cadastro
```

## Templates

Um template pode acessar os dados diretamente através do ponto:

```gotemplate
<h1>{{.Title}}</h1>

<p>{{.SiteName}}</p>
```

Para percorrer uma lista:

```gotemplate
<ul>
	{{range .MenuData}}
		<li>{{.}}</li>
	{{end}}
</ul>
```

## Layout

O arquivo padrão utilizado pelo MGRender é:

```text
layout/layout.html
```

Os arquivos adicionados através de:

```go
tpl.AddFile("menu.html")
```

ficam disponíveis para o layout.

Os arquivos específicos de uma página podem ser adicionados no handler:

```go
ctx.AddFile("lista.html")
```

Assim, uma requisição para `/` pode utilizar:

```text
layout/layout.html
layout/menu.html
layout/lista.html
```

Enquanto `/editar` pode utilizar:

```text
layout/layout.html
layout/menu.html
layout/editar.html
```

## Contexto por requisição

O `renderer` é criado uma única vez:

```go
tpl := mgrender.New()
```

Porém, cada chamada:

```go
tpl.Handle(cadastro.Lista)
```

cria um novo `Context` para a requisição.

Isso permite que requisições simultâneas tenham seus próprios arquivos e dados.

```text
Renderer
│
├── arquivos globais
├── dados globais
│
├── Request 1
│   └── Context
│       ├── arquivos da página
│       └── dados da página
│
├── Request 2
│   └── Context
│       ├── arquivos da página
│       └── dados da página
│
└── Request 3
    └── Context
        ├── arquivos da página
        └── dados da página
```

Dessa forma, uma página não altera o contexto de outra requisição.

---

## 🧩 Compatibilidade

* Go 1.26.5+

---

## 👤 Autor

**Murilo Gomes**

🔗 [https://www.profmugomes.com.br](https://www.profmugomes.com.br)

📺 [https://youtube.com/@profmugomes](https://youtube.com/@profmugomes)

---

## License

Copyright (c) 2026 Murilo Gomes <profmugomes.com.br>. All Rights Reserved.

This project is licensed under the PolyForm Perimeter License 1.0.1.

### Summary

This software is available for commercial and noncommercial use, subject to the terms of the PolyForm Perimeter License 1.0.1.

You may:

* ✔ Use the software for commercial and noncommercial purposes.
* ✔ Inspect and study the source code.
* ✔ Modify the software.
* ✔ Create derivative works based on the software.
* ✔ Redistribute the software and permitted modifications.

You may not:

* ✖ Provide a product that competes with the software.

See the full license terms at LICENSE.md.

This summary is provided for convenience only and does not replace or modify the full license terms.