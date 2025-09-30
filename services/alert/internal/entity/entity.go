package entity

import "time"

type Alert struct {
	AlertID     string    `json:"alert_id"`
	TrxID       string    `json:"trx_id"`
	ServiceName string    `json:"service_name"`
	Environment string    `json:"environment"`
	Category    string    `json:"category"`
	RuleName    string    `json:"rule_name"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}
