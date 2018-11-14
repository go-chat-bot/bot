
## Setup steps

* Create new project in google cloud console
  * This is Config.PubSubProject
* Create service credentials for this project
  * Path to this file should be in env variable GOOGLE_APPLICATION_CREDENTIALS
  * Choose "Pub/Sub Editor" role
* Enable Pub/Sub API in cloud console
* Create new topic in the Pub/Sub (say "hangouts-chat")
  * This is Config.TopicName
* Modify permissions on created topic so that
  "chat-api-push@system.gserviceaccount.com" has Pub/Sub Publisher permissions
* Enable hangouts chat api
* Go to hangouts chat API config and fill in info
  * Connection settings - use Pub/Sub and fill in topic string you created
    above
  * Verification token copied to your Config.Token

Config.SubscriptionName should be unique for each environment or you'll not
process messages correctly
