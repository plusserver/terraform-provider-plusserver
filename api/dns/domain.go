package dns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type DomainSearchList struct {
	NameList []string `json:"nameList"`
}

type SearchDomain struct {
	DnsDomainSearchList []DomainSearchList `json:"dnsDomainSearchList"`
}

type SearchDomainResponse struct {
	DnsDomainList []DomainListResponse `json:"dnsDomainList"`
}

type DomainGetResponse struct {
	DnsDomain struct {
		CompanyId                      string   `json:"companyId"`
		ContractId                     string   `json:"contractId"`
		CreateDateTime                 string   `json:"createDateTime"`
		DnsDomainId                    int      `json:"dnsDomainId"`
		DnsNameserverPairName          string   `json:"dnsNameserverPairName"`
		Name                           string   `json:"name"`
		Protected                      bool     `json:"protected"`
		ReplicationMasterIpAddressList []string `json:"replicationMasterIpAddressList"`
		ReplicationType                string   `json:"replicationType"`
		UnicodeName                    string   `json:"unicodeName"`
	} `json:"dnsDomain"`
	ResultSetProperties struct {
	} `json:"resultSetProperties"`
}

type DomainListResponse struct {
	UnicodeName                    string   `json:"unicodeName"`
	Name                           string   `json:"name"`
	DnsNameserverPairName          string   `json:"dnsNameserverPairName"`
	CompanyId                      string   `json:"companyId"`
	DnsDomainId                    int      `json:"dnsDomainId"`
	Protected                      bool     `json:"protected"`
	ReplicationType                string   `json:"replicationType"`
	ReplicationMasterIpAddressList []string `json:"replicationMasterIpAddressList"`
	CreateDateTime                 string   `json:"createDateTime"`
	ContractId                     string   `json:"contractId"`
}

type DeleteDomainResponse struct {
	ResultSetProperties struct {
	} `json:"resultSetProperties"`
}

type UpdateDomain struct {
	DnsDomain struct {
		ReplicationMasterIpAddressList []string `json:"replicationMasterIpAddressList"`
		Protected                      bool     `json:"protected"`
		CompanyId                      string   `json:"companyId"`
		ContractId                     string   `json:"contractId"`
	} `json:"dnsDomain"`
}

type UpdateDomainResponse struct {
	DnsDomain struct {
		CompanyId                      string        `json:"companyId"`
		ContractId                     string        `json:"contractId"`
		CreateDateTime                 string        `json:"createDateTime"`
		DnsDomainId                    int           `json:"dnsDomainId"`
		DnsNameserverPairName          string        `json:"dnsNameserverPairName"`
		Name                           string        `json:"name"`
		Protected                      bool          `json:"protected"`
		ReplicationMasterIpAddressList []interface{} `json:"replicationMasterIpAddressList"`
		ReplicationType                string        `json:"replicationType"`
		UnicodeName                    string        `json:"unicodeName"`
	} `json:"dnsDomain"`
	ResultSetProperties struct {
	} `json:"resultSetProperties"`
}

type CreateDomain struct {
	DnsDomain struct {
		UnicodeName                    string   `json:"unicodeName"`
		CompanyId                      string   `json:"companyId"`
		DnsNameserverPairName          string   `json:"dnsNameserverPairName"`
		Protected                      bool     `json:"protected"`
		ReplicationType                string   `json:"replicationType"`
		ReplicationMasterIpAddressList []string `json:"replicationMasterIpAddressList"`
		ContractId                     string   `json:"contractId"`
	} `json:"dnsDomain"`
}

type CreateDomainResponse struct {
	DnsDomain struct {
		UnicodeName                    string   `json:"unicodeName"`
		Name                           string   `json:"name"`
		DnsNameserverPairName          string   `json:"dnsNameserverPairName"`
		CompanyId                      string   `json:"companyId"`
		DnsDomainId                    int      `json:"dnsDomainId"`
		Protected                      bool     `json:"protected"`
		ReplicationType                string   `json:"replicationType"`
		ReplicationMasterIpAddressList []string `json:"replicationMasterIpAddressList"`
		CreateDateTime                 string   `json:"createDateTime"`
		ContractId                     string   `json:"contractId"`
	} `json:"dnsDomain"`
}

func debugPrint(body []byte) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		log.Println("[ERROR] JSON parse error: ", err)
		return
	}
	log.Println("[INFO] payload:", string(prettyJSON.Bytes()))

}

func (c *Client) SearchDomains(ctx context.Context, domain string) (*SearchDomainResponse, error) {
	data, err := json.Marshal(&SearchDomain{DnsDomainSearchList: []DomainSearchList{{NameList: []string{domain}}}})
	debugPrint(data)
	if err != nil {
		return nil, err
	}
	body, err := c.MakeRequest(ctx, http.MethodPost, "dnsDomains/search", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var result SearchDomainResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetDomainById(ctx context.Context, domainId string) (*DomainGetResponse, error) {
	body, err := c.MakeRequest(ctx, http.MethodGet, fmt.Sprintf("dnsDomains/%s", domainId), &bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var result DomainGetResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) DeleteDomain(ctx context.Context, domainId string) (*DeleteDomainResponse, error) {
	body, err := c.MakeRequest(ctx, http.MethodDelete, fmt.Sprintf("dnsDomains/%s", domainId), &bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var result DeleteDomainResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) UpdateDomain(ctx context.Context, domainId string, req *UpdateDomain) (*UpdateDomainResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	body, err := c.MakeRequest(ctx, http.MethodPut, fmt.Sprintf("dnsDomains/%s", domainId), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var result UpdateDomainResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateDomain(ctx context.Context, req *CreateDomain) (*CreateDomainResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	body, err := c.MakeRequest(ctx, http.MethodPost, "dnsDomains", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var result CreateDomainResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
