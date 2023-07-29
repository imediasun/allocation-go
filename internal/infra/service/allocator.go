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
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/repo"
	"gitlab.hotel.tools/backend-team/allocation-go/internal/domain/service"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

//460769 mo.63a975e4784d1 p.63a975e473fe7

//Должен быть
//461054 p.57eb86aecace1

//Пробуем mo.63d112187cd45 сделать available
type AllocateResult struct {
	Status        string
	BookingID     int
	GroupID       int32
	ItemID        []uint8
	AllocatedRoom ProductObject
	Reason        string
}

type History struct {
	User   Agent
	Action string
	Before Reservation
	After  Reservation
}

type ActionAbstractUpdate struct{}

type allocatorService struct {
	ctx                   context.Context
	reservation           Reservation
	logger                log.Logger
	db                    db.DB
	bookingRepoFactory    repo.BookingRepoFactory
	bookingAdapterFactory adapter.BookingAdapterFactory
}

type Money struct {
	Amount   float64
	Currency string
}

type ReservationStatus string

type ReservationPaymentOption string

type ReservationGroup struct {
	Item           BookingItems
	ID             int32
	BookingID      int32
	PaxNationality string
	StartDate      time.Time
	EndDate        time.Time
	ParentID       null.Int64
	Items          []BookingItems
}

type CurrencyRate struct {
	BookingID int64
	Source    string
	Target    string
	Rate      float64
	Date      time.Time
	Final     bool
}

type Reservation struct {
	ID                  int
	Creator             Agent
	Price               Money
	CreationDate        time.Time
	Status              ReservationStatus
	ProviderReference   null.String
	Channel             null.String
	Remark              string
	Client              null.String
	Manual              bool
	PaymentOption       null.String
	Groups              []ReservationGroup
	CancellationDate    []uint8
	StartDate           []uint8
	EndDate             []uint8
	Segment             null.String
	Source              null.String
	Logs                interface{} // Replace with actual type
	CurrencyRates       []CurrencyRate
	Foct                bool
	IsCityTaxToProvider bool
	MetaGroupID         null.Int64
	Customer            *Client
}

