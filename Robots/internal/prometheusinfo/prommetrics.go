package prometheusinfo

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	CreatedRobot = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "robot_add_total",
			Help: "Количество созданных роботов",
		},
	)

	GetRobot = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "robot_get_total",
			Help: "Количество обновлённых роботов",
		},
	)

	UpdateRobotCords = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "robot_updatecord_total",
			Help: "Количество обновлений координат роботов",
		},
	)

	UpdateRobotNames = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "robot_updatename_total",
			Help: "Количество обновлений имён роботов",
		},
	)

	UpdateRobotType = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "robot_updatetype_total",
			Help: "Количество обновлений типов роботов",
		},
	)

	DeletedRobot = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "robot_deleted_total",
			Help: "Количество удалений роботов",
		},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "handler", "status"},
	)

	CountOfRobotType = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "robot_types_total",
			Help: "Типы созданных роботов",
		},
		[]string{"robot_type"},
	)
)

func Register() {
	prometheus.MustRegister(CreatedRobot)
	prometheus.MustRegister(GetRobot)
	prometheus.MustRegister(UpdateRobotCords)
	prometheus.MustRegister(UpdateRobotNames)
	prometheus.MustRegister(UpdateRobotType)
	prometheus.MustRegister(DeletedRobot)
	prometheus.MustRegister(CountOfRobotType)

	prometheus.MustRegister(RequestDuration)
}
