# {{ .Title }}

{{ with .Description }}
{{ . }}
{{ end }}
{{ with .RawContent }}
{{ . }}
{{ end }}

{{ if eq .Kind "taxonomy" }}
## Terms

{{ range .Data.Terms.Alphabetical }}
{{ $term := .Page }}
- [{{ $term.LinkTitle }}]({{ with $term.OutputFormats.Get "markdown" }}{{ .RelPermalink }}{{ end }}) ({{ .Count }})
{{ end }}
{{ else }}
## Pages

{{ range .Pages.ByTitle }}
{{ $page := . }}
- [{{ $page.LinkTitle }}]({{ with $page.OutputFormats.Get "markdown" }}{{ .RelPermalink }}{{ end }}){{ with $page.Description }} — {{ . }}{{ end }}
{{ end }}
{{ end }}
