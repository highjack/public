# trivy csv template - highjack 23
# usage if this file is the folder you run the scan from
# trivy rootfs . --format template --template '@csv.tpl' --scanners vuln 

{{- if . }}
"Target","Package","CVE","Title","Severity","Package Version","Fixed Version","Link"
{{- range . }}
{{- $target := .Target -}}
{{- if (eq (len .Vulnerabilities) 0) }}
No vulnerabilities
{{- else }}
{{- range .Vulnerabilities }}
"{{$target}}","{{ escapeXML .PkgName }}","{{ escapeXML .VulnerabilityID }}","{{ escapeXML .Title }}","{{ escapeXML .Vulnerability.Severity }}","{{ escapeXML .InstalledVersion }}","{{ escapeXML .FixedVersion }}","{{ escapeXML .PrimaryURL}}"

{{- end }}
{{- end }}
{{- end }}
{{- end }}
