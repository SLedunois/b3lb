package balancer

// Mock is a mock implementation of the Balancer interface
type Mock struct{}

var (
	// BalancerMockProcessFunc is the function to be called when Process is called
	BalancerMockProcessFunc func(instances []string) (string, error)
	// BalancerMockClusterStatusFunc is the function to be called when ClusterStatus is called
	BalancerMockClusterStatusFunc func(instances []string) ([]InstanceStatus, error)
)

// Process is a mock implementation of the Process method
func (b *Mock) Process(instances []string) (string, error) {
	return BalancerMockProcessFunc(instances)
}

// ClusterStatus is a mock implementation of the ClusterStatus method
func (b *Mock) ClusterStatus(instances []string) ([]InstanceStatus, error) {
	return BalancerMockClusterStatusFunc(instances)
}
