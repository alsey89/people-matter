package memorial

import "github.com/alsey89/people-matter/internal/schema"

type postContributorRequest struct {
	Email        string                   `json:"email" validate:"required,email"`
	Relationship schema.RelationshipConst `json:"relationship" validate:"required"`
}

type getContribtionsResponse struct {
	CondolenceElements []schema.ContributionCondolenceElement `json:"condolenceElements"`
	GalleryElements    []schema.ContributionGalleryElement    `json:"galleryElements"`
	StoryElements      []schema.ContributionStoryElement      `json:"storyElements"`
	TimelineElements   []schema.ContributionTimelineElement   `json:"timelineElements"`
}

type postCondolenceElementRequest struct {
	ElementTitle       string `json:"elementTitle" validate:"required"`
	ElementDescription string `json:"elementDescription" validate:"required"`
	DesignElementID    string `json:"designElementId" validate:"required"`
	ElementAuthor      string `json:"elementAuthor" validate:"required"`
	IsImmutable        bool   `json:"isImmutable" gorm:"default:false"`
}
type postGalleryElementRequest struct {
	ElementTitle         string  `json:"elementTitle" validate:"required"`
	ElementDescription   string  `json:"elementDescription" validate:"required"`
	ElementMediaURL      string  `json:"elementMediaUrl" validate:"required"`
	ElementMediaType     string  `json:"elementMediaType" validate:"required"`
	ElementDate          string  `json:"elementDate" validate:"required"`
	ElementLocation      *string `json:"elementLocation"`
	ElementGooglePlaceID *string `json:"elementGooglePlaceId"`
	IsImmutable          bool    `json:"isImmutable" gorm:"default:false"`
}
type postTimelineElementRequest struct {
	ElementTitle         string                `json:"elementTitle" validate:"required"`
	ElementDescription   string                `json:"elementDescription" validate:"required"`
	ElementMediaURL      *string               `json:"elementMediaUrl"`
	ElementEventType     schema.EventTypeConst `json:"elementEventType" validate:"required"`
	ElementDate          string                `json:"elementDate" validate:"required"`
	ElementLocation      *string               `json:"elementLocation"`
	ElementGooglePlaceID *string               `json:"elementGooglePlaceId"`
	IsImmutable          bool                  `json:"isImmutable" gorm:"default:false"`
}
type postStoryElementRequest struct {
	ElementTitle       string  `json:"elementTitle" validate:"required"`
	ElementDescription string  `json:"elementDescription" validate:"required"`
	ElementMediaURL    *string `json:"elementMediaUrl"`
	ElementAuthor      string  `json:"elementAuthor" validate:"required"`
	IsImmutable        bool    `json:"isImmutable" gorm:"default:false"`
}

type putCondolenceElementRequest struct {
	ElementTitle       string `json:"elementTitle" validate:"required"`
	ElementDescription string `json:"elementDescription" validate:"required"`
	ElementAuthor      string `json:"elementAuthor" validate:"required"`
	DesignElementID    string `json:"designElementId" validate:"required"`
	IsImmutable        bool   `json:"isImmutable" gorm:"default:false"`
}
type putGalleryElementRequest struct {
	ElementTitle         string  `json:"elementTitle" validate:"required"`
	ElementDescription   string  `json:"elementDescription" validate:"required"`
	ElementDate          string  `json:"elementDate" validate:"required"`
	ElementLocation      *string `json:"elementLocation"`
	ElementGooglePlaceID *string `json:"elementGooglePlaceId"`
	IsImmutable          bool    `json:"isImmutable" gorm:"default:false"`

	// ElementMediaType     string  `json:"elementMediaType" validate:"required"` //not allowed to update media
	// ElementMediaURL      string  `json:"elementMediaUrl" validate:"required"` //not allowed to update media
}
type putTimelineElementRequest struct {
	ElementTitle         string                `json:"elementTitle" validate:"required"`
	ElementDescription   string                `json:"elementDescription" validate:"required"`
	ElementEventType     schema.EventTypeConst `json:"elementEventType" validate:"required"`
	ElementDate          string                `json:"elementDate" validate:"required"`
	ElementLocation      *string               `json:"elementLocation"`
	ElementGooglePlaceID *string               `json:"elementGooglePlaceId"`
	IsImmutable          bool                  `json:"isImmutable" gorm:"default:false"`

	// ElementMediaURL      *string               `json:"elementMediaUrl"` //not allowed to update media
}
type putStoryElementRequest struct {
	ElementTitle       string `json:"elementTitle" validate:"required"`
	ElementDescription string `json:"elementDescription" validate:"required"`
	ElementAuthor      string `json:"elementAuthor" validate:"required"`
	IsImmutable        bool   `json:"isImmutable" gorm:"default:false"`

	// ElementMediaURL    string `json:"elementMediaUrl" validate:"required"` //not allowed to update media
}

type patchExportStateRequest struct {
	ExportState schema.ExportStateConst `json:"exportState" validate:"required"`
}

type putCondolenceElementStateRequest struct {
	IsApproved bool `json:"isApproved" gorm:"default:false,not null"`
}
type putGalleryElementStateRequest struct {
	IsApproved bool `json:"isApproved" gorm:"default:false,not null"`
}
type putTimelineElementStateRequest struct {
	IsApproved bool `json:"isApproved" gorm:"default:false,not null"`
}
type putStoryElementStateRequest struct {
	IsApproved bool `json:"isApproved" gorm:"default:false,not null"`
}
