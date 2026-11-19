{{ $mode := .Site.Params.appMode | default "FULL" }}
{{ $blogMode := or (eq $mode "BLOG_ONLY") (eq .Site.Params.blogMode true) (eq .Site.Params.blogMode "true") }}
{{ $title := .Params.heading | default .Site.Params.homepageHeading | default (printf "About %s" (.Site.Params.author | default .Site.Title)) }}
# {{ $title }}

{{ with (.Params.subtitle | default .Site.Params.homepageSubtitle | default .Site.Params.description) }}
{{ . }}
{{ end }}
{{ with .RawContent }}
{{ . }}
{{ end }}

{{ if $blogMode }}
## Recent posts

{{ $posts := where .Site.RegularPages "Section" "post" }}
{{ $posts = where $posts "Kind" "page" }}
{{ $posts = where $posts "Draft" "ne" true }}
{{ $posts = where $posts "Params.hideFromBlogHome" "ne" true }}
{{ range first 12 (sort $posts "Date" "desc") }}
{{ $page := . }}
- [{{ $page.LinkTitle }}]({{ with $page.OutputFormats.Get "markdown" }}{{ .RelPermalink }}{{ end }}){{ with $page.Description }} — {{ . }}{{ end }}
{{ end }}
{{ else if eq $mode "SINGLE" }}
{{ $singlePath := .Site.Params.singleContent | default "content/single.md" | replaceRE "^content/" "" }}
{{ $singlePages := where .Site.RegularPages "File.Path" $singlePath }}
{{ if $singlePages }}
{{ with (index $singlePages 0) }}
## {{ .Title }}

{{ .RawContent }}
{{ end }}
{{ end }}
{{ else }}
## Sections

{{ $sections := where .Site.Pages "Kind" "section" }}
{{ range $sections.ByTitle }}
{{ $page := . }}
- [{{ $page.LinkTitle }}]({{ with $page.OutputFormats.Get "markdown" }}{{ .RelPermalink }}{{ end }}){{ with $page.Description }} — {{ . }}{{ end }}
{{ end }}
{{ end }}
