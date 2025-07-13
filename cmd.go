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
		setPillCmd,
		checkPillCmd,
		notifyCmd,
		setNotifyCmd,
	},
	Shortcuts:   Z.ArgMap{},
	Version:     `v0.0.2`,
	Source:      `https://github.com/espinosajuanma/pills`,
	Issues:      `https://github.com/espinosajuanma/pills/issues`,
	Summary:     `Add `,
	Description: ``,
}

var setPillCmd = &Z.Cmd{
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

var setNotifyCmd = &Z.Cmd{
	Name:        `set-notify`,
	Aliases:     []string{"config-notify"},
	Usage:       "",
	Commands:    []*Z.Cmd{help.Cmd},
	Summary:     `configure SMTP notification settings`,
	Description: `Prompts for SMTP configuration values needed for sending notification emails.`,
	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) > 0 {
			return x.UsageError()
		}

		// Host
		currentHost, _ := x.Root().Get("smtp.host")
		host := term.Prompt("SMTP Host [%s] ", currentHost)
		if host != "" {
			x.Root().Set("smtp.host", host)
		} else if currentHost == "" {
			return fmt.Errorf("smtp.host can't be empty")
		}

		// Port
		currentPort, _ := x.Root().Get("smtp.port")
		if currentPort == "" {
			currentPort = "587"
		}
		portStr := term.Prompt("SMTP Port [%s] ", currentPort)
		if portStr != "" {
			if _, err := strconv.Atoi(portStr); err != nil {
				return fmt.Errorf("invalid smtp.port: %w", err)
			}
			x.Root().Set("smtp.port", portStr)
		} else if currentPort == "" {
			return fmt.Errorf("smtp.port can't be empty")
		}

		// User
		currentUser, _ := x.Root().Get("smtp.user")
		user := term.Prompt("SMTP User [%s] ", currentUser)
		if user != "" {
			x.Root().Set("smtp.user", user)
		} else if currentUser == "" {
			return fmt.Errorf("smtp.user can't be empty")
		}

		// Password
		currentPass, _ := x.Root().Get("smtp.pass")
		pass := term.Prompt("SMTP Password (leave blank to keep current): ")
		if pass != "" {
			x.Root().Set("smtp.pass", pass)
		} else if currentPass == "" {
			return fmt.Errorf("smtp.pass can't be empty")
		}

		// From
		currentFrom, _ := x.Root().Get("smtp.from")
		from := term.Prompt("From address [%s] ", currentFrom)
		if from != "" {
			x.Root().Set("smtp.from", from)
		} else if currentFrom == "" {
			return fmt.Errorf("smtp.from can't be empty")
		}

		// To
		currentTo, _ := x.Root().Get("smtp.to")
		if currentTo == "" {
			currentTo = currentFrom
		}
		to := term.Prompt("To address [%s] ", currentTo)
		if to != "" {
			x.Root().Set("smtp.to", to)
		} else if currentTo == "" {
			return fmt.Errorf("smtp.to can't be empty")
		}

		term.Print("SMTP configuration saved.")
		return nil
	},
}

var checkPillCmd = &Z.Cmd{
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

var notifyCmd = &Z.Cmd{
	Name:        `notify`,
	Aliases:     []string{},
	Commands:    []*Z.Cmd{help.Cmd},
	Summary:     `sends a warning email if the pills are close to run out`,
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
			smtpConfig, err := GetSmtpConfig(x.Root())
			if err != nil {
				return fmt.Errorf("failed to get smtp config: %w", err)
			}

			subject := "💊 Refill Reminder"
			body := pill.ReminderMessage()

			err = smtpConfig.SendEmail(subject, body)
			if err != nil {
				return fmt.Errorf("failed to send email: %w", err)
			}
			term.Printf("Notification email sent for %s", pill.PillName)
		}
		return nil
	},
}
