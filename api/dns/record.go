package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RecordsResponseEntry struct {
	Content             string `json:"content"`
	DnsDomainId         int    `json:"dnsDomainId"`
	DnsResourceRecordId string `json:"dnsResourceRecordId"`
	Name                string `json:"name"`
	Ttl                 int    `json:"ttl"`
	Type                string `json:"type"`
}

type RecordsResponse struct {
	DnsResourceRecordList []*RecordsResponseEntry `json:"dnsResourceRecordList"`
	ResultSetProperties   struct {
	} `json:"resultSetProperties"`
}

type RecordUpdate struct {
	DnsResourceRecord struct {
		Content string `json:"content"`
		Ttl     int    `json:"ttl"`
	} `json:"dnsResourceRecord"`
}

type RecordUpdateResponse struct {
	ResultSetProperties struct {
	} `json:"resultSetProperties"`
}

type RecordCreateRequest struct {
	DnsResourceRecordList []struct {
		Content     string `json:"content"`
		DnsDomainId int    `json:"dnsDomainId"`
		Type        string `json:"type"`
		Name        string `json:"name"`
		Ttl         int    `json:"ttl"`
	} `json:"dnsResourceRecordList"`
}

func (c *Client) GetRecords(ctx context.Context, domainId int) (*RecordsResponse, error) {
	body, err := c.MakeRequest(ctx, http.MethodGet, fmt.Sprintf("dnsDomains/%d/dnsResourceRecords", domainId), &bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var result RecordsResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateRecord(ctx context.Context, records *RecordCreateRequest) (*RecordUpdateResponse, error){
	data, err := json.Marshal(records)
	if err != nil {
		return nil, err
	}
	body, err := c.MakeRequest(ctx, http.MethodPost, "/dnsResourceRecords", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var result RecordsResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Client) UpdateRecord(ctx context.Context, domainId int, recordId string, content string, ttl int) (*RecordUpdateResponse, error) {
	data, err := json.Marshal(&RecordUpdate{DnsResourceRecord: struct {
		Content string `json:"content"`
		Ttl     int    `json:"ttl"`
	}(struct {
		Content string
		Ttl     int
	}{Content: content, Ttl: ttl})})
	if err != nil {
		return nil, err
	}
	body, err := c.MakeRequest(ctx, http.MethodPut, fmt.Sprintf("dnsResourceRecords/%d/%s", domainId, recordId), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var result RecordUpdateResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteRecord(ctx context.Context, domainId int, recordId string) (*RecordUpdateResponse, error) {
	body, err := c.MakeRequest(ctx, http.MethodDelete, fmt.Sprintf("dnsResourceRecords/%d/%s", domainId, recordId), &bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var result RecordUpdateResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
