/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	tags "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/aws/tags"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSClusterProviderSpec is the providerConfig for AWS in the cluster
// object
// +k8s:openapi-gen=true
type AWSClusterProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// NetworkSpec encapsulates all things related to AWS network.
	NetworkSpec NetworkSpec `json:"networkSpec,omitempty"`

	// The AWS Region the cluster lives in.
	Region string `json:"region,omitempty"`

	// SSHKeyName is the name of the ssh key to attach to the bastion host.
	SSHKeyName string `json:"sshKeyName,omitempty"`

	// CAKeyPair is the key pair for ca certs.
	CAKeyPair KeyPair `json:"caKeyPair,omitempty"`

	//EtcdCAKeyPair is the key pair for etcd.
	EtcdCAKeyPair KeyPair `json:"etcdCAKeyPair,omitempty"`

	// FrontProxyCAKeyPair is the key pair for FrontProxyKeyPair.
	FrontProxyCAKeyPair KeyPair `json:"frontProxyCAKeyPair,omitempty"`

	// SAKeyPair is the service account key pair.
	SAKeyPair KeyPair `json:"saKeyPair,omitempty"`
}

// KeyPair is how operators can supply custom keypairs for kubeadm to use.
type KeyPair struct {
	// base64 encoded cert and key
	Cert []byte `json:"cert"`
	Key  []byte `json:"key"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

func init() {
	SchemeBuilder.Register(&AWSClusterProviderSpec{})
}

// HasCertAndKey returns whether a keypair contains cert and key of non-zero length.
func (kp *KeyPair) HasCertAndKey() bool {
	return len(kp.Cert) != 0 && len(kp.Key) != 0
}

// NetworkSpec encapsulates all things related to AWS network.
type NetworkSpec struct {
	// VPC configuration.
	// +optional
	VPC VPCSpec `json:"vpc,omitempty"`

	// Subnets configuration.
	// +optional
	Subnets Subnets `json:"subnets,omitempty"`
}

// VPCSpec configures an AWS VPC.
type VPCSpec struct {
	// ID is the vpc-id of the VPC this provider should use to create resources.
	ID string `json:"id,omitempty"`

	// CidrBlock is the CIDR block to be used when the provider creates a managed VPC.
	// Defaults to 10.0.0.0/16.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// InternetGatewayID is the id of the internet gateway associated with the VPC.
	// +optional
	InternetGatewayID *string `json:"internetGatewayId,omitempty"`

	// Tags is a collection of tags describing the resource.
	Tags tags.Map `json:"tags,omitempty"`
}

// String returns a string representation of the VPC.
func (v *VPCSpec) String() string {
	return fmt.Sprintf("id=%s", v.ID)
}

// IsProvided returns true if the VPC is unmanaged.
func (v *VPCSpec) IsProvided() bool {
	return v.ID != "" && !v.Tags.HasManaged()
}

// SubnetSpec configures an AWS Subnet.
type SubnetSpec struct {
	// ID defines a unique identifier to reference this resource.
	ID string `json:"id,omitempty"`

	// CidrBlock is the CIDR block to be used when the provider creates a managed VPC.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// AvailabilityZone defines the availability zone to use for this subnet in the cluster's region.
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// IsPublic defines the subnet as a public subnet. Refer to the AWS documentation for further information.
	// +optional
	IsPublic bool `json:"isPublic,omitempty"`

	// RouteTableID is the routing table id associated with the subnet.
	// +optional
	RouteTableID *string `json:"routeTableId"`

	// NatGatewayID is the NAT gateway id associated with the subnet.
	// Ignored if the subnet is public.
	// +optional
	NatGatewayID *string `json:"natGatewayId"`

	// Tags is a collection of tags describing the resource.
	Tags tags.Map `json:"tags,omitempty"`
}

// String returns a string representation of the subnet.
func (s *SubnetSpec) String() string {
	return fmt.Sprintf("id=%s/az=%s/public=%v", s.ID, s.AvailabilityZone, s.IsPublic)
}
