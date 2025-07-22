# WhitelistBlacklistApi

All URIs are relative to *http://localhost:55007*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**apiV1BlacklistRegexGet**](#apiv1blacklistregexget) | **GET** /api/v1/blacklist/regex | Get all blacklist regex filters for healthcheck ou group|
|[**apiV1BlacklistRegexIdDelete**](#apiv1blacklistregexiddelete) | **DELETE** /api/v1/blacklist/regex/{id} | Soft delete a blacklist regex filter|
|[**apiV1BlacklistRegexIdPut**](#apiv1blacklistregexidput) | **PUT** /api/v1/blacklist/regex/{id} | Update a blacklist regex filter|
|[**apiV1BlacklistRegexPost**](#apiv1blacklistregexpost) | **POST** /api/v1/blacklist/regex | Add a regex to blacklist for healthcheck ou group|
|[**apiV1WhitelistRegexGet**](#apiv1whitelistregexget) | **GET** /api/v1/whitelist/regex | Get all whitelist regex filters for healthcheck ou group|
|[**apiV1WhitelistRegexIdDelete**](#apiv1whitelistregexiddelete) | **DELETE** /api/v1/whitelist/regex/{id} | Soft delete a whitelist regex filter|
|[**apiV1WhitelistRegexIdPut**](#apiv1whitelistregexidput) | **PUT** /api/v1/whitelist/regex/{id} | Update a whitelist regex filter|
|[**apiV1WhitelistRegexPost**](#apiv1whitelistregexpost) | **POST** /api/v1/whitelist/regex | Add a regex to whitelist for healthcheck ou group|

# **apiV1BlacklistRegexGet**
> Array<RegexFilter> apiV1BlacklistRegexGet()


### Example

```typescript
import {
    WhitelistBlacklistApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new WhitelistBlacklistApi(configuration);

const { status, data } = await apiInstance.apiV1BlacklistRegexGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<RegexFilter>**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of blacklist regex filters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **apiV1BlacklistRegexIdDelete**
> apiV1BlacklistRegexIdDelete()


### Example

```typescript
import {
    WhitelistBlacklistApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new WhitelistBlacklistApi(configuration);

let id: number; // (default to undefined)

const { status, data } = await apiInstance.apiV1BlacklistRegexIdDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] |  | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Filter deleted successfully |  -  |
|**404** | Filter not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **apiV1BlacklistRegexIdPut**
> RegexFilter apiV1BlacklistRegexIdPut(apiV1WhitelistRegexIdPutRequest)


### Example

```typescript
import {
    WhitelistBlacklistApi,
    Configuration,
    ApiV1WhitelistRegexIdPutRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new WhitelistBlacklistApi(configuration);

let id: number; // (default to undefined)
let apiV1WhitelistRegexIdPutRequest: ApiV1WhitelistRegexIdPutRequest; //

const { status, data } = await apiInstance.apiV1BlacklistRegexIdPut(
    id,
    apiV1WhitelistRegexIdPutRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **apiV1WhitelistRegexIdPutRequest** | **ApiV1WhitelistRegexIdPutRequest**|  | |
| **id** | [**number**] |  | defaults to undefined|


### Return type

**RegexFilter**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Filter updated successfully |  -  |
|**400** | Invalid input |  -  |
|**404** | Filter not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **apiV1BlacklistRegexPost**
> RegexFilter apiV1BlacklistRegexPost(apiV1WhitelistRegexPostRequest)


### Example

```typescript
import {
    WhitelistBlacklistApi,
    Configuration,
    ApiV1WhitelistRegexPostRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new WhitelistBlacklistApi(configuration);

let apiV1WhitelistRegexPostRequest: ApiV1WhitelistRegexPostRequest; //

const { status, data } = await apiInstance.apiV1BlacklistRegexPost(
    apiV1WhitelistRegexPostRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **apiV1WhitelistRegexPostRequest** | **ApiV1WhitelistRegexPostRequest**|  | |


### Return type

**RegexFilter**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Filter created successfully |  -  |
|**400** | Invalid input (e.g., malformed regex) |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **apiV1WhitelistRegexGet**
> Array<RegexFilter> apiV1WhitelistRegexGet()


### Example

```typescript
import {
    WhitelistBlacklistApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new WhitelistBlacklistApi(configuration);

const { status, data } = await apiInstance.apiV1WhitelistRegexGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**Array<RegexFilter>**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of whitelist regex filters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **apiV1WhitelistRegexIdDelete**
> apiV1WhitelistRegexIdDelete()


### Example

```typescript
import {
    WhitelistBlacklistApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new WhitelistBlacklistApi(configuration);

let id: number; // (default to undefined)

const { status, data } = await apiInstance.apiV1WhitelistRegexIdDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] |  | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Filter deleted successfully |  -  |
|**404** | Filter not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **apiV1WhitelistRegexIdPut**
> RegexFilter apiV1WhitelistRegexIdPut(apiV1WhitelistRegexIdPutRequest)


### Example

```typescript
import {
    WhitelistBlacklistApi,
    Configuration,
    ApiV1WhitelistRegexIdPutRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new WhitelistBlacklistApi(configuration);

let id: number; // (default to undefined)
let apiV1WhitelistRegexIdPutRequest: ApiV1WhitelistRegexIdPutRequest; //

const { status, data } = await apiInstance.apiV1WhitelistRegexIdPut(
    id,
    apiV1WhitelistRegexIdPutRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **apiV1WhitelistRegexIdPutRequest** | **ApiV1WhitelistRegexIdPutRequest**|  | |
| **id** | [**number**] |  | defaults to undefined|


### Return type

**RegexFilter**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Filter updated successfully |  -  |
|**400** | Invalid input |  -  |
|**404** | Filter not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **apiV1WhitelistRegexPost**
> RegexFilter apiV1WhitelistRegexPost(apiV1WhitelistRegexPostRequest)


### Example

```typescript
import {
    WhitelistBlacklistApi,
    Configuration,
    ApiV1WhitelistRegexPostRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new WhitelistBlacklistApi(configuration);

let apiV1WhitelistRegexPostRequest: ApiV1WhitelistRegexPostRequest; //

const { status, data } = await apiInstance.apiV1WhitelistRegexPost(
    apiV1WhitelistRegexPostRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **apiV1WhitelistRegexPostRequest** | **ApiV1WhitelistRegexPostRequest**|  | |


### Return type

**RegexFilter**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Filter created successfully |  -  |
|**400** | Invalid input (e.g., malformed regex) |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

