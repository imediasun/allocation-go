package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/volatiletech/null/v8"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/db"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/common/adapter/log"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/adapter"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/model"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/repo"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/service"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type ActionAbstractUpdate struct{}

type allocatorService struct {
	ctx                   context.Context
	reservation           model.Reservation
	logger                log.Logger
	db                    db.DB
	bookingRepoFactory    repo.BookingRepoFactory
	bookingAdapterFactory adapter.BookingAdapterFactory
}

type ReservationStatus string

type ReservationPaymentOption string

func (s *allocatorService) getAgent(id int32) (*model.Agent, error) {
	//fmt.Printf("Value is: %d and type is Agent: %T\\n", id)
	row := s.db.QueryRow("SELECT id, name, AccountID FROM agents WHERE id = ?", id)

	agent := &model.Agent{}
	err := row.Scan(&agent.ID, &agent.Name, &agent.AccountID)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func (s *allocatorService) getCurrencyRates(bookingID int) ([]model.CurrencyRate, error) {
	rows, err := s.db.Query("SELECT BookingID, Source, Target, Rate, Date, Final FROM booking_currency_rates WHERE BookingID = ?", bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []model.CurrencyRate

	for rows.Next() {
		var rate model.CurrencyRate

		err := rows.Scan(&rate.BookingID, &rate.Source, &rate.Target, &rate.Rate, &rate.Date, &rate.Final)
		if err != nil {
			return nil, err
		}

		rates = append(rates, rate)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return rates, nil
}

func convertNullStringToString(ns null.String) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func (s *allocatorService) getItemsForGroup(ctx context.Context, groupID int32) ([]model.BookingItems, error) {
	logger := s.logger.WithMethod(ctx, "AllocateAll")
	rows, err := s.db.Query("SELECT ID,Type,VenueID,ProductID,Status FROM booking_items WHERE GroupID = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.BookingItems, 0)
	var count int
	for rows.Next() {
		count++
		var item model.BookingItems

		err := rows.Scan(&item.ID, &item.Type, &item.VenueID, &item.ProductID, &item.Status)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Total rows BeforeGetProduct: %d\n", item.ProductID)
		// Get the product information for the item
		product, err := s.getProduct(convertNullStringToString(item.ProductID))
		if err != nil {
			return nil, err
		}
		item.Product = product

		items = append(items, item)
	}
	fmt.Printf("Total rows getItems: %d\n", count)
	if err = rows.Err(); err != nil {
		return nil, err
	}

	itemsJSON, err := json.Marshal(items)
	if err != nil {
		logger.Error("failed to marshal user to JSON", zap.Error(err))
		return nil, err
	}

	fmt.Println(string(itemsJSON))

	return items, nil
}

func (s *allocatorService) getGroups(ctx context.Context, bookingID int) ([]model.ReservationGroup, error) {

	logger := s.logger.WithMethod(ctx, "AllocateAll")

	rows, err := s.db.Query("SELECT * FROM booking_groups WHERE BookingID = ?", bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]model.ReservationGroup, 0)
	var count int
	for rows.Next() {
		count++
		group := model.ReservationGroup{}
		var startDate, endDate string // Use string type to store the date and time as strings

		err = rows.Scan(&group.ID, &group.BookingID, &group.PaxNationality, &startDate, &endDate, &group.ParentID)
		if err != nil {
			return nil, err
		}

		// List of possible layouts for date and time
		dateLayouts := []string{
			"2006-01-02 15:04:05", // Use your known format, add more if needed
			"2006-01-02T15:04:05Z07:00",
			time.RFC3339,
			// Add more layouts as required based on possible formats in the database
		}

		// Parse the date and time using multiple layouts
		var startTime, endTime time.Time
		var errParse error
		for _, layout := range dateLayouts {
			if startTime, errParse = time.Parse(layout, startDate); errParse == nil {
				break
			}
		}
		if errParse != nil {
			return nil, errParse
		}
		group.StartDate = startTime

		for _, layout := range dateLayouts {
			if endTime, errParse = time.Parse(layout, endDate); errParse == nil {
				break
			}
		}
		if errParse != nil {
			return nil, errParse
		}
		group.EndDate = endTime

		//groups = append(groups, group)

		// get Items for this group
		items, err := s.getItemsForGroup(ctx, group.ID)
		if err != nil {
			return nil, err
		}
		group.Items = items

		groups = append(groups, group)
	}

	fmt.Printf("Total rows getGroups: %d\n", count)

	if err = rows.Err(); err != nil {
		return nil, err
	}
	groupsJSON, err := json.Marshal(groups)
	if err != nil {
		logger.Error("failed to marshal user to JSON", zap.Error(err))
		return nil, err
	}

	fmt.Println(string(groupsJSON))
	return groups, nil
}

func (s *allocatorService) getProduct(productID string) (model.Product, error) {

	//fmt.Printf("Value is: %d and type is PProductID: %T\\n", productID)
	row := s.db.QueryRow("SELECT ID, Status, product_type FROM products WHERE ID = ?", productID)

	var product model.Product
	err := row.Scan(&product.ID, &product.Status, &product.ProductType) // and so on for all fields in Product
	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (s *allocatorService) getReservations(ctx context.Context, bookingIDs []int32) ([]model.Reservation, error) {
	var reservationList []model.Reservation
	//fmt.Printf("Value is: %d and type is BookingID: %T\\n", bookingIDs)
	ids := bookingIDs
	placeholders := make([]string, len(ids))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	query := fmt.Sprintf("SELECT ID,ProviderReference,Channel,Client,PaymentOption,CancellationDate,Segment,Source,Foct,MetaGroupID	FROM bookings WHERE ID IN (%s)", strings.Join(placeholders, ","))
	var interfaceIDs []interface{}
	for _, id := range ids {
		interfaceIDs = append(interfaceIDs, id)
	}

	// Execute the query with the interfaceIDs as separate parameters
	rows, err := s.db.Query(query, interfaceIDs...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []model.Reservation{}, ErrUserNotFound
		}
		return []model.Reservation{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var reservation model.Reservation
		err := rows.Scan(
			&reservation.ID,
			&reservation.ProviderReference,
			&reservation.Channel,
			&reservation.Client,
			&reservation.PaymentOption,
			&reservation.CancellationDate,
			&reservation.Segment,
			&reservation.Source,
			&reservation.Foct,
			&reservation.MetaGroupID,
		)
		if err != nil {
			return nil, err
		}

		// get Agent
		agent, err := s.getAgent(1) // Replace with correct Agent ID
		if err != nil {
			return nil, err
		}
		reservation.Creator = *agent

		// get CurrencyRates
		rates, err := s.getCurrencyRates(reservation.ID)
		if err != nil {
			return nil, err
		}
		reservation.CurrencyRates = rates

		// get Groups
		groups, err := s.getGroups(ctx, reservation.ID)
		if err != nil {
			return nil, err
		}
		reservation.Groups = groups

		reservationList = append(reservationList, reservation)
	}
	return reservationList, nil
}

func (s *allocatorService) getReservation(ctx context.Context, bookingID int) (*model.Reservation, error) {

	//fmt.Printf("Value is: %d and type is BookingID: %T\\n", bookingID)
	row := s.db.QueryRow("SELECT ID,ProviderReference,Channel,Client,PaymentOption,CancellationDate,Segment,Source,Foct,MetaGroupID	FROM bookings WHERE ID = ?", bookingID)

	reservation := &model.Reservation{}
	err := row.Scan(
		&reservation.ID,
		&reservation.ProviderReference,
		&reservation.Channel,
		&reservation.Client,
		&reservation.PaymentOption,
		&reservation.CancellationDate,
		&reservation.Segment,
		&reservation.Source,
		&reservation.Foct,
		&reservation.MetaGroupID,
	)
	if err != nil {
		return nil, err
	}

	// get Agent
	agent, err := s.getAgent(1) // Replace with correct Agent ID
	if err != nil {
		return nil, err
	}
	reservation.Creator = *agent

	// get CurrencyRates
	rates, err := s.getCurrencyRates(bookingID)
	if err != nil {
		return nil, err
	}
	reservation.CurrencyRates = rates

	// get Groups
	groups, err := s.getGroups(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	reservation.Groups = groups

	return reservation, nil
}

type SegmentReservation struct {
	// Define fields for the SegmentReservation struct here
}

type Service struct {
	// Define fields for the Service struct here
}

type VenueAutoAllocate struct {
	AutoAllocate bool `db:"AutoAllocate"`
}

type ReservationCommissionType struct {
	// Define fields for the ReservationCommissionType struct here
}

type Time struct {
	// Define fields for the Time struct here
}

func NewAllocatorService(
	ctx context.Context,
	logger log.Logger,
	db db.DB,
	bookingRepoFactory repo.BookingRepoFactory,
	bookingAdapterFactory adapter.BookingAdapterFactory,
) service.AllocatorService {
	return &allocatorService{
		logger:                logger.WithComponent(ctx, "ConfiguratorService"),
		db:                    db,
		bookingRepoFactory:    bookingRepoFactory,
		bookingAdapterFactory: bookingAdapterFactory,
	}
}

func (s *allocatorService) getVenueAutoAllocate(ctx context.Context, reservationID int) (bool, error) {
	//fmt.Println("test")
	var venueAutoAllocate VenueAutoAllocate
	//fmt.Printf("Value is: %d and type is reservationID: %T\\n", reservationID)
	query := `SELECT MAX(venues.AutoAllocate) as AutoAllocate
FROM bookings
JOIN booking_groups ON bookings.ID = booking_groups.BookingID
JOIN booking_items ON booking_groups.ID = booking_items.GroupID
JOIN venues ON booking_items.VenueID = venues.ID
WHERE bookings.ID = ?;`
	// Use QueryRow to fetch a single row result directly
	err := s.db.QueryRow(query, reservationID).Scan(
		&venueAutoAllocate.AutoAllocate,
	)

	//fmt.Printf("Value is: %d and type is AutoAllocatable: %T\\n", venueAutoAllocate.AutoAllocate)

	if err != nil || !venueAutoAllocate.AutoAllocate {
		if err == sql.ErrNoRows {
			fmt.Println("Error=>")
			return false, ErrUserNotFound
		}
	}

	return true, err

}

func (s *allocatorService) getAllocatableRooms(ctx context.Context, venueID int32, productEntity model.Product, startDate, endDate time.Time) ([]model.MetaObjects, error) {
	// Prepare the criteria for the ProductObject query
	var productIDs []string
	productIDs = append(productIDs, productEntity.ID)

	productObjectCriteria := ProductObjectCriteria{
		VenueID:     venueID,
		ProductIDs:  productIDs,
		PeriodStart: startDate,
		PeriodEnd:   endDate,
		PeriodType:  "allocatable",
	}

	//fmt.Printf("Value is: %d and type is productObjectCriteria: %T\\n", productObjectCriteria)

	// Fetch the ProductObjects from the database using the criteria
	productObjects, err := s.fetchAllocatableProductObjects(ctx, productIDs, productObjectCriteria, 628044)
	if err != nil {
		return nil, err
	}

	// Filter out the excluded rooms, if any
	var allocatableRooms []model.MetaObjects
	for _, room := range productObjects {
		allocatableRooms = append(allocatableRooms, room)
	}

	return allocatableRooms, nil
}

func (s *allocatorService) getAllocatedObject(itemID int) (bool, error) {
	var status string
	query := "SELECT status FROM booking_allocations WHERE BookingProductID = ? "

	err := s.db.QueryRow(query, itemID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Return empty status if no row found
		}
		return false, err
	}

	return true, nil
}
func (s *allocatorService) getAllocatedObjectStatus(bookingProductID int) (string, error) {
	var status string
	query := "SELECT Status FROM booking_allocations WHERE BookingProductID = ? LIMIT 1"

	err := s.db.QueryRow(query, bookingProductID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Return empty status if no row found
		}
		return "", err
	}

	return status, nil
}

func (s *allocatorService) AllocateAll(ctx context.Context, reservationIDs []int32, userID *int32) ([]model.AllocateResult, error) {
	logger := s.logger.WithMethod(ctx, "AllocateAll")
	//needToAmend := false
	var results []model.AllocateResult

	if reservationIDs == nil || len(reservationIDs) == 0 {
		return nil, errors.New("Invalid parameter reservationIDs")
	}

	user, err := s.getUserFromDatabase(userID)
	if err != nil {
		logger.Error("failed to get user from database", zap.Error(err))
		return nil, err
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		logger.Error("failed to marshal user to JSON", zap.Error(err))
		return nil, err
	}

	// Print the JSON to the console
	fmt.Println(string(userJSON))

	reservations, err := s.getReservations(ctx, reservationIDs)

	reservationsJSON, err := json.Marshal(reservations)
	if err != nil {
		logger.Error("failed to marshal user to JSON", zap.Error(err))
		return nil, err
	}
	fmt.Println(string(reservationsJSON))
	allocatedStatus := "allocated"
	//var results []AllocateResult
	for _, reservation := range reservations {
		//needToAmend := false
		//bookingBeforeAmend := reservation

		for _, group := range reservation.Groups {
			startDate := group.StartDate
			endDate := group.EndDate
			//fmt.Printf("Value is: %d and type is group.ID: %T\\n", group.ID)
			for _, item := range group.Items {
				//fmt.Printf("Value is: %d and type is item.ID: %T\\n", item.ID)

				alreadyAllocated, err := s.allocation(ctx, item, allocatedStatus, startDate, endDate)

				if err != nil && alreadyAllocated != nil {
					logger.Error("failed to marshal user to JSON", zap.Error(err))
				} else if alreadyAllocated == nil {
					allocateResult := model.AllocateResult{
						Status:        "unallocated",
						BookingID:     strconv.Itoa(reservation.ID),
						GroupID:       strconv.Itoa(int(group.ID)),
						ItemID:        strconv.Itoa(item.ID),
						AllocatedRoom: nil,
						Reason:        "No available rooms for this date period!",
					}
					fmt.Println("NoResults")
					results = append(results, allocateResult)
				} else {
					for _, obj := range alreadyAllocated {
						allocateResult := model.AllocateResult{
							Status:        "allocated",
							BookingID:     strconv.Itoa(reservation.ID),
							GroupID:       strconv.Itoa(int(group.ID)),
							ItemID:        strconv.Itoa(item.ID),
							AllocatedRoom: &model.AllocatedRoom{Id: obj.MetaObjectsID, Name: "room"},
						}
						fmt.Println("Results")
						results = append(results, allocateResult)

					}
					fmt.Println("Return")
					//return results, nil

				}
			}
		}

		/*		if needToAmend {
				if err := tx.Commit(); err != nil {
					logger.Error("failed to commit transaction", zap.Error(err))
					return nil, err
				}
			}*/
	}

	fmt.Println("Json8=>")

	return results, nil
}

type Company struct {
	ID               int32       `db:"id"`
	Name             null.String `db:"Name"`
	ShortName        null.String `db:"ShortName"`
	CountryCode      null.String `db:"CountryCode"`
	State            null.String `db:"State"`
	City             null.String `db:"City"`
	Address          null.String `db:"Address"`
	Address2         null.String `db:"Address2"`
	Zip              null.String `db:"Zip"`
	CompanyNumber    null.String `db:"CompanyNumber"`
	CompanyVatNumber null.String `db:"CompanyVatNumber"`
	Iban             null.String `db:"Iban"`
	Swift            null.String `db:"Swift"`
	OnAccount        null.Bool   `db:"OnAccount"`
	Preferences      null.String `db:"Preferences"`
	Comments         null.String `db:"Comments"`
	ContactPerson    null.String `db:"ContactPerson"`
	Status           null.String `db:"Status"`
}

type Comment struct {
	ID          int32       `db:"id"`
	EntityType  null.String `db:"entityType"`
	EntityID    null.Int32  `db:"entityID"`
	AgentID     null.Int32  `db:"agentID"`
	CommentType null.String `db:"commentType"`
	Comment     null.String `db:"comment"`
	PostedAt    []uint8     `db:"postedAt"`
	UpdatedAt   []uint8     `db:"updatedAt"`
}

type Account struct {
	ID            null.Int32  `db:"BookingID"`
	Name          null.String `db:"Name"`
	ContactPerson null.String `db:"ContactPerson"`
	ContactEmail  null.String `db:"ContactEmail"`
	CountryCode   null.String `db:"CountryCode"`
	PhoneNumber   null.String `db:"PhoneNumber"`
	Active        null.Bool   `db:"Active"`
	ReportBug     null.Bool   `db:"ReportBug"`
	Type          null.String `db:"Type"`
	TwoFactorAuth null.Bool   `db:"TwoFactorAuth"`
	AppVersion    null.String `db:"AppVersion"`
}

type CasbinRuleAgents struct {
	ID                 null.Int32  `db:"ID"`
	AccountID          null.Int32  `db:"AccountID"`
	Name               null.String `db:"Name"`
	Password           null.String `db:"Password"`
	AccessLevel        null.String `db:"AccessLevel"`
	Active             null.Bool   `db:"Active"`
	Locale             null.String `db:"Locale"`
	DefaultVenueID     null.Int32  `db:"DefaultVenueID"`
	FirstTimeLogin     null.Bool   `db:"FirstTimeLogin"`
	AllowedVenues      null.String `db:"AllowedVenues"`
	Email              null.String `db:"Email"`
	SendEmail          null.String `db:"SendEmail"`
	EmailNotifications null.String `db:"EmailNotifications"`
	FavoritePages      null.String `db:"FavoritePages"`
	Avatar             null.String `db:"Avatar"`
	Role               null.String `db:"Role"`
}

type BookingRemarks struct {
	BookingID null.Int32 `db:"BookingID"`
	Remark    null.Int32 `db:"Remark"`
}

type BookingChangesLog struct {
	ID            null.Int32  `db:"id"`
	ReservationID null.Int32  `db:"reservation_id"`
	Changes       null.String `db:"changes"`
	Agent         null.String `db:"agent"`
	Time          null.String `db:"time"`
}

type BookingCurrencyRates struct {
	BookingID null.Int32   `db:"BookingID"`
	Source    null.String  `db:"Source"`
	Target    null.String  `db:"Target"`
	Rate      null.Float32 `db:"Rate"`
	Date      []uint8      `db:"Date"`
	Final     null.Bool    `db:"Final"`
}

type BookingPax struct {
	GroupID  null.Int32 `db:"GroupID"`
	ClientID null.Int32 `db:"ClientID"`
}

type BookingGroupRatePlan struct {
	RatePlanID int32       `db:"RatePlanID"`
	GroupId    null.Int32  `db:"GroupId"`
	JsonData   null.String `db:"JsonData"`
}

type BookingItemRateDetails struct {
	ItemID          int32        `db:"id"`
	Date            []uint8      `db:"Date"`
	Amount          null.Float32 `db:"Amount"`
	AmountBeforeTax null.Float32 `db:"AmountBeforeTax"`
	Currency        null.String  `db:"Currency"`
	OriginalAmount  null.Float32 `db:"OriginalAmount"`
}

type ProductAffected struct {
	ID                int32       `db:"id"`
	ProductID         null.String `db:"product_id"`
	ClientID          null.Int32  `db:"client_id"`
	UsedAt            []uint8     `db:"used_at"`
	IsUsed            null.Int32  `db:"is_used"`
	GroupID           null.Int32  `db:"group_id"`
	UseAt             []uint8     `db:"use_at"`
	ReservationItemID null.Int32  `db:"reservation_item_id"`
}

type BookingItemCancellationPolicies struct {
	ItemID            int32        `db:"id"`
	DaysBeforeCheckIn null.Int32   `db:"DaysBeforeCheckIn"`
	PenaltyType       null.String  `db:"PenaltyType"`
	PenaltyValue      null.Float32 `db:"PenaltyValue"`
}

type BookingAllocationAndItems struct {
	ItemID       int32       `db:"id"`
	Status       null.String `db:"Status"`
	MetaObjectID null.String `db:"MetaObjectID"`
	StatusTimes  null.String `db:"StatusTimes"`
	LockedBy     null.Int32  `db:"LockedBy"`
}

func (s *allocatorService) getUserFromDatabase(userID *int32) (model.Agent, error) {
	query := "SELECT id, name, AccountID FROM agents WHERE id = ?"
	var agent model.Agent

	rows, err := s.db.Query(query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Agent{}, ErrUserNotFound
		}
		return model.Agent{}, err
	}
	defer rows.Close()

	if rows.Next() {
		// Scan the values from the row into the user variable
		err := rows.Scan(&agent.ID, &agent.Name, &agent.AccountID)
		if err != nil {
			return model.Agent{}, err
		}
	} else {
		return model.Agent{}, ErrUserNotFound
	}

	return agent, nil
}

type ProductObject struct {
	ID     string `db:"id"`
	RoomID int32  `db:"room_id"`
	// Другие поля объекта продукта
}

type HistorySaver interface {
	Save(history model.History) error
}

type historySaver struct {
	db db.DB
}

type ProductObjectCriteria struct {
	PeriodStart      time.Time
	PeriodEnd        time.Time
	ProductIDs       []string
	PeriodType       string
	VenueID          int32
	ViewingAccountID int
	Active           bool
	IDs              []int
}

func (s *allocatorService) AutoAllocate(ctx context.Context, reservationID int, isNotify bool) {
	//fmt.Printf("Value is: %d and type is reservationID: %T\\n", reservationID)

	logger := s.logger.WithMethod(ctx, "AllocateAll")
	venueAutoAllocate, err := s.getVenueAutoAllocate(ctx, reservationID)
	if err != nil {
		logger.Error("Error getting venueAutoAllocate:", zap.Error(err))
		//return nil, err
	}

	if venueAutoAllocate {
		s.autoAllocateReservation(ctx, reservationID, isNotify)
	}

}

func buildQuery(productObjectCriteria ProductObjectCriteria, metaObjectsList []string) (string, []interface{}) {
	var query strings.Builder
	var params []interface{}

	mysqlDateFormatPeriodStart := productObjectCriteria.PeriodStart.Format("2006-01-02")
	mysqlDateFormatPeriodEnd := productObjectCriteria.PeriodEnd.Format("2006-01-02")
	query.WriteString("SELECT DISTINCT MAX(po.ID) ")
	query.WriteString("FROM product_objects AS po ")
	query.WriteString("WHERE po.ID IN (")
	for i := range metaObjectsList {
		if i > 0 {
			query.WriteString(", ")
		}
		query.WriteString("?")
		params = append(params, metaObjectsList[i])
	}
	query.WriteString(") ")
	query.WriteString("AND NOT po.ID IN (SELECT DISTINCT pos.MetaObjectID ")
	query.WriteString("FROM product_object_statuses AS pos ")
	query.WriteString("WHERE pos.Status IN ('out_of_order', 'out_of_service') ")
	query.WriteString("AND Date BETWEEN DATE(%s) AND DATE(%s) - INTERVAL 1 DAY) ")
	//params = append(params, mysqlDateFormatPeriodStart, mysqlDateFormatPeriodEnd)
	query.WriteString("AND NOT po.ID IN (SELECT ba.MetaObjectID AS ID ")
	query.WriteString("FROM booking_groups AS bg ")
	query.WriteString("INNER JOIN booking_items AS bi ON bi.GroupID = bg.ID ")
	query.WriteString("LEFT JOIN booking_allocations AS ba ON bi.ID = ba.BookingProductID ")
	query.WriteString("WHERE bi.ProductID IN (")
	for i := range productObjectCriteria.ProductIDs {
		if i > 0 {
			query.WriteString(", ")
		}
		query.WriteString("?")
		params = append(params, productObjectCriteria.ProductIDs[i])
	}
	query.WriteString(") ")

	query.WriteString("AND bi.ProductType = 'room' ")

	query.WriteString("AND DATE(bg.EndDate) - INTERVAL 1 DAY >= DATE(%s) ")
	query.WriteString("AND DATE(bg.StartDate) <= DATE(%s) - INTERVAL 1 DAY);")

	//params = append(params, mysqlDateFormatPeriodStart, mysqlDateFormatPeriodEnd)
	resultQuery := fmt.Sprintf(query.String(), mysqlDateFormatPeriodStart, mysqlDateFormatPeriodEnd, mysqlDateFormatPeriodStart, mysqlDateFormatPeriodEnd)
	return resultQuery, params
}

func (s *allocatorService) fetchAllocatableProductObjects(ctx context.Context, bookingProductIDs []string, criteria ProductObjectCriteria, bookingProductID int) ([]model.MetaObjects, error) {

	logger := s.logger.WithMethod(ctx, "AllocateAll")
	//fmt.Printf("Value is: %d and type is hashCriteria: %T\\n", hashCriteria)
	mysqlDateFormatPeriodStart := criteria.PeriodStart.Format("2006-01-02")
	mysqlDateFormatPeriodEnd := criteria.PeriodEnd.Format("2006-01-02")
	productIds := bookingProductIDs
	productIdsPlaceholders := make([]string, len(productIds))
	for i := range productIdsPlaceholders {
		productIdsPlaceholders[i] = "?"
	}

	productObjectsQuery := fmt.Sprintf("SELECT DISTINCT po.ID FROM product_objects AS po INNER JOIN product_objects AS poActive ON po.ID = poActive.ID AND poActive.Key = 'active'  INNER JOIN product_objects AS poProductID ON po.ID = poProductID.ID AND poProductID.Key = 'product_id' INNER JOIN product_objects AS poRoomNumber ON po.ID = poRoomNumber.ID AND poRoomNumber.Key = 'name' INNER JOIN booking_allocations AS ba ON po.ID = ba.MetaObjectID WHERE  poActive.Value = '1' AND poProductID.Value IN  (%s) AND po.ID NOT IN (SELECT DISTINCT ba.MetaObjectID AS ID  FROM booking_groups AS bg INNER JOIN booking_items AS bi ON bi.GroupID = bg.ID LEFT JOIN booking_allocations AS ba ON bi.ID = ba.BookingProductID WHERE bi.ProductID IN (%s) AND ba.MetaObjectID IS NOT NULL AND bi.ProductType = 'room'  AND DATE(bg.EndDate) - INTERVAL 1 DAY >= DATE('%s') AND DATE(bg.StartDate) <= DATE('%s') - INTERVAL 1 DAY )", strings.Join(productIdsPlaceholders, ","), strings.Join(productIdsPlaceholders, ","), mysqlDateFormatPeriodStart, mysqlDateFormatPeriodEnd)

	var productObjectsInterfaceIDs []interface{}
	for _, id := range bookingProductIDs {
		productObjectsInterfaceIDs = append(productObjectsInterfaceIDs, id)
		productObjectsInterfaceIDs = append(productObjectsInterfaceIDs, id)
	}

	rows, err := s.db.Query(productObjectsQuery, productObjectsInterfaceIDs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var metaObjectsList []string
	for rows.Next() {
		var productObject ProductObject
		err := rows.Scan(
			&productObject.ID,
		)
		if err != nil {
			return nil, err
		}

		metaObjectsList = append(metaObjectsList, productObject.ID)
	}

	fmt.Println("Json109=>")
	metaObjectsListJson, err := json.Marshal(metaObjectsList)
	if err != nil {
		logger.Error("failed to marshal user to JSON", zap.Error(err))
	}

	fmt.Println(string(metaObjectsListJson))

	ids := metaObjectsList
	availableMetaObjectsPlaceholders := make([]string, len(ids))
	for i := range availableMetaObjectsPlaceholders {
		availableMetaObjectsPlaceholders[i] = "?"
	}
	availableMetaObjectsQuery, params := buildQuery(criteria, metaObjectsList)

	// Execute the query with the interfaceIDs as separate parameters
	availableMetaObjectsRows, err := s.db.Query(availableMetaObjectsQuery, params...)
	if err != nil {
		logger.Error("failed to marshal user to JSON", zap.Error(err))
	}
	defer availableMetaObjectsRows.Close()
	var allocatableProductObjects []model.MetaObjects
	for availableMetaObjectsRows.Next() {
		var availableMetaObjects model.MetaObjects
		err := availableMetaObjectsRows.Scan(
			&availableMetaObjects.MetaObjectsID,
		)
		if err != nil {
			return nil, err
		}

		allocatableProductObjects = append(allocatableProductObjects, availableMetaObjects)
	}

	//fmt.Println("Json7=>")
	allocatableProductObjectsJson, err := json.Marshal(allocatableProductObjects)
	if err != nil {
		logger.Error("failed to marshal user to JSON", zap.Error(err))
	}

	fmt.Println(string(allocatableProductObjectsJson))

	return allocatableProductObjects, nil
}

func (s *allocatorService) allocation(ctx context.Context, item model.BookingItems, allocatedStatus string, startDate time.Time, endDate time.Time) ([]model.MetaObjects, error) {
	logger := s.logger.WithMethod(ctx, "AllocateAll")
	var product model.Product
	//fmt.Printf("Value is: %d and type is item.Type: %T\\n", item.Type)

	//fmt.Printf("Value is: %d and type is item.Product.ProductType: %T\\n", item.Product.ProductType)
	//fmt.Printf("Value is: %d and type is ProductID: %T\\n", item.Product.ID)
	isAllocatedObject, err := s.getAllocatedObject(item.ID)
	if err != nil {
		logger.Error("failed on getting AllocatedObject", zap.Error(err))
	}
	if item.Type == "product" && item.Product.ProductType == "room" &&
		!isAllocatedObject {
		fmt.Println("Point")

		product = item.Product

		// Create the productObjectCriteria
		productObjectCriteria := ProductObjectCriteria{
			PeriodStart: startDate,
			PeriodEnd:   endDate,
			ProductIDs:  []string{product.ID}, // Assuming product.ID is int
			// ... set other criteria fields ...
		}

		// Fetch allocatable product objects using criteria
		allocatableProductObjects, err := s.fetchAllocatableProductObjects(ctx, []string{product.ID}, productObjectCriteria, item.ID)
		if err != nil {
			logger.Error("failed to marshal user to JSON", zap.Error(err))
		}
		if len(allocatableProductObjects) > 0 {
			fmt.Println("InCheck")
			err := s.updateAllocationStatus(ctx, item.ID, allocatedStatus, allocatableProductObjects)
			if err != nil {
				logger.Error("failed to marshal user to JSON", zap.Error(err))
			}

			/*if len(allocatableProductObjects) > 1 {
				return allocatableProductObjects[1:], nil
			} else{*/
			return allocatableProductObjects, err
			/*	}*/

		}

	} else {
		fmt.Println("Not Point")
	}

	return nil, err

}

func (s *allocatorService) autoAllocateReservation(ctx context.Context, reservationID int, isNotify bool) ([]model.MetaObjects, error) {
	logger := s.logger.WithMethod(ctx, "AllocateAll")
	reservationToEdit, err := s.getReservation(ctx, reservationID)
	if err != nil {
		// Handle the error
		logger.Error("failed to marshal user to JSON", zap.Error(err))
		return nil, err
	}

	allocatedStatus := "allocated"
	reservationToEditJSON, err := json.Marshal(reservationToEdit)
	if err != nil {
		logger.Error("failed to marshal user to JSON", zap.Error(err))
	}
	//fmt.Println("reservationToEditJSON")
	var allocatableProductObjects []model.MetaObjects
	fmt.Println(string(reservationToEditJSON))
	for _, group := range reservationToEdit.Groups {
		startDate := group.StartDate
		endDate := group.EndDate
		groupItemsJSON, err := json.Marshal(group.Items)
		if err != nil {
			logger.Error("failed to marshal user to JSON", zap.Error(err))
		}
		//fmt.Println("groupItemsJSON")
		fmt.Println(string(groupItemsJSON))

		for _, item := range group.Items {
			tempAllocated, _ := s.allocation(ctx, item, allocatedStatus, startDate, endDate)
			if err != nil {
				logger.Error("failed to marshal user to JSON", zap.Error(err))
			}
			for _, obj := range tempAllocated {
				allocatableProductObjects = append(allocatableProductObjects, obj)
			}
		}

		return allocatableProductObjects, nil
	}

	return nil, err

}

func (s *allocatorService) updateAllocationStatus(ctx context.Context, bookingProductID int, status string, productObjects []model.MetaObjects) error {
	for _, layout := range productObjects {
		var bookingProductIdUpdated int

		err := s.db.QueryRow("SELECT BookingProductID FROM booking_allocations WHERE MetaObjectID = ? AND BookingProductID = ?", layout.MetaObjectsID, bookingProductID).Scan(&bookingProductIdUpdated)
		if err != nil {
			_, err = s.db.Exec("INSERT INTO booking_allocations (BookingProductID,MetaObjectID, Status, StatusTimes, LockedBy) VALUES (?,?, ?, ?, ?)", bookingProductID, layout.MetaObjectsID, "allocated", "[]", nil)
			if err != nil {
				return fmt.Errorf("failed to update allocation status: %w", err)
			}
			fmt.Println("New row inserted successfully!")
		} else {
			_, err = s.db.Exec("UPDATE booking_allocations SET BookingProductID = ?, Status = ?, StatusTimes = ?, LockedBy = ? WHERE MetaObjectID = ?", bookingProductID, "allocated", "[]", nil, layout.MetaObjectsID)
			if err != nil {
				return fmt.Errorf("failed to update allocation status: %w", err)
			}
			fmt.Println("Row updated successfully!")
		}

	}

	return nil
}
