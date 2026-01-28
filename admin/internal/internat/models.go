package internat

type Messages_ModelsI18nModuleKey struct {
	Admin string
	New   string

	MediaLibrary string
	Advertisers  string
	Campaigns    string
	Bgroups      string
	Banners      string
	Users        string
	Audiences    string
	Additions    string
	Settings     string

	// Banner
	Banner           string
	BannersInfo      string
	BannersTargeting string
	BannersCapping   string
	BannersBudget    string
	BannersTracking  string
	BannersMarker    string

	BannersID           string
	BannersImage        string
	BannersIcon         string
	BannersBgroupID     string
	BannersBgroup       string
	BannersPrice        string
	BannersActive       string
	BannersTitle        string
	BannersLabel        string
	BannersClicktracker string
	BannersImptracker   string
	BannersTarget       string
	BannersErid         string
	BannersCpm          string
	BannersDescription  string
	BannersTimetable    string
	BannersMacros       string

	// Bgroup
	Bgroup           string
	BgroupsInfo      string
	BgroupsTargeting string
	BgroupsTimetable string
	BgroupsCapping   string
	BgroupsBudget    string

	BgroupsID         string
	BgroupsPrice      string
	BgroupsCampaign   string
	BgroupsCampaignID string
	BgroupsActive     string
	BgroupsTitle      string
	BgroupsStart      string
	BgroupsEnd        string
	BgroupsCpm        string

	// Campaign
	Campaign           string
	CampaignsInfo      string
	CampaignsTargeting string
	CampaignsTimetable string
	CampaignsCapping   string
	CampaignsBudget    string

	CampaignsID           string
	CampaignsTitle        string
	CampaignsActive       string
	CampaignsBundle       string
	CampaignsStart        string
	CampaignsEnd          string
	CampaignsAdvertiser   string
	CampaignsAdvertiserID string

	// Advertiser
	Advertiser           string
	AdvertisersInfo      string
	AdvertisersTargeting string
	AdvertisersTimetable string
	AdvertisersCapping   string
	AdvertisersBudget    string

	AdvertisersID          string
	AdvertisersTitle       string
	AdvertisersActive      string
	AdvertisersStart       string
	AdvertisersEnd         string
	AdvertisersOrdContract string
	AdvertisersOrdEnable   string

	// User
	User          string
	UsersName     string
	UsersAccount  string
	UsersPassword string

	// Audience
	AudiencesId    string
	Audience       string
	AudiencesTitle string
	AudiencesName  string
	AudiencesFile  string
	AudiencesInfo  string

	// Workers
	Workers          string
	WorkersId        string
	WorkersJob       string
	WorkersStatus    string
	WorkersCreatedAt string

	// Network
	Network       string
	Networks      string
	NetworksName  string
	NetworksTitle string
}

var Messages_en_EN_ModelsI18nModuleKey = &Messages_ModelsI18nModuleKey{
	Admin: "AdEx",
	New:   "New",

	MediaLibrary: "Media Library",
	Advertisers:  "Advertisers",
	Campaigns:    "Campaigns",
	Bgroups:      "Groups",
	Banners:      "Banners",
	Users:        "Users",

	BannersTracking: "Tracking",
	BannersTitle:    "Title",
	BannersLabel:    "Label",
}

var Messages_ru_RU_ModelsI18nModuleKey = &Messages_ModelsI18nModuleKey{
	Admin: "AdEx",
	New:   "Новый",

	MediaLibrary: "Библиотека медиа",
	Advertisers:  "Рекламодатели",
	Campaigns:    "Кампании",
	Bgroups:      "Группы",
	Banners:      "Креативы",
	Users:        "Пользователи",
	Audiences:    "Аудитории",
	Additions:    "Дополнения",
	Settings:     "Настройки",

	Banner:           "Креатив",
	BannersInfo:      "Общая информация",
	BannersTargeting: "Таргетинг",
	BannersBudget:    "Бюджет",
	BannersCapping:   "Каппинг",
	BannersTracking:  "Трекинг",
	BannersMarker:    "Маркировка",

	BannersID:           "ID",
	BannersImage:        "Картинка",
	BannersIcon:         "Иконка",
	BannersBgroupID:     "Группа",
	BannersBgroup:       "Группа",
	BannersPrice:        "Цена",
	BannersActive:       "Включен",
	BannersTitle:        "Заголовок",
	BannersLabel:        "Лейбл",
	BannersClicktracker: "Кликтрекер",
	BannersImptracker:   "Имптрекер",
	BannersTarget:       "Таргет",
	BannersErid:         "Ерид",
	BannersCpm:          "CPM",
	BannersDescription:  "Описание",
	BannersTimetable:    "Расписание активности",
	BannersMacros:       "Макросы",

	Bgroup:           "Группа",
	BgroupsInfo:      "Общая информация",
	BgroupsTargeting: "Таргетинг",
	BgroupsTimetable: "Расписание активности",
	BgroupsBudget:    "Бюджет",
	BgroupsCapping:   "Каппинг",

	BgroupsID:         "ID",
	BgroupsCampaign:   "Кампания",
	BgroupsCampaignID: "Кампания",
	BgroupsPrice:      "Цена",
	BgroupsActive:     "Включен",
	BgroupsTitle:      "Заголовок",
	BgroupsStart:      "Начало",
	BgroupsEnd:        "Конец",
	BgroupsCpm:        "CPM",

	Campaign:           "Кампания",
	CampaignsInfo:      "Общая информация",
	CampaignsTargeting: "Таргетинг",
	CampaignsTimetable: "Расписание активности",
	CampaignsBudget:    "Бюджет",
	CampaignsCapping:   "Каппинг",

	CampaignsID:           "ID",
	CampaignsTitle:        "Заголовок",
	CampaignsActive:       "Включен",
	CampaignsBundle:       "Бандл приложения",
	CampaignsStart:        "Начало",
	CampaignsEnd:          "Конец",
	CampaignsAdvertiser:   "Рекламодатель",
	CampaignsAdvertiserID: "Рекламодатель",

	Advertiser:           "Рекламодатель",
	AdvertisersInfo:      "Общая информация",
	AdvertisersTargeting: "Таргетинг",
	AdvertisersTimetable: "Расписание активности",
	AdvertisersCapping:   "Каппинг",
	AdvertisersBudget:    "Бюджет",

	AdvertisersID:          "ID",
	AdvertisersTitle:       "Заголовок",
	AdvertisersActive:      "Включен",
	AdvertisersStart:       "Начало",
	AdvertisersEnd:         "Конец",
	AdvertisersOrdContract: "ОРД Контракт",
	AdvertisersOrdEnable:   "Использовать ОРД",

	User:          "Пользователь",
	UsersName:     "Имя",
	UsersAccount:  "Аккаунт",
	UsersPassword: "Пароль",

	AudiencesId:    "ID",
	Audience:       "Аудитория",
	AudiencesTitle: "Заголовок",
	AudiencesName:  "Название",
	AudiencesFile:  "Файл",
	AudiencesInfo:  "Информация",

	Workers:          "Workers",
	WorkersId:        "ID",
	WorkersJob:       "Job",
	WorkersStatus:    "Status",
	WorkersCreatedAt: "Created At",

	Network:       "Ендпоинт",
	Networks:      "Ендпоинты",
	NetworksName:  "Название",
	NetworksTitle: "Заголовок",
}
