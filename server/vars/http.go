package vars

const (
	ChatHost = "http://chat:8888"

	ChatGetCustomerConfigUri         = "/api/config/get-customer-config"
	ChatFindCustomerConfigUri        = "/api/config/find-get-customer-config"
	ChatUpdateCustomerConfigUri      = "/api/config/update-customer-config"
	ChatSaveCustomerConfigUri        = "/api/config/create-customer-config"
	ChatDeleteCustomerConfigUri      = "/api/config/delete-customer-config"
	ChatDeleteCustomerConfigByIdsUri = "/api/config/delete-customer-config-by-ids"

	ChatGetApplicationConfigUri         = "/api/config/get-application-config"
	ChatFindApplicationConfigUri        = "/api/config/find-get-application-config"
	ChatSaveApplicationConfigUri        = "/api/config/create-application-config"
	ChatUpdateApplicationConfigUri      = "/api/config/update-application-config"
	ChatDeleteApplicationConfigUri      = "/api/config/delete-application-config"
	ChatDeleteApplicationConfigByIdsUri = "/api/config/delete-application-config-by-ids"

	ChatGetChatUri = "/api/msg/get-chat"
)
