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
	// GetAuthTrendsByInterval fetches authorization activity trends aggregated over a specified interval
	// Tracks Authorization_Granted (success) vs Authorization_Denied (failure) attempts
	// Supports intervals: "minute" or "hour"
	// Returns time-series data points with granted and denied counts
	GetAuthTrendsByInterval(ctx context.Context, interval string, lastHours int) ([]dto.AuthTrendDataPoint, error)

	// GetAuthTrendsMetrics returns aggregated authorization metrics for a summary view
	GetAuthTrendsMetrics(ctx context.Context, lastHours int) (*dto.DashboardMetrics, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository creates a new dashboard repository instance
func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

// GetAuthTrendsByInterval fetches authorization trends aggregated by the specified interval
// Uses proper GORM queries to fetch activity logs, then aggregates in Go code
// Tracks: Authorization_Granted (success) and Authorization_Denied (failure)
func (r *dashboardRepository) GetAuthTrendsByInterval(ctx context.Context, interval string, lastHours int) ([]dto.AuthTrendDataPoint, error) {
	// Validate interval
	if interval != "minute" && interval != "hour" {
		interval = "hour"
	}

	// Calculate cutoff time
	cutoffTime := time.Now().UTC().Add(-time.Duration(lastHours) * time.Hour)

	// Fetch all activity logs for the specified time period using pure GORM
	// Filter for Authorization_Granted and Authorization_Denied types only
	var activityLogs []database.ActivityLog
	err := r.db.WithContext(ctx).
		Where("created_at >= ?", cutoffTime).
		Where("type IN (?, ?)", database.ActivityAuthorizationGranted, database.ActivityAuthorizationDenied).
		Order("created_at ASC").
		Find(&activityLogs).Error

	if err != nil {
		return nil, err
	}

	// Aggregate in Go by truncating timestamps to the specified interval
	// Use a map to group by timestamp bucket and activity type
	trendMap := make(map[time.Time]*dto.AuthTrendDataPoint)

	for _, log := range activityLogs {
		// Truncate timestamp to interval (minute or hour)
		var timeBucket time.Time
		if interval == "minute" {
			timeBucket = log.CreatedAt.UTC().Truncate(time.Minute)
		} else {
			timeBucket = log.CreatedAt.UTC().Truncate(time.Hour)
		}

		// Initialize the bucket if it doesn't exist
		if _, exists := trendMap[timeBucket]; !exists {
			trendMap[timeBucket] = &dto.AuthTrendDataPoint{
				Timestamp:    timeBucket,
				SuccessCount: 0,
				FailureCount: 0,
			}
		}

		// Increment count based on activity type
		if log.Type == database.ActivityAuthorizationGranted {
			trendMap[timeBucket].SuccessCount++
		} else if log.Type == database.ActivityAuthorizationDenied {
			trendMap[timeBucket].FailureCount++
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

// GetAuthTrendsMetrics returns aggregated metrics for authorization trends
func (r *dashboardRepository) GetAuthTrendsMetrics(ctx context.Context, lastHours int) (*dto.DashboardMetrics, error) {
	var grantedCount, deniedCount int64

	// Calculate cutoff time using proper Go time operations (not SQL)
	cutoffTime := time.Now().UTC().Add(-time.Duration(lastHours) * time.Hour)

	// Count total authorization grants in the period using pure GORM
	if err := r.db.WithContext(ctx).
		Where("created_at >= ?", cutoffTime).
		Where("type = ?", database.ActivityAuthorizationGranted).
		Model(&database.ActivityLog{}).
		Count(&grantedCount).Error; err != nil {
		return nil, err
	}

	// Count total authorization denials in the period using pure GORM
	if err := r.db.WithContext(ctx).
		Where("created_at >= ?", cutoffTime).
		Where("type = ?", database.ActivityAuthorizationDenied).
		Model(&database.ActivityLog{}).
		Count(&deniedCount).Error; err != nil {
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
		avgSuccess = float64(grantedCount) / float64(len(trends))
		avgFailure = float64(deniedCount) / float64(len(trends))

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
		TotalSuccessfulAuthentications: grantedCount,
		TotalFailedAuthentications:     deniedCount,
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
