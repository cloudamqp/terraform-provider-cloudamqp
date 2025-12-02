## Coding style

* Use `any` instead of `interface{}`
* Use `int64` instead of `int` for Terraform attributes (framework uses `types.Int64`)
* Mark sensitive attributes with `Sensitive: true` in schema
* Use `types.String`, `types.Int64`, `types.Bool`, etc. for attribute types
* Prefer `ctx context.Context` as first parameter in all functions
* Use `tflog` package for logging instead of standard log package

## Client library

### API Models

* Create request/response models in `api/models/` subdirectories organized by category
* Use descriptive struct names: `VpcRequest`, `VpcResponse`, `CustomCertificateRequest`
* Use JSON tags for all fields: `` `json:"field_name"` ``
* Keep models simple - no business logic in model structs
* Response models should match API response structure exactly
* Request models should contain only fields accepted by the API

### API Client Methods

* All API methods should accept `ctx context.Context` as first parameter
* Use `callWithRetry()` for all HTTP requests via `retryRequest` struct
* Method naming convention:
  * `Create{Resource}(ctx, params)` - POST requests
  * `Read{Resource}(ctx, id)` - GET requests
  * `Update{Resource}(ctx, id, params)` - PUT/PATCH requests
  * `Delete{Resource}(ctx, id)` - DELETE requests
* Return appropriate model types and error
* Use `tflog.Debug()` to log method, path, and parameters

### Using callWithRetry

* Build `retryRequest` struct with:
  * `functionName` - name of API method for logging (e.g., "CreateVPC")
  * `resourceName` - resource type for error messages (e.g., "VPC")
  * `attempt` - always start at 1
  * `sleep` - delay between retries (typically `5 * time.Second`)
  * `data` - pointer to response model struct or `nil`
  * `failed` - pointer to `map[string]any` for error responses
  * `customRetryCode` - optional, for custom retry logic (e.g., 400 for polling)
* Chain sling request: `api.sling.New().Post(path).BodyJSON(params)`
* Pass to `callWithRetry()`: `api.callWithRetry(ctx, sling, retryRequest{...})`

### Example API Method

```go
func (api *API) CreateVPC(ctx context.Context, params model.VpcRequest) (model.VpcResponse, error) {
    var (
        data   model.VpcResponse
        failed map[string]any
        path   = "/api/vpcs"
    )

    tflog.Debug(ctx, fmt.Sprintf("method=POST path=%s params=%v", path, params))
    err := api.callWithRetry(ctx, api.sling.New().Post(path).BodyJSON(params), retryRequest{
        functionName: "CreateVPC",
        resourceName: "VPC",
        attempt:      1,
        sleep:        5 * time.Second,
        data:         &data,
        failed:       &failed,
    })
    if err != nil {
        return model.VpcResponse{}, err
    }
    return data, nil
}
```

### Polling Patterns

* For async operations, create separate polling methods
* Use `customRetryCode` to retry on non-success status codes
* Poll until resource is ready (typically returns 200)
* Example: `pollForVpcReady()` retries on 400 until VPC returns 200

### Error Handling in API Methods

* Return wrapped errors with context: `fmt.Errorf("failed to read VPC: %w", err)`
* Handle resource drift by checking for zero values in response
* Return `nil` response (not error) when resource is gone (404, 410)

## Terraform Framework Guidelines

### Resource Structure

* Implement required interfaces:
  * `resource.Resource` - base interface
  * `resource.ResourceWithConfigure` - for provider configuration
  * `resource.ResourceWithImportState` - for import support
* Use private struct with lowercase name: `type myResourceResource struct`
* Implement these methods:
  * `Metadata()` - set resource type name
  * `Schema()` - define resource schema
  * `Configure()` - receive API client from provider
  * `Create()`, `Read()`, `Update()`, `Delete()` - CRUD operations
  * `ImportState()` - for resource import

### Schema Definition

* Use `schema.Schema` with `Attributes` map
* Common attribute modifiers:
  * `Required: true` - attribute must be set
  * `Optional: true` - attribute can be omitted
  * `Computed: true` - attribute is set by provider
  * `Sensitive: true` - hide value in logs and output
* Plan modifiers:
  * `RequiresReplace()` - force new resource on change
  * `UseStateForUnknown()` - use state value when unknown
  * Default values: `int64default.StaticInt64(1)`

### Model Structs

* Use separate struct for resource data model
* Suffix with `ResourceModel`: `type myResourceResourceModel struct`
* Use struct tags: `` `tfsdk:"attribute_name"` ``
* Use framework types: `types.String`, `types.Int64`, `types.Bool`, etc.

### Error Handling

* Add errors to diagnostics: `resp.Diagnostics.AddError("Title", err.Error())`
* Check diagnostics: `if resp.Diagnostics.HasError() { return }`
* Use descriptive error titles for user clarity

### Logging

* Use `tflog` package for structured logging
* Available levels: `Debug()`, `Info()`, `Warn()`, `Error()`
* Always pass context: `tflog.Debug(ctx, "message", map[string]any{"key": "value"})`
* Use for debugging API calls and resource operations

### Async Operations

* Use `context.WithTimeout()` for long-running operations
* Poll for job completion using `PollForJobCompleted()`
* Make sleep and timeout configurable via resource attributes
* Default values: sleep=10s, timeout=1800s (30 minutes)

### Best Practices

* Set `ID` attribute in Create method
* Handle imports by parsing import ID string
* Use `path.Root("attribute")` for setting nested attributes
* Always defer `cancel()` when using timeout contexts
* Return early on diagnostic errors
* Keep methods focused and single-purpose

## Changelog

* Section Headers (in order):
  * NOTES
  * FEATURES
  * IMPROVEMENTS
  * BUG FIXES
  * DEPENDENCIES
  * DEPRECATED

* Formatting Rules:
  * Past tense verbs: "Added", "Fixed", "Updated", "Removed", "Deprecated", "Bumped"
  * PR references: Always include ([#XXX]) at the end of each entry
  * PR links: Add reference links at the bottom of each release section
  * Consistent structure: Bullet points for all entries
  * Clear descriptions: Brief but informative

### Example Changelog Entry

```markdown
## 1.x.x (Unreleased)

FEATURES

* **New Resource:** `cloudamqp_custom_certificate` - Upload custom certificates to cluster ([#XXX])

IMPROVEMENTS

* resource/cloudamqp_instance: Added support for new plan types ([#XXX])

BUG FIXES

* resource/cloudamqp_alarm: Fixed panic when alarm was deleted outside Terraform ([#XXX])

DEPENDENCIES

* Bumped github.com/hashicorp/terraform-plugin-framework to v1.4.2 ([#XXX])

[#XXX]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/XXX
```
