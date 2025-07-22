# Connection


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [optional] [default to undefined]
**user_id** | **string** |  | [optional] [default to undefined]
**pty_session_id** | **string** |  | [optional] [default to undefined]
**start_role** | [**StartRole**](StartRole.md) |  | [optional] [default to undefined]
**status** | [**ConnectionStatus**](ConnectionStatus.md) |  | [optional] [default to undefined]
**join_time** | **string** |  | [optional] [default to undefined]
**leave_time** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { Connection } from './api';

const instance: Connection = {
    id,
    user_id,
    pty_session_id,
    start_role,
    status,
    join_time,
    leave_time,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
