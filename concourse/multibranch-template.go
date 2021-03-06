package concourse

var repositoryPipelineTemplate = `

{{define "full-pipeline" -}}
  {{template "header" . -}}
  {{template "branch-resources" . -}}
  {{template "jobs-header" . -}}
  {{template "self-update-job" . -}}
  {{template "build-jobs-passed" .}}
{{end}}
{{define "self-update-pipeline" -}}
  {{template "header" . -}}
  {{template "jobs-header" . -}}
  {{template "self-update-job" .}}
{{end}}
{{define "build-pipeline" -}}
  {{template "header" . -}}
  {{template "branch-resources" . -}}
  {{template "jobs-header" . -}}
  {{template "build-jobs" .}}
{{end}}


{{define "header" -}}
resource_types:
  - name: branches-resource-type
    type: registry-image
    source:
      repository: gofabian/git-branches-resource
resources:
  - name: branches
    type: branches-resource-type
    source:
      uri: ((GIT_URL))
{{- end}}


{{define "branch-resources" -}}
  {{range .Branches}}
  - name: "checkout-{{.HarmonizedName}}"
    type: git
    source:
      uri: "((GIT_URL))"
      branch: "{{.Name}}"
  {{- end}}
{{- end}}


{{define "jobs-header"}}
jobs:
{{- end}}


{{define "self-update-job"}}
  - name: self-update
    plan:
      - get: branches
        trigger: true
      - task: generate-multibranch-pipeline
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
                b=$(sort < branches | tr '\n' ',' | sed -e 's/,*$//')
                flo generate-pipeline -s multibranch -j self-update,build -b "$b" -i "{{.DroneFile}}" -o ../flo/pipeline.yml
                cat ../flo/pipeline.yml
      - set_pipeline: self
        file: flo/pipeline.yml
        vars:
          GIT_URL: ((GIT_URL))
{{- end}}


{{define "build-jobs-passed" -}}
  {{template "build-job-header-passed" . -}}
  {{range .Branches -}}
    {{template "build-job-content" . -}}
  {{- end}}
{{- end}}


{{define "build-jobs" -}}
  {{template "build-job-header" . -}}
  {{range .Branches -}}
    {{template "build-job-content" . -}}
  {{- end}}
{{- end}}


{{define "build-job-header"}}
  - name: update-branch-pipelines
    plan:
      - get: branches
        trigger: true
{{- end}}

{{define "build-job-header-passed"}}
  - name: update-branch-pipelines
    plan:
      - get: branches
        trigger: true
        passed:
          - self-update
{{- end}}


{{define "build-job-content"}}
      - get: "checkout-{{.HarmonizedName}}"
      - task: "generate-branch-pipeline-{{.HarmonizedName}}"
        input_mapping:
          workspace: "checkout-{{.HarmonizedName}}"
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
                flo generate-pipeline -s branch -j self-update,build -i "{{.DroneFile}}" -o ../flo/pipeline.yml
                cat ../flo/pipeline.yml
      - set_pipeline: "branch-{{.HarmonizedName}}"
        file: flo/pipeline.yml
        vars:
          GIT_URL: ((GIT_URL))
          GIT_BRANCH: "{{.Name}}"
{{- end}}

`
