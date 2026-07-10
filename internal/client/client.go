package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"tracker/internal/config"
	"tracker/internal/service"
)

var httpClient = &http.Client{Timeout: service.HTTPTimeout}

func isRetryableError(err error, statusCode int) bool {
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return true
		}
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			if urlErr.Timeout() {
				return true
			}
			var opErr *net.OpError
			if errors.As(urlErr.Err, &opErr) {
				return true
			}
		}
		return false
	}

	switch statusCode {
	case http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusRequestTimeout,
		http.StatusTooManyRequests:
		return true
	default:
		return false
	}
}

func parseRetryAfter(value string) time.Duration {
	if value == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		if seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
		return 0
	}
	if t, err := http.ParseTime(value); err == nil {
		delay := time.Until(t)
		if delay > 0 {
			return delay
		}
	}
	return 0
}

func doRequestWithCtx(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	apiURL := config.GetAPIURL()
	fullURL := apiURL + path

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("ошибка сериализации: %w", err)
		}
	}

	var lastErr error
	backoff := service.InitialBackoff

	for attempt := 0; attempt <= service.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 2
			if backoff > service.MaxBackoff {
				backoff = service.MaxBackoff
			}
		}

		var reqBody io.Reader
		if bodyBytes != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
		if err != nil {
			return err
		}

		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		token := config.LoadToken()
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("ошибка сети: %w", err)
			if isRetryableError(err, 0) && attempt < service.MaxRetries {
				continue
			}
			return lastErr
		}

		if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close()
			return fmt.Errorf("сессия истекла. Выполните: tracker login")
		}
		if resp.StatusCode == http.StatusForbidden {
			resp.Body.Close()
			return fmt.Errorf("доступ запрещён")
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
			resp.Body.Close()

			if attempt < service.MaxRetries {
				if retryAfter > 0 {
					backoff = retryAfter
				}
				lastErr = fmt.Errorf("превышен rate limit (429)")
				continue
			}
			return fmt.Errorf("превышен rate limit после %d попыток", service.MaxRetries+1)
		}

		if resp.StatusCode >= 400 {
			bodyBytesResp, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if isRetryableError(nil, resp.StatusCode) && attempt < service.MaxRetries {
				lastErr = fmt.Errorf("ошибка %d: %s", resp.StatusCode, string(bodyBytesResp))
				continue
			}

			var errResp struct {
				Detail string `json:"detail"`
			}
			if err := json.Unmarshal(bodyBytesResp, &errResp); err == nil && errResp.Detail != "" {
				return fmt.Errorf("ошибка %d: %s", resp.StatusCode, errResp.Detail)
			}
			return fmt.Errorf("ошибка %d: %s", resp.StatusCode, string(bodyBytesResp))
		}

		if result != nil {
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				resp.Body.Close()
				return err
			}
		}
		resp.Body.Close()
		return nil
	}

	return fmt.Errorf("превышено число попыток: %w", lastErr)
}

func doRequest(method, path string, body interface{}, result interface{}) error {
	return doRequestWithCtx(context.Background(), method, path, body, result)
}

func doRawRequestWithCtx(ctx context.Context, method, path string, body interface{}) ([]byte, http.Header, error) {
	apiURL := config.GetAPIURL()
	fullURL := apiURL + path

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка сериализации: %w", err)
		}
	}

	var lastErr error
	backoff := service.InitialBackoff

	for attempt := 0; attempt <= service.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, nil, ctx.Err()
			case <-time.After(backoff):
			}
			backoff *= 2
			if backoff > service.MaxBackoff {
				backoff = service.MaxBackoff
			}
		}

		var reqBody io.Reader
		if bodyBytes != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
		if err != nil {
			return nil, nil, err
		}

		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		token := config.LoadToken()
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("ошибка сети: %w", err)
			if isRetryableError(err, 0) && attempt < service.MaxRetries {
				continue
			}
			return nil, nil, lastErr
		}

		if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close()
			return nil, nil, fmt.Errorf("сессия истекла. Выполните: tracker login")
		}
		if resp.StatusCode == http.StatusForbidden {
			resp.Body.Close()
			return nil, nil, fmt.Errorf("доступ запрещён")
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
			resp.Body.Close()

			if attempt < service.MaxRetries {
				if retryAfter > 0 {
					backoff = retryAfter
				}
				lastErr = fmt.Errorf("превышен rate limit (429)")
				continue
			}
			return nil, nil, fmt.Errorf("превышен rate limit после %d попыток", service.MaxRetries+1)
		}

		if resp.StatusCode >= 400 {
			bodyBytesResp, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if isRetryableError(nil, resp.StatusCode) && attempt < service.MaxRetries {
				lastErr = fmt.Errorf("ошибка %d: %s", resp.StatusCode, string(bodyBytesResp))
				continue
			}

			var errResp struct {
				Detail string `json:"detail"`
			}
			if err := json.Unmarshal(bodyBytesResp, &errResp); err == nil && errResp.Detail != "" {
				return nil, nil, fmt.Errorf("ошибка %d: %s", resp.StatusCode, errResp.Detail)
			}
			return nil, nil, fmt.Errorf("ошибка %d: %s", resp.StatusCode, string(bodyBytesResp))
		}

		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, nil, err
		}

		return data, resp.Header, nil
	}

	return nil, nil, fmt.Errorf("превышено число попыток: %w", lastErr)
}

func doRawRequest(method, path string, body interface{}) ([]byte, http.Header, error) {
	return doRawRequestWithCtx(context.Background(), method, path, body)
}
