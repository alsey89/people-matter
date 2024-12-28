package schema

import (
	"time"
)

type ContributionStateConst string

const (
	ContributionStatePending ContributionStateConst = "pending" // contributor has submitted this contribution
	ContributionStatePrivate ContributionStateConst = "private" // curator has marked this contribution for publishing, it's in a private state
	ContributionStatePublic  ContributionStateConst = "public"  // curator has published this contribution, it's public in the memorial
)

type ContributionCondolenceElement struct {
	BaseModelWithTime
	FSPID             uint                   `json:"fspId"`
	MemorialID        uint                   `json:"memorialId"`
	Contributor       *User                  `json:"contributor" gorm:"foreignKey:ContributorID"`
	ContributorID     uint                   `json:"contributorId"`
	ContributionState ContributionStateConst `json:"contributionState" sql:"type:ENUM('pending', 'private', 'public')"`
	IsImmutable       bool                   `json:"isImmutable" gorm:"default:false"` //contributor has specified that their contribution cannot be modified
	IsPinned          bool                   `json:"isPinned" gorm:"default:false"`

	ElementTitle       string `json:"elementTitle"`
	ElementDescription string `json:"elementDescription"`
	ElementAuthor      string `json:"elementAuthor"`
	DesignElementID    string `json:"designElementId"`

	//relationship between ContributorID and user_memorial_roles.user_id needs to be explicitly defined
	ContributorMemorialRole *UserMemorialRole `json:"contributorMemorialRole" gorm:"foreignKey:ContributorID;references:UserID"`
}

type ContributionGalleryElement struct {
	BaseModelWithTime
	FSPID             uint                   `json:"fspId"`
	MemorialID        uint                   `json:"memorialId"`
	Contributor       *User                  `json:"contributor" gorm:"foreignKey:ContributorID"`
	ContributorID     uint                   `json:"contributorId"`
	ContributionState ContributionStateConst `json:"contributionState" sql:"type:ENUM('pending', 'private', 'public')"`
	IsImmutable       bool                   `json:"isImmutable" gorm:"default:false"` //contributor has specified that their contribution cannot be modified

	ElementTitle         string    `json:"elementTitle"`
	ElementDescription   string    `json:"elementDescription"`
	ElementDate          time.Time `json:"elementDate"`
	ElementMediaURL      string    `json:"elementMediaUrl"`
	ElementMediaType     string    `json:"elementType"`
	ElementLocation      *string   `json:"elementLocation"`
	ElementGooglePlaceID *string   `json:"elementGooglePlaceId"`
	HasEXIF              bool      `json:"hasEXIF"`
	UseEXIF              bool      `json:"useEXIF"`

	//relationship between ContributorID and user_memorial_roles.user_id needs to be explicitly defined
	ContributorMemorialRole *UserMemorialRole `json:"contributorMemorialRole" gorm:"foreignKey:ContributorID;references:UserID"`
}

type ContributionStoryElement struct {
	BaseModelWithTime
	FSPID             uint                   `json:"fspId"`
	MemorialID        uint                   `json:"memorialId"`
	Contributor       *User                  `json:"contributor" gorm:"foreignKey:ContributorID"`
	ContributorID     uint                   `json:"contributorId"`
	ContributionState ContributionStateConst `json:"contributionState" sql:"type:ENUM('pending', 'private', 'public')"`
	IsImmutable       bool                   `json:"isImmutable" gorm:"default:false"` //contributor has specified that their contribution cannot be modified

	ElementTitle       string    `json:"elementTitle"`
	ElementDescription string    `json:"elementDescription"`
	ElementDate        time.Time `json:"elementDate"`
	ElementMediaURL    string    `json:"elementMediaUrl"`
	ElementAuthor      string    `json:"elementAuthor"`

	//relationship between ContributorID and user_memorial_roles.user_id needs to be explicitly defined
	ContributorMemorialRole *UserMemorialRole `json:"contributorMemorialRole" gorm:"foreignKey:ContributorID;references:UserID"`
}

type EventTypeConst string

const (
	EventTypeBirth       EventTypeConst = "birth"
	EventTypeDeath       EventTypeConst = "death"
	EventTypeMarriage    EventTypeConst = "marriage"
	EventTypeDivorce     EventTypeConst = "divorce"
	EventTypeEducation   EventTypeConst = "education"
	EventTypeGraduation  EventTypeConst = "graduation"
	EventTypeEmployment  EventTypeConst = "employment"
	EventTypeRetirement  EventTypeConst = "retirement"
	EventTypePromotion   EventTypeConst = "promotion"
	EventTypeAward       EventTypeConst = "award"
	EventTypeAchievement EventTypeConst = "achievement"
	EventTypeMilestone   EventTypeConst = "milestone"
	EventTypeOther       EventTypeConst = "other"
)

type ContributionTimelineElement struct {
	BaseModelWithTime
	FSPID             uint                   `json:"fspId"`
	MemorialID        uint                   `json:"memorialId"`
	Contributor       *User                  `json:"contributor" gorm:"foreignKey:ContributorID"`
	ContributorID     uint                   `json:"contributorId"`
	ContributionState ContributionStateConst `json:"contributionState" sql:"type:ENUM('pending', 'private', 'public')"`
	IsImmutable       bool                   `json:"isImmutable" gorm:"default:false"` //contributor has specified that their contribution cannot be modified

	ElementTitle         string         `json:"elementTitle"`
	ElementDescription   string         `json:"elementDescription"`
	ElementDate          time.Time      `json:"elementDate"`
	ElementMediaURL      *string        `json:"elementMediaUrl"`
	ElementEventType     EventTypeConst `json:"elementEventType" sql:"type:ENUM('birth', 'death', 'marriage', 'divorce', 'education', 'graduation', 'employment', 'retirement', 'promotion', 'award', 'achievement', 'milestone', 'other')"`
	ElementLocation      *string        `json:"elementLocation"`
	ElementGooglePlaceID *string        `json:"elementGooglePlaceId"`

	//relationship between ContributorID and user_memorial_roles.user_id needs to be explicitly defined
	ContributorMemorialRole *UserMemorialRole `json:"contributorMemorialRole" gorm:"foreignKey:ContributorID;references:UserID"`
}
