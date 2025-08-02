package queue_ranked

import (
	"fmt"
	"io"

	"github.com/junioryono/Riot-API-Golang/constants/queue"
)

type ID queue.ID
type String string

const (
	RankedSolo5x5 ID = 420
	RankedFlexSR  ID = 440
	RankedFlexTT  ID = 470
)

var stringToIDMap = map[String]ID{
	"RANKED_SOLO_5x5": RankedSolo5x5,
	"RANKED_FLEX_SR":  RankedFlexSR,
	"RANKED_FLEX_TT":  RankedFlexTT,
}

var idToStringMap = map[ID]String{
	RankedSolo5x5: "RANKED_SOLO_5x5",
	RankedFlexSR:  "RANKED_FLEX_SR",
	RankedFlexTT:  "RANKED_FLEX_TT",
}

func (q ID) String() String {
	return idToStringMap[q]
}

func (q String) ID() ID {
	return stringToIDMap[q]
}

func (q String) PrettyString() queue.PrettyString {
	return queue.PrettyString(q)
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (q *ID) UnmarshalGQL(v interface{}) error {
	intValue, ok := v.(int16)
	if !ok {
		return fmt.Errorf("rank must be an int16")
	}

	*q = ID(intValue)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (q ID) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, q)
}
