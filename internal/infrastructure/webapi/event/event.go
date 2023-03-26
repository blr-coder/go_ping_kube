package event

type CreateEvent struct {
	Title      string
	CampaignID int64
	Country    string
	Lang       string
	Device     string
}
