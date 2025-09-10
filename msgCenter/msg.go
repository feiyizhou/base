package msgcenter

type Msg struct {
	MsgID     string `json:"msgId"`
	MsgType   string `json:"msgType"`
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

func (m Msg) ToMap() map[string]any {
	return map[string]any{
		"msgId":     m.MsgID,
		"msgType":   m.MsgType,
		"source":    m.Source,
		"timestamp": m.Timestamp,
	}
}

const (
	TopicAgentInstall   = "agent_install"
	TopicAgentUninstall = "agent_uninstall"

	GroupAgentInstall   = "agent_install_group"
	GroupAgentUninstall = "agent_uninstall_group"

	ProducerAgent  = "producer_agent"
	ProducerScript = "producer_script"
)
