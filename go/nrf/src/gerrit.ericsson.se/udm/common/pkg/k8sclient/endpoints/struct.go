package endpoints

// TEndpoints defines the struct of endpoints json data received from k8s API
type TEndpoints struct {
	Subsets []TSubset `json:"subsets"`
}

// TSubset defines the struct of attribute subsets
type TSubset struct {
	Addresses []TAddress `json:"addresses"`
}

// TAddress defines the struct of attribute addresses
type TAddress struct {
	IP string `json:"ip"`
}
