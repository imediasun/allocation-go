package infra

import (
	"gitlab.hotel.tools/backend-team/allocation-go/internal/infra/adapter"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/infra/repo"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/infra/service"
	"go.uber.org/fx"
)

var Constructors = fx.Provide(
	repo.NewBookingRepoFactory,

	adapter.NewBookingAdapterFactory,

	service.NewAllocatorService,
)
