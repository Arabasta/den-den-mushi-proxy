# TokenResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**token** | **string** | JWT token for proxy | [default to undefined]
**proxyUrl** | **string** | Load balancer URL for proxy group (e.g., &#x60;https://proxy.os.com&#x60;, &#x60;https://proxy.db.com&#x60;) | [default to undefined]

## Example

```typescript
import { TokenResponse } from './api';

const instance: TokenResponse = {
    token,
    proxyUrl,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
