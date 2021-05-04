package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics will return a global variable Metrics which will contain all metrics collectors
var Metrics *MetricsCollectors

// MetricsUsers will return a set of users' related prometheus metrics
type MetricsUsers struct {
	UsersCreated prometheus.Counter
}

// MetricsCards will return a set of cards' related prometheus metrics
type MetricsCards struct {
	CardsCreated prometheus.Counter
}

// MetricsBalances will return a set of balances' related prometheus metrics
type MetricsBalances struct {
	BalancesCreated prometheus.Counter
}

// MetricsSpends will return a set of spends' related prometheus metrics
type MetricsSpends struct {
	SpendsCreated prometheus.Counter
}

// MetricsCollectors will return a struct with all metrics collectors
type MetricsCollectors struct {
	Users    *MetricsUsers
	Cards    *MetricsCards
	Balances *MetricsBalances
	Spends   *MetricsSpends
}

// InitMetrics will
func InitMetrics() (m *MetricsCollectors) {
	usersCreated := promauto.NewCounter(prometheus.CounterOpts{
		Name: "budget_tracker_users_created_total",
		Help: "The total number of created users",
	})

	cardsCreated := promauto.NewCounter(prometheus.CounterOpts{
		Name: "budget_tracker_cards_created_total",
		Help: "The total number of created cards",
	})

	balancesCreated := promauto.NewCounter(prometheus.CounterOpts{
		Name: "budget_tracker_balances_created_total",
		Help: "The total number of created balances",
	})

	spendsCreated := promauto.NewCounter(prometheus.CounterOpts{
		Name: "budget_tracker_spends_created_total",
		Help: "The total number of created spends",
	})

	Metrics = &MetricsCollectors{
		&MetricsUsers{
			UsersCreated: usersCreated,
		},
		&MetricsCards{
			CardsCreated: cardsCreated,
		},
		&MetricsBalances{
			BalancesCreated: balancesCreated,
		},
		&MetricsSpends{
			SpendsCreated: spendsCreated,
		},
	}

	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.Register(Metrics.Users.UsersCreated)
	prometheus.Register(Metrics.Cards.CardsCreated)
	prometheus.Register(Metrics.Balances.BalancesCreated)
	prometheus.Register(Metrics.Spends.SpendsCreated)

	return m
}
