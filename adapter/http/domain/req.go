package domain

type Req struct {
	TicketID string      `json:"ticket_id"`
	Msg      string      `json:"msg"`
	Code     int         `json:"code"`
	Data     interface{} `json:"data"`
}

type Rsp struct {
	TicketID string      `json:"ticket_id"`
	Msg      string      `json:"msg"`
	Code     int         `json:"code"`
	Data     interface{} `json:"data"`
}

func NewRsp(ticketID string, code int, msg string, data interface{}) *Rsp {
	return &Rsp{
		TicketID: ticketID,
		Msg:      msg,
		Code:     code,
		Data:     data,
	}
}
