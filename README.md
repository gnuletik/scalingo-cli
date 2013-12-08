Appsdeck-CLI v0.2.1
===================

Command line utility to manage its appsdeck application.

```
NAME:
   Appsdeck Client - Manage your apps and containers

USAGE:
   Appsdeck Client [global options] command [command options] [arguments...]

VERSION:
   0.2.1

COMMANDS:
   logs, l	[-n <nblines> | --stream]
   run, r	Run any command for your app
   apps, a	Manage your apps
   logout	Logout from Appsdeck
   create, c	appsdeck create <name>
   destroy, d	appsdeck destroy <id or canonical name>
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --app '<name>'	Name of the app
   --version		print the version
   --help, -h		show help

```

Dev usage
---------

Define (example) :

* `APPSDECK_LOG=http://127.0.0.1:10004`
* `APPSDECK_API=http://127.0.0.1`
* `DEBUG=1`
