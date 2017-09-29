# kongctl

CLI tool to version & manage your Kong APIs, Plugins and Consumers. Supporting multiple contexts, exports and diffs.

## Roadmap
 - [] Connect to Kong services
 - [] Export Kong apis, plugins and settings into versionable .yaml files
 - [] Add way of applying .yaml config files to update Kong
 - [] Store a diff of local config files and what's on the remote so we can only make the requests for what has changed, if something has been removed then we need to remove it, https://github.com/go-yaml/yaml/, https://github.com/kylelemons/godebug/
 - [] Refactor & hook up tests into a CI tool
 - [] Keep track of the Kong hosts as "contexts" and allow switching
 - [] Document tools api & options
 - [] Produce pretty landing page & branding
 - [] Add to Homebrew