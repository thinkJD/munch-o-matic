# munch-o-matic

My children receive their school meals from a local caterer. Every day, there are three dishes to choose from. The dishes are announced about four weeks in advance, and you can select what you would like to eat using a web UI. Sometimes, everyday life gets really busy, and it can be a close call with placing the order.

Munch-o-matic can order dishes automatically based on multiple strategies. Use it from your terminal or run it in daemon mode on your server. You can change the orders whenever you like using the official web UI.

## Synopsis

```
TasteNext API Client

Usage:
  munch-o-matic [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  daemon      Run munch-o-matic in daemon mode
  help        Help about any command
  info        Information and statistics
  menu        Work with menus
  order       Order or cancel a dish from the menu

Flags:
      --config string   config file (default is $HOME/.munch-o-matic.yaml)
  -h, --help            help for munch-o-matic

Use "munch-o-matic [command] --help" for more information about a command.
```
