package concourse

var repositoryPipelineTemplate = `

{{define "full-pipeline" -}}
  {{template "header" . -}}
  {{template "jobs-header" . -}}
  {{template "refresh-job" . -}}
  {{template "build-job-header-passed" . -}}
  {{template "build-job-content" . -}}
{{end}}
{{define "refresh-pipeline" -}}
  {{template "header" . -}}
  {{template "jobs-header" . -}}
  {{template "refresh-job" . -}}
{{end}}
{{define "build-pipeline" -}}
  {{template "header" . -}}
  {{template "jobs-header" . -}}
  {{template "build-job-header" . -}}
  {{template "build-job-content" . -}}
{{end}}


{{define "header" -}}
resource_types:
  - name: git-branches
    type: registry-image
    source:
      repository: vito/git-branches-resource
resources:
  - name: branches
    type: git-branches
    source:
      uri: ((GIT_URL))
{{- end}}


{{define "jobs-header"}}
jobs:
{{- end}}


{{define "refresh-job"}}
  - name: refresh
    plan:
      - get: branches
        trigger: true
      - task: generate
        input_mapping:
          workspace: branches
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
                bs=$(tr '\n' ',' < branches | sed -e 's/,*$//' | sed -e 's/,/ -b /g')
                flo generate repository -g "((GIT_URL))" -b $bs \
                  -i .drone.yml -o ../flo/pipeline.yml -j all
                cat ../flo/pipeline.yml
      - set_pipeline: self
        file: flo/pipeline.yml
        vars:
          GIT_URL: ((GIT_URL))
{{- end}}


{{define "build-job-header"}}
  - name: pipelines
    plan:
      - get: branches
        trigger: true
{{- end}}

{{define "build-job-header-passed"}}
  - name: pipelines
    plan:
      - get: branches
        trigger: true
        passed:
          - refresh
{{- end}}


{{define "build-job-content"}}
      - task: generate
        input_mapping:
          workspace: branches
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
                flo generate branch -g "((GIT_URL))" -b dummy \
                  -i .drone.yml -o ../flo/pipeline.yml -j refresh
                cat ../flo/pipeline.yml
  {{- range .Branches}}
      - set_pipeline: "{{.HarmonizedName}}"
        file: flo/pipeline.yml
        vars:
          GIT_URL: ((GIT_URL))
          GIT_BRANCH: "{{.Name}}"
  {{- end}}
{{- end}}

`
