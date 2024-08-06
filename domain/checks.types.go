package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

// Defines values for CheckBaseLastRunStatus.
const (
	CheckBaseLastRunStatusCanceled CheckBaseLastRunStatus = "canceled"

	CheckBaseLastRunStatusFailed CheckBaseLastRunStatus = "failed"

	CheckBaseLastRunStatusSuccess CheckBaseLastRunStatus = "success"
)

// Defines values for CheckPatchStatus.
const (
	CheckPatchStatusActive CheckPatchStatus = "active"

	CheckPatchStatusInactive CheckPatchStatus = "inactive"
)

// Defines values for CheckStatusLevel.
const (
	CheckStatusLevelCRIT CheckStatusLevel = "CRIT"

	CheckStatusLevelINFO CheckStatusLevel = "INFO"

	CheckStatusLevelOK CheckStatusLevel = "OK"

	CheckStatusLevelUNKNOWN CheckStatusLevel = "UNKNOWN"

	CheckStatusLevelWARN CheckStatusLevel = "WARN"
)

// Defines values for DeadmanCheckType.
const (
	DeadmanCheckTypeDeadman DeadmanCheckType = "deadman"
)

// Defines values for ThresholdCheckType.
const (
	ThresholdCheckTypeThreshold ThresholdCheckType = "threshold"
)

// Defines values for CustomCheckType.
const (
	CustomCheckTypeCustom CustomCheckType = "custom"
)

// CheckPatch defines model for CheckPatch.
type CheckPatch struct {
	Description *string           `json:"description,omitempty"`
	Name        *string           `json:"name,omitempty"`
	Status      *CheckPatchStatus `json:"status,omitempty"`
}

// CheckPatchStatus defines model for CheckPatch.Status.
type CheckPatchStatus string

// The state to record if check matches a criteria.
type CheckStatusLevel string

// CustomCheckType defines model for CustomCheck.Type.
type CustomCheckType string

// DeadmanCheckType defines model for DeadmanCheck.Type.
type DeadmanCheckType string

// ThresholdCheckType defines model for ThresholdCheck.Type.
type ThresholdCheckType string

// CreateCheckJSONRequestBody defines body for CreateCheck for application/json ContentType.
type CreateCheckJSONRequestBody CreateCheckJSONBody

// PatchChecksIDJSONRequestBody defines body for PatchChecksID for application/json ContentType.
type PatchChecksIDJSONRequestBody PatchChecksIDJSONBody

// PutChecksIDJSONRequestBody defines body for PutChecksID for application/json ContentType.
type PutChecksIDJSONRequestBody PutChecksIDJSONBody

// PostChecksIDLabelsJSONRequestBody defines body for PostChecksIDLabels for application/json ContentType.
type PostChecksIDLabelsJSONRequestBody PostChecksIDLabelsJSONBody

// Check defines model for Check.
type Check interface {
	Type() string
}

// Checks defines model for Checks.
type Checks struct {
	Checks *[]Check `json:"checks,omitempty"`

	// URI pointers for additional paged results.
	Links *Links `json:"links,omitempty"`
}

