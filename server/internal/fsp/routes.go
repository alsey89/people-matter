package fsp

func (d *Domain) registerRoutes() {
	d.logger.Info("Registering routes.")

	e := d.params.Server.GetServer()

	FSPGroup := e.Group(
		"/api/v1/fsp",
		d.params.Identity.MustResolveFSPID(),
		d.params.Identity.MustBeAuthenticated(),
		d.params.Identity.MustMatchTenantIdentifierAndToken(),
		d.params.Identity.MustBeFSPAdminOrHigher(),
	)
	FSPGroup.GET("",
		d.handlerGetDashboard,
	)

	// forms -------------------------------------------------------

	// Tenant SUPERADMIN ----------------------------------------------

	// account

	FSPGroup.GET("/superadmin/account",
		d.handlerGetFSPAccount,
		d.params.Identity.MustBeFSPSuperAdmin(),
	)
	FSPGroup.PUT("/superadmin/account",
		d.handlerUpdateFSPAccount,
		d.params.Identity.MustBeFSPSuperAdmin(),
	)

	// team

	FSPGroup.GET("/superadmin/team",
		d.handlerGetFSPTeam,
		d.params.Identity.MustBeFSPSuperAdmin(),
	)
	FSPGroup.POST("/superadmin/team",
		d.handlerAddFSPTeam,
		d.params.Identity.MustBeFSPSuperAdmin(),
	)
	FSPGroup.PUT("/superadmin/team/:teamMemberID",
		d.handlerUpdateFSPTeam,
		d.params.Identity.MustBeFSPSuperAdmin(),
	)
	FSPGroup.DELETE("/superadmin/team/:teamMemberID",
		d.handlerDeleteFSPTeam,
		d.params.Identity.MustBeFSPSuperAdmin(),
	)

	// memorial

	FSPGroup.GET("/superadmin/memorial",
		d.handlerGetMemorials,
		d.params.Identity.MustBeFSPSuperAdmin(),
	)
	FSPGroup.POST("/superadmin/memorial",
		d.handlerAddMemorial,
		d.params.Identity.MustBeFSPSuperAdmin(),
	)
	FSPGroup.PUT("/superadmin/memorial/:memorialID",
		d.handlerUpdateMemorial,
		d.params.Identity.MustBeFSPSuperAdmin(),
	)
	FSPGroup.DELETE("/superadmin/memorial/:memorialID",
		d.handlerDeleteMemorial,
		d.params.Identity.MustBeFSPSuperAdmin(),
	)

	//! Tenant Admin ---------------------------------------------------

	// memorial
	FSPGroup.GET("/admin/memorial",
		d.handlerGetMemorials,
		d.params.Identity.MustBeFSPAdminOrHigher(),
	)
	FSPGroup.POST("/admin/memorial",
		d.handlerAddMemorial,
		d.params.Identity.MustBeFSPAdminOrHigher(),
	)
	FSPGroup.PUT("/admin/memorial/:memorialID",
		d.handlerUpdateMemorial,
		d.params.Identity.MustBeFSPAdminOrHigher(),
	)
	FSPGroup.DELETE("/admin/memorial/:memorialID",
		d.handlerDeleteMemorial,
		d.params.Identity.MustBeFSPAdminOrHigher(),
	)

	//! TBD ---------------------------------------------------------

	FSPGroup.GET("/partner",
		d.handlerGetPartners,
	)

}
