package model

import "time"

type EventMapper struct {
	MapperID           string            `json:"mapper_id"`
	Name               string            `json:"name"`
	CreatedAt          time.Time         `json:"created_at"`
	MappingInstruction map[string]string `json:"mapping_instruction"`
}

type Event struct {
	MapperID  string                 `json:"mapper_id"`
	Drcr      string                 `json:"drcr"`
	BalanceID string                 `json:"balance_id"`
	Data      map[string]interface{} `json:"data"`
}
