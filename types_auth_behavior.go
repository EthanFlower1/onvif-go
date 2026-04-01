package onvif

// AuthBehaviorServiceCapabilities represents the capabilities of the Authentication Behavior service.
type AuthBehaviorServiceCapabilities struct {
	MaxLimit                                uint
	MaxAuthenticationProfiles               uint
	MaxPoliciesPerAuthenticationProfile     uint
	MaxSecurityLevels                       uint
	MaxRecognitionGroupsPerSecurityLevel    uint
	MaxRecognitionMethodsPerRecognitionGroup uint
	ClientSuppliedTokenSupported            bool
	SupportedAuthenticationModes            string
}

// AuthenticationProfileInfo contains basic information about an authentication profile.
type AuthenticationProfileInfo struct {
	Token       string
	Name        string
	Description string
}

// AuthenticationPolicy associates a security level with a schedule.
type AuthenticationPolicy struct {
	ScheduleToken              string
	SecurityLevelConstraints   []SecurityLevelConstraint
}

// SecurityLevelConstraint defines what security level should be active depending on the schedule state.
type SecurityLevelConstraint struct {
	ActiveRegularSchedule    bool
	ActiveSpecialDaySchedule bool
	AuthenticationMode       string
	SecurityLevelToken       string
}

// AuthenticationProfile contains all properties of AuthenticationProfileInfo plus policies.
type AuthenticationProfile struct {
	AuthenticationProfileInfo
	DefaultSecurityLevelToken string
	AuthenticationPolicies    []AuthenticationPolicy
}

// SecurityLevelInfo contains basic information about a security level instance.
type SecurityLevelInfo struct {
	Token       string
	Name        string
	Priority    int
	Description string
}

// RecognitionMethod defines a recognition method with type and order.
type RecognitionMethod struct {
	RecognitionType string
	Order           int
}

// RecognitionGroup contains a list of recognition methods.
type RecognitionGroup struct {
	RecognitionMethods []RecognitionMethod
}

// SecurityLevel contains all properties of SecurityLevelInfo plus recognition groups.
type SecurityLevel struct {
	SecurityLevelInfo
	RecognitionGroups []RecognitionGroup
}
