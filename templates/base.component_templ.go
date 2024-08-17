// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/nikola-susa/pigeon-box/model"
	"github.com/slack-go/slack"
)

func logo() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<svg viewBox=\"0 0 128 128\" fill=\"none\"><g fill=\"currentColor\"><path d=\"m7.042 114.751 41.536-41.536.047-19.096L7.279 12.773a2.204 2.204 0 0 0-1.131-.598 2.222 2.222 0 0 0-2.272.944c-.245.363-.377.79-.38 1.226l-.04 18.477a2.183 2.183 0 0 0 .745 1.702l.8.683 27.25 22.52-26.384-2.855a2.234 2.234 0 0 0-2.302 1.317c-.128.29-.194.602-.192.918l.05 13.387a2.213 2.213 0 0 0 .73 1.68 2.203 2.203 0 0 0 1.747.55l26.398-3.008L4.897 92.352l-.81.693a2.221 2.221 0 0 0-.777 1.682l-.045 18.47a2.21 2.21 0 0 0 1.366 2.035c.403.166.847.209 1.276.123a2.234 2.234 0 0 0 1.135-.604z\"></path> <path d=\"M120.958 15.61 79.422 57.145l-.047 19.096 41.346 41.346c.31.307.703.515 1.131.598a2.222 2.222 0 0 0 2.272-.944c.245-.363.377-.79.38-1.226l.04-18.477a2.183 2.183 0 0 0-.745-1.703l-.8-.682-27.25-22.52 26.384 2.855a2.234 2.234 0 0 0 2.302-1.317c.128-.29.194-.603.192-.919l-.05-13.387a2.213 2.213 0 0 0-.73-1.68 2.203 2.203 0 0 0-1.747-.549l-26.398 3.008 27.401-22.637.81-.692a2.221 2.221 0 0 0 .777-1.682l.045-18.47a2.21 2.21 0 0 0-1.366-2.035 2.216 2.216 0 0 0-1.276-.124 2.234 2.234 0 0 0-1.134.604z\"></path> <path d=\"M113.237 121.384 71.7 79.848 52.604 79.8l-41.346 41.346c-.307.309-.515.703-.598 1.131a2.222 2.222 0 0 0 .944 2.272c.363.245.79.377 1.226.38l18.477.04a2.183 2.183 0 0 0 1.703-.745l.682-.8 22.52-27.25-2.855 26.384a2.234 2.234 0 0 0 1.317 2.302c.29.128.603.194.919.192l13.387-.05a2.213 2.213 0 0 0 1.68-.73 2.203 2.203 0 0 0 .549-1.747L68.2 96.128l22.637 27.401.692.81a2.221 2.221 0 0 0 1.682.777l18.47.045a2.21 2.21 0 0 0 2.035-1.366c.166-.403.21-.847.124-1.276a2.234 2.234 0 0 0-.604-1.135z\"></path></g> <path class=\"text-primary\" fill=\"currentColor\" d=\"m13.561 6.616 41.537 41.536 19.096.047L115.54 6.853c.307-.31.515-.703.598-1.131a2.222 2.222 0 0 0-.944-2.272 2.227 2.227 0 0 0-1.226-.38L95.49 3.03a2.183 2.183 0 0 0-1.703.745l-.682.8-22.52 27.25L73.44 5.44a2.234 2.234 0 0 0-1.317-2.302 2.234 2.234 0 0 0-.919-.192l-13.387.05a2.213 2.213 0 0 0-1.68.73 2.203 2.203 0 0 0-.549 1.747l3.008 26.398L35.96 4.471l-.692-.81a2.221 2.221 0 0 0-1.682-.777L15.115 2.84a2.21 2.21 0 0 0-2.034 1.366 2.216 2.216 0 0 0-.124 1.276c.086.429.296.823.604 1.135z\"></path></svg>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func textareaRichEdit(label string, name string, message model.RenderMessage) templ.Component {
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
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div x-data=\"Editor\" x-init=\"init(true)\" class=\"relative\"><div class=\"flex flex-col w-full border rounded bg-alt focus-within:border-primary shadow\"><div class=\"flex-1 relative\"><label for=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(name)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.component.templ`, Line: 23, Col: 32}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><span class=\"label hidden\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(label)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.component.templ`, Line: 24, Col: 53}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span> <textarea x-ref=\"textarea\" tabindex=\"1\" id=\"message-input\" name=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(name)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.component.templ`, Line: 29, Col: 30}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" x-model=\"value\" @input=\"handleInputEvent(); resize()\" rows=\"1\" title=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(label)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.component.templ`, Line: 33, Col: 32}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" placeholder=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var7 string
		templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(label)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.component.templ`, Line: 34, Col: 38}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"block w-full min-h-[30px] border-0 p-4 pb-11 bg-inherit resize-none placeholder:text-sidebar-foreground/50\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var8 string
		templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(message.Text)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.component.templ`, Line: 34, Col: 169}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</textarea></label></div></div><div class=\"absolute bottom-2 right-2 flex gap-1.5\"><button tabindex=\"1\" type=\"button\" class=\"button button-basic\" title=\"Cancel editing\" x-on:htmx:after-request=\"forceFocus($el)\" hx-post=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var9 string
		templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs("/t/" + message.ThreadID + "/m/" + message.ID + "/cancel")
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.component.templ`, Line: 44, Col: 83}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-swap=\"outerHTML\" hx-target=\"closest li\" hx-trigger=\"click, submit, keyup-escape from:body\">Cancel</button> <button tabindex=\"1\" type=\"submit\" class=\"button button-success shadow\" x-bind:disabled=\"value.trim() === &#39;&#39;\" title=\"Save message\">Save</button></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func userAvatar(user slack.User) templ.Component {
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
		templ_7745c5c3_Var10 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var10 == nil {
			templ_7745c5c3_Var10 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"w-5 h-5 rounded shadow bg-alt relative\" x-data=\"{ open: false }\" @mouseenter=\"open = true\" @mouseleave=\"open = false\" @click=\"open = ! open\"><img x-ref=\"avatar\" src=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var11 string
		templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(user.Profile.Image24)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.component.templ`, Line: 62, Col: 53}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"absolute inset-0 rounded opacity-90\" alt=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var12 string
		templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(user.Profile.RealName)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.component.templ`, Line: 62, Col: 126}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div x-show=\"open\" x-anchor=\"$refs.avatar\" @click.outside=\"open = false\" class=\"z-10 w-40 bg-background text-foreground border rounded shadow\"><div class=\"p-2 pb-0\"><p class=\"text-base\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var13 string
		templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(user.Profile.RealName)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.component.templ`, Line: 65, Col: 59}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p><p class=\"text-sm text-alt-foreground\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var14 string
		templ_7745c5c3_Var14, templ_7745c5c3_Err = templ.JoinStringErrs(user.Profile.Title)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.component.templ`, Line: 66, Col: 74}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var14))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p></div><div class=\"divider my-2\"></div><div class=\"relative rounded shadow border h-[142px] w-[142px] bg-primary/10 m-2\"><img src=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var15 string
		templ_7745c5c3_Var15, templ_7745c5c3_Err = templ.JoinStringErrs(user.Profile.ImageOriginal)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.component.templ`, Line: 70, Col: 52}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var15))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"absolute inset-0 rounded flex justify-center items-center text-sm text-primary\" alt=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var16 string
		templ_7745c5c3_Var16, templ_7745c5c3_Err = templ.JoinStringErrs(user.Profile.RealName)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/base.component.templ`, Line: 70, Col: 167}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var16))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}