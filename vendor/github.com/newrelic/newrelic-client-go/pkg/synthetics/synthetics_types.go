package synthetics

// MonitorScriptLocation represents a New Relic Synthetics monitor script location.
type MonitorScriptLocation struct {
	Name string `json:"name"`
	HMAC string `json:"hmac"`
}

// MonitorScript represents a New Relic Synthetics monitor script.
type MonitorScript struct {
	Text      string                  `json:"scriptText"`
	Locations []MonitorScriptLocation `json:"scriptLocations"`
}

// MonitorType represents a Synthetics monitor type.
type MonitorType string

// MonitorStatusType represents a Synthetics monitor status type.
type MonitorStatusType string

var (
	// MonitorTypes specifies the possible types for a Synthetics monitor.
	MonitorTypes = struct {
		Ping            MonitorType
		Browser         MonitorType
		ScriptedBrowser MonitorType
		APITest         MonitorType
	}{
		Ping:            "SIMPLE",
		Browser:         "BROWSER",
		ScriptedBrowser: "SCRIPT_BROWSER",
		APITest:         "SCRIPT_API",
	}

	// MonitorStatus specifies the possible Synthetics monitor status types.
	MonitorStatus = struct {
		Enabled  MonitorStatusType
		Muted    MonitorStatusType
		Disabled MonitorStatusType
	}{
		Enabled:  "ENABLED",
		Muted:    "MUTED",
		Disabled: "DISABLED",
	}
)

// MonitorOptions represents the options for a New Relic Synthetics monitor.
type MonitorOptions struct {
	ValidationString       string `json:"validationString,omitempty"`
	VerifySSL              bool   `json:"verifySSL,omitempty"`
	BypassHEADRequest      bool   `json:"bypassHEADRequest,omitempty"`
	TreatRedirectAsFailure bool   `json:"treatRedirectAsFailure,omitempty"`
}

// Monitor represents a New Relic Synthetics monitor.
type Monitor struct {
	ID           string            `json:"id,omitempty"`
	Name         string            `json:"name"`
	Type         MonitorType       `json:"type"`
	Frequency    uint              `json:"frequency"`
	URI          string            `json:"uri"`
	Locations    []string          `json:"locations"`
	Status       MonitorStatusType `json:"status"`
	SLAThreshold float64           `json:"slaThreshold"`
	UserID       uint              `json:"userId,omitempty"`
	APIVersion   string            `json:"apiVersion,omitempty"`
	ModifiedAt   *Time             `json:"modifiedAt,omitempty"`
	CreatedAt    *Time             `json:"createdAt,omitempty"`
	Options      MonitorOptions    `json:"options,omitempty"`
}

// MonitorLabel represents a single label for a New Relic Synthetics monitor.
type MonitorLabel struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Href  string `json:"href"`
}

// SecureCredential represents a Synthetics secure credential.
type SecureCredential struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Value       string `json:"value"`
	CreatedAt   *Time  `json:"createdAt"`
	LastUpdated *Time  `json:"lastUpdated"`
}
