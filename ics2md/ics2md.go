package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	ics "github.com/arran4/golang-ical"
)

const (
	icsDateFormat  = "20060102T150405Z"
	dateFormat     = "2006/01/02"
	dateTimeFormat = "2006/01/02 15:04:05"
)

func getCal(fp string) (*ics.Calendar, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	c, err := ics.ParseCalendar(f)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func main() {
	fp := os.Args[1]
	c, err := getCal(fp)
	if err != nil {
		panic(err)
	}
	dateEventDict := make(map[string][]*ics.VEvent)
	for _, ev := range c.Events() {
		st, err := ev.GetAllDayStartAt()
		if err != nil {
			panic(err)
		}
		dateStr := st.Format(dateFormat)
		dtEvs, ok := dateEventDict[dateStr]
		if ok {
			dateEventDict[dateStr] = append(dtEvs, ev)
		} else {
			dateEventDict[dateStr] = []*ics.VEvent{ev}
		}
	}
	dateStrs := make([]string, 0, len(dateEventDict))
	for dtstr, _ := range dateEventDict {
		dateStrs = append(dateStrs, dtstr)
	}
	sort.Strings(dateStrs)
	var sb strings.Builder
	for _, dtstr := range dateStrs {
		sb.WriteString("## ")
		sb.WriteString(dtstr)
		sb.WriteString("\n\n")
		evs := dateEventDict[dtstr]
		sort.Slice(evs, func(i, j int) bool {
			// should no be error this time
			ist, _ := evs[i].GetStartAt()
			jst, _ := evs[j].GetStartAt()
			return ist.Before(jst)
		})
		for _, ev := range evs {
			st, _ := ev.GetStartAt()
			if len(evs) > 1 {
				sb.WriteString("### ")
				sb.WriteString(st.Format(dateTimeFormat))
				sb.WriteString("\n\n")
			}
			desc := ev.GetProperty(ics.ComponentPropertyDescription)
			if desc != nil {
				// ics format, newline is written with \n string
				descNL := strings.ReplaceAll(desc.Value, `\n`, "\n")
				sb.WriteString(descNL)
			}
			sb.WriteString("\n\n")
		}
	}
	fmt.Println(sb.String())
}
