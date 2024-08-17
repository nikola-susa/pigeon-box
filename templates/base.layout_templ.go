// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func BaseLayout(title string, description string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\" data-theme=\"dark\"><head><title>Pigeon box ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if title != "" {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("| ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(title)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.layout.templ`, Line: 10, Col: 20}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</title><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><link rel=\"icon\" href=\"/static/favicon.ico\" sizes=\"any\"><meta name=\"color-scheme\" content=\"dark\"><meta name=\"theme-color\" content=\"#f5c0c0\"><link rel=\"stylesheet\" href=\"/static/build.css\"><script defer src=\"https://cdn.jsdelivr.net/npm/@alpinejs/anchor@3.x.x/dist/cdn.min.js\"></script><script defer src=\"https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js\"></script></head><body hx-indicator=\".loading-bar\" x-data=\"{ dragover: false }\" @dragover.prevent=\"dragover = true\" @dragleave.prevent=\"dragover = false\" @drop.prevent=\"drop\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<section id=\"toast-container\" class=\"fixed top-12 right-3 w-[400px] text-sm tracking-wide\" style=\"display: none;\"><div id=\"toast-container-body\" class=\"font-mono overflow-auto max-h-32 z-[99999] relative flex flex-col gap-1.5\"><div class=\"toast z-[99999] bg-background/70 backdrop-blur border rounded shadow-lg p-1.5 relative\" id=\"toast-template\" style=\"display: none;\"><div class=\"toast-title whitespace-nowrap\"><span class=\"text-alt-foreground/70\"></span> <strong class=\"uppercase font-normal\"></strong> <span></span></div></div></div></section></body><script src=\"/static/vendor/htmx.min.js\"></script><script src=\"https://unpkg.com/htmx.org@1.9.12/dist/ext/sse.js\"></script><script src=\"/static/main.js\"></script></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}