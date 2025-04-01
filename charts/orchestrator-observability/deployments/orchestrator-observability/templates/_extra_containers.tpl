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
  resources:
    requests:
      cpu: 10m
      memory: 32Mi
    limits:
      cpu: 500m
      memory: 256Mi
  {{- if .Values.grafana_proxy.env }}
  args:
    - --log-level={{ .Values.grafana_proxy.env.logLevel }}
    {{- if eq .Values.grafana_proxy.env.includeDeleted true }}
    - --include-deleted
    {{- end }}
    {{- if eq .Values.grafana_proxy.env.includeEdgenodeSystem true}}
    - --include-edgenode-system
    {{- end }}
  {{- end }}
{{- end }}
{{- end }}