func (s *allocatorService) getAgent(id int32) (*Agent, error) {
	fmt.Printf("Value is: %d and type is Agent: %T\\n", id)
	row := s.db.QueryRow("SELECT id, name, AccountID FROM agents WHERE id = ?", id)

	agent := &Agent{}
	err := row.Scan(&agent.ID, &agent.Name, &agent.AccountID)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func (s *allocatorService) getCurrencyRates(bookingID int) ([]CurrencyRate, error) {
	rows, err := s.db.Query("SELECT BookingID, Source, Target, Rate, Date, Final FROM booking_currency_rates WHERE BookingID = ?", bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []CurrencyRate

	for rows.Next() {
		var rate CurrencyRate

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

func (s *allocatorService) getItemsForGroup(ctx context.Context, groupID int32) ([]BookingItems, error) {
	logger := s.logger.WithMethod(ctx, "AllocateAll")
	rows, err := s.db.Query("SELECT ID,Type,VenueID,ProductID,Status FROM booking_items WHERE GroupID = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]BookingItems, 0)
	var count int
	for rows.Next() {
		count++
		var item BookingItems

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

func (s *allocatorService) getGroups(ctx context.Context, bookingID int) ([]ReservationGroup, error) {

	logger := s.logger.WithMethod(ctx, "AllocateAll")

	rows, err := s.db.Query("SELECT * FROM booking_groups WHERE BookingID = ?", bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]ReservationGroup, 0)
	var count int
	for rows.Next() {
		count++
		group := ReservationGroup{}
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

func (s *allocatorService) getProduct(productID string) (Product, error) {

	fmt.Printf("Value is: %d and type is PProductID: %T\\n", productID)
	row := s.db.QueryRow("SELECT ID, Status, product_type FROM products WHERE ID = ?", productID)

	var product Product
	err := row.Scan(&product.ID, &product.Status, &product.ProductType) // and so on for all fields in Product
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func (s *allocatorService) getReservations(ctx context.Context, bookingIDs []int32) ([]Reservation, error) {
	var reservationList []Reservation
	fmt.Printf("Value is: %d and type is BookingID: %T\\n", bookingIDs)
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
			return []Reservation{}, ErrUserNotFound
		}
		return []Reservation{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var reservation Reservation
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

func (s *allocatorService) getReservation(ctx context.Context, bookingID int) (*Reservation, error) {

	fmt.Printf("Value is: %d and type is BookingID: %T\\n", bookingID)
	row := s.db.QueryRow("SELECT ID,ProviderReference,Channel,Client,PaymentOption,CancellationDate,Segment,Source,Foct,MetaGroupID	FROM bookings WHERE ID = ?", bookingID)

	reservation := &Reservation{}
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

type Client struct {
	ID             null.Int32  `db:"ID"`
	AccountID      null.Int32  `db:"AccountID"`
	Email          null.String `db:"Email"`
	Phone          null.String `db:"Phone"`
	Title          null.String `db:"Title"`
	Gender         null.String `db:"Gender"`
	Nationality    null.String `db:"Nationality"`
	LanguageID     null.Int32  `db:"LanguageID"`
	Identification null.String `db:"Identification"`
	LastName       null.String `db:"LastName"`
	BirthDate      []uint8     `db:"BirthDate"`
	Address        null.String `db:"Address"`
	AdditionalInfo null.String `db:"AdditionalInfo"`
	AgentID        null.Int32  `db:"AgentID"`
	CreatedAt      []uint8     `db:"CreatedAt"`
	Status         null.String `db:"Status"`
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
    fmt.Println("test")
	var venueAutoAllocate VenueAutoAllocate
	fmt.Printf("Value is: %d and type is reservationID: %T\\n", reservationID)
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

	fmt.Printf("Value is: %d and type is AutoAllocatable: %T\\n", venueAutoAllocate.AutoAllocate)

	if err != nil || !venueAutoAllocate.AutoAllocate {
		if err == sql.ErrNoRows {
			fmt.Println("Error=>")
			return false, ErrUserNotFound
		}
	}

	return true, err

}

func (s *allocatorService) getAllocatableRooms(ctx context.Context, venueID int32, productEntity Product, startDate, endDate time.Time) ([]ProductObject, error) {
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

	fmt.Printf("Value is: %d and type is productObjectCriteria: %T\\n", productObjectCriteria)

	// Fetch the ProductObjects from the database using the criteria
	productObjects, err := s.fetchAllocatableProductObjects(ctx, productObjectCriteria)
	if err != nil {
		return nil, err
	}

	// Filter out the excluded rooms, if any
	var allocatableRooms []ProductObject
	for _, room := range productObjects {
		allocatableRooms = append(allocatableRooms, room)
	}

	return allocatableRooms, nil
}

func (s *allocatorService) getAllocatedObject(itemID []uint8) (bool, error) {
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

func (s *allocatorService) AllocateAll(ctx context.Context, reservationIDs []int32, userID *int32) ([]byte, error) {
	logger := s.logger.WithMethod(ctx, "AllocateAll")
	//needToAmend := false
	var results []AllocateResult

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

	//var results []AllocateResult
	for _, reservation := range reservations {
		needToAmend := false
		//bookingBeforeAmend := reservation

		tx, err := s.db.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer tx.Rollback(ctx)

		for _, group := range reservation.Groups {
			for _, item := range group.Items {
				isAllocatedObject, err := s.getAllocatedObject(item.ID)
				if err != nil {
					return nil, err
				}
				if item.Product.ProductType == "room" &&
					item.Status == "confirmed" &&
					!isAllocatedObject {
					fmt.Printf("getAllocatableRooms")
					roomsByProduct, err := s.getAllocatableRooms(ctx, item.VenueID, item.Product, group.StartDate, group.EndDate)
					if err != nil {
						return nil, err
					}

					if len(roomsByProduct) == 0 {
						results = append(results, AllocateResult{
							Status:    "unallocated",
							BookingID: reservation.ID,
							GroupID:   group.ID,
							ItemID:    item.ID,
							Reason:    "No available rooms for this date period!",
						})
						continue
					}

					room := roomsByProduct[0]
					//roomsByProduct = roomsByProduct[1:]
					//needToAmend, err = s.allocateRoom(item.Product.ID)
					fmt.Println("Json095=>")
					roomsByProductJson, err := json.Marshal(roomsByProduct)
					if err != nil {
						logger.Error("failed to marshal user to JSON", zap.Error(err))
					}

					fmt.Println(string(roomsByProductJson))
					err = s.updateAllocationStatus(ctx, reservation.ID, item.Product.ID, "allocated", roomsByProduct)
					if err != nil {
						return nil, err
					}

					results = append(results, AllocateResult{
						Status:        "allocated",
						BookingID:     reservation.ID,
						GroupID:       group.ID,
						ItemID:        item.ID,
						AllocatedRoom: room,
					})

					/*					history := History{
											User:   user,
											Action: "update",
											Before: bookingBeforeAmend, // клон объекта reservation
											After:  reservation,
										}
										s.historySaver.Save(history)*/
				}
			}
		}

		if needToAmend {
			if err := tx.Commit(); err != nil {
				logger.Error("failed to commit transaction", zap.Error(err))
				return nil, err
			}
		}
	}

	fmt.Println("Json8=>")
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		logger.Error("failed to marshal user to JSON", zap.Error(err))
		return nil, err
	}

	fmt.Println(string(resultsJSON))

	return resultsJSON, nil
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

type BookingItems struct {
	ID        []uint8     `db:"id"`
	Type      string      `db:"Type"`
	VenueID   int32       `db:"VenueID"`
	ProductID null.String `db:"ProductID"`
	Status    string      `db:"Status"`
	Product   Product
}

type BookingGroups struct {
	ID             int32       `db:"ID"`
	BookingID      int32       `db:"BookingID"`
	PaxNationality string      `db:"PaxNationality"`
	StartDate      []uint8     `db:"StartDate"`
	EndDate        []uint8     `db:"EndDate"`
	ParentID       null.String `db:"ParentID"`
	Remark         null.String `db:"Remark"`
}

type Agent struct {
	ID        int32  `db:"id"`
	Name      string `db:"name"`
	AccountID int32  `db:"accountID"`
}

func (s *allocatorService) getUserFromDatabase(userID *int32) (Agent, error) {
	query := "SELECT id, name, AccountID FROM agents WHERE id = ?"
	var agent Agent

	rows, err := s.db.Query(query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return Agent{}, ErrUserNotFound
		}
		return Agent{}, err
	}
	defer rows.Close()

	if rows.Next() {
		// Scan the values from the row into the user variable
		err := rows.Scan(&agent.ID, &agent.Name, &agent.AccountID)
		if err != nil {
			return Agent{}, err
		}
	} else {
		return Agent{}, ErrUserNotFound
	}

	return agent, nil
}

type Product struct {
	ID          string
	Status      string
	ProductType string `db:"product_type"`
}

type ProductObject struct {
	ID     string `db:"id"`
	RoomID int32  `db:"room_id"`
	// Другие поля объекта продукта
}

type HistorySaver interface {
	Save(history History) error
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

func joinInts(ints []int) string {
	var strInts []string
	for _, i := range ints {
		strInts = append(strInts, strconv.Itoa(i))
	}
	return strings.Join(strInts, "|")
}

func convertUint8ToInt(int32Slice [][]uint8) [][]int {
	intSlice := make([][]int, len(int32Slice))
	for i, innerSlice := range int32Slice {
		intSlice[i] = make([]int, len(innerSlice))
		for j, val := range innerSlice {
			intSlice[i][j] = int(val)
		}
	}
	return intSlice
}

func joinTimes(t1, t2 []uint8) string {
	// Convert []uint8 to strings
	t1Str := string(t1)
	t2Str := string(t2)

	// Parse the strings into time.Time objects
	t1Time, err := time.Parse(time.RFC3339, t1Str)
	if err != nil {
		// Handle the error if necessary
	}
	t2Time, err := time.Parse(time.RFC3339, t2Str)
	if err != nil {
		// Handle the error if necessary
	}

	// Format the time objects and return the result
	return fmt.Sprintf("%s|%s", t1Time.Format(time.RFC3339), t2Time.Format(time.RFC3339))
}

func joinArrayInts(ints [][]int) string {
	var strInts []string
	for _, innerInts := range ints {
		strInnerInts := make([]string, len(innerInts))
		for i, val := range innerInts {
			strInnerInts[i] = strconv.Itoa(val)
		}
		strInts = append(strInts, strings.Join(strInnerInts, "|"))
	}
	return strings.Join(strInts, "|")
}

/*func (c *ProductObjectCriteria) Hash() string {
	result := []string{
		joinArrayInts(convertUint8ToInt(c.ProductIDs)),
		joinTimes(c.PeriodStart, c.PeriodEnd),
		strconv.Itoa(int(c.VenueID)),
		strconv.Itoa(c.ViewingAccountID),
		strconv.FormatBool(c.Active),
		joinInts(c.IDs),
		c.PeriodType,
	}

	for i, v := range result {
		if v == "" {
			result[i] = "\x00" // Change empty string to binary zero
		}
	}
	return strings.Join(result, "|")
}*/

func (s *allocatorService) AutoAllocate(ctx context.Context, agentID *int32, reservationID int, isNotify bool) {
	fmt.Printf("Value is: %d and type is reservationID: %T\\n", reservationID)

	logger := s.logger.WithMethod(ctx, "AllocateAll")
	venueAutoAllocate, err := s.getVenueAutoAllocate(ctx, reservationID)
	if err != nil {
		logger.Error("Error getting venueAutoAllocate:", zap.Error(err))
		//return nil, err
	}
	fmt.Println("venueAutoAllocate")
	fmt.Println(venueAutoAllocate)
	if venueAutoAllocate {
		fmt.Println("getVenueAutoAllocate==true")
		s.autoAllocateReservation(ctx, reservationID, isNotify)
	}

}

func buildQuery(productObjectCriteria ProductObjectCriteria) (string, []interface{}) {
	// Create placeholders for ProductIDs
	placeholders := make([]string, len(productObjectCriteria.ProductIDs))
	values := make([]interface{}, len(productObjectCriteria.ProductIDs))

	for i, id := range productObjectCriteria.ProductIDs {
		placeholders[i] = "?"
		values[i] = id
		fmt.Printf("Value is: %d and type is 44: %T\\n", id)
	}

	// Create the WHERE clause for ProductIDs
	whereProductIDs := fmt.Sprintf("`Key`='product_id' AND `Value` IN (%s)", strings.Join(placeholders, ","))

	// Create the WHERE clause for Date range
	whereDateRange := fmt.Sprintf("pos.Date BETWEEN ? AND ?")

	// Combine all the WHERE clauses
	where := fmt.Sprintf("%s AND %s AND pos.Status = 'available'", whereProductIDs, whereDateRange)

	// Create the full SQL query
	query := fmt.Sprintf("SELECT DISTINCT ID FROM product_objects AS po LEFT JOIN product_object_statuses AS pos ON pos.MetaObjectID = po.ID WHERE %s", where)

	// Create the list of parameters for the SQL query
	params := append(values, productObjectCriteria.PeriodStart, productObjectCriteria.PeriodEnd)

	return query, params
}

func (s *allocatorService) fetchAllocatableProductObjects(ctx context.Context, criteria ProductObjectCriteria) ([]ProductObject, error) {
	//hashCriteria := criteria.Hash()
	logger := s.logger.WithMethod(ctx, "AllocateAll")
	//fmt.Printf("Value is: %d and type is hashCriteria: %T\\n", hashCriteria)
	query, params := buildQuery(criteria)
	fmt.Println(query)

	rows, err := s.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allocatableProductObjects []ProductObject
	for rows.Next() {
		var productObject ProductObject

		err := rows.Scan(&productObject.ID)
		if err != nil {
			return nil, err
		}

		allocatableProductObjects = append(allocatableProductObjects, productObject)
	}

	fmt.Println("Json7=>")
	allocatableProductObjectsJson, err := json.Marshal(allocatableProductObjects)
	if err != nil {
		logger.Error("failed to marshal user to JSON", zap.Error(err))
	}

	fmt.Println(string(allocatableProductObjectsJson))

	return allocatableProductObjects, nil
}

func (s *allocatorService) autoAllocateReservation(ctx context.Context, reservationID int, isNotify bool) {
	logger := s.logger.WithMethod(ctx, "AllocateAll")
	fmt.Printf("Value is: %d and type is ReservationID: %T\\n", reservationID)
	reservationToEdit, err := s.getReservation(ctx, reservationID)
	if err != nil {
		// Handle the error
		fmt.Println("Error fetching reservation:", err)
		return
	}

	allocatedStatus := "allocated"
	reservationToEditJSON, err := json.Marshal(reservationToEdit)
	if err != nil {
		logger.Error("failed to marshal user to JSON", zap.Error(err))
	}
	fmt.Println("reservationToEditJSON")
	fmt.Println(string(reservationToEditJSON))
	for _, group := range reservationToEdit.Groups {
		startDate := group.StartDate
		endDate := group.EndDate
		groupItemsJSON, err := json.Marshal(group.Items)
		if err != nil {
			logger.Error("failed to marshal user to JSON", zap.Error(err))
		}
		fmt.Println("groupItemsJSON")
		fmt.Println(string(groupItemsJSON))
		var product Product
		for _, item := range group.Items {
			fmt.Printf("Value is: %d and type is item.Type: %T\\n", item.Type)

			fmt.Printf("Value is: %d and type is item.Product.ProductType: %T\\n", item.Product.ProductType)
			fmt.Printf("Value is: %d and type is ProductID: %T\\n", item.Product.ID)
			//Здесь ProductID должен прилетать ввиде строки типа p.56fa2a0e51daf
			if item.Type == "product" && item.Product.ProductType == "room" {
				fmt.Println("Point")

				product = item.Product

				// Create the productObjectCriteria
				productObjectCriteria := ProductObjectCriteria{
					PeriodStart: startDate,
					PeriodEnd:   endDate,
					ProductIDs:  []string{product.ID}, // Assuming product.ID is int
					// ... set other criteria fields ...
				}

				fmt.Printf("Value is: %d and type is productObjectCriteria2: %T\\n", productObjectCriteria)

				// Fetch allocatable product objects using criteria
				allocatableProductObjects, err := s.fetchAllocatableProductObjects(ctx, productObjectCriteria)
				if err != nil {
					// Handle the error
				}
				fmt.Println("beforeCheck")

				fmt.Println(len(allocatableProductObjects))
				if len(allocatableProductObjects) > 0 {
					fmt.Println("InCheck")
					err := s.updateAllocationStatus(ctx, reservationID, item.Product.ID, allocatedStatus, allocatableProductObjects)
					if err != nil {
						// Handle the error
					}
					allocatableProductObjects = allocatableProductObjects[1:]
				}
			}
		}
	}
}

/*func (s *allocatorService) allocateRoom(bookingProductID string) (bool, error) {
	// Construct the SQL query
	var productObject ProductObject
	metaObjectQuery := "SELECT * FROM product_objects as po WHERE po.Key='product_id' AND Value= ? "
	err := s.db.QueryRow(metaObjectQuery, bookingProductID).Scan(
		&productObject.ID,
	)

	if err != nil {
		return false, fmt.Errorf("failed to get metObjct: %w", err)
	}

	query := "INSERT INTO booking_allocations (BookingProductID, MetaObjectID, Status, StatusTimes, LockedBy) VALUES (?, ?, 'allocated', '[]', NULL) ON DUPLICATE KEY UPDATE MetaObjectID = VALUES(`MetaObjectID`), Status = VALUES(`Status`), StatusTimes = VALUES(`StatusTimes`), LockedBy = VALUES(`LockedBy`)"

	// Execute the SQL query with the provided parameters
	_, err = s.db.Exec(query, bookingProductID, productObject.ID)
	if err != nil {
		return false, fmt.Errorf("failed to update allocation status: %w", err)
	}

	return true, nil
}*/

func convertUint8ToInt32(uint8Slice []uint8) []int32 {
	int32Slice := make([]int32, len(uint8Slice))
	for i, val := range uint8Slice {
		int32Slice[i] = int32(val)
	}
	return int32Slice
}

func (s *allocatorService) updateAllocationStatus(ctx context.Context, reservationID int, bookingProductID string, status string, productObjects []ProductObject) error {

	logger := s.logger.WithMethod(ctx, "AllocateAll")
	fmt.Println("updateAllocationStatus")
	fmt.Printf("Value is: %d and type is 33: %T\\n", bookingProductID)
	bookingProductIDResults := bookingProductID
	fmt.Println("Json098=>")
	bookingProductIDJson, err := json.Marshal(bookingProductIDResults)
	if err != nil {
		logger.Error("failed to marshal user to JSON", zap.Error(err))
	}

	fmt.Println(string(bookingProductIDJson))

	fmt.Println("Json099=>")
	productObjectsJson, err := json.Marshal(productObjects)
	if err != nil {
		logger.Error("failed to marshal user to JSON", zap.Error(err))
	}
	fmt.Println(string(productObjectsJson))
	for _, layout := range productObjects {
		var bookingProductID int32
		err = s.db.QueryRow("SELECT BookingProductID FROM booking_allocations WHERE MetaObjectID = ?", layout.ID).Scan(&bookingProductID)
		if err != nil {
			// If the row doesn't exist, insert a new row with the provided data
			_, err = s.db.Exec("INSERT INTO booking_allocations (BookingProductID,MetaObjectID, Status, StatusTimes, LockedBy) VALUES (?,?, ?, ?, ?)", reservationID, layout.ID, "allocated", "[]", nil)
			if err != nil {
				return fmt.Errorf("failed to update allocation status: %w", err)
			}
			fmt.Println("New row inserted successfully!")
		} else {
			// If the row already exists, update it with the provided data
			_, err = s.db.Exec("UPDATE booking_allocations SET Status = ?, StatusTimes = ?, LockedBy = ? WHERE MetaObjectID = ?", "allocated", "[]", nil, layout.ID)
			if err != nil {
				return fmt.Errorf("failed to update allocation status: %w", err)
			}
			fmt.Println("Row updated successfully!")
		}

	}

	return nil
}
