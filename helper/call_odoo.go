package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rifkyfu32/odoo-client/model/web"
)

func CallOdoo(url string, service string, method string, args []any) (any, error) {
	reqBody := web.JsonRpcRequest{
		JsonRpc: "2.0",
		Method:  "call",
		ID:      1, // ID bisa apa saja, untuk mencocokkan request & response
		Params: map[string]any{
			"service": service, // 'common' untuk login, 'object' untuk data
			"method":  method,  // Metode Odoo yang ingin dipanggil
			"args":    args,    // Argumen untuk metode Odoo
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("gagal marshal request JSON: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal mengirim request ke Odoo: %w", err)
	}

	defer func() {
		err := resp.Body.Close()
		PanicIfError(err)
	}()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca response body: %w", err)
	}

	var rpcResponse web.JsonRpcResponse
	if err := json.Unmarshal(respBody, &rpcResponse); err != nil {
		return nil, fmt.Errorf("gagal unmarshal response JSON: %w", err)
	}

	if rpcResponse.Error != nil {
		return nil, fmt.Errorf("error dari Odoo (Code: %d): %s", rpcResponse.Error.Code, rpcResponse.Error.Message)
	}

	return rpcResponse.Result, nil
}
