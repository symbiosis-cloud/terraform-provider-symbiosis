package symbiosis

import (
	"fmt"
)

type SymbiosisApiError struct {
	Status    int32
	ErrorType string `json:"error"`
	Message   string
	Path      string
}

type NodePool struct {
	Id              string
	ClusterName     string
	NodeTypeName    string
	IsMaster        bool
	DesiredQuantity int
}

type Cluster struct {
	Name      string
	State     string
	NodePools []NodePool
}

type TeamMember struct {
	Email string
	Role  string
}

func (error *SymbiosisApiError) Error() string {
	return fmt.Sprintf("Symbiosis: %v (type=%v, path=%v)", error.Message, error.ErrorType, error.Path)
}

func (client *SymbiosisClient) describeCluster(name string) (*Cluster, error) {
	api := client.symbiosisApi
	resp, err := api.R().SetResult(Cluster{}).ForceContentType("application/json").Get(fmt.Sprintf("rest/v1/cluster/%s", name))
	if err != nil {
		var cluster Cluster
		return &cluster, err
	}
	if resp.StatusCode() == 404 {
		return nil, nil
	}
	if resp.StatusCode() != 200 {
		symbiosisErr := resp.Error().(*SymbiosisApiError)
		return nil, symbiosisErr
	}
	cluster := resp.Result().(*Cluster)
	return cluster, nil
}

func (client *SymbiosisClient) describeNodePool(id string) (*NodePool, error) {
	api := client.symbiosisApi
	resp, err := api.R().SetResult(NodePool{}).ForceContentType("application/json").Get(fmt.Sprintf("rest/v1/node-pool/%v", id))
	if err != nil {
		var NodePool NodePool
		return &NodePool, err
	}
	if resp.StatusCode() == 404 {
		return nil, nil
	}
	if resp.StatusCode() != 200 {
		symbiosisErr := resp.Error().(*SymbiosisApiError)
		return nil, symbiosisErr
	}
	nodePool := resp.Result().(*NodePool)
	return nodePool, nil
}

func (client *SymbiosisClient) describeTeamMember(email string) (*TeamMember, error) {
	api := client.symbiosisApi
	resp, err := api.R().SetResult(TeamMember{}).ForceContentType("application/json").Get(fmt.Sprintf("rest/v1/team/member/%s", email))
	if err != nil {
		var teamMember TeamMember
		return &teamMember, err
	}
	if resp.StatusCode() == 404 {
		return nil, nil
	}
	if resp.StatusCode() != 200 {
		symbiosisErr := resp.Error().(*SymbiosisApiError)
		return nil, symbiosisErr
	}
	teamMember := resp.Result().(*TeamMember)
	return teamMember, nil
}

func (client *SymbiosisClient) describeTeamMemberInvitation(email string) (*TeamMember, error) {
	api := client.symbiosisApi
	resp, err := api.R().SetResult(TeamMember{}).ForceContentType("application/json").Get(fmt.Sprintf("rest/v1/team/member/invite/%s", email))
	if err != nil {
		var teamMember TeamMember
		return &teamMember, err
	}
	if resp.StatusCode() == 404 {
		return nil, nil
	}
	if resp.StatusCode() != 200 {
		symbiosisErr := resp.Error().(*SymbiosisApiError)
		return nil, symbiosisErr
	}
	teamMember := resp.Result().(*TeamMember)
	return teamMember, nil
}
