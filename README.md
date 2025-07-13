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
- `pills set-notify`
  - Set SMTP Configuration
- `pills notify <alias>`
  - Sends a notification email if a reminder is due for the given pill alias.

## SMTP Configuration

To use the `notify` command, you need to configure your SMTP server details.
You can do this using the `set-notify` command 

```bash
pills set-noitify
```

### Cronjob Alarm

Set up an alarm every day at one hour that you prefer. Edit your cron jobs using
`crontab -e` and add this line:

```bash
0 8 * * * pills notify <alias>
```

## Tab completion

To activate bash completion use the `complete -C` option from your `.bashrc` or command line.
Tehere is no messy sourcin required. All the completion is done by the program itself.

```bash
complete -C pills pills
```