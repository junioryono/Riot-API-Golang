package staticdata

import (
	"fmt"

	"github.com/junioryono/Riot-API-Golang/constants/language"
	"github.com/junioryono/Riot-API-Golang/constants/patch"
	"github.com/junioryono/Riot-API-Golang/constants/queue"
)

type Queues []Queue

type Queue struct {
	QueueID     queue.ID `json:"queueId"`
	Map         string   `json:"map"`
	Description string   `json:"description"`
	Notes       string   `json:"notes"`
}

func GetQueues(v patch.Patch, lang language.Language) (Queues, error) {
	var res Queues
	err := getJSON("https://static.developer.riotgames.com/docs/lol/queues.json", &res)
	return res, err
}

func (queues Queues) Queue(queueId int) (Queue, error) {
	for _, q := range queues {
		if q.QueueID == queue.ID(queueId) {
			return q, nil
		}
	}

	return Queue{}, fmt.Errorf("queue %d not found", queueId)
}
