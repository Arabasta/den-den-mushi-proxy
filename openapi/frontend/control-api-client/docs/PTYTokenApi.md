# PTYTokenApi

All URIs are relative to *http://localhost:55007*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**apiV1PtyTokenJoinPost**](#apiv1ptytokenjoinpost) | **POST** /api/v1/pty_token/join | Mint a join token for an existing PTY session|
|[**apiV1PtyTokenStartPost**](#apiv1ptytokenstartpost) | **POST** /api/v1/pty_token/start | Mint a start token for a new PTY session|

# **apiV1PtyTokenJoinPost**
> TokenResponse apiV1PtyTokenJoinPost(joinRequest)


### Example

```typescript
import {
    PTYTokenApi,
    Configuration,
    JoinRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new PTYTokenApi(configuration);

let joinRequest: JoinRequest; //

const { status, data } = await apiInstance.apiV1PtyTokenJoinPost(
    joinRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **joinRequest** | **JoinRequest**|  | |


### Return type

**TokenResponse**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Authorization token and proxy load balancer URL to connect to |  -  |
|**4XX** | Invalid request |  -  |
|**5XX** | Server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **apiV1PtyTokenStartPost**
> TokenResponse apiV1PtyTokenStartPost(startRequest)


### Example

```typescript
import {
    PTYTokenApi,
    Configuration,
    StartRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new PTYTokenApi(configuration);

let startRequest: StartRequest; //

const { status, data } = await apiInstance.apiV1PtyTokenStartPost(
    startRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **startRequest** | **StartRequest**|  | |


### Return type

**TokenResponse**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Authorization token and proxy load balancer URL to connect to |  -  |
|**4XX** | Invalid request |  -  |
|**5XX** | Server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

