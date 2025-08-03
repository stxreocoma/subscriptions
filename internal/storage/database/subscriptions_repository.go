package database

import (
	"context"
	"fmt"
	"subscriptions/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (d *Database) GetSubscription(userID uuid.UUID, serviceName string) (*models.Subscription, error) {
	row := d.pool.QueryRow(context.Background(),
		`SELECT service_name, user_id, price, start_date, end_date 
		 FROM subscriptions 
		 WHERE user_id = @user_id and service_name = @service_name;`,
		pgx.NamedArgs{"user_id": userID, "service_name": serviceName})

	sub := &models.Subscription{}
	var startDate, endDate string
	err := row.Scan(&sub.ServiceName, &sub.UserID, &startDate, &endDate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No subscription found
		}
		return nil, fmt.Errorf("error scanning row: %v", err)
	}

	sub.StartDate, err = time.Parse("2006-01", startDate)
	if err != nil {
		return nil, fmt.Errorf("error parsing start date: %v", err)
	}
	sub.EndDate, err = time.Parse("2006-01", endDate)
	if err != nil {
		return nil, fmt.Errorf("error parsing end date: %v", err)
	}
	return sub, nil
}

func (d *Database) CreateSubscription(subscription *models.Subscription) (*models.Subscription, error) {
	_, err := d.pool.Exec(context.Background(),
		`INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date)
		 VALUES (@user_id, @service_name, @price, @start_date, @end_date)`,
		pgx.NamedArgs{
			"user_id":      subscription.UserID,
			"service_name": subscription.ServiceName,
			"price":        subscription.Price,
			"start_date":   subscription.StartDate.Format("2006-01"),
			"end_date":     subscription.EndDate.Format("2006-01"),
		})
	if err != nil {
		return nil, fmt.Errorf("error creating subscription: %v", err)
	}
	return subscription, nil
}

func (d *Database) UpdateSubscription(subscription *models.Subscription) (*models.Subscription, error) {
	_, err := d.pool.Exec(context.Background(),
		`UPDATE subscriptions
		 SET price = @price, start_date = @start_date, end_date = @end_date 
		 WHERE user_id = @user_id AND service_name @service_name`,
		pgx.NamedArgs{
			"user_id":      subscription.UserID,
			"service_name": subscription.ServiceName,
			"price":        subscription.Price,
			"start_date":   subscription.StartDate.Format("2006-01"),
			"end_date":     subscription.EndDate.Format("2006-01"),
		})
	if err != nil {
		return nil, fmt.Errorf("error updating subscription: %v", err)
	}

	return subscription, nil
}

func (d *Database) DeleteSubscription(userID uuid.UUID, serviceName string) error {
	_, err := d.pool.Exec(context.Background(),
		`DELETE FROM subscriptions
		 WHERE user_id = @user_id AND service_name = @service_name`,
		pgx.NamedArgs{
			"user_id":      userID,
			"service_name": serviceName,
		})
	if err != nil {
		return fmt.Errorf("error deleting subscription: %v", err)
	}

	return nil
}

func (d *Database) ListSubscriptions(userID uuid.UUID, page int) ([]*models.Subscription, error) {
	rows, err := d.pool.Query(context.Background(),
		`SELECT service_name, user_id, price, start_date, end_date
		 FROM subscriptions
		 WHERE user_id = @user_id
		 ORDER BY start_date DESC
		 LIMIT 10 OFFSET 10*(@page-1)`,
		pgx.NamedArgs{"user_id": userID, "page": page})
	if err != nil {
		return nil, fmt.Errorf("error listing subscriptions: %v", err)
	}
	defer rows.Close()
	var subscriptions []*models.Subscription
	for rows.Next() {
		sub := &models.Subscription{}
		var startDate, endDate string
		if err := rows.Scan(&sub.ServiceName, &sub.UserID, &sub.Price, &startDate, &endDate); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		sub.StartDate, err = time.Parse("2006-01", startDate)
		if err != nil {
			return nil, fmt.Errorf("error parsing start date: %v", err)
		}
		sub.EndDate, err = time.Parse("2006-01", endDate)
		if err != nil {
			return nil, fmt.Errorf("error parsing end date: %v", err)
		}
		subscriptions = append(subscriptions, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}
	return subscriptions, nil
}

func (d *Database) SubscriptionTotalCost(userID uuid.UUID, serviceName, startDate, endDate string) (*models.TotalSubscriptionCost, error) {
	row := d.pool.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(
		 	price *
		 		EXTRACT(YEAR FROM AGE(
		 			LEAST(end_date, @end_date), end_date),
					GREATEST(start_date, @start_date)
				)) * 12 +
				EXTRACT(MONTH FROM AGE(
					LEAST(COALESCE(end_date, @end_date), @end_date),
					GREATEST(start_date, @start_date)
				))
		 	)
		 ), 0) AS total_cost
		 FROM subscriptions
		 WHERE user_id = @user_id
		   AND (service_name IS NULL OR service_name = @service_name)
		   AND start_date <= @start_date
		   AND (end_date IS NULL OR end_date >= @end_date);`,
		pgx.NamedArgs{
			"user_id":      userID,
			"start_date":   startDate,
			"end_date":     endDate,
			"service_name": serviceName,
		})
	var totalCost int
	if err := row.Scan(&totalCost); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error scanning row: %v", err)
	}
	return &models.TotalSubscriptionCost{TotalCost: totalCost}, nil
}

//апдейт: спец ручка эта скорее всего должна тока число возвращать, а не структуру, сделаю обязательный параметр userID и доп service_name, и получается получаем там просто старт енд и прайс берем разницк по месяцам и умножаем на нее сумму, если сервис нейм не придет то груп баим по айди и сумируем прайс а затем также делаем
