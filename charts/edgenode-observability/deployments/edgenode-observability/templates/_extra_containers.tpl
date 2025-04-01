{{- define "orchestrator.extraContainers" -}}
{{- if .Values.grafana_proxy -}}
- name: {{ .Values.grafana_proxy.name }}
  image: {{ .Values.grafana_proxy.registry }}/{{ .Values.grafana_proxy.repository }}:{{ .Values.grafana_proxy.tag }}
  imagePullPolicy: {{ .Values.grafana_proxy.pullPolicy }}
  ports:
    - containerPort: 9190
  securityContext:
    capabilities:
      drop:
      - ALL
    allowPrivilegeEscalation: false
    readOnlyRootFilesystem: true
    seccompProfile:
      type: RuntimeDefault
  {{- with .Values.grafana_proxy.resources }}
  resources:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  {{- if .Values.grafana_proxy.env }}
  args:
    - --log-level={{ .Values.grafana_proxy.env.logLevel }}
    {{- if eq .Values.grafana_proxy.env.includeDeleted true }}
    - --include-deleted
    {{- end }}
    {{- if eq .Values.grafana_proxy.env.includeEdgenodeSystem true }}
    - --include-edgenode-system
    {{- end }}
    {{- with .Values.grafana_proxy.env.namespace }}
    - {{ print "--mimir-url=http://edgenode-observability-mimir-gateway." (.edgenode | default "orch-infra") ".svc.cluster.local:8181/prometheus" | quote }}
    - {{ print "--loki-url=http://edgenode-observability-loki-gateway." (.edgenode | default "orch-infra") ".svc.cluster.local:80" | quote }}
    - {{ print "--otc-url=observability-tenant-controller." (.platform | default "orch-platform") ".svc.cluster.local:50051" | quote }}
    {{- end }}
  {{- end }}
{{- end }}
{{- end }}
