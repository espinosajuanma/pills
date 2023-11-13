# Pills

Simple CLI tool to remind to buy your medicine

## Motivation

I've been sick since almost a year now and I have to take an annoying amount of pills 
I know myself and I've been cautious about making alarms to take the pills on time.
I mean, at the start I am pretty aware of taken them. But once I got used to them I can forget them.

Now is time to make my refill reminders so I don't forget to buy my pills.

## Installation

```bash
go install github.com/espinosajuanma/pills/cmd/pills@latest
```

## Commands

- `pills add <alias>`
  - Warning if alias already exists
  - Prompt to setup the pill configurations
    - Pill Name
    - Stock Amount
    - Frequency Time
    - Doses Amount
    - Date (Today) 
    - Remind At
- `pills peek <alias>`
  - Prints current pill configurations
- `pills check <alias>`
  - Prints reminder 

## Cronjob Alarm

Set up an alarm every day at one hour that you prefer.

```bash
#!/bin/bash
msg="$(pills check ibuprofen)"
if [[ ! -z "$msg" ]]; then
  notify-send "💊 Refill Reminder" "$msg"
fi
```

The last example would just send a notification. The preffered way you want the notification.
My intention is connect it to my phone push notifications.

## Tab completion

To activate bash completion use the `complete -C` option from your `.bashrc` or command line.
Tehere is no messy sourcin required. All the completion is done by the program itself.

```bash
complete -C pills pills
```