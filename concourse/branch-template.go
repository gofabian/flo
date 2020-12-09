package concourse

var branchPipelineTemplate = `

{{define "full-pipeline" -}}
  {{template "header" . -}}
  {{template "jobs-header" . -}}
  {{template "self-update-job" . -}}
  {{template "build-job-header-passed" . -}}
  {{template "build-job-content" . -}}
{{end}}
{{define "self-update-pipeline" -}}
  {{template "header" . -}}
  {{template "jobs-header" . -}}
  {{template "self-update-job" . -}}
{{end}}
{{define "build-pipeline" -}}
  {{template "header" . -}}
  {{template "jobs-header" . -}}
  {{template "build-job-header" . -}}
  {{template "build-job-content" . -}}
{{end}}


{{define "header" -}}
resources:
  - name: checkout
    type: git
    source:
      branch: "((GIT_BRANCH))"
      uri: "((GIT_URL))"
{{- end}}


{{define "jobs-header"}}
jobs:
{{- end}}


{{define "self-update-job"}}
  - name: self-update
    plan:
      - get: checkout
        trigger: true
      - task: generate-branch-pipeline
        input_mapping:
          workspace: checkout
        config:
          platform: linux
          image_resource:
            type: registry-image
            source: {repository: gofabian/flo, tag: "0"}
          inputs:
            - name: workspace
          outputs:
            - name: workspace
            - name: flo
          run:
            dir: workspace
            path: sh
            args:
                - -exc
                - |-
                  flo generate branch -g "((GIT_URL))" -b "((GIT_BRANCH))" \
                      -i .drone.yml -o ../flo/pipeline.yml -j all
                  cat ../flo/pipeline.yml
      - set_pipeline: self
        file: flo/pipeline.yml
        vars:
          GIT_BRANCH: "((GIT_BRANCH))"
          GIT_URL: "((GIT_URL))"
{{- end}}


{{define "build-job-header"}}
  - name: "{{.Name}}"
    plan:
      - get: checkout
        trigger: true
{{- end}}

{{define "build-job-header-passed"}}
  - name: "{{.Name}}"
    plan:
      - get: checkout
        trigger: true
        passed:
          - self-update
{{- end}}


{{define "build-job-content"}}
  {{- range .Steps}}
      - task: "{{.Name}}"
        input_mapping:
          workspace: checkout
        config:
          platform: linux
          image_resource:
            type: registry-image
            source: {repository: "{{.Repository}}" {{- if .Tag -}}, tag: "{{.Tag}}" {{- end -}} }
          inputs:
            - name: workspace
          outputs:
            - name: workspace
          run:
            dir: workspace
    {{- if .Command}}
            path: "{{.Command}}"
      {{- if .CommandArgs}}
            args:
        {{- range .CommandArgs}}
              - "{{.}}"
        {{- end}}
      {{- end}}
    {{- end}}
    {{- if .Commands}}
            path: sh
            args:
                - -exc
                - |-
      {{- range .Commands}}
                  {{.}}
      {{- end}}
    {{- end}}
  {{- end}}
{{- end}}

`
