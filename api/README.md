# API 계약 관리

`openapi.yaml`이 HTTP API 계약의 기준입니다. Go 구조체와 주석에서 문서를 역으로 만들지 않습니다.

1. `openapi.yaml`에서 경로, operationId, 요청·응답 스키마, 오류 상태를 수정합니다.
2. `go generate ./api`로 Go 타입과 strict 서버 인터페이스를 다시 생성합니다.
3. 각 모듈의 `transport/http/handler.go`에서 변경된 인터페이스를 구현합니다.
4. `response.go`에서 도메인 모델을 공개 응답 DTO로 변환합니다.
5. `make test`로 생성 코드 일치 여부와 계약·단위 테스트를 검증합니다.

생성 파일 `internal/api/api.gen.go`는 직접 수정하지 않으며 명세 변경과 함께 커밋합니다. 생성기 버전은 go.mod의 tool 의존성으로 고정되어 있습니다. `oapi-codegen.yaml`이 생성 옵션을 관리합니다.

## 문서 보기

```sh
go run ./cmd/apidocs
```

- 문서 UI: http://127.0.0.1:8081/docs/
- 원본 명세: http://127.0.0.1:8081/openapi.yaml
- 문서 전용 서버는 DB에 연결하지 않고 업무 API도 실행하지 않으며 Try it out을 비활성화합니다.
- 실제 API 서버에서는 동일한 문서를 `/docs/`, `/openapi.yaml`로 제공합니다. 실제 서버의 Swagger UI에서 Try it out을 사용하면 데이터를 변경할 수 있습니다.
- UI의 CSS·JavaScript는 Go 바이너리에 포함되어 외부 CDN 연결이 필요 없습니다.

## 계약과 업무 규칙

OpenAPI는 필수 필드, 타입, 길이, 수치 범위, nullable, 상태 enum, 응답 상태를 정의합니다. 서버는 명세에 따라 요청을 검사합니다. 서비스는 회원·상품 존재 여부, 중복, 판매 기간, 재고, 자기 판매 구매 금지 등 업무 규칙을 검사합니다.

회원 요청의 모르는 필드는 무시하고, 상품·판매·구매 요청은 `additionalProperties: false`로 거부합니다. JSON 키는 명세의 이름을 정확하게 사용합니다. 전체 수정에서 nullable 필드 생략과 null은 모두 NULL로 처리합니다. 응답의 nullable 필드는 생략하지 않고 null로 반환합니다.

기존 요청 정책을 유지하여 첫 번째 JSON 값만 읽고 뒤의 값이나 문자열은 무시합니다. 첫 JSON을 읽을 때 1 MiB 제한을 적용합니다. Content-Type 생략 시 application/json으로 간주하며, 명시한 Content-Type이 지원되지 않으면 400을 반환합니다.

실제 응답이 명세에 맞는지도 계약 테스트에서 검증합니다. `make check-api`는 임시 파일로 재생성해 결과를 비교하므로, 검증 과정에서 작업 중인 생성 파일을 덮어쓰지 않습니다.
