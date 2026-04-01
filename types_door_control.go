package onvif

// DoorControlServiceCapabilities represents the capabilities of the Door Control service.
type DoorControlServiceCapabilities struct {
	MaxLimit                    uint
	MaxDoors                    uint
	ClientSuppliedTokenSupported bool
	DoorManagementSupported     bool
}

// DoorCapabilities reflects optional functionality of a particular door instance.
type DoorCapabilities struct {
	Access              *bool
	AccessTimingOverride *bool
	Lock                *bool
	Unlock              *bool
	Block               *bool
	DoubleLock          *bool
	LockDown            *bool
	LockOpen            *bool
	DoorMonitor         *bool
	LockMonitor         *bool
	DoubleLockMonitor   *bool
	Alarm               *bool
	Tamper              *bool
	Fault               *bool
}

// DoorInfo represents a Door as a physical object with capabilities.
type DoorInfo struct {
	Token        string
	Name         string
	Description  string
	Capabilities DoorCapabilities
}

// DoorTimings defines timing parameters for a door.
type DoorTimings struct {
	ReleaseTime            string
	OpenTime               string
	ExtendedReleaseTime    string
	DelayTimeBeforeRelock  string
	ExtendedOpenTime       string
	PreAlarmTime           string
}

// Door includes all properties of DoorInfo plus type and timing configuration.
type Door struct {
	DoorInfo
	DoorType string
	Timings  DoorTimings
}

// DoorPhysicalState represents the physical state of a door.
type DoorPhysicalState string

// Door physical state constants.
const (
	DoorPhysicalStateUnknown DoorPhysicalState = "Unknown"
	DoorPhysicalStateOpen    DoorPhysicalState = "Open"
	DoorPhysicalStateClosed  DoorPhysicalState = "Closed"
	DoorPhysicalStateFault   DoorPhysicalState = "Fault"
)

// LockPhysicalState represents the physical state of a lock.
type LockPhysicalState string

// Lock physical state constants.
const (
	LockPhysicalStateUnknown  LockPhysicalState = "Unknown"
	LockPhysicalStateLocked   LockPhysicalState = "Locked"
	LockPhysicalStateUnlocked LockPhysicalState = "Unlocked"
	LockPhysicalStateFault    LockPhysicalState = "Fault"
)

// DoorAlarmState describes the alarm state of a door.
type DoorAlarmState string

// Door alarm state constants.
const (
	DoorAlarmStateNormal        DoorAlarmState = "Normal"
	DoorAlarmStateDoorForcedOpen DoorAlarmState = "DoorForcedOpen"
	DoorAlarmStateDoorOpenTooLong DoorAlarmState = "DoorOpenTooLong"
)

// DoorTamperState describes the state of a tamper detector.
type DoorTamperState string

// Door tamper state constants.
const (
	DoorTamperStateUnknown       DoorTamperState = "Unknown"
	DoorTamperStateNotInTamper   DoorTamperState = "NotInTamper"
	DoorTamperStateTamperDetected DoorTamperState = "TamperDetected"
)

// DoorTamper contains tampering information for a door.
type DoorTamper struct {
	Reason string
	State  DoorTamperState
}

// DoorFaultState describes the state of a door fault.
type DoorFaultState string

// Door fault state constants.
const (
	DoorFaultStateUnknown      DoorFaultState = "Unknown"
	DoorFaultStateNotInFault   DoorFaultState = "NotInFault"
	DoorFaultStateFaultDetected DoorFaultState = "FaultDetected"
)

// DoorFault contains fault information for a door.
type DoorFault struct {
	Reason string
	State  DoorFaultState
}

// DoorMode represents the logical operating mode of a door.
type DoorMode string

// Door mode constants.
const (
	DoorModeUnknown     DoorMode = "Unknown"
	DoorModeLocked      DoorMode = "Locked"
	DoorModeUnlocked    DoorMode = "Unlocked"
	DoorModeAccessed    DoorMode = "Accessed"
	DoorModeBlocked     DoorMode = "Blocked"
	DoorModeLockedDown  DoorMode = "LockedDown"
	DoorModeLockedOpen  DoorMode = "LockedOpen"
	DoorModeDoubleLocked DoorMode = "DoubleLocked"
)

// DoorState contains the current aggregate runtime status of a door.
type DoorState struct {
	DoorPhysicalState       *DoorPhysicalState
	LockPhysicalState       *LockPhysicalState
	DoubleLockPhysicalState *LockPhysicalState
	Alarm                   *DoorAlarmState
	Tamper                  *DoorTamper
	Fault                   *DoorFault
	DoorMode                DoorMode
}
