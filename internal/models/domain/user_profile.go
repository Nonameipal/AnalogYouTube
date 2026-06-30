package domain 
import "time"


type UserProfile struct {
	ID int `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Description string `json:"description"`
	SubscribersCount int `json:"subscribers_count"`
	SubscriptionsCount int `json:"subscriptions_count"`
	Subscribers *[]User `json:"subcribers,omitempty"`
	Subscriptions *[]User `json:"subscriptions,omitempty"`
	Videos []Video `json:"videos"`
	Role string `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"` 
}