// CheckBase defines model for CheckBase.
type CheckBase struct {
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// An optional description of the check.
	Description   *string                 `json:"description,omitempty"`
	Id            *string                 `json:"id,omitempty"`
	Labels        *Labels                 `json:"labels,omitempty"`
	LastRunError  *string                 `json:"lastRunError,omitempty"`
	LastRunStatus *CheckBaseLastRunStatus `json:"lastRunStatus,omitempty"`

	// A timestamp ([RFC3339 date/time format](https://docs.influxdata.com/influxdb/v2.3/reference/glossary/#rfc3339-timestamp)) of the latest scheduled and completed run.
	LatestCompleted *time.Time `json:"latestCompleted,omitempty"`
	Links           *struct {
		// URI of resource.
		Labels *Link `json:"labels,omitempty"`

		// URI of resource.
		Members *Link `json:"members,omitempty"`

		// URI of resource.
		Owners *Link `json:"owners,omitempty"`

		// URI of resource.
		Query *Link `json:"query,omitempty"`

		// URI of resource.
		Self *Link `json:"self,omitempty"`
	} `json:"links,omitempty"`
	Name string `json:"name"`

	// The ID of the organization that owns this check.
	OrgID string `json:"orgID"`

	// The ID of creator used to create this check.
	OwnerID *string        `json:"ownerID,omitempty"`
	Query   DashboardQuery `json:"query"`

	// `inactive` cancels scheduled runs and prevents manual runs of the task.
	Status TaskStatusType `json:"status"`

	// The ID of the task associated with this check.
	TaskID    *string    `json:"taskID,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type CheckBaseExtend struct {
	CheckBase
	// Embedded fields due to inline allOf schema
	// Check repetition interval.
	Every *string `json:"every,omitempty"`

	// Duration to delay after the schedule, before executing check.
	Offset *string `json:"offset,omitempty"`

	// The template used to generate and write a status message.
	StatusMessageTemplate *string `json:"statusMessageTemplate,omitempty"`

	// List of tags to write to each status.
	Tags *[]struct {
		Key   *string `json:"key,omitempty"`
		Value *string `json:"value,omitempty"`
	} `json:"tags,omitempty"`
}

// DeadmanCheck defines model for DeadmanCheck.
type DeadmanCheck struct {
	CheckBaseExtend
	// The state to record if check matches a criteria.
	Level *CheckStatusLevel `json:"level,omitempty"`

	// If only zero values reported since time, trigger an alert
	ReportZero *bool `json:"reportZero,omitempty"`

	// String duration for time that a series is considered stale and should not trigger deadman.
	StaleTime *string `json:"staleTime,omitempty"`

	// String duration before deadman triggers.
	TimeSince *string `json:"timeSince,omitempty"`
}

func (d DeadmanCheck) Type() string {
	return string(DeadmanCheckTypeDeadman)
}

// MarshalJSON implement json.Marshaler interface.
func (d DeadmanCheck) MarshalJSON() ([]byte, error) {
	type deadmanCheckAlias DeadmanCheck
	return json.Marshal(
		struct {
			deadmanCheckAlias
			Type string `json:"type"`
		}{
			deadmanCheckAlias: deadmanCheckAlias(d),
			Type:              d.Type(),
		})
}

// ThresholdCheck defines model for ThresholdCheck.
type ThresholdCheck struct {
	// Embedded struct due to allOf(#/components/schemas/CheckBase)
	CheckBaseExtend
	Thresholds *[]Threshold `json:"thresholds,omitempty"`
}

// MarshalJSON implement json.Marshaler interface.
func (t *ThresholdCheck) MarshalJSON() ([]byte, error) {
	type thresholdCheckAlias ThresholdCheck
	return json.Marshal(
		struct {
			thresholdCheckAlias
			Type string `json:"type"`
		}{
			thresholdCheckAlias: thresholdCheckAlias(*t),
			Type:                t.Type(),
		})
}

type thresholdCheckDecode struct {
	CheckBaseExtend
	Thresholds []thresholdDecode `json:"thresholds"`
}

type thresholdDecode struct {
	ThresholdBase
	Type   string  `json:"type"`
	Value  float32 `json:"value"`
	Min    float32 `json:"min"`
	Max    float32 `json:"max"`
	Within bool    `json:"within"`
}

// UnmarshalJSON implement json.Unmarshaler interface.
func (t *ThresholdCheck) UnmarshalJSON(b []byte) error {
	var tdRaws thresholdCheckDecode
	if err := json.Unmarshal(b, &tdRaws); err != nil {
		return err
	}
	t.CheckBaseExtend = tdRaws.CheckBaseExtend
	a := make([]Threshold, 0, len(tdRaws.Thresholds))
	t.Thresholds = &a
	for _, tdRaw := range tdRaws.Thresholds {
		switch tdRaw.Type {
		case "lesser":
			td := &LesserThreshold{
				ThresholdBase: tdRaw.ThresholdBase,
				Value:         tdRaw.Value,
			}
			*t.Thresholds = append(*t.Thresholds, td)
		case "greater":
			td := &GreaterThreshold{
				ThresholdBase: tdRaw.ThresholdBase,
				Value:         tdRaw.Value,
			}
			*t.Thresholds = append(*t.Thresholds, td)
		case "range":
			td := &RangeThreshold{
				ThresholdBase: tdRaw.ThresholdBase,
				Min:           tdRaw.Min,
				Max:           tdRaw.Max,
				Within:        tdRaw.Within,
			}
			*t.Thresholds = append(*t.Thresholds, td)
		default:
			return fmt.Errorf("invalid threshold type %s", tdRaw.Type)
		}
	}
	return nil
}

func (t ThresholdCheck) Type() string {
	return string(ThresholdCheckTypeThreshold)
}

// Threshold defines model for Threshold.
type Threshold interface {
	Type() string
}

// ThresholdBase defines model for ThresholdBase.
type ThresholdBase struct {
	// If true, only alert if all values meet threshold.
	AllValues *bool `json:"allValues,omitempty"`

	// The state to record if check matches a criteria.
	Level *CheckStatusLevel `json:"level,omitempty"`
}

// LesserThreshold defines model for LesserThreshold.
type LesserThreshold struct {
	// Embedded struct due to allOf(#/components/schemas/ThresholdBase)
	ThresholdBase
	// Embedded fields due to inline allOf schema
	Typ   LesserThresholdType `json:"type"`
	Value float32             `json:"value"`
}

func (t LesserThreshold) Type() string {
	return string(LesserThresholdTypeLesser)
}

// MarshalJSON implement json.Marshaler interface.
func (t LesserThreshold) MarshalJSON() ([]byte, error) {
	type lesserThresholdAlias LesserThreshold
	return json.Marshal(
		struct {
			lesserThresholdAlias
			Type string `json:"type"`
		}{
			lesserThresholdAlias: lesserThresholdAlias(t),
			Type:                 t.Type(),
		})
}

// GreaterThreshold defines model for GreaterThreshold.
type GreaterThreshold struct {
	// Embedded struct due to allOf(#/components/schemas/ThresholdBase)
	ThresholdBase
	// Embedded fields due to inline allOf schema
	Typ   GreaterThresholdType `json:"type"`
	Value float32              `json:"value"`
}

func (t GreaterThreshold) Type() string {
	return string(GreaterThresholdTypeGreater)
}

// MarshalJSON implement json.Marshaler interface.
func (t GreaterThreshold) MarshalJSON() ([]byte, error) {
	type greaterThresholdAlias GreaterThreshold
	return json.Marshal(
		struct {
			greaterThresholdAlias
			Type string `json:"type"`
		}{
			greaterThresholdAlias: greaterThresholdAlias(t),
			Type:                  t.Type(),
		})
}

// RangeThreshold defines model for RangeThreshold.
type RangeThreshold struct {
	// Embedded struct due to allOf(#/components/schemas/ThresholdBase)
	ThresholdBase
	// Embedded fields due to inline allOf schema
	Max    float32            `json:"max"`
	Min    float32            `json:"min"`
	Typ    RangeThresholdType `json:"type"`
	Within bool               `json:"within"`
}

func (r RangeThreshold) Type() string {
	return string(RangeThresholdTypeRange)
}

// MarshalJSON implement json.Marshaler interface.
func (r RangeThreshold) MarshalJSON() ([]byte, error) {
	type rangeThresholdAlias RangeThreshold
	return json.Marshal(
		struct {
			rangeThresholdAlias
			Type string `json:"type"`
		}{
			rangeThresholdAlias: rangeThresholdAlias(r),
			Type:                r.Type(),
		})
}

// CustomCheck defines model for CustomCheck.
type CustomCheck struct {
	// Embedded struct due to allOf(#/components/schemas/CheckBase)
	CheckBase
}

func (c CustomCheck) Type() string {
	return string(CustomCheckTypeCustom)
}

// MarshalJSON implement json.Marshaler interface.
func (c CustomCheck) MarshalJSON() ([]byte, error) {
	type customCheckAlias CustomCheck
	return json.Marshal(
		struct {
			customCheckAlias
			Type string `json:"type"`
		}{
			customCheckAlias: customCheckAlias(c),
			Type:             c.Type(),
		})
}

// GetChecksParams defines parameters for GetChecks.
type GetChecksParams struct {
	// The offset for pagination.
	// The number of records to skip.
	Offset *Offset `json:"offset,omitempty"`

	// Limits the number of records returned. Default is `20`.
	Limit *Limit `json:"limit,omitempty"`

	// Only show checks that belong to a specific organization ID.
	OrgID string `json:"orgID"`

	// OpenTracing span context
	ZapTraceSpan *TraceSpan `json:"Zap-Trace-Span,omitempty"`
}

// CreateCheckJSONBody defines parameters for CreateCheck.
type CreateCheckJSONBody Check

// CreateCheckAllParams defines type for all parameters for CreateCheck.
type CreateCheckAllParams struct {
	Body CreateCheckJSONRequestBody
}

// DeleteChecksIDParams defines parameters for DeleteChecksID.
type DeleteChecksIDParams struct {
	// OpenTracing span context
	ZapTraceSpan *TraceSpan `json:"Zap-Trace-Span,omitempty"`
}

// DeleteChecksIDAllParams defines type for all parameters for DeleteChecksID.
type DeleteChecksIDAllParams struct {
	DeleteChecksIDParams

	CheckID string
}

// GetChecksIDParams defines parameters for GetChecksID.
type GetChecksIDParams struct {
	// OpenTracing span context
	ZapTraceSpan *TraceSpan `json:"Zap-Trace-Span,omitempty"`
}

// GetChecksIDAllParams defines type for all parameters for GetChecksID.
type GetChecksIDAllParams struct {
	GetChecksIDParams

	CheckID string
}

// PatchChecksIDJSONBody defines parameters for PatchChecksID.
type PatchChecksIDJSONBody CheckPatch

// PatchChecksIDParams defines parameters for PatchChecksID.
type PatchChecksIDParams struct {
	// OpenTracing span context
	ZapTraceSpan *TraceSpan `json:"Zap-Trace-Span,omitempty"`
}

// PatchChecksIDAllParams defines type for all parameters for PatchChecksID.
type PatchChecksIDAllParams struct {
	PatchChecksIDParams

	CheckID string

	Body PatchChecksIDJSONRequestBody
}

// PutChecksIDJSONBody defines parameters for PutChecksID.
type PutChecksIDJSONBody Check

// PutChecksIDParams defines parameters for PutChecksID.
type PutChecksIDParams struct {
	// OpenTracing span context
	ZapTraceSpan *TraceSpan `json:"Zap-Trace-Span,omitempty"`
}

// PutChecksIDAllParams defines type for all parameters for PutChecksID.
type PutChecksIDAllParams struct {
	PutChecksIDParams

	CheckID string

	Body PutChecksIDJSONRequestBody
}

// GetChecksIDLabelsParams defines parameters for GetChecksIDLabels.
type GetChecksIDLabelsParams struct {
	// OpenTracing span context
	ZapTraceSpan *TraceSpan `json:"Zap-Trace-Span,omitempty"`
}

// GetChecksIDLabelsAllParams defines type for all parameters for GetChecksIDLabels.
type GetChecksIDLabelsAllParams struct {
	GetChecksIDLabelsParams

	CheckID string
}

// PostChecksIDLabelsJSONBody defines parameters for PostChecksIDLabels.
type PostChecksIDLabelsJSONBody LabelMapping

// PostChecksIDLabelsParams defines parameters for PostChecksIDLabels.
type PostChecksIDLabelsParams struct {
	// OpenTracing span context
	ZapTraceSpan *TraceSpan `json:"Zap-Trace-Span,omitempty"`
}

// PostChecksIDLabelsAllParams defines type for all parameters for PostChecksIDLabels.
type PostChecksIDLabelsAllParams struct {
	PostChecksIDLabelsParams

	CheckID string

	Body PostChecksIDLabelsJSONRequestBody
}

// DeleteChecksIDLabelsIDParams defines parameters for DeleteChecksIDLabelsID.
type DeleteChecksIDLabelsIDParams struct {
	// OpenTracing span context
	ZapTraceSpan *TraceSpan `json:"Zap-Trace-Span,omitempty"`
}

// DeleteChecksIDLabelsIDAllParams defines type for all parameters for DeleteChecksIDLabelsID.
type DeleteChecksIDLabelsIDAllParams struct {
	DeleteChecksIDLabelsIDParams

	CheckID string

	LabelID string
}

// GetChecksIDQueryParams defines parameters for GetChecksIDQuery.
type GetChecksIDQueryParams struct {
	// OpenTracing span context
	ZapTraceSpan *TraceSpan `json:"Zap-Trace-Span,omitempty"`
}

// GetChecksIDQueryAllParams defines type for all parameters for GetChecksIDQuery.
type GetChecksIDQueryAllParams struct {
	GetChecksIDQueryParams

	CheckID string
}
