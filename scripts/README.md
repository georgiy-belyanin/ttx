# TTX scripts

This directory contains supplementary scripts used during the TTX development.

## Automatically generate Go models from Tarantool jsonschema

This section describes the details on how you might generate Go models
automatically using Tarantool config jsonschemas.

You might need the following tools.

* [Tarantool](https://github.com/tarantool/tarantool)
* [go-jsonschema](https://github.com/omissis/go-jsonschema)
* [json-schema-defaults](https://www.npmjs.com/package/json-schema-defaults)
* [yq](https://github.com/mikefarah/yq)

TL;DR just install the mentioned dependencies and execute the script:

```bash
# Install the dependencies other than Tarantool
go get github.com/atombender/go-jsonschema/...
go install github.com/atombender/go-jsonschema@latest
npm install -g json-schema-defaults
pip3 install yq

# Generate `instance_config.go` and `instance.config.default.yaml`
./gen-instance-config
```

For more information on each step see the sections below.

### Get Tarantool instance config jsonschema

Execute the following Lua script in Tarantool to get JSONschema for instance
configuration.

```lua
do
    local fio = require('fio')
    local json = require('json')
    local instance_config = require('internal.config.instance_config')

    f = fio.open('instance.config', {'O_CREAT', 'O_WRONLY'})
    f:write(json.encode(instance_config:jsonschema()))
end
```

### Generate Go models from the jsonschema

You might use [go-jsonschema](https://github.com/omissis/go-jsonschema) to
generate Go models for unmarshalling YAML files.

```bash
go-jsonschema --only-models --tags yaml -e -p config instance.config -o instance_config.go
```

### Generate defaults for JSONschema

You might use [json-schema-defaults]
(https://www.npmjs.com/package/json-schema-defaults) to generate a default
configuration for JSONschema. [yq] might be useful for converting the received
json into a YAML.


```bash
json-schema-defaults instance.config | yq -P
```
