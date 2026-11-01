# Go 모듈러 모노리스

회원·상품·판매·구매를 한 프로세스에서 실행합니다. `dealer` 테이블 구조에 맞춘 도메인 모델과 메모리 저장소, HTTP API를 제공합니다. 표준 라이브러리만 사용합니다.

```sh
go run ./cmd/server
# 기본 :8080, 변경: HTTP_ADDR=:9090 go run ./cmd/server
go test -race ./...
go vet ./...
```

## 구조와 책임

```text
cmd/server/                         서버 실행·종료
internal/app/                       모듈 조립, 모듈 간 어댑터, API 접두사
internal/modules/{member,product,sale,purchase}/
  model.go                          테이블에 대응하는 도메인 모델
  service.go                        서비스 입력, 검증, 업무 흐름, 저장소 포트
  memory.go                         저장소 구현과 원자적 데이터 변경
  transport/http/
    routes.go                       경로 등록
    request.go                      모듈별 DTO, JSON 읽기, 숫자 경로 해석
    handler.go                      HTTP 입력 → 서비스 호출 → 응답
    response.go                     모듈별 오류·HTTP 응답
internal/modules/sale/price.go       DECIMAL(10,2) 금액 타입
internal/modules/purchase/history.go 구매 이력 모델
internal/platform/                  오류 분류, nullable 값 복사 유틸리티
database/                          스키마와 수동 마이그레이션 SQL
```

다른 모듈의 저장소를 직접 사용하지 않습니다. 소비자가 정의한 인터페이스를 `app`의 어댑터로 연결합니다. 공통 CRUD 핸들러는 사용하지 않으며, HTTP 없이 서비스를 호출해도 동일한 업무 검증이 적용됩니다.

## 테이블 반영

| 모듈 | 테이블 | 주요 변경 |
| --- | --- | --- |
| 회원 | `dealer.dt_member` | 숫자 `member_no`와 사용자 `id` 구분, `id` 유일성, nullable 연락처 |
| 상품 | `dealer.dt_product` | 숫자 `product_no`, 유일한 `product_code`, nullable 이름·설명 |
| 판매 | `dealer.dt_sale` | 판매자·상품 번호, 수량, 정확한 소수 금액, 상태·판매 기간 |
| 구매 | `dealer.dt_purchase` | 판매·판매자·구매자 번호, 수량, 구매 상태·시각 |
| 구매 이력 | `dealer.dt_purchase_history` | 구매 당시 데이터와 수량 스냅샷, nullable `purchase_no` |

- `BIGINT`는 `int64`, `INT`는 `int32`, `DATETIME`은 `time.Time`입니다.
- nullable 컬럼은 `*string`, `*time.Time`, `*int64`로 표현합니다. JSON에서는 `null`을 반환하며 빈 문자열과 구분합니다.
- `*_no`는 저장소가 생성하고, `create_at`·`update_at`·`purchase_at`은 서버에서 관리합니다. 회원 수정 시 생성 시각은 유지합니다.
- 회원 이메일은 유일하지 않습니다. 같은 이메일을 가진 회원을 허용합니다. 유일성 검사는 `id`와 `product_code`에 적용합니다.
- 메모리 저장소의 문자열 중복 검사는 대소문자를 구분합니다. 실제 DB 도입 시 테이블 collation과 일치시켜야 합니다.
- `price`는 내부적으로 1/100 단위 정수로 보관하며 JSON에서는 `"12000.50"`처럼 문자열로 받습니다. 실수와 지수 표기는 받지 않으며 양수 및 최대 `99999999.99`를 검증합니다.
- 판매 생성 시 수량은 양수이며 `ACTIVE` 또는 `CLOSED` 상태를 받습니다. 구매 시 `ACTIVE`, 시작 시각 도달, 종료 시각 이전, 충분한 재고를 검사합니다.
- 구매의 판매자는 판매 정보에서 가져옵니다. 자기 판매 구매를 거부하고 요청 수량만큼 재고를 원자적으로 차감합니다.
- 구매 상태는 현재 `COMPLETED`입니다. 결제 완료를 뜻하는 것은 아니며, 현 단계의 구매 생성 완료 상태입니다. 상태 전이·취소는 아직 제공하지 않습니다.
- 구매와 최초 이력은 저장소의 `CreateWithHistory`에서 함께 저장합니다. 저장 실패 시 재고를 복구합니다. 이력은 구매 모듈이 소유합니다.

## 승인된 스키마 추가

원래 테이블에 없던 다음 컬럼을 추가하는 구조입니다.

