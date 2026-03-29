package repository

import (
	"context"
	"time"

	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/database"
	"github.com/Z-TAS-Solutions/Z-QryptGIN/internal/app/dto"
	"gorm.io/gorm"
)

// DashboardRepository handles database queries for admin dashboard analytics
type DashboardRepository interface {
	// GetAuthTrendsByInterval fetches authentication activity trends aggregated over a specified interval
	// Supports intervals: "minute" or "hour"
	// Returns time-series data points with success and failure counts
	GetAuthTrendsByInterval(ctx context.Context, interval string, lastHours int) ([]dto.AuthTrendDataPoint, error)

	// GetAuthTrendsMetrics returns aggregated authentication metrics for a summary view
	GetAuthTrendsMetrics(ctx context.Context, lastHours int) (*dto.DashboardMetrics, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository creates a new dashboard repository instance
func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

// GetAuthTrendsByInterval fetches authentication trends aggregated by the specified interval
// Fetches activity logs for successful and failed logins over the last N hours
func (r *dashboardRepository) GetAuthTrendsByInterval(ctx context.Context, interval string, lastHours int) ([]dto.AuthTrendDataPoint, error) {
	var results []struct {
		TimeBucket time.Time
		Type       string
		Count      int64
	}

	// Determine the PostgreSQL date_trunc interval format
	intervalFormat := interval
	if interval != "minute" && interval != "hour" {
		intervalFormat = "hour" // Default to hour if invalid
	}

	// Query: Aggregate activity logs by time bucket, counting successes and failures separately
	// DATE_TRUNC groups logs by the specified interval (minute or hour)
	query := `
		SELECT
			DATE_TRUNC(?, activity_logs.created_at) AS time_bucket,
			activity_logs.type,
			COUNT(*) AS count
		FROM activity_logs
		WHERE activity_logs.created_at >= NOW() - INTERVAL '1 hour' * ?
			AND activity_logs.type IN (?, ?)
		GROUP BY DATE_TRUNC(?, activity_logs.created_at), activity_logs.type
		ORDER BY time_bucket ASC
	`

	if err := r.db.WithContext(ctx).Raw(
		query,
		intervalFormat,
		lastHours,
		string(database.ActivityLoginSuccess),
		string(database.ActivityFailedLogin),
		intervalFormat,
	).Scan(&results).Error; err != nil {
		return nil, err
	}

	// Transform flat results into time-series data points
	// Use a map to group by timestamp
	trendMap := make(map[time.Time]*dto.AuthTrendDataPoint)

	for _, result := range results {
		if _, exists := trendMap[result.TimeBucket]; !exists {
			trendMap[result.TimeBucket] = &dto.AuthTrendDataPoint{
				Timestamp:    result.TimeBucket,
				SuccessCount: 0,
				FailureCount: 0,
			}
		}

		if result.Type == string(database.ActivityLoginSuccess) {
			trendMap[result.TimeBucket].SuccessCount = result.Count
		} else if result.Type == string(database.ActivityFailedLogin) {
			trendMap[result.TimeBucket].FailureCount = result.Count
		}
	}

	// Convert map to sorted slice
	var trends []dto.AuthTrendDataPoint
	for _, point := range trendMap {
		trends = append(trends, *point)
	}

	// Sort by timestamp ascending
	if len(trends) > 1 {
		for i := 0; i < len(trends)-1; i++ {
			for j := i + 1; j < len(trends); j++ {
				if trends[i].Timestamp.After(trends[j].Timestamp) {
					trends[i], trends[j] = trends[j], trends[i]
				}
			}
		}
	}

	// Ensure empty intervals are filled
	// This prevents gaps in the time series visualization
	trends = r.fillMissingIntervals(trends, interval, lastHours)

	return trends, nil
}

// GetAuthTrendsMetrics returns aggregated metrics for authentication trends
func (r *dashboardRepository) GetAuthTrendsMetrics(ctx context.Context, lastHours int) (*dto.DashboardMetrics, error) {
	var successCount, failureCount int64

	// Count total successful logins in the period
	if err := r.db.WithContext(ctx).
		Where("created_at >= NOW() - INTERVAL '1 hour' * ?", lastHours).
		Where("type = ?", string(database.ActivityLoginSuccess)).
		Model(&database.ActivityLog{}).
		Count(&successCount).Error; err != nil {
		return nil, err
	}

	// Count total failed logins in the period
	if err := r.db.WithContext(ctx).
		Where("created_at >= NOW() - INTERVAL '1 hour' * ?", lastHours).
		Where("type = ?", string(database.ActivityFailedLogin)).
		Model(&database.ActivityLog{}).
		Count(&failureCount).Error; err != nil {
		return nil, err
	}

	// Fetch trends to calculate averages and peaks
	trends, err := r.GetAuthTrendsByInterval(ctx, "hour", lastHours)
	if err != nil {
		return nil, err
	}

	var avgSuccess, avgFailure float64
	var peakSuccess, peakFailure int64
	var peakSuccessTime, peakFailureTime time.Time

	if len(trends) > 0 {
		avgSuccess = float64(successCount) / float64(len(trends))
		avgFailure = float64(failureCount) / float64(len(trends))

		for _, trend := range trends {
			if trend.SuccessCount > peakSuccess {
				peakSuccess = trend.SuccessCount
				peakSuccessTime = trend.Timestamp
			}
			if trend.FailureCount > peakFailure {
				peakFailure = trend.FailureCount
				peakFailureTime = trend.Timestamp
			}
		}
	}

	return &dto.DashboardMetrics{
		TotalSuccessfulAuthentications: successCount,
		TotalFailedAuthentications:     failureCount,
		AverageSuccessPerInterval:      avgSuccess,
		AverageFailurePerInterval:      avgFailure,
		PeakSuccessCount:               peakSuccess,
		PeakFailureCount:               peakFailure,
		PeakSuccessTime:                peakSuccessTime,
		PeakFailureTime:                peakFailureTime,
	}, nil
}

// fillMissingIntervals ensures that all time buckets in the range are represented
// This prevents gaps in time-series visualizations when no activity occurred in an interval
func (r *dashboardRepository) fillMissingIntervals(trends []dto.AuthTrendDataPoint, interval string, lastHours int) []dto.AuthTrendDataPoint {
	if len(trends) == 0 {
		return trends
	}

	// Create a map for quick lookup
	trendMap := make(map[time.Time]dto.AuthTrendDataPoint)
	for _, trend := range trends {
		trendMap[trend.Timestamp] = trend
	}

	// Determine the interval duration
	var intervalDuration time.Duration
	if interval == "minute" {
		intervalDuration = time.Minute
	} else {
		intervalDuration = time.Hour
	}

	// Generate all expected timestamps from lastHours ago to now
	endTime := time.Now().UTC().Truncate(intervalDuration)
	startTime := endTime.Add(-time.Duration(lastHours) * time.Hour).Truncate(intervalDuration)

	var filledTrends []dto.AuthTrendDataPoint
	for current := startTime; !current.After(endTime); current = current.Add(intervalDuration) {
		if trend, exists := trendMap[current]; exists {
			filledTrends = append(filledTrends, trend)
		} else {
			// Add a zero-value data point for this interval
			filledTrends = append(filledTrends, dto.AuthTrendDataPoint{
				Timestamp:    current,
				SuccessCount: 0,
				FailureCount: 0,
			})
		}
	}

	return filledTrends
}
