{{/*
Expand the name of the chart.
*/}}
{{- define "pf-backend.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create chart labels
*/}}
{{- define "pf-backend.labels" -}}
app.kubernetes.io/name: {{ include "pf-backend.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "pf-backend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "pf-backend.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Generate the default fully qualified app name.
*/}}
{{- define "pf-backend.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "pf-backend.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
