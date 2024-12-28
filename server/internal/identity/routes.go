package identity

func (d *Domain) registerRoutes() {
	d.logger.Info("Registering routes.")

	e := d.params.Server.GetServer()

	authGroup := e.Group(
		"/api/v1/auth",
		d.MustResolveFSPID(),
	)

	// Forms ----------------------------------------------------------

	// user story: As a curator, I want to be able to select correct roles in forms
	authGroup.GET("/form/relationship",
		d.handlerGetRelationships,
	)

	// Auth -----------------------------------------------------------

	// this route is called for *forms that are reached in isolation* (e.g. user goes to the sign up address directly)
	authGroup.GET("/csrf",
		d.csrfHandler,
	)
	// this is a general route that checks if the user is authenticated
	authGroup.GET("/check",
		d.authCheckHandler,
		d.MustBeAuthenticated(),
		d.MustMatchTenantIdentifierAndToken(),
	)
	// user story: As a user, I want to sign in
	authGroup.POST("/signin",
		d.signinHandler,
	)
	// user story: As a user with roles in multiple memorials, I want to switch my active memorial & role
	authGroup.PATCH("/memorial/:memorialID/active",
		d.handlerSwitchActiveMemorial,
		d.MustBeAuthenticated(),
	)
	// user story: As a user, I want to confirm my email address after signing up
	authGroup.POST("/signup/confirm",
		d.confirmEmailHandler,
		d.MustHaveValidConfirmationToken(d.config.JWTEmailConfirmationScope),
	)
	// user story: As a user, I want to signout
	authGroup.POST("/signout",
		d.signoutHandler,
	)
	// user story: As a user, I want to request a password reset
	authGroup.POST("/password/reset/request",
		d.requestResetPasswordHandler,
	)
	// user story: As a user, I want to confirm my password reset and set a new password
	authGroup.POST("/password/reset/confirm",
		d.handlerResetPassword,
		d.MustHaveValidConfirmationToken(d.config.JWTPasswordResetScope),
	)
	// user story: As a user invited by an Tenant admin, I want to set my password and confirm my email address
	authGroup.POST("/password/set/confirm",
		d.handlerSetPasswordAndConfirmEmail,
		d.MustHaveValidConfirmationToken(d.config.JWTPasswordResetScope),
	)

	// Applications & Invitations for contributors -------------------------------------

	// user story: as a new user without an account in the fsp, I want to apply to be a contributor to a memorial and create an account
	authGroup.POST("/application/signup",
		d.handlerApplicationSignup,
	)
	// user story: as a user with an account in the fsp, I want to apply to be a contributor to a memorial without creating a new account
	authGroup.POST("/application",
		d.handlerApplication,
		d.MustBeAuthenticated(),
	)
	// user story: as a user without an account in the fsp, I want to accept an invitation to be a contributor to a memorial and create an account
	authGroup.POST("/invitation/accept/signup",
		d.handlerAcceptInvitationSignup,
	)
	// user story: as a user with an account in the fsp, I want to accept an invitation to be a contributor to a memorial without creating a new account
	authGroup.POST("/invitation/accept",
		d.handlerAcceptInvitation,
		d.MustBeAuthenticated(),
	)
}
