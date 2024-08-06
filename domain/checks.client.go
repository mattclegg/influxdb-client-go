package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/oapi-codegen/runtime"
	"io"
	"net/http"
	"net/url"
)

var typeToCheck = map[string]func() Check{
	"deadman":   func() Check { return &DeadmanCheck{} },
	"threshold": func() Check { return &ThresholdCheck{} },
	"custom":    func() Check { return &CustomCheck{} },
}

// UnmarshalJSON will convert
func unmarshalCheckJSON(b []byte) (Check, error) {
	var raw struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		m := "unable to detect the check type from json"
		e := &Error{
			Code:    ErrorCodeInvalid,
			Message: &m,
		}
		return nil, e.Error()
	}
	factoryFunc, ok := typeToCheck[raw.Type]
	if !ok {
		return nil, fmt.Errorf("invalid check type %s", raw.Type)
	}
	check := factoryFunc()
	err := json.Unmarshal(b, check)
	return check, err
}

// GetChecks calls the GET on /checks
// List all checks
func (c *Client) GetChecks(ctx context.Context, params *GetChecksParams) (*Checks, error) {
	var err error

	serverURL, err := url.Parse(c.APIEndpoint)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("./checks")

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	queryValues := queryURL.Query()

	if params.Offset != nil {

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "offset", runtime.ParamLocationQuery, *params.Offset); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

	}

	if params.Limit != nil {

		if queryFrag, err := runtime.StyleParamWithLocation("form", true, "limit", runtime.ParamLocationQuery, *params.Limit); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

	}

	if queryFrag, err := runtime.StyleParamWithLocation("form", true, "orgID", runtime.ParamLocationQuery, params.OrgID); err != nil {
		return nil, err
	} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
		return nil, err
	} else {
		for k, v := range parsed {
			for _, v2 := range v {
				queryValues.Add(k, v2)
			}
		}
	}

	queryURL.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params.ZapTraceSpan != nil {
		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "Zap-Trace-Span", runtime.ParamLocationHeader, *params.ZapTraceSpan)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Zap-Trace-Span", headerParam0)
	}

	req = req.WithContext(ctx)
	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(rsp.Body)

	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	type raw struct {
		Checks *[]json.RawMessage `json:"checks,omitempty"`

		// URI pointers for additional paged results.
		Links *Links `json:"links,omitempty"`
	}
	response := &Checks{}

	switch rsp.StatusCode {
	case 200:
		var a raw
		if err := unmarshalJSONResponse(bodyBytes, &a); err != nil {
			return nil, err
		}
		if a.Checks != nil && len(*a.Checks) > 0 {
			c := make([]Check, len(*a.Checks))
			response.Checks = &c
			for i, m := range *a.Checks {
				check, err := unmarshalCheckJSON(m)
				if err != nil {
					return nil, err
				}
				(*response.Checks)[i] = check
			}
		}
		response.Links = a.Links
	default:
		return nil, decodeError(bodyBytes, rsp)
	}
	return response, nil

}

// CreateCheck calls the POST on /checks
// Add new check
func (c *Client) CreateCheck(ctx context.Context, params *CreateCheckAllParams) (Check, error) {
	var err error
	var bodyReader io.Reader
	buf, err := json.Marshal(params.Body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)

	serverURL, err := url.Parse(c.APIEndpoint)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("./checks")

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	req = req.WithContext(ctx)
	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(rsp.Body)

	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	switch rsp.StatusCode {
	case 201:
		check, err := unmarshalCheckJSON(bodyBytes)
		if err != nil {
			return nil, err
		}
		return check, nil
	default:
		return nil, decodeError(bodyBytes, rsp)
	}
}

// DeleteChecksID calls the DELETE on /checks/{checkID}
// Delete a check
func (c *Client) DeleteChecksID(ctx context.Context, params *DeleteChecksIDAllParams) error {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "checkID", runtime.ParamLocationPath, params.CheckID)
	if err != nil {
		return err
	}

	serverURL, err := url.Parse(c.APIEndpoint)
	if err != nil {
		return err
	}

	operationPath := fmt.Sprintf("./checks/%s", pathParam0)

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return err
	}

	if params.ZapTraceSpan != nil {
		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "Zap-Trace-Span", runtime.ParamLocationHeader, *params.ZapTraceSpan)
		if err != nil {
			return err
		}

		req.Header.Set("Zap-Trace-Span", headerParam0)
	}

	req = req.WithContext(ctx)
	rsp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = rsp.Body.Close() }()

	if rsp.StatusCode > 299 {
		bodyBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		return decodeError(bodyBytes, rsp)
	}
	return nil

}

// GetChecksID calls the GET on /checks/{checkID}
// Retrieve a check
func (c *Client) GetChecksID(ctx context.Context, params *GetChecksIDAllParams) (Check, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "checkID", runtime.ParamLocationPath, params.CheckID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(c.APIEndpoint)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("./checks/%s", pathParam0)

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params.ZapTraceSpan != nil {
		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "Zap-Trace-Span", runtime.ParamLocationHeader, *params.ZapTraceSpan)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Zap-Trace-Span", headerParam0)
	}

	req = req.WithContext(ctx)
	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(rsp.Body)

	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	switch rsp.StatusCode {
	case 200:
		check, err := unmarshalCheckJSON(bodyBytes)
		if err != nil {
			return nil, err
		}
		return check, nil
	default:
		return nil, decodeError(bodyBytes, rsp)
	}
}

