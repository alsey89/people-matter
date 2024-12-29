package company

func (d *Domain) registerRoutes() {
	d.logger.Info("Registering routes.")

	e := d.params.Server.GetServer()

	TenantGroup := e.Group(
		"/api/v1/company",
		// d.params.Identity.MustResolveFSPID(),
		// d.params.Identity.MustBeAuthenticated(),
	)
	TenantGroup.GET("",
		d.handlerGetDashboard,
	)

	// forms -------------------------------------------------------

	// account
	TenantGroup.GET("/account",
		d.handlerGetCompanyAccount,
	)
	TenantGroup.PUT("/account",
		d.handlerUpdateCompanyAccount,
	)
	// positions
	TenantGroup.GET("/positions",
		d.handlerGetPositions,
	)
	TenantGroup.POST("/positions",
		d.handlerCreatePosition,
	)
	TenantGroup.PUT("/positions/:positionID",
		d.handlerUpdatePosition,
	)
	TenantGroup.DELETE("/positions/:positionID",
		d.handlerDeletePosition,
	)

}
