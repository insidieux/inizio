# User guide

## Installing

### Github Release

Visit the [releases page](https://github.com/insidieux/inizio/releases/latest) to download one of the pre-built binaries
for your platform.

### Docker

Use the [Docker image](https://hub.docker.com/repository/docker/insidieux/inizio)

```shell
docker pull insidieux/inizio:v1.0.0
```

### go get

Alternatively, you can use the go get method:

```shell
go get github.com/insidieux/pinchy/cmd/inizio
```

Ensure that `GOPATH/bin` is added to your `PATH`.

## Usage

Example config for running binary you can find [here](./../configs/inizio/plugins.yaml)

### Binary

```shell
inizio 
  --plugins.config /etc/inizio/plugins.yaml
  --plugins.path /usr/local/bin/inizio-plugins
  ./working-directory
```

### Docker

```shell
docker run
  -v ./:/projects
  -w /projects 
  insidieux/inizio:v1.0.0
  /projects/working-directory
```

## Command flags

```shell
--layout.cleanup                      cleanup working directory before generation
--layout.template.dockerfile string   path to custom Dockerfile template (must have "gotmpl" extension)
--layout.template.makefile   string   path to custom Makefile template (must have "gotmpl" extension)
--logger.level               string   log level (default "info")
--plugins.config             string   path to plugins config yaml file
--plugins.fail-fast                   stop after first plugin failure
--plugins.path               string   path to plugins directory (default "/usr/local/bin/inizio-plugins")
```

## Plugins config

### Example

```yaml
# An example config with full definition of plugin executable arguments, flags and environment
# Important notes:
# 1. Each env.Name will be processed by uppercase
# 2. Each env.Value will be passed to envsubst process (substitutes environment variables in shell format strings)
# 3. Each flags.Value will be passed to envsubst process (substitutes environment variables in shell format strings)
- name: plugin-first
  env:
    - name: ENV_VARIABLE
      value: "${SOME_ENV_VARIABLE}"
  flags:
    - name: some.flag
      value: "${SOME_ENV_VARIABLE}"
  args:
    - first
    - second
```
