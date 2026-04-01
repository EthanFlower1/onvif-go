package onvif

// ScheduleServiceCapabilities represents the capabilities of the Schedule service.
type ScheduleServiceCapabilities struct {
	MaxLimit                     uint
	MaxSchedules                 uint
	MaxTimePeriodsPerDay         uint
	MaxSpecialDayGroups          uint
	MaxDaysInSpecialDayGroup     uint
	MaxSpecialDaysSchedules      uint
	ExtendedRecurrenceSupported  bool
	SpecialDaysSupported         bool
	StateReportingSupported      bool
	ClientSuppliedTokenSupported bool
}

// ScheduleInfo contains basic information about a schedule instance.
type ScheduleInfo struct {
	Token       string
	Name        string
	Description string
}

// Schedule includes all properties of ScheduleInfo plus iCalendar data and special days.
type Schedule struct {
	ScheduleInfo
	Standard    string
	SpecialDays []SpecialDaysSchedule
}

// TimePeriod defines a start and optional end time within a day.
type TimePeriod struct {
	From  string
	Until string
}

// SpecialDaysSchedule defines alternate time periods for a group of special days.
type SpecialDaysSchedule struct {
	GroupToken string
	TimeRange  []TimePeriod
}

// ScheduleState contains state information for a schedule.
type ScheduleState struct {
	Active     bool
	SpecialDay *bool
}

// SpecialDayGroupInfo contains basic information about a special day group.
type SpecialDayGroupInfo struct {
	Token       string
	Name        string
	Description string
}

// SpecialDayGroup includes all properties of SpecialDayGroupInfo plus iCalendar day data.
type SpecialDayGroup struct {
	SpecialDayGroupInfo
	Days string
}
