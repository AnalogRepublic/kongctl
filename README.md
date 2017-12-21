# kongctl

*(This tool is currently a work in progress)*

CLI tool to version & manage your Kong APIs, Plugins and Consumers. Supporting multiple contexts, exports and diffs.

## Configuration

Copy the `example.kongctl.yaml` file into the root of your project and rename it to `kongctl.yaml`. You can also place this file in `$HOME/.kongctl/` if you want to keep your kongctl settings global.

## Roadmap
 - [x] Connect to Kong services
 - [x] Export Kong apis, plugins and settings into versionable .yaml files
 - [x] Add way of applying .yaml config files to update Kong
 - [x] Keep track of the Kong hosts as "contexts" and allow switching
 - [ ] Refactor & hook up tests into a CI tool
 - [ ] Document tools api & options
 - [ ] Produce pretty landing page & branding
 - [ ] Add to Homebrew


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

This command will take a file & diff it with the Kong service, it will then make any updates, additions or removals.

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