package filter

import "den-den-mushi-Go/pkg/dto/proxy_host"

type Func func(*proxy_host.Record2) bool
