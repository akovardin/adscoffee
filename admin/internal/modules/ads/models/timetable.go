package models

import "encoding/json"

type Timetable map[int]map[int]bool // 7 дней, 24 часа

func NewTimetable(data string) (Timetable, error) {
	t := Timetable{}

	if data == "" {
		return t, nil
	}

	if err := json.Unmarshal([]byte(data), &t); err != nil {
		return Timetable{}, err
	}

	return t, nil
}

func (t Timetable) Validate(day, hour int) bool {
	if len(t) == 0 {
		return true
	}

	if day < 1 || day > 7 || hour < 0 || hour > 23 {
		return false
	}

	val, ok := t[day][hour]
	if !ok {
		return false
	}

	return val
}

func (t Timetable) Set(day, hour int, value bool) {
	if t[day] == nil {
		t[day] = map[int]bool{}
	}
	t[day][hour] = value
}

func (t Timetable) String() string {
	if len(t) == 0 {
		return ""
	}

	data, _ := json.Marshal(t)

	return string(data)
}
