{{ define "schema" }}
  {{ range .SchemaSlice }}
    <h1 class="page-header">
      {{ .Id }} <a name="{{ .Id }}" class="anchorjs-link" href="#{{ .Id }}"> <small><span class="glyphicon glyphicon-link xx-small" aria-hidden="true"></span></small></a> 
    </h1>
    <p>{{ .Description }}</p>
    <h2>Attributes</h2>
    <div class="table-responsive">
      <table class="table table-striped">
        <thead>
          <tr>
            <th>#</th>
            <th>Name</th>
            <th>Type</th>
            <th>Format</th>
            <th>Description</th>
            <th>Example</th>
          </tr>
        </thead>
        <tbody>
          {{ range $n, $d := .Properties }}
          <tr>
            <td></td>
            <td>{{ $n }}</td>
            <td>{{ $d.ResolveType }}</td>
            <td>{{ $d.ResolveFormat }}</td>
            <td>{{ $d.ResolveDescription }}</td>
            <td>
              <code>{{ $d.ExampleJSON }}</code>
            </td>
          </tr>
          {{ end }}
        </tbody>
      </table>
    </div>
    {{ range .Links }}
      <h2>
        {{ .Method }} {{ .Href }} - {{ .Title }} <a name="{{ .Method }}-{{ .Href }}" class="anchorjs-link" href="#{{ .Method }}-{{ .Href }}"><small><span class="glyphicon glyphicon-link" aria-hidden="true"></span></small></a> 
      </h2>

      <p>{{ .Description }}</p>

      {{ template "curl_example" . }}
      {{ template "response_example" . }}

    {{ end }}
  {{ end }}
{{ end }}
