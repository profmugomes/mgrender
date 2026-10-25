# MGRender

[![License](https://img.shields.io/badge/license-MIT-5351FB)](LICENSE)

Biblioteca para renderização de templates HTML em aplicações Go.

O MGRender permite configurar arquivos e dados globais no `main` e criar um `Context` independente para cada requisição HTTP.

## Características

* Renderização de templates HTML
* Layout padrão definido na criação do renderer
* Dados globais compartilhados entre páginas
* Arquivos globais compartilhados entre páginas
* Contexto independente por requisição
* Compatível com `net/http`
* Compatível com roteadores como Chi
* Seguro para uso concorrente
* Sem dependências externas

## Instalação

```bash
go get github.com/profmugomes/mgrender
```

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

O layout principal é informado diretamente através do `New()`:

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/profmugomes/mgrender"
)

func main() {
	tpl := mgrender.New("layout/layout.html")

	tpl.AddFile("layout/menu.html")

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
func Lista(
	ctx *mgrender.Context,
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx.AddFile("layout/lista.html")

	ctx.SetData("Title", "Lista de Cadastro")

	ctx.SetData("Data", []map[string]string{
		{
			"Nome":  "João",
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

func Lista(
	ctx *mgrender.Context,
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx.AddFile("layout/lista.html")

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

func Editar(
	ctx *mgrender.Context,
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx.AddFile("layout/editar.html")

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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/profmugomes/mgcadastro/cadastro"
	"github.com/profmugomes/mgcadastro/menu"
	"github.com/profmugomes/mgrender"
)

func main() {
	tpl := mgrender.New("layout/layout.html")

	tpl.AddFile("layout/menu.html")

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

O layout principal é definido na criação do renderer:

```go
tpl := mgrender.New("layout/layout.html")
```

O caminho é definido pelo desenvolvedor e não é imposto pela biblioteca.

Os arquivos adicionados através de:

```go
tpl.AddFile("layout/menu.html")
```

ficam disponíveis para o layout.

Os arquivos específicos de uma página podem ser adicionados no handler:

```go
ctx.AddFile("layout/lista.html")
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

A biblioteca não exige que os templates estejam em um diretório específico. Os caminhos são definidos pela aplicação.

## Contexto por requisição

O `renderer` é criado uma única vez:

```go
tpl := mgrender.New("layout/layout.html")
```

Porém, cada chamada:

```go
tpl.Handle(cadastro.Lista)
```

cria um novo `Context` para a requisição.

Os dados globais são clonados para o `Context`, permitindo que uma página altere ou adicione dados sem modificar os dados globais.

Os arquivos globais também são copiados para o contexto da requisição, permitindo adicionar os arquivos específicos da página sem alterar o estado do `renderer`.

Isso permite que requisições simultâneas tenham seus próprios arquivos e dados.

```text
Renderer

├── layout padrão
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

## Compatibilidade

* Go 1.27.2+
* `net/http`
* Roteadores compatíveis com `net/http`
* Chi e outros roteadores que utilizem `http.Handler` ou `http.HandlerFunc`

O MGRender não possui dependências externas.

## 👤 Autor

**Murilo Gomes**

🔗 [https://www.profmugomes.com.br](https://www.profmugomes.com.br)

📺 [https://youtube.com/@profmugomes](https://youtube.com/@profmugomes)

---

## License

Copyright (c) 2026 Murilo Gomes <profmugomes.com.br>

Licensed under the [MIT](https://github.com/profmugomes/mgrender/blob/main/LICENSE) license.

All contributions to the MGRender are subject to this license.
