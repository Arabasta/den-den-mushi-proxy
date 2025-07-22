# HostSessionDetails


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**host** | [**Host**](Host.md) |  | [optional] [default to undefined]
**os_users** | **Array&lt;string&gt;** |  | [optional] [default to undefined]
**pty_sessions** | [**Array&lt;PtySessionSummary&gt;**](PtySessionSummary.md) |  | [optional] [default to undefined]

## Example

```typescript
import { HostSessionDetails } from './api';

const instance: HostSessionDetails = {
    host,
    os_users,
    pty_sessions,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
