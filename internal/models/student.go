package models

import "time"

const (
	StatusPending = "pending"
	StatusPassed  = "passed"
	StatusFailed  = "failed"

	StageGitHub = "github"
	StageDocker = "docker"
	StageK8s    = "k8s"
)

type Student struct {
	ID                uint   `json:"id" gorm:"primaryKey"`
	Name              string `json:"name" gorm:"not null"`
	Email             string `json:"email" gorm:"not null;uniqueIndex"`
	GitHubUsername    string `json:"githubUsername" gorm:"not null"`
	DockerHubUsername string `json:"dockerHubUsername" gorm:"not null"`
	Approved          bool   `json:"approved" gorm:"not null;default:false"`

	GitHubStatus        string     `json:"githubStatus" gorm:"not null;default:pending"`
	GitHubErrorMessage  string     `json:"githubErrorMessage"`
	GitHubLastCheckedAt *time.Time `json:"githubLastCheckedAt"`

	DockerStatus        string     `json:"dockerStatus" gorm:"not null;default:pending"`
	DockerErrorMessage  string     `json:"dockerErrorMessage"`
	DockerLastCheckedAt *time.Time `json:"dockerLastCheckedAt"`

	K8sStatus        string     `json:"k8sStatus" gorm:"not null;default:pending"`
	K8sErrorMessage  string     `json:"k8sErrorMessage"`
	K8sLastCheckedAt *time.Time `json:"k8sLastCheckedAt"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *Student) AllPassed() bool {
	return s.GitHubStatus == StatusPassed &&
		s.DockerStatus == StatusPassed &&
		s.K8sStatus == StatusPassed
}
