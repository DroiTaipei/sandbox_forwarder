# v0.6.1
- util/slack
   - Add `CustomClient` to accept http client. Useful for httpmock on other package
   - Add `ThumbURL` in `Payload` to show image on slack message
   - Remove `ProxyURL` and `Timeout` from `WebhookInfo` because these can be included in `CustomeClient`
   - Upgrade to gopkg.in/jarcoal/httpmock.v1

# v0.6.0
- trace
    - Init commit for droitrace which depends on open tracing
- util/slack
    - Define Payload struct which can be used for rich content webhook
    - Remove Content field on WebhookInfo
    - Refine `AddField` by using `sort.Slice` which is supported after Golang 1.8
- file
    - Add error code

# v0.5.0
- errors
    - Add TraceDroiError for recording tracable error stacks
    - Remove NewApiError
- rdb
    - Migrate NewApiError to ConstDroiError

# v0.4.3
- rdb
    - Refactor rdb struct & error
    - Add QueryAppPrefixPayload for query app datas with appID prefix

# v0.4.2
- rdb
    - Add rdb struct & error into droipkg

# v0.4.1
- file
    - Add ParseCdnURL() to get fqdn, qppid, fid, and filename from CDN url

# v0.4.0
- file
    - Modify file.GetFileBasename CDN link rule

# v0.3.0
- file
    - Change the Op field of majesty's CallBackPayload to map[string]interface{}
    - Add log fields
    - Add OpCodes ID and mapping

# v0.2.0
- file
    - Add file.GetFileBasename to get basename for CDN link

# v0.1.0
- init. commit
