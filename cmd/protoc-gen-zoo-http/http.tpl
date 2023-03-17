{{$svcType := .ServiceType}}

type {{$svcType}}HTTPService interface {
{{- range .MethodSets}}
	{{.Comment}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}

func Register{{$svcType}}HTTPService(g *gin.RouterGroup, svc {{$svcType}}HTTPService) {
    r := g.Group("")
	{{- range .Methods}}
	r.{{.Method}}("{{.Path}}", _{{$svcType}}_{{.Name}}{{.Num}}_HTTP_Handler(svc))
	{{- end}}
}

{{range .Methods}}
func _{{$svcType}}_{{.Name}}{{.Num}}_HTTP_Handler(svc {{$svcType}}HTTPService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
        c := &http.Context{Context: ctx}
		shouldBind := func(req *{{.Request}}) error {
			{{- if .HasBody}}
			if err := c.ShouldBind(req{{.Body}}); err != nil {
				return err
			}
			{{- if not (eq .Body "")}}
			if err := c.ShouldBindQuery(req); err != nil {
				return err
			}
			{{- end}}
			{{- else}}
			{{- if not (eq .Method "PATCH")}}
			if err := c.ShouldBindQuery(req{{.Body}}); err != nil {
				return err
			}
			{{- end}}
			{{- end}}
			{{- if .HasParams}}
			if err := c.ShouldBindUri(req); err != nil {
				return err
			}
			{{- end}}
			return nil
		}

		var err error
		var req {{.Request}}
		var rsp *{{.Reply}}

		if err = shouldBind(&req); err != nil {
		    c.SetError(err)
			return
		}
		rsp, err = svc.{{.Name}}(c.Request.Context(), &req)
		if err != nil {
		    c.SetError(err)
			return
		}
        c.JSON(200, rsp{{.ResponseBody}})
	}
}
{{- end}}

