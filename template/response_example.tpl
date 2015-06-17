{{ define "response_example" }}

<h3>Response Example</h3>

{{ if eq .Rel "create" }}
<pre><code>HTTP/1.1 201 Created</code></pre>
{{ else if eq .Rel "empty" }}
<pre><code>HTTP/1.1 202 Accepted</code></pre>
{{ else if eq .Rel "update" }}
<pre><code>HTTP/1.1 204 No Content</code></pre>
{{ else if eq .Rel "destroy" }}
<pre><code>HTTP/1.1 204 No Content</code></pre>
{{ else if eq .Rel "notImpremented" }}
<pre><code>HTTP/1.1 501 Not Implemented</code></pre>
{{ else }}
<pre><code>HTTP/1.1 200 OK</code></pre>
<pre><code class="json">{{ .TargetSchema.ExampleJSON }}</code></pre> 
{{ end }}

{{ end }}
