package grpc_client

import "testing"

func TestResolveTarget(t *testing.T) {
	tests := []struct {
		name       string
		service    string
		address    string
		targetName string
		want       string
	}{
		{name: "service discovery", service: "greeter", targetName: "greeter", want: "greeter"},
		{name: "direct address", address: "127.0.0.1:50051", targetName: "greeter", want: "127.0.0.1:50051"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := resolveTarget(test.service, test.address, test.targetName)
			if err != nil {
				t.Fatalf("resolveTarget() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("resolveTarget() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestResolveTargetRequiresConfig(t *testing.T) {
	if _, err := resolveTarget("", "", "calculator"); err == nil {
		t.Fatal("resolveTarget() error = nil, want an error")
	}
}
