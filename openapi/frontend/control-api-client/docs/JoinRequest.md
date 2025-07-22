# JoinRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**pty_session_id** | **string** | ID of the existing PTY session to join, user will be validated against this session\&#39;s initial connection details | [default to undefined]
**start_role** | [**StartRole**](StartRole.md) |  | [default to undefined]

## Example

```typescript
import { JoinRequest } from './api';

const instance: JoinRequest = {
    pty_session_id,
    start_role,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
