# munch-o-matic

My children receive their school meals from a local caterer. Every day, there are three dishes to choose from. The dishes are announced about four weeks in advance, and you can select what you would like to eat using a web UI. Sometimes, everyday life gets really busy, and it can be a close call with placing the order.

munch-o-matic can order dishes automatically based on multiple strategies. Use it from your terminal or run it in daemon mode on your server. You can change the orders whenever you like using the official web UI.

## Jobs

checkBalance
Weekly on Wednesday 17:00
If balance is bellow 20€ (5 Days)

orderFood
Weekly on Wednesday 17:00
auto-order --weeks 2 --strategy SchoolFav

updateMetrics
Get all upcoming dishes and update the counter
Hourly

## Demo

![](https://github.com/thinkJD/munch-o-matic/blob/main/assets/render1700930307166.gif)

## Synopsis

```
                                _                                             _    _
 _ __ ___   _   _  _ __    ___ | |__           ___          _ __ ___    __ _ | |_ (_)  ___
| '_ ` _ \ | | | || '_ \  / __|| '_ \  _____  / _ \  _____ | '_ ` _ \  / _` || __|| | / __|
| | | | | || |_| || | | || (__ | | | ||_____|| (_) ||_____|| | | | | || (_| || |_ | || (__
|_| |_| |_| \__,_||_| |_| \___||_| |_|        \___/        |_| |_| |_| \__,_| \__||_| \___|

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
