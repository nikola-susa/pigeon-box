package md

import (
	"fmt"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"io"
)

func Parse(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)

	opts := mdhtml.RendererOptions{
		Flags:          mdhtml.CommonFlags | mdhtml.HrefTargetBlank,
		RenderNodeHook: renderHook,
	}
	renderer := mdhtml.NewRenderer(opts)

	maybeUnsafeHTML := markdown.ToHTML(md, p, renderer)

	policy := createPolicy()
	safeHtml := policy.SanitizeBytes(maybeUnsafeHTML)

	return safeHtml
}

func createPolicy() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("style").OnElements("code", "pre", "div", "span", "p", "a", "ul", "ol", "li", "table", "tr", "td", "th", "h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowAttrs("class").OnElements("code", "pre", "div", "span", "p", "a", "ul", "ol", "li", "table", "tr", "td", "th", "h1", "h2", "h3", "h4", "h5", "h6")
	return p
}

func htmlHighlighter(w io.Writer, source, lang, defaultLang string) error {
	if lang == "" {
		lang = defaultLang
	}
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}

	htmlFormatter := html.New(html.TabWidth(2), html.WithClasses(true), html.ClassPrefix("hl-"))
	if htmlFormatter == nil {
		fmt.Printf("couldn't create html formatter")
		panic("couldn't create html formatter")
	}
	styleName := "base16-snazzy"
	highlightStyle := styles.Get(styleName)

	if highlightStyle == nil {
		fmt.Printf("didn't find style '%s'", styleName)
	}

	return htmlFormatter.Format(w, highlightStyle, it)
}

func renderCode(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	defaultLang := "bash"
	lang := string(codeBlock.Info)

	err := htmlHighlighter(w, string(codeBlock.Literal), lang, defaultLang)
	if err != nil {
		fmt.Printf("error highlighting code: %s", err)
		return
	}
}

func renderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok {
		renderCode(w, code, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}
