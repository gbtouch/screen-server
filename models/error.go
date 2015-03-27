package models

//Error make server forwarding screen device debug info
type Error struct {
	Device    string `json:"device"`
	Message   string `json:"message"`
	DebugInfo string `json:"debuginfo"`
}
