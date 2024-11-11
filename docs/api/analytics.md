# Analytics Service Documentation

## Financial Analytics API

The analytics service provides comprehensive financial reporting capabilities for administrators. It includes detailed tournament income analysis, school performance metrics, and expense tracking.

### GetFinancialReports

Endpoint: `AnalyticsService.GetFinancialReports`
Authorization: Admin users only

Request:
```json
{
  "token": "your_auth_token_here",
  "date_range": {
    "start_date": "2024-01-01",
    "end_date": "2024-12-31"
  },
  "tournament_id": "123",  // Optional
  "report_type": "income_overview",
  "group_by": "category"  // Required for school_performance report type
}
```

#### Report Types

1. **Income Overview** (`report_type: "income_overview"`)
   - Shows tournament-wise financial performance
   - Metrics include:
     - Total income per tournament
     - Net revenue (after discounts)
     - Net profit (revenue minus expenses)
   - Grouped by league and tournament
   - Filterable by date range and specific tournament

2. **School Performance** (`report_type: "school_performance"`)
   - Financial performance metrics for schools
   - Two grouping options:
     - By Category (`group_by: "category"`): Groups schools by their type (Private, Public, etc.)
     - By Location (`group_by: "location"`):
       - For Rwandan schools: Groups by province
       - For international schools: Groups by country
   - Metrics include:
     - Total amount paid
     - Number of schools in each group

3. **Expenses** (`report_type: "expenses"`)
   - Detailed breakdown of tournament expenses
   - Categories include:
     - Food expenses
     - Transport expenses
     - Per diem expenses
     - Awarding expenses
     - Stationary expenses
     - Other expenses
   - Available as:
     - Per tournament breakdown (when tournament_id is provided)
     - Summary across all tournaments (when tournament_id is not provided)

## Testing Analytics Features

To test the analytics features:

1. Authentication and Authorization:
   - Verify that only admin users can access the reports
   - Test with invalid tokens and non-admin users
   - Verify appropriate error messages are shown

2. Income Overview Report:
   - Test with different date ranges
   - Test with specific tournament IDs
   - Verify calculations for:
     - Total income
     - Net revenue
     - Net profit

3. School Performance Report:
   - Test both grouping options:
     - Category grouping
     - Location grouping
   - Verify correct aggregation of:
     - Total amounts
     - School counts

4. Expenses Report:
   - Test tournament-specific expense breakdown
   - Test overall expense summary
   - Verify all expense categories are correctly calculated

5. Date Range Filtering:
   - Test with various date ranges
   - Test with invalid date formats
   - Test with future dates
   - Test with very large date ranges

6. Error Handling:
   - Test with missing required fields
   - Test with invalid tournament IDs
   - Test with invalid report types
   - Test with invalid group by options

7. Performance Testing:
   - Test with large datasets
   - Test concurrent requests
   - Verify response times are within acceptable limits

## Notes

- All monetary values are returned as float64 for precision
- Dates should be provided in "YYYY-MM-DD" format
- The service handles null values and provides appropriate defaults
- Responses include metadata about the report generation