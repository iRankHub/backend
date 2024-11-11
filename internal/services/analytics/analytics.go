package services

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/iRankHub/backend/internal/grpc/proto/analytics"
	"github.com/iRankHub/backend/internal/models"
	"github.com/iRankHub/backend/internal/utils"
)

type AnalyticsService struct {
	db *sql.DB
}

func NewAnalyticsService(db *sql.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

func (s *AnalyticsService) GetFinancialReports(ctx context.Context, req *analytics.FinancialReportRequest) (*analytics.FinancialReportResponse, error) {
	// Validate admin access
	claims, err := utils.ValidateToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can access financial reports")
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.DateRange.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", req.DateRange.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %v", err)
	}

	// Convert tournament ID
	var tournamentID int32
	if req.TournamentId != nil && *req.TournamentId != "" {
		id, err := strconv.ParseInt(*req.TournamentId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid tournament ID: %v", err)
		}
		tournamentID = int32(id)
	}

	queries := models.New(s.db)
	response := &analytics.FinancialReportResponse{}
	if req.ReportType != nil {
		response.ReportType = *req.ReportType
	}

	// Handle report types
	if req.ReportType == nil {
		return nil, fmt.Errorf("report type is required")
	}

	switch *req.ReportType {
	case "income_overview":
		incomes, err := queries.GetTournamentIncomeOverview(ctx, models.GetTournamentIncomeOverviewParams{
			Startdate:   startDate,
			Startdate_2: endDate,
			Column3:     tournamentID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get income overview: %v", err)
		}

		for _, income := range incomes {
			var leagueID string
			var leagueName string
			if income.Leagueid.Valid {
				leagueID = fmt.Sprintf("%d", income.Leagueid.Int32)
			}
			if income.LeagueName.Valid {
				leagueName = income.LeagueName.String
			}

			totalIncome := 0.0
			if v, ok := income.TotalIncome.(float64); ok {
				totalIncome = v
			}
			netRevenue := 0.0
			if v, ok := income.NetRevenue.(float64); ok {
				netRevenue = v
			}

			response.TournamentIncomes = append(response.TournamentIncomes, &analytics.TournamentIncome{
				TournamentId:   fmt.Sprintf("%d", income.Tournamentid),
				TournamentName: income.TournamentName,
				LeagueId:       leagueID,
				LeagueName:     leagueName,
				TotalIncome:    totalIncome,
				NetRevenue:     netRevenue,
				NetProfit:      float64(income.NetProfit),
				TournamentDate: income.Startdate.Format("2006-01-02"),
			})
		}

	case "school_financial_performance":
		if req.GroupBy == nil {
			return nil, fmt.Errorf("group_by parameter is required for school performance report")
		}

		var performanceData interface{}
		switch *req.GroupBy {
		case "category":
			performanceData, err = queries.GetSchoolPerformanceByCategory(ctx, models.GetSchoolPerformanceByCategoryParams{
				Startdate:   startDate,
				Startdate_2: endDate,
				Column3:     tournamentID,
			})
		case "location":
			performanceData, err = queries.GetSchoolPerformanceByLocation(ctx, models.GetSchoolPerformanceByLocationParams{
				Startdate:   startDate,
				Startdate_2: endDate,
				Column3:     tournamentID,
			})
		default:
			return nil, fmt.Errorf("invalid group_by parameter: %s", *req.GroupBy)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to get school performance: %v", err)
		}

		switch data := performanceData.(type) {
		case []models.GetSchoolPerformanceByCategoryRow:
			for _, row := range data {
				totalAmount := 0.0
				if v, ok := row.TotalAmount.(float64); ok {
					totalAmount = v
				}
				response.SchoolPerformance = append(response.SchoolPerformance, &analytics.SchoolPerformanceData{
					GroupName:   row.GroupName,
					TotalAmount: totalAmount,
					SchoolCount: int32(row.SchoolCount),
				})
			}
		case []models.GetSchoolPerformanceByLocationRow:
			for _, row := range data {
				groupName := ""
				if v, ok := row.GroupName.(string); ok {
					groupName = v
				}
				totalAmount := 0.0
				if v, ok := row.TotalAmount.(float64); ok {
					totalAmount = v
				}
				response.SchoolPerformance = append(response.SchoolPerformance, &analytics.SchoolPerformanceData{
					GroupName:   groupName,
					TotalAmount: totalAmount,
					SchoolCount: int32(row.SchoolCount),
				})
			}
		}

	case "expenses":
		if req.TournamentId != nil && *req.TournamentId != "" {
			expenses, err := queries.GetExpensesByTournament(ctx, models.GetExpensesByTournamentParams{
				Startdate:   startDate,
				Startdate_2: endDate,
				Column3:     tournamentID,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to get expenses by tournament: %v", err)
			}

			for _, expense := range expenses {
				totalExpense := 0.0
				if expense.Totalexpense.Valid {
					total, _ := strconv.ParseFloat(expense.Totalexpense.String, 64)
					totalExpense = total
				}

				foodExpense, _ := strconv.ParseFloat(expense.Foodexpense, 64)
				transportExpense, _ := strconv.ParseFloat(expense.Transportexpense, 64)
				perDiemExpense, _ := strconv.ParseFloat(expense.Perdiemexpense, 64)
				awardingExpense, _ := strconv.ParseFloat(expense.Awardingexpense, 64)
				stationaryExpense, _ := strconv.ParseFloat(expense.Stationaryexpense, 64)
				otherExpenses, _ := strconv.ParseFloat(expense.Otherexpenses, 64)

				response.ExpenseCategories = append(response.ExpenseCategories, &analytics.ExpenseCategory{
					TournamentId:      fmt.Sprintf("%d", expense.Tournamentid),
					TournamentName:    expense.TournamentName,
					FoodExpense:       foodExpense,
					TransportExpense:  transportExpense,
					PerDiemExpense:    perDiemExpense,
					AwardingExpense:   awardingExpense,
					StationaryExpense: stationaryExpense,
					OtherExpenses:     otherExpenses,
					TotalExpense:      totalExpense,
				})
			}
		} else {
			summary, err := queries.GetExpensesSummary(ctx, models.GetExpensesSummaryParams{
				Startdate:   startDate,
				Startdate_2: endDate,
				Column3:     tournamentID,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to get expenses summary: %v", err)
			}

			// Convert interface{} values to float64
			foodExpense := getFloat64FromInterface(summary.FoodExpense)
			transportExpense := getFloat64FromInterface(summary.TransportExpense)
			perDiemExpense := getFloat64FromInterface(summary.PerDiemExpense)
			awardingExpense := getFloat64FromInterface(summary.AwardingExpense)
			stationaryExpense := getFloat64FromInterface(summary.StationaryExpense)
			otherExpenses := getFloat64FromInterface(summary.OtherExpenses)
			totalExpense := getFloat64FromInterface(summary.TotalExpense)

			response.ExpenseCategories = append(response.ExpenseCategories, &analytics.ExpenseCategory{
				FoodExpense:       foodExpense,
				TransportExpense:  transportExpense,
				PerDiemExpense:    perDiemExpense,
				AwardingExpense:   awardingExpense,
				StationaryExpense: stationaryExpense,
				OtherExpenses:     otherExpenses,
				TotalExpense:      totalExpense,
			})
		}
	}

	return response, nil
}

func (s *AnalyticsService) GetAttendanceReports(ctx context.Context, req *analytics.AttendanceReportRequest) (*analytics.AttendanceReportResponse, error) {
	// Validate admin access
	claims, err := utils.ValidateToken(req.Token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	userRole, ok := claims["user_role"].(string)
	if !ok || userRole != "admin" {
		return nil, fmt.Errorf("unauthorized: only admins can access attendance reports")
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.DateRange.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", req.DateRange.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %v", err)
	}

	// Convert tournament ID
	var tournamentID int32
	if req.TournamentId != nil && *req.TournamentId != "" {
		id, err := strconv.ParseInt(*req.TournamentId, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid tournament ID: %v", err)
		}
		tournamentID = int32(id)
	}

	queries := models.New(s.db)
	response := &analytics.AttendanceReportResponse{
		ReportType: req.ReportType,
	}

	switch req.ReportType {
	case "category":
		attendance, err := queries.GetSchoolAttendanceByCategory(ctx, models.GetSchoolAttendanceByCategoryParams{
			Startdate:   startDate,
			Startdate_2: endDate,
			Column3:     tournamentID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get category attendance: %v", err)
		}

		var totalSchools int64 = 0
		var weightedTotalChange float64 = 0.0

		for _, data := range attendance {
			totalSchools += data.SchoolCount

			// Convert interface{} to float64 for percentage change
			percentageChange := 0.0
			if pc, ok := data.PercentageChange.(float64); ok {
				percentageChange = pc
			}

			response.CategoryAttendance = append(response.CategoryAttendance, &analytics.CategoryAttendance{
				Category:         data.Category,
				SchoolCount:      int32(data.SchoolCount),
				PercentageChange: percentageChange,
			})

			weightedTotalChange += float64(data.SchoolCount) * percentageChange
		}

		response.TotalSchools = int32(totalSchools)
		if totalSchools > 0 {
			response.TotalPercentageChange = weightedTotalChange / float64(totalSchools)
		}

	case "location":
		showRwandaProvinces := len(req.Countries) == 1 && req.Countries[0] == "Rwanda"

		attendance, err := queries.GetSchoolAttendanceByLocation(ctx, models.GetSchoolAttendanceByLocationParams{
			Startdate:   startDate,
			Startdate_2: endDate,
			Column3:     tournamentID,
			Column4:     showRwandaProvinces,
			Column5:     req.Countries,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get location attendance: %v", err)
		}

		var totalSchools int64 = 0
		var weightedTotalChange float64 = 0.0

		for _, data := range attendance {
			totalSchools += data.SchoolCount

			// Convert interface{} to string for location
			location := ""
			if loc, ok := data.Location.(string); ok {
				location = loc
			}

			// Convert interface{} to float64 for percentage change
			percentageChange := 0.0
			if pc, ok := data.PercentageChange.(float64); ok {
				percentageChange = pc
			}

			response.LocationAttendance = append(response.LocationAttendance, &analytics.LocationAttendance{
				Location:         location,
				LocationType:     data.LocationType,
				SchoolCount:      int32(data.SchoolCount),
				PercentageChange: percentageChange,
			})

			weightedTotalChange += float64(data.SchoolCount) * percentageChange
		}

		response.TotalSchools = int32(totalSchools)
		if totalSchools > 0 {
			response.TotalPercentageChange = weightedTotalChange / float64(totalSchools)
		}

	default:
		return nil, fmt.Errorf("invalid report type: %s", req.ReportType)
	}

	return response, nil
}

// Helper function to convert interface{} to float64
func getFloat64FromInterface(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case int64:
		return float64(v)
	case int32:
		return float64(v)
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	default:
		return 0
	}
}
