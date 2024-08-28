package kvdbv1

// TPartitioner struct
type TPartitioner struct {
	Name                 string   `json:"name"`
	UploadDirectory      string   `json:"uploadDirectory"`
	UploadEndTimestamp   string   `json:"uploadEndTimestamp"`
	UploadStartTimestamp string   `json:"uploadStartTimestamp"`
	Deployed             bool     `json:"deployed"`
	DeployedOnMembers    []string `json:"deployedOnMembers"`
}

// TProblemDetails struct
type TProblemDetails struct {
	Cause   string `json:"cause,omitempty"`
	Summary string `json:"summary,omitempty"`
}

// TGfshCommand struct
type TGfshCommand struct {
	Output            string `json:"output,omitempty"`
	StartTimestamp    string `json:"startTimestamp,omitempty"`
	StatusCode        int    `json:"statusCode,omitempty"`
	Id                string `json:"id,omitempty"`
	ReceivedTimestamp string `json:"receivedTimestamp,omitempty"`
	Command           string `json:"command"`
	ExecutionStatus   string `json:"executionStatus,omitempty"`
	EndTimestamp      string `json:"endTimestamp,omitempty"`
	Metadata          string `json:"metadata,omitempty"`
}

// TGfshCommandId struct
type TGfshCommandId struct {
	CommandId string `json:"commandId"`
}

// TInitialPartitionersDeploymentStatus struct
type TInitialPartitionersDeploymentStatus struct {
	InitialPartitionersDeploymentDone bool `json:"initialPartitionersDeploymentDone"`
}
