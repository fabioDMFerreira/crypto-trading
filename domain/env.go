package domain

type Env struct {
	MongoURL                    string
	MongoDB                     string
	NotificationsReceiver       string
	NotificationsSender         string
	NotificationsSenderPassword string
	AppEnv                      string
	AppID                       string
}
