// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"assistant/api/assets/templ/components"
	"assistant/types"
)

func Home(projectName string, contracts []types.Contract, conversation []types.Message, errorMessages []int, isSidebarOpen bool, selectedModel string) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<html lang=\"en\"><head><meta charset=\"UTF-8\"><link rel=\"stylesheet\" href=\"/static/css/style.css\"><link rel=\"stylesheet\" href=\"/static/css/vendor/notiflix.css\"><link rel=\"stylesheet\" href=\"https://cdn.jsdelivr.net/npm/@shoelace-style/shoelace@2.16.0/cdn/themes/dark.css\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><!--\n        * 204 No Content by default does nothing, but is not an error\n        * 2xx, 3xx and 422 responses are non-errors and are swapped\n        * 4xx & 5xx responses are swapped but are errors\n        * 503 responses are not swapped and is an error\n        * all other responses are swapped using \"...\" as a catch-all\n        --><meta name=\"htmx-config\" content=\"{\n                &#34;timeout&#34;:300000,\n                &#34;responseHandling&#34;:[\n                    {&#34;code&#34;:&#34;204&#34;, &#34;swap&#34;: false},\n                    {&#34;code&#34;:&#34;[23]..&#34;, &#34;swap&#34;: true},\n                    {&#34;code&#34;:&#34;422&#34;, &#34;swap&#34;: true},\n                    {&#34;code&#34;:&#34;503&#34;, &#34;swap&#34;: false, &#34;error&#34;:true},\n                    {&#34;code&#34;:&#34;[45]..&#34;, &#34;swap&#34;: true, &#34;error&#34;:true},\n                    {&#34;code&#34;:&#34;...&#34;, &#34;swap&#34;: true}\n                ]\n            }\"><script src=\"https://cdn.tailwindcss.com\"></script><title>Audit Assistant</title></head><body><div class=\"flex flex-col h-screen bg-gray-900 text-white\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.HeaderTemplate(projectName, isSidebarOpen, selectedModel).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.MainContent(contracts, conversation, errorMessages, isSidebarOpen).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.GenerateReportDialog().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><script type=\"module\" src=\"https://cdn.jsdelivr.net/npm/@shoelace-style/shoelace@2.16.0/cdn/shoelace-autoloader.js\"></script><script src=\"https://unpkg.com/htmx.org@2.0.2\" integrity=\"sha384-Y7hw+L/jvKeWIRRkqWYfPcvVxHzVzn5REgzbawhxAuQGwX1XWe70vji+VSeHOThJ\" crossorigin=\"anonymous\"></script><script src=\"https://unpkg.com/showdown/dist/showdown.min.js\"></script><script src=\"/static/js/vendor/notiflix.js\"></script><script src=\"/static/js/script.js\"></script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
