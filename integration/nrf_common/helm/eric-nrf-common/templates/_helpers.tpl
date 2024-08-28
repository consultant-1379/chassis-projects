{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "eric-nrf-common.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "eric-nrf-common.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart version as used by the chart label.
*/}}
{{- define "eric-nrf-common.version" -}}
{{- printf "%s" .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "eric-nrf-common.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "eric-nrf-common.product-info" }}
ericsson.com/product-name: "Ericsson NRF common Service"
ericsson.com/product-number: ""
ericsson.com/product-revision: "{{.Values.productInfo.rstate}}"
{{- end}}

{{- define "eric-nrf-common.labels" }}
app.kubernetes.io/name: {{ include "eric-nrf-common.name" . }}
app.kubernetes.io/version: {{ include "eric-nrf-common.version" . }}
helm.sh/chart: {{ template "eric-nrf-common.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end}}

{{- define "eric-nrf-common.joblabels" }}
app.kubernetes.io/name: {{ include "eric-nrf-common.name" . }}
app.kubernetes.io/version: {{ include "eric-nrf-common.version" . }}
helm.sh/chart: {{ template "eric-nrf-common.chart" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end}}