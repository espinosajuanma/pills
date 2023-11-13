package pills

import (
	"fmt"
	"strconv"
	"time"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/rwxrob/term"
	"github.com/rwxrob/vars"
)

func init() {
	Z.Vars.SoftInit()
}

var Cmd = &Z.Cmd{
	Name: `pills`,
	Commands: []*Z.Cmd{
		help.Cmd, vars.Cmd, conf.Cmd, // common
		setCmd,
		checkCmd,
	},
	Shortcuts:   Z.ArgMap{},
	Version:     `v0.0.1`,
	Source:      `https://github.com/espinosajuanma/pills`,
	Issues:      `https://github.com/espinosajuanma/pills/issues`,
	Summary:     `Add `,
	Description: ``,
}

var setCmd = &Z.Cmd{
	Name:        `set`,
	Aliases:     []string{"add", "update"},
	Usage:       "<alias>",
	Commands:    []*Z.Cmd{help.Cmd},
	Summary:     `configure a new pill`,
	Description: ``,
	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) == 0 {
			return x.UsageError()
		}
		alias := args[0]

		conf, err := x.Root().Get("pill." + alias)
		if err != nil {
			return err
		}
		var pill Pill
		if conf != "" {
			pill, err = FromJson(conf)
			if err != nil {
				return err
			}
		}

		pillName := term.Prompt("Pill Name [%s] ", pill.PillName)
		if pillName == "" && pill.PillName == "" {
			return fmt.Errorf("pill name can't be empty")
		}
		if pillName != "" {
			pill.PillName = pillName
		}

		stockAmount := term.Prompt("Stock Amount [%d] ", pill.StockAmount)
		if stockAmount == "" && pill.StockAmount == 0 {
			return fmt.Errorf("stock amount can't be 0")
		}
		if stockAmount != "" {
			pill.StockAmount, err = strconv.Atoi(stockAmount)
			if err != nil {
				return err
			}
		}

		dosageFrequency := term.Prompt("Dosage Frequency [%s] ", pill.DosageFrequency)
		if dosageFrequency == "" && pill.DosageFrequency == 0 {
			return fmt.Errorf("dosage frequency can't be empty")
		}
		if dosageFrequency != "" {
			d, err := time.ParseDuration(dosageFrequency)
			if err != nil {
				return err
			}
			pill.DosageFrequency = d
		}

		dosesPerIntake := term.Prompt("Doses per intake [%.2f] ", pill.DosesPerIntake)
		if dosesPerIntake == "" && pill.DosesPerIntake == 0 {
			return fmt.Errorf("doses per intake can't be 0")
		}
		if dosesPerIntake != "" {
			pill.DosesPerIntake, err = strconv.ParseFloat(dosesPerIntake, 64)
			if err != nil {
				return err
			}
		}

		if pill.StockDate.IsZero() {
			pill.StockDate = time.Now()
		}
		stockDate := term.Prompt("Stock date [%s] ", pill.StockDate.Format(DATE_FORMAT))
		if stockDate != "" {
			_, err := pill.SetStockDate(stockDate)
			if err != nil {
				return err
			}
		}

		reminderDays := term.Prompt("Refill reminder days [%d] ", pill.RefillReminderDays)
		if reminderDays == "" && pill.RefillReminderDays == 0 {
			return fmt.Errorf("refill reminder days can't be 0")
		}
		if reminderDays != "" {
			pill.RefillReminderDays, err = strconv.Atoi(reminderDays)
			if err != nil {
				return err
			}
		}

		pill.CalculateDates()
		term.Printf("Exhaustion Date: %s", pill.ExhaustionDate.Format(DATE_FORMAT))
		term.Printf("Reminder Date: %s", pill.ReminderDate.Format(DATE_FORMAT))

		parsed, err := pill.ToJson()
		if err != nil {
			return err
		}
		x.Root().Set("pill."+alias, parsed)

		return nil
	},
}

var checkCmd = &Z.Cmd{
	Name:        `check`,
	Aliases:     []string{"alarm"},
	Commands:    []*Z.Cmd{help.Cmd},
	Summary:     `returns a warning message if the pills are close to run out`,
	Description: ``,
	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) == 0 {
			return x.UsageError()
		}
		alias := args[0]
		t, err := x.Root().Get("pill." + alias)
		if err != nil {
			return err
		}
		if t == "" {
			return fmt.Errorf("%s is not a pill", alias)
		}
		pill, err := FromJson(t)
		if err != nil {
			return err
		}

		if pill.ShouldRemind() {
			term.Print(pill.ReminderMessage())
		}
		return nil
	},
}
