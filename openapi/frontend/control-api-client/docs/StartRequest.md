# StartRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**purpose** | [**ConnectionPurpose**](ConnectionPurpose.md) |  | [default to undefined]
**change_id** | **string** | Only required if purpose is \&quot;change_request\&quot;, ID of the change request to connect to | [optional] [default to undefined]
**server** | [**ServerInfo**](ServerInfo.md) |  | [default to undefined]

## Example

```typescript
import { StartRequest } from './api';

const instance: StartRequest = {
    purpose,
    change_id,
    server,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
