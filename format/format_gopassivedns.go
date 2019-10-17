// balboa
// Copyright (c) 2018, DCSO GmbH

package format

import (
	"encoding/json"
	"net"
	"time"

	"balboa/observation"

	log "github.com/sirupsen/logrus"
)

type dnsLogEntry struct {
	QueryID             uint16 `json:"query_id"`
	ResponseCode        int    `json:"rcode"`
	Question            string `json:"q"`
	QuestionType        string `json:"qtype"`
	Answer              string `json:"a"`
	AnswerType          string `json:"atype"`
	TTL                 uint32 `json:"ttl"`
	Server              net.IP `json:"dst"`
	Client              net.IP `json:"src"`
	Timestamp           string `json:"tstamp"`
	Elapsed             int64  `json:"elapsed"`
	ClientPort          string `json:"sport"`
	Level               string `json:"level"`
	Length              int    `json:"bytes"`
	Proto               string `json:"protocol"`
	Truncated           bool   `json:"truncated"`
	AuthoritativeAnswer bool   `json:"aa"`
	RecursionDesired    bool   `json:"rd"`
	RecursionAvailable  bool   `json:"ra"`
}

// MakeGopassivednsInputObservations is a MakeObservationFunc that accepts
// input in the format as generated by https://github.com/Phillipmartin/gopassivedns.
func MakeGopassivednsInputObservations(inputJSON []byte, sensorID string, out chan observation.InputObservation, stop chan bool) error {
	var in dnsLogEntry
	err := json.Unmarshal(inputJSON, &in)
	if err != nil {
		log.Warn(err)
		return nil
	}
	tst, err := time.Parse("2006-01-02 15:04:05.999999 -0700 MST", in.Timestamp)
	if err != nil {
		log.Info(err)
		return nil
	}
	o := observation.InputObservation{
		Count:          1,
		Rdata:          in.Answer,
		Rrname:         in.Question,
		Rrtype:         in.AnswerType,
		SensorID:       sensorID,
		TimestampEnd:   tst,
		TimestampStart: tst,
	}
	out <- o
	log.Debug("enqueued 1 observation")
	return nil
}
