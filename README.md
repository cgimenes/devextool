# devextool

A simple task runner and automation tool to provide a better developer experience.

## Introduction and Rationale

Good UX (in my mind and a bit simplified) is defined by a few traits:

  * Clear/Common sense of how to accomplish a goal/task
  * Discoverable when it's a new task (or not as clear/common)
  * Consistency
  
There are lots of existing tools that provide task running/automation (Make, Thor/Rakefiles, Mage, Taskfile.dev, a number in JS land, etc) so why a new one?  I feel they are all directory focused (there are ways around that but seem like hacks). Instead, primarily I want a tool that centralises as much as possible to DRY up all those little shell scripts/functions/aliases.

## Using dx

dx runs programs from a fixed directory (DXHome by default) The location of this folder would searched for in order:  

  * specified using an env variable (DXHOME - eg to use with direnv).
  * In the current directory, 
  * $HOME/DXHome

There are only a few things to know about the contents of DXHome:

  * The only required directory is cmd.  It contains executables/scripts and subdirectories.  Subdirectories create 'namespace' comamnds to group related commands together.
  * The scripts/programs in cmd must be executable (eg chmod 755) but can be written in any programming lanuage (but probably shell scripts).
  * Folders and files inside cmd that start with an underscore are ignored (eg cmd/_lib to house shared functions in your shell scripts)
  * $DXHOME is set (eg so you can "source $DXHOME/utils.sh" or cat $DXHOME/README.me) for commands
  * DXHome can contain anything else you might need (Docs, config files, etc)
  * $GOOS_cmd can be used as the command folder to provide platform specific set of commands (eg for Windows, note this overrides cmd)

Commands can have options specified for them by creating a json file int he same directory as the command. This should help make the dev tools you are running be self documenting. The file will be named _meta_COMMANDNAME.json.  The format should be as follows (each is optional):

```
{
  "use": "cmd --usage --flags --etc",
  "short": "Short description that is displayed in the list of commands",
  "long": "Block of text shown when dx <command> --help is run",
  "aliases": ["an", "array", "of", "alt", "names"],
  "example": "a string that shows an example of running command"
}
```

## Examples

TBD

## Building

Prebuilt binaries are (will be) generated and availabe under Releases but a simple 

``` 
  git clone https://github.com/vhodges/devextool.git
  cd devextool
  go build ./...
```

should suffice.

## License

MIT

