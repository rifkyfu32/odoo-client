package web

type JsonRpcResponse struct {
	JsonRpc string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Result  any           `json:"result"` // Data sukses
	Error   *JsonRpcError `json:"error"`  // Data error
}
