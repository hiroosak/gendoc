{{ define "curl_example" }}

<h3>Curl Example</h3>

{{ $hs := headers }}
{{ $encType := .EncType }}

{{ if eq .Method "GET" }}
<pre><code class="bash">curl {{ baseURL }}{{ .Href }} -X GET{{ range $hs }} \
       -H "{{ . }}"{{ end }}{{ if eq $encType "application/json" }} \
       -H "Content-Type: application/json" \
       -d '{{ .Schema.ExampleJSON }}'{{else if eq $encType "application/x-www-form-urlencoded" }} \
       -H "Content-Type: application/x-www-form-urlencoded"{{ $ps := .Schema.ExampleGetData }}{{ range $ps }} \
       -d '{{ . }}'{{ end }}{{ end }}
</code></pre>
{{ else if eq .Method "POST" }}
<pre><code class="bash">curl {{ baseURL }}{{ .Href }} -X POST{{ range $hs }} \
       -H "{{ . }}"{{ end }} \{{ if eq $encType "application/json" }}
       -H "Content-Type: application/json" \{{ end }}
       -d '{{ .Schema.ExampleJSON }}'
</code></pre>
{{ else if eq .Method "PUT" }}
<pre><code class="bash">curl {{ baseURL }}{{ .Href }} -X PUT{{ range $hs }} \
       -H "{{ . }}"{{ end }} \{{ if eq $encType "application/json" }}
       -H "Content-Type: application/json" \{{ end }}
       -d '{{ .Schema.ExampleJSON }}'
</code></pre>
{{ else if eq .Method "PATCH" }}
<pre><code class="bash">curl {{ baseURL }}{{ .Href }} -X PATCH{{ range $hs }} \
       -H "{{ . }}"{{ end }} \{{ if eq $encType "application/json" }}
       -H "Content-Type: application/json" \{{ end }}
       -d '{{ .Schema.ExampleJSON }}'
</code></pre>
{{ else if eq .Method "DELETE" }}
<pre><code class="bash">curl {{ baseURL }}{{ .Href }} -X DELETE{{ range $hs }} \
       -H "{{ . }}"{{ end }} \{{ if eq $encType "application/json" }}
       -H "Content-Type: application/json" \{{ end }}
       -d '{{ .Schema.ExampleJSON }}'
</code></pre>
{{ end }}

{{ end }}
