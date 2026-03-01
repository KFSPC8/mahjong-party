package parlour

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricOpenRooms = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "parlour_open_rooms",
		Help: "Current number of open rooms",
	})
	metricRoomSubscriptions = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "parlour_room_subscriptions",
		Help: "Current number of room subscriptions",
	})
)
