package memorial

func (d *Domain) registerRoutes() {
	d.logger.Info("Registering routes.")

	e := d.params.Server.GetServer()

	MemorialGroup := e.Group(
		"/api/v1/memorial/:memorialID",
		d.params.Identity.MustResolveFSPID(),
		d.params.Identity.MustBeAuthenticated(),
		d.params.Identity.MustMatchTenantIdentifierAndToken(),
	)

	//! FE Application ----------------------------------------------
	// the nuxt3 app will use this route to fetch the memorial data
	e.GET(
		"/api/v1/memorial/:memorialID/app/contribution",
		d.handlerGetMemorialContributions,
		d.params.Identity.MustResolveFSPID(),
	)
	// the nuxt3 app will use this route to update the export status
	e.PATCH(
		"/api/v1/memorial/:memorialID/app/export/:exportID",
		d.handlerUpdateExportState,
		d.params.Identity.MustResolveFSPID(),
	)

	//! Forms --------------------------------------------------------
	// user story: As a curator, I want to be able to select correct roles in forms
	MemorialGroup.GET(
		"/curator/form/role",
		d.handlerGetMemorialFormRole,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a contributor, when contributing timeline elements, I want to be able to select correct event types in forms
	MemorialGroup.GET(
		"/contributor/form/timeline/type",
		d.handlerGetTimelineElementType,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)

	//! Curator --------------------------------------------

	// user story: As a curator, I want to view the memorial dashboard
	MemorialGroup.GET(
		"/curator/dashboard",
		d.handlerGetDashboard,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to see a list of all the contributors
	MemorialGroup.GET(
		"/curator/contributor",
		d.handlerGetContributors,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to see a list of all contributor applicants
	MemorialGroup.GET(
		"/curator/contributor/application",
		d.handlerGetContributorApplications,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to approve a contributor application
	MemorialGroup.PUT(
		"/curator/contributor/application/:applicationID",
		d.handlerAcceptContributorApplication,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to reject a contributor application
	MemorialGroup.DELETE(
		"/curator/contributor/application/:applicationID",
		d.handlerRejectContributorApplication,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to see a list of all contributors I have invited
	MemorialGroup.GET(
		"/curator/contributor/invitation",
		d.handlerGetContributorInvitations,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to invite a contributor to the memorial
	MemorialGroup.POST(
		"/curator/contributor",
		d.handlerInviteContributor,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to be able to delete existing contributor invitations
	MemorialGroup.DELETE(
		"/curator/contributor/invitation/:invitationID",
		d.handlerDeleteContributorInvitation,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to reinvite a contributor to the memorial if they have not accepted the invitation or if their invitation has expired
	MemorialGroup.PUT(
		"/curator/contributor/invitation/:invitationID",
		d.handlerReinviteContributor,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to update a contributor's role
	MemorialGroup.PUT(
		"/curator/contributor/:contributorMemorialRoleID",
		d.handlerUpdateContributor,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to remove a contributor from the memorial
	MemorialGroup.DELETE(
		"/curator/contributor/:contributorMemorialRoleID",
		d.handlerDeleteContributor,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)

	// user story: As a curator, I want to update info on a condolenceElement before it is approved
	MemorialGroup.PUT(
		"/curator/contributions/condolence/:contributionCondolenceElementID",
		d.handlerUpdateContributionCondolenceElement,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to update info on a galleryElement before it is approved
	MemorialGroup.PUT(
		"/curator/contributions/gallery/:contributionGalleryElementID",
		d.handlerUpdateContributionGalleryElement,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to update info on a timelineElement before it is approved
	MemorialGroup.PUT(
		"/curator/contributions/timeline/:contributionTimelineElementID",
		d.handlerUpdateContributionTimelineElement,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to update info on a storyElement before it is approved
	MemorialGroup.PUT(
		"/curator/contributions/story/:contributionStoryElementID",
		d.handlerUpdateContributionStoryElement,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)

	// user story: As a curator, I want to view a list of all contributions from all contributors
	MemorialGroup.GET(
		"/curator/contribution",
		d.handlerGetMemorialContributions,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to approve or unapprove a condolenceElement contribution
	MemorialGroup.PATCH(
		"/curator/contribution/condolence/:contributionCondolenceElementID",
		d.handlerUpdateContributionCondolenceElementState,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to approve or unapprove a galleryElement contribution
	MemorialGroup.PATCH(
		"/curator/contribution/gallery/:contributionGalleryElementID",
		d.handlerUpdateContributionGalleryElementState,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to approve or unapprove a timelineElement contribution
	MemorialGroup.PATCH(
		"/curator/contribution/timeline/:contributionTimelineElementID",
		d.handlerUpdateContributionTimelineElementState,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to approve or unapprove a storyElement contribution
	MemorialGroup.PATCH(
		"/curator/contribution/story/:contributionStoryElementID",
		d.handlerUpdateContributionStoryElementState,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)

	// user story: As a curator, I want to fetch a list of all exports related to the memorial
	MemorialGroup.GET(
		"/curator/export",
		d.handlerGetAllExports,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to fetch the exportState of the memorial
	MemorialGroup.GET(
		"/curator/export/:exportID",
		d.handlerGetExportState,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)
	// user story: As a curator, I want to export/publish the preview as a static site when I'm happy with the contributions
	MemorialGroup.POST(
		"/curator/export",
		d.handlerExportMemorial,
		d.params.Identity.MustBeCuratorOfCurrentMemorial(),
	)

	//! Contributor -----------------------------------------

	//todo: unify the route naming, it's a bit of a mess now

	// user story: As a contributor, I want to view a list of all my contributions
	MemorialGroup.GET(
		"/contributor/contribution",
		d.handlerGetContributions,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)
	// user story: As a contributor, I want to fetch a presigned upload URL and use it to upload media
	MemorialGroup.GET(
		"/contributor/upload/presigned",
		d.handlerGetPresignedUploadURL,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)

	// user story: As a contributor, I want to create a contributionCondolenceElement record
	MemorialGroup.POST(
		"/contributor/condolence",
		d.handlerCreateContributionCondolenceElement,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)
	// user story: As a contributor, I want to create a contributionGalleryElement record and store the uploaded mediaURL
	MemorialGroup.POST(
		"/contributor/gallery",
		d.handlerCreateContributionGalleryElement,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)
	// user story: As a contributor, I want to create a contributionTimelineElement record and optionally store the uploaded mediaURL
	MemorialGroup.POST(
		"/contributor/timeline",
		d.handlerCreateContributionTimelineElement,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)
	// user story: As a contributor, I want to create a contributionStoryElement record and optionally store the uploaded mediaURL
	MemorialGroup.POST(
		"/contributor/story",
		d.handlerCreateContributionStoryElement,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)

	//user story: As a contributor, I want to update a contributionGalleryElement record
	MemorialGroup.PUT(
		"/contributor/contributions/gallery/:contributionGalleryElementID",
		d.handlerUpdateContributionGalleryElement,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)
	//user story: As a contributor, I want to update a contributionTimelineElement record
	MemorialGroup.PUT(
		"/contributor/contributions/timeline/:contributionTimelineElementID",
		d.handlerUpdateContributionTimelineElement,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)
	//user story: As a contributor, I want to update a contributionStoryElement record
	MemorialGroup.PUT(
		"/contributor/contributions/story/:contributionStoryElementID",
		d.handlerUpdateContributionStoryElement,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)
	// user story: As a contributor, I want to update a contributionCondolenceElement record
	MemorialGroup.PUT(
		"/contributor/contributions/condolence/:contributionCondolenceElementID",
		d.handlerUpdateContributionCondolenceElement,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)

	//user story: As a contributor, I want to delete a contributionGalleryElement record
	MemorialGroup.DELETE(
		"/contributor/contributions/gallery/:contributionGalleryElementID",
		d.handlerDeleteContributionGalleryElement,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)
	//user story: As a contributor, I want to delete a contributionTimelineElement record
	MemorialGroup.DELETE(
		"/contributor/contributions/timeline/:contributionTimelineElementID",
		d.handlerDeleteContributionTimelineElement,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)
	//user story: As a contributor, I want to delete a contributionStoryElement record
	MemorialGroup.DELETE(
		"/contributor/contributions/story/:contributionStoryElementID",
		d.handlerDeleteContributionStoryElement,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)
	// user story: As a contributor, I want to delete a contributionCondolenceElement record
	MemorialGroup.DELETE(
		"/contributor/contributions/condolence/:contributionCondolenceElementID",
		d.handlerDeleteContributionCondolenceElement,
		d.params.Identity.MustBeContributorOrCuratorOfCurrentMemorial(),
	)

	//! User ---------------------------------------------------------
	// user story: As **non-contributor**, I want to apply to be a contributor
	MemorialGroup.POST(
		"/contributor/application",
		d.handlerApplyToBeContributor,
	)
	// user story: As a **non-contributor**, I want to accept an invitation to be a contributor
	MemorialGroup.POST(
		"/contributor/invitation/:invitationID/accept",
		d.handlerAcceptContributorInvitation,
	)
}
