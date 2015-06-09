{{ define "sidemenu" }}
<ul class="nav nav-sidebar">
  <li><a href="#"><strong>{{ .Meta.Title }}</strong></a></li>
{{ range .SchemaSlice }}
  <li><a href="#{{ .Id }}">{{ .Id }}</a></li>
  <li>
    <ul class="nav nav-sidebar-submenu">
      {{ range .Links }}
        <li><a href="#{{ .Method }}-{{ .Href }}">
            {{ if eq .Rel "notImplemented" }}
            <span class="label label-warning">Not Implemented</span>
            {{ end }}
            {{ .Title }}
          </a>
        </li>
      {{ end }}
    </ul>
  </li>
{{ end }}
</ul>
{{ end }}
