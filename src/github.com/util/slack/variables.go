package slack

import (
	"net/http"
)

const (
	host               = "hooks.slack.com"
	port               = "443"
	defaultWebhookPath = "/services/T0BL2BW3G/B3ZUSL677/X9yHXvWoEm3ocqKyqNEPgviH"
	defaultText        = "Droi Slack Webhook"
	defaultChannel     = "bot_home"
)

// WebhookInfo - basic information for slack webhook
// Required fields: Payload and Payload.Channel
type WebhookInfo struct {
	CustomClient *http.Client
	Payload      *Payload
}

// Payload - slack basic formatted message
// See https://api.slack.com/docs/message-formatting
type Payload struct {
	Parse       string       `json:"parse,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	LinkNames   string       `json:"link_names,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment - rich content of slack message
// See https://api.slack.com/docs/message-attachments
type Attachment struct {
	Fallback   string   `json:"fallback"`
	Color      string   `json:"color"`
	PreText    string   `json:"pretext"`
	AuthorName string   `json:"author_name"`
	AuthorLink string   `json:"author_link"`
	AuthorIcon string   `json:"author_icon"`
	Title      string   `json:"title"`
	TitleLink  string   `json:"title_link"`
	Text       string   `json:"text"`
	ImageURL   string   `json:"image_url"`
	ThumbURL   string   `json:"thumb_url"`
	Fields     []*Field `json:"fields"`
	Footer     string   `json:"footer"`
	FooterIcon string   `json:"footer_icon"`
}

// Field - the field of attachment
type Field struct {
	Title string      `json:"title"`
	Value interface{} `json:"value"`
	Short bool        `json:"short"`
}
