{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "eric-nef-nf-registration.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create chart version as used by the chart label.
*/}}
{{- define "eric-nef-nf-registration.version" -}}
{{- printf "%s" .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "eric-nef-nf-registration.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create image registry url
*/}}
{{- define "eric-nef-nf-registration.registryUrl" -}}
{{- if .Values.imageCredentials.registry -}}
{{- if .Values.imageCredentials.registry.url -}}
{{- print .Values.imageCredentials.registry.url -}}
{{- else if .Values.global.registry.url -}}
{{- print .Values.global.registry.url -}}
{{- else -}}
""
{{- end -}}
{{- else if .Values.global.registry.url -}}
{{- print .Values.global.registry.url -}}
{{- else -}}
""
{{- end -}}
{{- end -}}

{{/*
Create image pull secrets
*/}}
{{- define "eric-nef-nf-registration.pullSecrets" -}}
{{- if .Values.imageCredentials.registry -}}
{{- if .Values.imageCredentials.registry.pullSecret -}}
{{- print .Values.imageCredentials.registry.pullSecret -}}
{{- else if .Values.global.registry.pullSecret -}}
{{- print .Values.global.registry.pullSecret -}}
{{- end -}}
{{- else if .Values.global.registry.pullSecret -}}
{{- print .Values.global.registry.pullSecret -}}
{{- end -}}
{{- end -}}

{{/*
Full image name including image registry and tag
*/}}

{{/*
Ericsson product information
*/}}
{{- define "eric-nef-nf-registration.product-info" -}}
ericsson.com/product-name: ""
ericsson.com/product-number: ""
ericsson.com/product-revision: "{{.Values.productInfo.rstate}}"
{{- end -}}
