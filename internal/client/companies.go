package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"tracker/internal/models"
)

func ListCompanies(limit, offset int) (*models.CompanyListResponse, error) {
	values := url.Values{}
	if limit > 0 {
		values.Set("limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		values.Set("offset", fmt.Sprintf("%d", offset))
	}

	path := "/companies"
	if len(values) > 0 {
		path += "?" + values.Encode()
	}

	data, _, err := doRawRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	trimmed := bytes.TrimLeft(data, " \t\r\n")
	if len(trimmed) == 0 {
		return &models.CompanyListResponse{}, nil
	}

	var resp models.CompanyListResponse

	if trimmed[0] == '[' {
		var companies []models.Company
		if err := json.Unmarshal(data, &companies); err != nil {
			return nil, fmt.Errorf("ошибка парсинга массива компаний: %w", err)
		}
		resp.Companies = companies
		resp.Total = len(companies)
		resp.Limit = limit
		resp.Offset = offset
	} else if trimmed[0] == '{' {
		if err := json.Unmarshal(data, &resp); err != nil {
			return nil, fmt.Errorf("ошибка парсинга структуры компаний: %w", err)
		}
		if resp.Total == 0 {
			resp.Total = len(resp.Companies)
		}
		if resp.Limit == 0 && limit > 0 {
			resp.Limit = limit
		}
		if resp.Offset == 0 && offset > 0 {
			resp.Offset = offset
		}
	} else {
		return nil, fmt.Errorf("неожиданный формат ответа сервера")
	}

	return &resp, nil
}

func CreateCompany(name, description string) (*models.Company, error) {
	payload := map[string]string{"name": name}
	if description != "" {
		payload["description"] = description
	}
	var resp models.Company
	err := doRequest("POST", "/companies", payload, &resp)
	return &resp, err
}

func DeleteCompany(name string) error {
	return doRequest("DELETE", fmt.Sprintf("/companies/%s", url.PathEscape(name)), nil, nil)
}
