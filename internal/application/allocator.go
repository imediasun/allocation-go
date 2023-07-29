package application

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/db"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/log"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/adapter"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/service"
	"gitlab.hotel.tools/backend-team/allocation-go/pkg/api/openapi"
	"go.uber.org/zap"
)

type AllocatorService struct {
	ctx    context.Context
	logger log.Logger
	db     db.DB

	bookingAdapterFactory adapter.BookingAdapterFactory
	allocatorService      service.AllocatorService
}

func NewAllocatorServer(
	ctx context.Context,
	logger log.Logger,
	db db.DB,
	bookingAdapterFactory adapter.BookingAdapterFactory,
	allocatorService service.AllocatorService,
) openapi.ServerInterface {
	server := &AllocatorService{
		ctx:                   ctx,
		logger:                logger.Named("allocator"),
		db:                    db,
		bookingAdapterFactory: bookingAdapterFactory,
		allocatorService:      allocatorService,
	}

	return server
}

func (s *AllocatorService) FindPets(ctx echo.Context, params openapi.FindPetsParams) error {
	s.logger.Info("FindPets")
	return domain.ErrUnimplemented
}

func (s *AllocatorService) AllocateAll(ctx echo.Context, reservationIDs []int32, userID *int32) error {

	// Call the AllocateAll method and capture the slice of AllocateResult and error
	_, err := s.allocatorService.AllocateAll(s.ctx, reservationIDs, userID)
	if err != nil {
		s.logger.Error("failed to allocate all", zap.Error(err))

		// If there was an error, return it
		return err
	}

	s.logger.Info("AddPet")

	// If there was no error, return nil
	return nil
}

func (s *AllocatorService) AutoAllocate(ctx echo.Context, reservationID int, isNotify bool) error {
	fmt.Printf("Value is: %d and type is: %T\\n", reservationID)
	// Call the AllocateAll method and capture the slice of AllocateResult and error
	s.allocatorService.AutoAllocate(s.ctx, reservationID, true)

	s.logger.Info("AddPet")

	// If there was no error, return nil
	return nil
}

func (s *AllocatorService) DeletePet(ctx echo.Context, id int64) error {
	s.logger.Info("DeletePet")
	return domain.ErrUnimplemented
}

func (s *AllocatorService) FindPetByID(ctx echo.Context, id int64) error {
	s.logger.Info("FindPetByID")
	return domain.ErrUnimplemented
}
