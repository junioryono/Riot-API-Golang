package patch

import (
	"fmt"
	"strings"
	"time"

	"github.com/junioryono/Riot-API-Golang/constants/region"
)

type ShortPatch string

func (v ShortPatch) String() string {
	return string(v)
}

type Patch string

func (v Patch) String() string {
	return string(v)
}

func (v *Patch) ShortPatch() ShortPatch {
	vSplit := strings.Split(v.String(), ".")
	if len(vSplit) < 2 {
		return ShortPatch(*v)
	}

	return ShortPatch(fmt.Sprintf("%s.%s", vSplit[0], vSplit[1]))
}

type PatchWithStartTime struct {
	Patch     ShortPatch
	StartTime time.Time
	Shifts    map[region.Region]int
}

func (v *PatchWithStartTime) GetRegionStartDate(region region.Region) time.Time {
	// Add the shift to the start date
	shift, ok := v.Shifts[region]
	if !ok {
		return v.StartTime
	}

	// Shift is in seconds
	return v.StartTime.Add(time.Duration(shift) * time.Second)
}
