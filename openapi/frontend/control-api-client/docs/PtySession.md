# PtySession


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [optional] [default to undefined]
**created_by** | **string** |  | [optional] [default to undefined]
**start_time** | **string** |  | [optional] [default to undefined]
**end_time** | **string** |  | [optional] [default to undefined]
**state** | [**PtySessionState**](PtySessionState.md) |  | [optional] [default to undefined]
**last_activity** | **string** |  | [optional] [default to undefined]
**purpose** | [**ConnectionPurpose**](ConnectionPurpose.md) |  | [optional] [default to undefined]
**change_id** | **string** | ID of the change request associated with this PTY session, if applicable | [optional] [default to undefined]
**connections** | [**Array&lt;Connection&gt;**](Connection.md) | List of connections to this PTY session\&#39;s life time | [optional] [default to undefined]

## Example

```typescript
import { PtySession } from './api';

const instance: PtySession = {
    id,
    created_by,
    start_time,
    end_time,
    state,
    last_activity,
    purpose,
    change_id,
    connections,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
