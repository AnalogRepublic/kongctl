# kongctl

*(This tool is currently a work in progress)*

CLI tool to version & manage your Kong APIs, Plugins and Consumers. Supporting multiple contexts, exports and diffs.

## Configuration

Copy the `example.kongctl.yaml` file into the root of your project and rename it to `kongctl.yaml`. You can also place this file in `$HOME/.kongctl/` or `/etc/kongctl/` if you want to keep your kongctl settings global.

## Roadmap
 - [x] Connect to Kong services
 - [ ] Export Kong apis, plugins and settings into versionable .yaml files
 - [ ] Add way of applying .yaml config files to update Kong
 - [ ] Store a diff of local config files and what's on the remote so we can only make the requests for what has changed, if something has been removed then we need to remove it, https://github.com/go-yaml/yaml/, https://github.com/kylelemons/godebug/
 - [ ] Refactor & hook up tests into a CI tool
 - [ ] Keep track of the Kong hosts as "contexts" and allow switching
 - [ ] Document tools api & options
 - [ ] Produce pretty landing page & branding
 - [ ] Add to Homebrew


## Commands (Work in progress)

```
// Describe commands; should list out the resource items or a specific one of that resource
// e.g. "kongctl describe apis"
kongctl describe [resource]
kongctl describe [resource] [id/name]

// Export command should export a resource by it's ID and optionally allow for
// manually outputting the export filename.
kongctl export [resource] [id] -o [filename]

// Allow the user to export the entire Kong
// service to a single file.
kongctl export -o [filename]

// Apply command should apply the file(s) to the Kong service
// This should have the option to recursively look for files to apply
// "Applying" requires us to diff the Kong service with all of the information
// We've got from the files. A backup flag can be used to export the current
// live configuration before importing a new one.
kongctl apply [filename/folder]

// Remove a resource by it's ID or name. A "disable" flag can be used
// for resources which can also be disabled, like a soft delete.
kongctl remove [resource] [id/name]

// Should allow a user to add a new context to the
// kongctl config. Should prompt user with questions.
kongctl context add [name]
    
// Switch the context to another
kongctl context switch [name]

// List all available contexts
kongctl context list

// Remove a context
kongctl context remove [name]

// Tell the user what their current context is
kongctl context get
```