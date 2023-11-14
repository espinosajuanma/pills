package pills

import (
	"encoding/json"
	"fmt"
	"time"
)

const DATE_FORMAT = "02-01-2006"

type Pill struct {
	PillName           string        `json:"pillName"`
	StockAmount        int           `json:"stockAmount"`
	DosageFrequency    time.Duration `json:"dosageFrequency"`
	DosesPerIntake     float64       `json:"dosesPerIntake"`
	StockDate          time.Time     `json:"stockDate"`
	RefillReminderDays int           `json:"refillReminderDays"`
	ExhaustionDate     time.Time     `json:"exhaustionDate"`
	ReminderDate       time.Time     `json:"reminderDate"`
}

func (p *Pill) SetStockDate(dateString string) (time.Time, error) {
	if dateString == "" {
		dateString = time.Now().Local().Format(DATE_FORMAT)
	}
	t, err := time.Parse(DATE_FORMAT, dateString)
	if err != nil {
		return t, err
	}
	p.StockDate = t
	return t, nil
}

func (p *Pill) CalculateDates() {
	totalDoses := float64(p.StockAmount) / p.DosesPerIntake
	totalDuration := time.Duration(totalDoses) * p.DosageFrequency
	p.ExhaustionDate = p.StockDate.Add(totalDuration)

	remindAt := time.Duration(p.RefillReminderDays*24) * time.Hour
	remindDate := p.ExhaustionDate.Add(-remindAt)
	p.ReminderDate = remindDate
}

func (p Pill) ShouldRemind() bool {
	if p.ExhaustionDate.IsZero() || p.ReminderDate.IsZero() {
		p.CalculateDates()
	}
	today := time.Now()
	return today.After(p.ReminderDate)
}

func (p Pill) ReminderMessage() string {
	return fmt.Sprintf("Running out of [%s] on [%s]", p.PillName, p.ExhaustionDate.Format(DATE_FORMAT))
}

func FromJson(str string) (Pill, error) {
	var pill Pill
	err := json.Unmarshal([]byte(str), &pill)
	if err != nil {
		return Pill{}, err
	}
	return pill, nil
}

func (p Pill) ToJson() (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
