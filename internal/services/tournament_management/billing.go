package services

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/iRankHub/backend/internal/grpc/proto/tournament_management"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type BillingService struct {
	db *sql.DB
}

func NewBillingService(db *sql.DB) *BillingService {
	return &BillingService{db: db}
}

// Tournament Expenses Methods

func (s *BillingService) CreateTournamentExpenses(ctx context.Context, req *tournament_management.CreateExpensesRequest) (*tournament_management.ExpensesResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can manage expenses")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	queries := models.New(s.db)

	// Get tournament currency based on league type
	currency, err := queries.GetRegistrationCurrency(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament currency: %v", err)
	}

	expenses, err := queries.CreateTournamentExpenses(ctx, models.CreateTournamentExpensesParams{
		Tournamentid:      req.GetTournamentId(),
		Foodexpense:       float64ToString(req.GetFoodExpense()),
		Transportexpense:  float64ToString(req.GetTransportExpense()),
		Perdiemexpense:    float64ToString(req.GetPerDiemExpense()),
		Awardingexpense:   float64ToString(req.GetAwardingExpense()),
		Stationaryexpense: float64ToString(req.GetStationaryExpense()),
		Otherexpenses:     float64ToString(req.GetOtherExpenses()),
		Currency:          currency,
		Notes:             sql.NullString{String: req.GetNotes(), Valid: req.Notes != ""},
		Createdby:         sql.NullInt32{Int32: int32(userID), Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create expenses: %v", err)
	}

	return expensesToProto(expenses), nil
}

func (s *BillingService) UpdateTournamentExpenses(ctx context.Context, req *tournament_management.UpdateExpensesRequest) (*tournament_management.ExpensesResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can update expenses")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	queries := models.New(s.db)

	expenses, err := queries.UpdateTournamentExpenses(ctx, models.UpdateTournamentExpensesParams{
		Tournamentid:      req.GetTournamentId(),
		Foodexpense:       float64ToString(req.GetFoodExpense()),
		Transportexpense:  float64ToString(req.GetTransportExpense()),
		Perdiemexpense:    float64ToString(req.GetPerDiemExpense()),
		Awardingexpense:   float64ToString(req.GetAwardingExpense()),
		Stationaryexpense: float64ToString(req.GetStationaryExpense()),
		Otherexpenses:     float64ToString(req.GetOtherExpenses()),
		Currency:          req.GetCurrency(),
		Notes:             sql.NullString{String: req.GetNotes(), Valid: req.Notes != ""},
		Updatedby:         sql.NullInt32{Int32: int32(userID), Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update expenses: %v", err)
	}

	return expensesToProto(expenses), nil
}

func (s *BillingService) GetTournamentExpenses(ctx context.Context, req *tournament_management.GetExpensesRequest) (*tournament_management.ExpensesResponse, error) {
    if err := validateAuthentication(req.GetToken()); err != nil {
        return nil, err
    }

    queries := models.New(s.db)

    expenses, err := queries.GetTournamentExpenses(ctx, req.GetTournamentId())
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("no expenses found for tournament ID %d", req.GetTournamentId())
        }
        return nil, fmt.Errorf("failed to get tournament expenses: %v", err)
    }

    return expensesToProto(expenses), nil
}

// School Registration Methods

func (s *BillingService) CreateSchoolRegistration(ctx context.Context, req *tournament_management.CreateRegistrationRequest) (*tournament_management.RegistrationResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	queries := models.New(s.db)

	// Get tournament currency based on league type
	currency, err := queries.GetRegistrationCurrency(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament currency: %v", err)
	}

	// Get tournament fee
	tournament, err := queries.GetTournamentByID(ctx, req.GetTournamentId())
	if err != nil {
		return nil, fmt.Errorf("failed to get tournament: %v", err)
	}

	registration, err := queries.CreateSchoolRegistration(ctx, models.CreateSchoolRegistrationParams{
		Schoolid:          req.GetSchoolId(),
		Tournamentid:      req.GetTournamentId(),
		Plannedteamscount: req.GetPlannedTeamsCount(),
		Amountperteam:     tournament.Tournamentfee,
		Currency:          currency,
		Createdby:         sql.NullInt32{Int32: int32(userID), Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create registration: %v", err)
	}

	return registrationToProto(registration), nil
}

func (s *BillingService) UpdateSchoolRegistration(ctx context.Context, req *tournament_management.UpdateRegistrationRequest) (*tournament_management.RegistrationResponse, error) {
	claims, err := utils.ValidateToken(req.GetToken())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can update payment status")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	queries := models.New(s.db)

	registration, err := queries.UpdateSchoolRegistration(ctx, models.UpdateSchoolRegistrationParams{
		Schoolid:         req.GetSchoolId(),
		Tournamentid:     req.GetTournamentId(),
		Actualteamscount: sql.NullInt32{Int32: req.GetActualTeamsCount(), Valid: true},
		Discountamount:   sql.NullString{String: float64ToString(req.GetDiscountAmount()), Valid: true},
		Actualpaidamount: sql.NullString{String: float64ToString(req.GetActualPaidAmount()), Valid: true},
		Paymentstatus:    req.GetPaymentStatus(),
		Updatedby:        sql.NullInt32{Int32: int32(userID), Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update registration: %v", err)
	}

	return registrationToProto(registration), nil
}

func (s *BillingService) GetSchoolRegistration(ctx context.Context, req *tournament_management.GetRegistrationRequest) (*tournament_management.DetailedRegistrationResponse, error) {
	if err := validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)

	registration, err := queries.GetSchoolRegistration(ctx, models.GetSchoolRegistrationParams{
		Schoolid:     req.GetSchoolId(),
		Tournamentid: req.GetTournamentId(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get registration: %v", err)
	}

	return &tournament_management.DetailedRegistrationResponse{
		RegistrationId:    registration.Registrationid,
		SchoolId:          registration.Schoolid,
		TournamentId:      registration.Tournamentid,
		SchoolName:        registration.Schoolname,
		SchoolEmail:       registration.Schoolemail,
		SchoolType:        registration.Schooltype,
		ContactEmail:      registration.Contactemail,
		ContactPersonName: registration.Contactpersonname,
		Country:           registration.Country.String,
		Province:          registration.Province.String,
		District:          registration.District.String,
		Address:           registration.Address,
		PlannedTeamsCount: registration.Plannedteamscount,
		ActualTeamsCount:  int32(registration.Actualteamscount.Int32),
		AmountPerTeam:     stringToFloat64(registration.Amountperteam),
		TotalAmount:       nullStringToFloat64(registration.Totalamount),
		DiscountAmount:    nullStringToFloat64(registration.Discountamount),
		ActualPaidAmount:  nullStringToFloat64(registration.Actualpaidamount),
		PaymentStatus:     registration.Paymentstatus,
		Currency:          registration.Currency,
	}, nil
}

func (s *BillingService) ListTournamentRegistrations(ctx context.Context, req *tournament_management.ListRegistrationsRequest) (*tournament_management.ListRegistrationsResponse, error) {
	if err := validateAuthentication(req.GetToken()); err != nil {
		return nil, err
	}

	queries := models.New(s.db)

	registrations, err := queries.ListTournamentRegistrations(ctx, models.ListTournamentRegistrationsParams{
		Tournamentid: req.GetTournamentId(),
		Limit:        req.GetPageSize(),
		Offset:       req.GetPageSize() * req.GetPageToken(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list registrations: %v", err)
	}

	response := &tournament_management.ListRegistrationsResponse{
		Registrations: make([]*tournament_management.ListRegistrationItem, len(registrations)),
		NextPageToken: req.GetPageToken() + 1,
	}

	for i, reg := range registrations {
		response.Registrations[i] = &tournament_management.ListRegistrationItem{
			RegistrationId:    reg.Registrationid,
			IDebateSchoolId:   reg.Idebateschoolid.String,
			SchoolName:        reg.Schoolname,
			SchoolEmail:       reg.Schoolemail,
			PlannedTeamsCount: reg.Plannedteamscount,
			ActualTeamsCount:  int32(reg.Actualteamscount.Int32),
			TotalAmount:       nullStringToFloat64(reg.Totalamount),
			PaymentStatus:     reg.Paymentstatus,
			Currency:          reg.Currency,
		}
	}

	return response, nil
}

// Helper functions for converting between model and proto types
func expensesToProto(e models.Tournamentexpense) *tournament_management.ExpensesResponse {
	return &tournament_management.ExpensesResponse{
		ExpenseId:         e.Expenseid,
		TournamentId:      e.Tournamentid,
		FoodExpense:       stringToFloat64(e.Foodexpense),
		TransportExpense:  stringToFloat64(e.Transportexpense),
		PerDiemExpense:    stringToFloat64(e.Perdiemexpense),
		AwardingExpense:   stringToFloat64(e.Awardingexpense),
		StationaryExpense: stringToFloat64(e.Stationaryexpense),
		OtherExpenses:     stringToFloat64(e.Otherexpenses),
		TotalExpense:      nullStringToFloat64(e.Totalexpense),
		Currency:          e.Currency,
		Notes:             e.Notes.String,
		CreatedAt:         nullTimeToString(e.Createdat),
		UpdatedAt:         nullTimeToString(e.Updatedat),
	}
}

func registrationToProto(r models.Schooltournamentregistration) *tournament_management.RegistrationResponse {
	return &tournament_management.RegistrationResponse{
		RegistrationId:    r.Registrationid,
		SchoolId:          r.Schoolid,
		TournamentId:      r.Tournamentid,
		PlannedTeamsCount: r.Plannedteamscount,
		ActualTeamsCount:  int32(r.Actualteamscount.Int32),
		AmountPerTeam:     stringToFloat64(r.Amountperteam),
		TotalAmount:       nullStringToFloat64(r.Totalamount),
		DiscountAmount:    nullStringToFloat64(r.Discountamount),
		ActualPaidAmount:  nullStringToFloat64(r.Actualpaidamount),
		PaymentStatus:     r.Paymentstatus,
		Currency:          r.Currency,
		CreatedAt:         nullTimeToString(r.Createdat),
		UpdatedAt:         nullTimeToString(r.Updatedat),
	}
}

// Helper function to convert float64 to string with 2 decimal places
func float64ToString(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

// Helper function to convert string to float64
func stringToFloat64(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// Helper function to handle sql.NullString to float64 conversion
func nullStringToFloat64(s sql.NullString) float64 {
	if !s.Valid {
		return 0
	}
	return stringToFloat64(s.String)
}

// Helper function to handle sql.NullTime to string conversion
func nullTimeToString(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(time.RFC3339)
}

func validateAuthentication(token string) error {
	_, err := utils.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}
	return nil
}
