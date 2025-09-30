package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"research-apm/services/alert/internal/entity"
	"time"
)

type AlertAPM struct {
	ServiceName string    `json:"service.name"`
	Environment string    `json:"service.environment"`
	Category    string    `json:"kibana.alert.rule.category"`
	RuleName    string    `json:"kibana.alert.rule.name"`
	Message     string    `json:"kibana.alert.reason"`
	Status      string    `json:"kibana.alert.status"`
	Timestamp   time.Time `json:"@timestamp"`
}

type AlertAPMHit struct {
	ID     string   `json:"_id"`
	Source AlertAPM `json:"_source"`
}

func (a AlertAPMHit) ToEntity() entity.Alert {
	h := sha256.New()
	h.Write([]byte(a.ID + a.Source.Status))
	trxID := hex.EncodeToString(h.Sum(nil))
	return entity.Alert{
		AlertID:     a.ID,
		TrxID:       trxID,
		ServiceName: a.Source.ServiceName,
		Environment: a.Source.Environment,
		Category:    a.Source.Category,
		RuleName:    a.Source.RuleName,
		Message:     a.Source.Message,
		Status:      a.Source.Status,
		Timestamp:   a.Source.Timestamp,
	}
}

func NewMessage(a entity.Alert) string {
	icon := "❓"
	switch a.Status {
	case "active":
		icon = "🚨"
	case "recovered":
		icon = "✅"
	}
	return fmt.Sprintf(`
<b>%s [%s]</b>
<b>%s</b>
%s
%s (UTC)

%s`,
		icon,
		a.Status,
		a.ServiceName,
		a.RuleName,
		a.Timestamp.Format("02-01-2006 15:04:05"),
		a.Message,
	)
}
