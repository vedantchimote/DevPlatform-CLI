{{/*
Expand the name of the chart.
*/}}
{{- define "devplatform-base.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "devplatform-base.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "devplatform-base.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "devplatform-base.labels" -}}
helm.sh/chart: {{ include "devplatform-base.chart" . }}
{{ include "devplatform-base.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app: {{ .Values.app.name }}
environment: {{ .Values.app.environment }}
managed-by: devplatform-cli
{{- end }}

{{/*
Selector labels
*/}}
{{- define "devplatform-base.selectorLabels" -}}
app.kubernetes.io/name: {{ include "devplatform-base.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "devplatform-base.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "devplatform-base.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Get resource limits and requests based on environment
*/}}
{{- define "devplatform-base.resources" -}}
{{- $env := .Values.app.environment | default "dev" }}
{{- if hasKey .Values.resources $env }}
{{- toYaml (index .Values.resources $env) }}
{{- else }}
{{- toYaml .Values.resources.dev }}
{{- end }}
{{- end }}
