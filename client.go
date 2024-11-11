package tapestry

type TapestryClient struct {
	tapestryApiBaseUrl string
	apiKey             string
	execution          Execution
	blockchain         string
}

type Execution string

const (
	ExecutionFastUnconfirmed Execution = "FAST_UNCONFIRMED"
	ExecutionQuickSignature  Execution = "QUICK_SIGNATURE"
	ExecutionConfirmedParsed Execution = "CONFIRMED_AND_PARSED"
)

func NewTapestryClient(apiKey string, tapestryApiBaseUrl string, execution Execution, blockchain string) TapestryClient {
	return TapestryClient{
		tapestryApiBaseUrl: tapestryApiBaseUrl,
		apiKey:             apiKey,
		execution:          execution,
		blockchain:         blockchain,
	}
}