- `dt_sale.sell_member_no BIGINT NOT NULL`
- `dt_purchase.quantity INT NOT NULL`
- `dt_purchase_history.quantity INT NOT NULL` — 이력에도 구매 당시 수량 보존

새 DB의 전체 구조는 `database/schema.sql`에 있습니다. 이미 만든 테이블에는 `database/migrations/001_add_sale_owner_and_purchase_quantity.sql`로 nullable 컬럼을 추가한 뒤, 기존 데이터의 실제 판매자와 수량을 채우고 `002_require_sale_owner_and_purchase_quantity.sql`로 NOT NULL을 적용합니다. 기존 값은 임의로 추정하지 않습니다. 테이블이 비어 있다면 백필 없이 두 단계를 진행할 수 있습니다.

**SQL은 실제 DB에 실행하지 않았습니다. 현재 서버는 계속 메모리 저장소를 사용하며 재시작하면 데이터가 사라집니다.** DB 드라이버·접속 설정·SQL 저장소는 아직 없습니다. DB를 연결할 때 모듈별 Repository를 구현하고, 구매·이력의 원자적 저장 및 재고의 조건부 갱신을 보장해야 합니다. 현재 프로세스 내부 보상은 서버 중단에 대한 영속적 트랜잭션을 대신하지 않습니다.

## API

| 메서드 | 경로 | 용도 |
| --- | --- | --- |
| GET | `/health` | 상태 확인 |
| POST | `/api/v1/members` | 회원 등록 |
| GET | `/api/v1/members/{member_no}` | 회원 조회 |
| PUT | `/api/v1/members/{member_no}` | 회원 전체 수정 |
| DELETE | `/api/v1/members/{member_no}` | 회원 삭제 |
| POST | `/api/v1/products` | 상품 등록 |
| GET | `/api/v1/products/{product_no}` | 상품 조회 |
| POST | `/api/v1/sales` | 판매 등록 |
| GET | `/api/v1/sales/{sale_no}` | 판매 조회 |
| POST | `/api/v1/purchases` | 구매 생성 |
| GET | `/api/v1/purchases/{purchase_no}` | 구매 조회 |
| GET | `/api/v1/purchases/{purchase_no}/histories` | 구매 이력 조회 |

생성 201, 조회·수정 200, 삭제 204(본문 없음), 입력 오류 400, 없는 데이터 404, 중복·재고 부족·판매 불가 409를 반환합니다. 경로는 양수 정수 PK입니다. 기존 UUID 경로와 `seller_id`, `buyer_id`, `unit_price`, `stock` 입력은 새 번호·컬럼명으로 변경되었습니다.

날짜 입력은 시간대가 포함된 RFC3339 문자열입니다. 첫 JSON 값만 읽으며 뒤에 붙은 데이터는 무시합니다. JSON 읽기에는 1 MiB 제한이 적용되지만 읽지 않은 뒷부분까지 포함한 전체 본문 크기를 검사하지는 않습니다. 모르는 필드 처리 정책은 각 모듈의 request.go에 있습니다.

### 호출 예시

응답의 번호를 다음 요청에 사용합니다. 아래 1·2는 새 메모리 저장소에서 순서대로 실행한 예시입니다.

```sh
curl localhost:8080/api/v1/members -H 'Content-Type: application/json' \
  -d '{"id":"seller","email":null,"address":"Seoul","phone_number":null}'

curl localhost:8080/api/v1/members -H 'Content-Type: application/json' \
  -d '{"id":"buyer","email":"buyer@example.com"}'

curl localhost:8080/api/v1/products -H 'Content-Type: application/json' \
  -d '{"product_code":"KEYBOARD-001","name":"키보드","description":null}'

curl localhost:8080/api/v1/sales -H 'Content-Type: application/json' \
  -d '{"sell_member_no":1,"product_no":1,"name":"키보드 판매","quantity":10,"price":"12000.50","status":"ACTIVE","start_at":"2020-01-01T00:00:00Z","end_at":null}'

curl localhost:8080/api/v1/purchases -H 'Content-Type: application/json' \
  -d '{"sale_no":1,"buy_member_no":2,"quantity":2}'

curl localhost:8080/api/v1/purchases/1/histories
curl localhost:8080/api/v1/sales/1
```

회원 PUT은 `id`와 nullable 연락처 필드를 받아 전체 교체합니다. 생략한 nullable 필드는 NULL이 됩니다. 삭제 시 다른 모듈의 데이터는 연쇄 삭제하지 않습니다. 인증·인가, 결제, 배송, 목록 검색, 구매 요청 멱등성은 아직 구현하지 않았습니다.
