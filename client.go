package tapestry

type TapestryClient struct {
	tapestryApiBaseUrl string
	apiKey             string
	execution          TapestryExecutionType
	blockchain         string
}

type TapestryExecutionType string

const (
	FastUnconfirmed TapestryExecutionType = "FAST_UNCONFIRMED"
	QuickSignature  TapestryExecutionType = "QUICK_SIGNATURE"
	ConfirmedParsed TapestryExecutionType = "CONFIRMED_AND_PARSED"
)

func NewTapestryClient(apiKey string, tapestryApiBaseUrl string, execution TapestryExecutionType, blockchain string) TapestryClient {
	return TapestryClient{
		tapestryApiBaseUrl: tapestryApiBaseUrl,
		apiKey:             apiKey,
		execution:          execution,
		blockchain:         blockchain,
	}
}
