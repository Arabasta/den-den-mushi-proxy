# MakeChangeApi

All URIs are relative to *http://localhost:55007*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**apiV1ChangeRequestsGet**](#apiv1changerequestsget) | **GET** /api/v1/change_requests/ | Get change request and associated PTY sessions|

# **apiV1ChangeRequestsGet**
> Array<ChangeRequestSessionsResponse> apiV1ChangeRequestsGet()

Returns all APPROVED Change Requests associated with the user\'s implementor group and associated PTY sessions and connections Implementor group is retrieved from the user\'s id provided by the authentication token.

### Example

```typescript
import {
    MakeChangeApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new MakeChangeApi(configuration);

let ticketIds: Array<string>; // (optional) (default to undefined)
let implementorGroups: Array<string>; // (optional) (default to undefined)
let lob: string; // (optional) (default to undefined)
let country: string; // (optional) (default to undefined)
let startTime: string; // (optional) (default to undefined)
let endTime: string; // (optional) (default to undefined)
let ptySessionState: PtySessionState; //Does not filter out CRs, field not for users (optional) (default to undefined)
let page: number; // (optional) (default to 1)
let pageSize: number; // (optional) (default to 20)

const { status, data } = await apiInstance.apiV1ChangeRequestsGet(
    ticketIds,
    implementorGroups,
    lob,
    country,
    startTime,
    endTime,
    ptySessionState,
    page,
    pageSize
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **ticketIds** | **Array&lt;string&gt;** |  | (optional) defaults to undefined|
| **implementorGroups** | **Array&lt;string&gt;** |  | (optional) defaults to undefined|
| **lob** | [**string**] |  | (optional) defaults to undefined|
| **country** | [**string**] |  | (optional) defaults to undefined|
| **startTime** | [**string**] |  | (optional) defaults to undefined|
| **endTime** | [**string**] |  | (optional) defaults to undefined|
| **ptySessionState** | **PtySessionState** | Does not filter out CRs, field not for users | (optional) defaults to undefined|
| **page** | [**number**] |  | (optional) defaults to 1|
| **pageSize** | [**number**] |  | (optional) defaults to 20|


### Return type

**Array<ChangeRequestSessionsResponse>**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List change request details and their associated PTY sessions |  -  |
|**404** | Not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

