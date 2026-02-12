package controllers

import (
	"log/slog"

	repoProvider "tsv-service/internal/providers/repository"
	svcProvider "tsv-service/internal/providers/service"
)

type Base interface {
	Logger() *slog.Logger
	ServiceProvider() svcProvider.Provider
	RepositoryProvider() repoProvider.Provider
}

type BaseController struct {
	logger *slog.Logger
	sp     svcProvider.Provider
	rp     repoProvider.Provider
}

func NewBaseController(
	logger *slog.Logger,
	sp svcProvider.Provider,
	rp repoProvider.Provider,
) *BaseController {
	return &BaseController{logger: logger, sp: sp, rp: rp}
}

func (c *BaseController) Logger() *slog.Logger                      { return c.logger }
func (c *BaseController) ServiceProvider() svcProvider.Provider     { return c.sp }
func (c *BaseController) RepositoryProvider() repoProvider.Provider { return c.rp }
