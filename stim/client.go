package stim

import (
	"fmt"
)

type StimApiError struct {
	Status    int32
	ErrorType string `json:"error"`
	Message   string
	Path      string
}

func (error *StimApiError) Error() string {
	return fmt.Sprintf("Stim: %v (%v)", error.Message, error.ErrorType)
}

func (client *StimClient) describeCluster(name string) (*Cluster, error) {
  api := client.stimApi
  resp, err := api.R().SetResult(Cluster{}).ForceContentType("application/json").Get(fmt.Sprintf("rest/v1/cluster/%s", name))
  if err != nil {
    var cluster Cluster
    return &cluster, err
  }
  if resp.StatusCode() == 404 {
    return nil, nil
  }
  if resp.StatusCode() != 200 {
		stimErr := resp.Error().(*StimApiError)
    var cluster Cluster
    return &cluster, stimErr
  }
	cluster := resp.Result().(*Cluster)
  return cluster, nil
}

func (client *StimClient) describeTeamMember(email string) (*TeamMember, error) {
  api := client.stimApi
	resp, err := api.R().SetResult(TeamMember{}).ForceContentType("application/json").Get(fmt.Sprintf("rest/v1/team/member/%s", email))
  if err != nil {
    var teamMember TeamMember
    return &teamMember, err
  }
  if resp.StatusCode() == 404 {
    return nil, nil
  }
  if resp.StatusCode() != 200 {
		stimErr := resp.Error().(*StimApiError)
    return nil, stimErr
  }
	teamMember := resp.Result().(*TeamMember)
  return teamMember, nil
}

func (client *StimClient) describeTeamMemberInvitation(email string) (*TeamMember, error) {
  api := client.stimApi
	resp, err := api.R().SetResult(TeamMember{}).ForceContentType("application/json").Get(fmt.Sprintf("rest/v1/team/member/invite/%s", email))
  if err != nil {
    var teamMember TeamMember
    return &teamMember, err
  }
  if resp.StatusCode() == 404 {
    return nil, nil
  }
  if resp.StatusCode() != 200 {
		stimErr := resp.Error().(*StimApiError)
    return nil, stimErr
  }
	teamMember := resp.Result().(*TeamMember)
  return teamMember, nil
}
