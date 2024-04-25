// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.663
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Register(banner templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<section class=\"registration-section\"><section class=\"reptr-heading\"><span>Reptr</span></section><section class=\"reptr-description\"><p>Create an account to start building your study decks.</p></section>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if banner != nil {
			templ_7745c5c3_Err = banner.Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<section class=\"form-container\"><form action=\"/register\" method=\"POST\"><section class=\"input-container\"><input type=\"text\" id=\"email\" name=\"email\" placeholder=\"Email\"><br></section><section class=\"input-container\"><input type=\"password\" id=\"password\" name=\"password\" placeholder=\"Password\"><br></section><section class=\"input-container-last\"><input type=\"password\" id=\"repassword\" name=\"repassword\" placeholder=\"Confirm Password\"><br></section><input type=\"submit\" value=\"Register\"><section class=\"or\"><span>Or</span></section><section class=\"login-link\"><a class=\"button\" href=\"/login\">Login</a></section></form></section></section>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
