package capabilities

// BindFlags define the flags used to generate the bundle report
type BindFlags struct {
	IndexImage      string `json:"image"`
	Limit           int32  `json:"limit"`
	HeadOnly        bool   `json:"headOnly"`
	S3Bucket        string `json:"s3Bucket"`
	Filter          string `json:"filter"`
	FilterBundle    string `json:"filterBundle"`
	OutputPath      string `json:"outputPath"`
	OutputFormat    string `json:"outputFormat"`
	ContainerEngine string `json:"containerEngine"`
	Namespace       string `json:"namespace"`
	PullSecretName  string `json:"pullSecretName"`
	ServiceAccount  string `json:"serviceAccount"`
}