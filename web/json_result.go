package web

// JSONResult is a json result.
type JSONResult struct {
	StatusCode int
	Response   interface{}
}

// Render renders the result
func (jr *JSONResult) Render(ctx *Ctx) error {
	return WriteJSON(ctx.Response(), ctx.Request(), jr.StatusCode, jr.Response)
}
