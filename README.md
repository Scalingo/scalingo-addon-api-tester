# Scalingo Addon API Tester tool

The goal of this tool is to help developers to test their Addon
web service before deploying it in production to ensure it answers
correctly to requests and to check if the manifest is correctly
written.

The documentation of the expected API can be found on our API
documentation website: http://developers.scalingo.com/addons

## Install

You need go to be installed on your computer

```
go get github.com/Scalingo/scalingo-addon-api-tester
```

## The addon manifest

The tool is expecting a `manifest.json` in the current directory.
The file should respect the format of the addon manifest documented
here: http://developers.scalingo.com/addons/manifest.html

You can specify another path with the global flag: `--manifest`

```sh
scalingo-addon-api-tester --manifest provision
```

## Commands

### Provision an addon

```sh
scalingo-addon-api-tester provision [--plan <plan>] [--app <app>]
```

Both flags are optional, the tool is generating random app name if
it is not specified on the command line, and the default plan is the
first defined in your manifest.

### Update an addon

```sh
scalingo-addon-api-tester update <id> --plan <plan>
```

Use an existing addon ID (use `list` command to get them) and make the
request to the addon web server to update the plan of the resource.

### Deprovision an addon

```sh
scalingo-addon-api-tester deprovision <id>
```

Make a request to deprovision an addon with the given ID

## Helper commands

The command line is saving the history of your provisionning in a file
(`$HOME/.scalingo-addon-tester`), and it implements different command
to display those data.

Example:

```sh
scalingo-addon-api-tester provision
→ OK
scalingo-addon-api-tester list
- addon-1: free
scalingo-addon-api-tester update addon-1 --plan premium
→ OK
scalingo-addon-api-tester list
- addon-1: premium
scalingo-addon-api-tester deprovision addon-1
→ OK
```

### List addons

```sh
scalingo-addon-api-tester list
```

Display all the addons which have been saved.

### Purge addons

```sh
scalingo-addon-api-tester purge
```

If the database is completely invalid, just run the purge command,
then the `list` command won't list anything anymore. Addons are
not deprovisioned, this command is just a helper for this command
line tool.