// PatchChecksID calls the PATCH on /checks/{checkID}
// Update a check
func (c *Client) PatchChecksID(ctx context.Context, params *PatchChecksIDAllParams) (Check, error) {
	var err error
	var bodyReader io.Reader
	buf, err := json.Marshal(params.Body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "checkID", runtime.ParamLocationPath, params.CheckID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(c.APIEndpoint)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("./checks/%s", pathParam0)

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if params.ZapTraceSpan != nil {
		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "Zap-Trace-Span", runtime.ParamLocationHeader, *params.ZapTraceSpan)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Zap-Trace-Span", headerParam0)
	}

	req = req.WithContext(ctx)
	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(rsp.Body)

	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	switch rsp.StatusCode {
	case 200:
		check, err := unmarshalCheckJSON(bodyBytes)
		if err != nil {
			return nil, err
		}
		return check, nil
	default:
		return nil, decodeError(bodyBytes, rsp)
	}
}

// PutChecksID calls the PUT on /checks/{checkID}
// Update a check
func (c *Client) PutChecksID(ctx context.Context, params *PutChecksIDAllParams) (Check, error) {
	var err error
	var bodyReader io.Reader
	buf, err := json.Marshal(params.Body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "checkID", runtime.ParamLocationPath, params.CheckID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(c.APIEndpoint)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("./checks/%s", pathParam0)

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if params.ZapTraceSpan != nil {
		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "Zap-Trace-Span", runtime.ParamLocationHeader, *params.ZapTraceSpan)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Zap-Trace-Span", headerParam0)
	}

	req = req.WithContext(ctx)
	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(rsp.Body)

	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	switch rsp.StatusCode {
	case 200:
		check, err := unmarshalCheckJSON(bodyBytes)
		if err != nil {
			return nil, err
		}
		return check, nil
	default:
		return nil, decodeError(bodyBytes, rsp)
	}
}

// GetChecksIDLabels calls the GET on /checks/{checkID}/labels
// List all labels for a check
func (c *Client) GetChecksIDLabels(ctx context.Context, params *GetChecksIDLabelsAllParams) (*LabelsResponse, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "checkID", runtime.ParamLocationPath, params.CheckID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(c.APIEndpoint)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("./checks/%s/labels", pathParam0)

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params.ZapTraceSpan != nil {
		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "Zap-Trace-Span", runtime.ParamLocationHeader, *params.ZapTraceSpan)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Zap-Trace-Span", headerParam0)
	}

	req = req.WithContext(ctx)
	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(rsp.Body)

	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &LabelsResponse{}

	switch rsp.StatusCode {
	case 200:
		if err := unmarshalJSONResponse(bodyBytes, &response); err != nil {
			return nil, err
		}
	default:
		return nil, decodeError(bodyBytes, rsp)
	}
	return response, nil

}

// PostChecksIDLabels calls the POST on /checks/{checkID}/labels
// Add a label to a check
func (c *Client) PostChecksIDLabels(ctx context.Context, params *PostChecksIDLabelsAllParams) (*LabelResponse, error) {
	var err error
	var bodyReader io.Reader
	buf, err := json.Marshal(params.Body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "checkID", runtime.ParamLocationPath, params.CheckID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(c.APIEndpoint)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("./checks/%s/labels", pathParam0)

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if params.ZapTraceSpan != nil {
		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "Zap-Trace-Span", runtime.ParamLocationHeader, *params.ZapTraceSpan)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Zap-Trace-Span", headerParam0)
	}

	req = req.WithContext(ctx)
	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(rsp.Body)

	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &LabelResponse{}

	switch rsp.StatusCode {
	case 201:
		if err := unmarshalJSONResponse(bodyBytes, &response); err != nil {
			return nil, err
		}
	default:
		return nil, decodeError(bodyBytes, rsp)
	}
	return response, nil

}

// DeleteChecksIDLabelsID calls the DELETE on /checks/{checkID}/labels/{labelID}
// Delete label from a check
func (c *Client) DeleteChecksIDLabelsID(ctx context.Context, params *DeleteChecksIDLabelsIDAllParams) error {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "checkID", runtime.ParamLocationPath, params.CheckID)
	if err != nil {
		return err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "labelID", runtime.ParamLocationPath, params.LabelID)
	if err != nil {
		return err
	}

	serverURL, err := url.Parse(c.APIEndpoint)
	if err != nil {
		return err
	}

	operationPath := fmt.Sprintf("./checks/%s/labels/%s", pathParam0, pathParam1)

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return err
	}

	if params.ZapTraceSpan != nil {
		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "Zap-Trace-Span", runtime.ParamLocationHeader, *params.ZapTraceSpan)
		if err != nil {
			return err
		}

		req.Header.Set("Zap-Trace-Span", headerParam0)
	}

	req = req.WithContext(ctx)
	rsp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = rsp.Body.Close() }()

	if rsp.StatusCode > 299 {
		bodyBytes, err := io.ReadAll(rsp.Body)
		if err != nil {
			return err
		}
		return decodeError(bodyBytes, rsp)
	}
	return nil

}

// GetChecksIDQuery calls the GET on /checks/{checkID}/query
// Retrieve a check query
func (c *Client) GetChecksIDQuery(ctx context.Context, params *GetChecksIDQueryAllParams) (*FluxResponse, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "checkID", runtime.ParamLocationPath, params.CheckID)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(c.APIEndpoint)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("./checks/%s/query", pathParam0)

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if params.ZapTraceSpan != nil {
		var headerParam0 string

		headerParam0, err = runtime.StyleParamWithLocation("simple", false, "Zap-Trace-Span", runtime.ParamLocationHeader, *params.ZapTraceSpan)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Zap-Trace-Span", headerParam0)
	}

	req = req.WithContext(ctx)
	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(rsp.Body)

	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &FluxResponse{}

	switch rsp.StatusCode {
	case 200:
		if err := unmarshalJSONResponse(bodyBytes, &response); err != nil {
			return nil, err
		}
	default:
		return nil, decodeError(bodyBytes, rsp)
	}
	return response, nil

}
