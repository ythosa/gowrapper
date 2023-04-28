# gowrapper

## Example
* Before:
    ```
    my-app/
    ├─ cmd/
    ├─ internal/
    │  ├─ protos/
    │  │  ├─ some_proto_interface.go
    │  ├─ usecases/
    │  │  ├─ file_with_usecase_interface.go
    │  │  ├─ excluded_file1.go
    │  │  ├─ excluded_file2.go
    ```
* Command:
    ```shell
    gowrapper gen -d ./internal --excluded-files="/usecase/excluded_file1.go,/usecase/excluded_file2.go" \ 
        --excluded-dirs="protos" -t templates/opentracing -p tracing
    ```
* After:
    ```
    my-app/
    ├─ cmd/
    ├─ internal/
    │  ├─ protos/
    │  │  ├─ some_proto_interface.go
    │  ├─ usecases/
    │  │  ├─ file_with_usecase_interface.go
    │  │  ├─ excluded_file1.go
    │  │  ├─ excluded_file2.go
    │  │  ├─ usecase_tracing.go   <-- generated file
    ```
