# library

```
[ResourceConfig]
    ServiceName = "trace_test"
	ServiceVersion = "1.0"
[ExportConfig]
	Name ="jeajer"
	[ExportConfig.Attribute]
        FileName = "trace.txt"
        PrettyPrint = true
        Timestamps = true
        Url = "http://localhost:14268/api/traces"

```