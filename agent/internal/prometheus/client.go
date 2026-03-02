package prometheus
import (
	"context"; "encoding/json"; "fmt"; "net/http"; "net/url"; "time"
)
type Client struct { baseURL, token string; http *http.Client }
type QueryResult struct { Metric map[string]string; Value float64 }
func NewClient(baseURL, token string) *Client {
	return &Client{baseURL: baseURL, token: token, http: &http.Client{Timeout: 30 * time.Second}}
}
func (c *Client) Query(ctx context.Context, query string) ([]QueryResult, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s/api/v1/query?%s", c.baseURL, url.Values{"query": {query}}.Encode()), nil)
	if c.token != "" { req.Header.Set("Authorization", "Bearer "+c.token) }
	resp, err := c.http.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	var r struct { Data struct { Result []struct {
		Metric map[string]string; Value []interface{}
	} } }
	json.NewDecoder(resp.Body).Decode(&r)
	var out []QueryResult
	for _, item := range r.Data.Result {
		var v float64
		if len(item.Value) == 2 { fmt.Sscanf(fmt.Sprint(item.Value[1]), "%f", &v) }
		out = append(out, QueryResult{Metric: item.Metric, Value: v})
	}
	return out, nil
}
