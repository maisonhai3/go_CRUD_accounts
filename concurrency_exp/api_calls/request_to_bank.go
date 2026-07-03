package api_calls

import (
	"context"
	"fmt"
	"net/http"
)

func CallBank(ctx context.Context, bankURL string) (string, error) {
	// Gắn context chứa tín hiệu cancel/timeout vào trực tiếp Request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bankURL, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bank error: %d", resp.StatusCode)
	}

	return "Thành công từ: " + bankURL, nil
}
