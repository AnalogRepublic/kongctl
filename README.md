# kongctl

*(This tool is currently a work in progress)*

CLI tool to version & manage your Kong APIs, Plugins and Consumers. Supporting multiple contexts, exports and diffs.

## Kong

Currently supporting Kong version `0.10.x`

## Configuration

Copy the `example.kongctl.yaml` file into the root of your project and rename it to `kongctl.yaml`. You can also place this file in `$HOME/.kongctl/` if you want to keep your kongctl settings global.

The example config file:
```
current_context: default
contexts:
  default:
    host: http://127.0.0.1:8001
```

## Defining APIs, Plugins & Other Resources

Currently we only support defining your Kong resources in a single YAML file. This file is a YAML representation of the data that Kong returns in it's APIs; so the fields you can use should match up with their API values.

Note that you can omit some fields if you want them to have their default value, this is defined as per the Kong API. However, omiting a Boolean value will make it false. You will need to provide it in order to set it to true.

Example file which defines two APIs:
```
apis:
  # The first example endpoint
  - name: Endpoint-1
    upstream_url: http://endpoint.kongctl.io
    uris:
      - /testing-1
    strip_uri: true

  # The second example endpoint
  - name: Endpoint-2
    upstream_url: http://endpoint.kongctl.io
    uris:
      - /testing-2
    strip_uri: false
    methods:
      - "GET"
      - "POST"
    preserve_host: true
    retries: 10
    https_only: true
```

See the Kong API docs for more information, [here](https://getkong.org/docs/0.10.x/admin-api/).

## Commands

### Describe

These commands give you a list of a resource or the information about a specific resource.

#### Usage

Get a list of all API resources
```
kongctl describe apis
```

Get information about an API called "Testing" (You can also provide the ID of the resource)
```
kongctl describe apis Testing
```

## Export

This command will allow you to export your Kong services to a YAML file

### Usage

Provide a filename & path e.g. `/folder/export.yaml`
```
kongctl export /folder/export.yaml
```

## Apply

This command will take a file & diff it with the Kong service, it will then make any updates, additions or removals. Note that if your Kong service has a resource that isn't found in the YAML file, it will be removed. We recommend running an export before applying any files so that you can revert any unwanted changes.

### Usage

Just supply the name of the file
```
kongctl apply my_kong_setup.yaml
```

## Context

Sometimes you'll want to manage a Kong configuration across many hosts, you can modify your `kongctl.yaml` file and add more contexts. Here are some commands to help you switch and list the contexts available to you.

### Usage
```
// COMING SOON
kongctl context add [name]
    
// COMING SOON
kongctl context remove [name]

// Switch the context to another
kongctl context switch my-other-context

// List all available contexts
kongctl context list

// Tell the user what their current context is
kongctl context
```

## Roadmap
 - [x] Connect to Kong services
 - [x] Export Kong apis, plugins and settings into versionable .yaml files
 - [x] Add way of applying .yaml config files to update Kong
 - [x] Keep track of the Kong hosts as "contexts" and allow switching
 - [ ] Refactor & hook up tests into a CI tool
 - [ ] Document tools api & options
 - [ ] Produce pretty landing page & branding
 - [ ] Add to Homebrew