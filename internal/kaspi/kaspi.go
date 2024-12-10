package kaspi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	kaspiBaseUrl = "https://kaspi.kz/shop/api/v2"
)

var (
	ErrKaspiRequestFail = errors.New("error request status code is not 2xx")
)

func do(ctx context.Context, token, method, link string, parsedRes any) error {
	req, err := http.NewRequestWithContext(ctx, method, kaspiBaseUrl+link, nil)
	if err != nil {
		return fmt.Errorf("creating request %s %s %s \n%w", method, link, token, err)
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Add("X-Auth-Token", token)
	req.Header.Add("Content-Type", "application/vnd.api+json")
	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil || res.StatusCode < 200 || res.StatusCode >= 300 {
		contents, _ := io.ReadAll(res.Body)
		return fmt.Errorf("request fail %s %s %s \n%w \n%s", method, link, token, err, contents)
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(data, parsedRes)
	if err != nil {
		return fmt.Errorf("response json parsing: %w", err)
	}
	return nil
}
