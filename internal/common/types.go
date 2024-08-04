package commons

type Queue struct {
	Name       string
	RoutingKey string
	Exchange   string
	Consumer   string
	UseDelay   bool
}
