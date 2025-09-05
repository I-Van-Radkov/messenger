package websocket

const (
	actionSendMessage   = "send_message"
	actionNewMessage    = "new_message"
	actionMessageStatus = "message_status"
	actionError         = "error_message"

	statusSent      = "sent"
	statusDelivered = "delivered"
)

type IncomingMessage struct {
	Action       string `json:"action"` // "send_message"
	ID           int64  `json:"id"`
	DialogID     int64  `json:"dialog_id"`
	RecipientID  int64  `json:"recipient_id"`
	Text         string `json:"text"`
	IsReplyToMsg bool   `json:"is_reply_to_msg"`
	ReplyToMsgID int64  `json:"reply_to_msg_id,omitempty"`
}

type OutgoingMessage struct {
	Action    string `json:"action"` // "new_message"
	MessageID int64  `json:"message_id"`
	Status    string `json:"status"`
}

type ErrorMessage struct {
	Action  string `json:"action"` // "error_message"
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}
