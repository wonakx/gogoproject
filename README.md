# Go 모듈러 모노리스

회원, 상품, 판매, 구매를 하나의 프로세스로 실행하는 HTTP API입니다. Go 표준 라이브러리만 사용합니다.

## 실행 및 검증

`go.mod`에 지정된 Go 버전으로 실행합니다.

```sh
go run ./cmd/server
# 기본 :8080, 변경 예: HTTP_ADDR=:9090 go run ./cmd/server
go test -race ./...
go vet ./...
```

## 디렉터리 구조

```text
cmd/server/                    서버 실행, 종료 처리
internal/app/                  모듈 조립, 모듈 간 어댑터, HTTP 라우팅
internal/modules/
  member/                      회원 등록, 조회, 이메일 중복 검사
  product/                     상품 카탈로그 등록, 조회
  sale/                        판매자·상품 연결, 판매가, 재고 관리
  purchase/                    구매 생성, 가격 스냅샷, 구매 조회
internal/platform/
  fault/                       공통 오류 분류
  identity/                    식별자 생성
  httpjson/                    공통 JSON 요청·응답 및 오류 처리
```

각 모듈의 `transport/http/handler.go`가 리소스 경로와 요청 DTO를 소유합니다. `internal/app/http.go`는 모듈 라우터를 조립하고 `/api/v1` 접두사를 한 번만 적용합니다. `/health`는 버전 없는 운영 경로입니다.

각 모듈의 `service.go`는 모델, 유스케이스, 소비자 관점의 의존 인터페이스를 정의합니다. `memory.go`는 해당 모듈의 저장소 구현입니다. 외부 의존성이 추가되면 모델과 포트를 별도 파일로 분리할 수 있습니다.

### 모듈 경계

- 다른 모듈의 패키지나 저장소를 직접 참조하지 않습니다. 회원·상품의 존재 확인, 판매 재고 예약 등 필요한 기능을 소비자 모듈의 인터페이스로 정의합니다.
- `internal/app/app.go`만 모든 모듈을 조립하고 어댑터로 연결합니다. 각 모듈의 HTTP 어댑터도 저장소를 직접 사용하지 않고 자신의 서비스를 호출합니다.
- 회원은 회원 정보, 상품은 카탈로그, 판매는 판매가와 재고, 구매는 주문을 소유합니다. 같은 상품을 여러 판매자가 서로 다른 가격과 재고로 판매할 수 있습니다.
- 구매는 회원 존재 확인 → 판매 재고 예약 → 주문 저장 순서입니다. 주문에는 판매자, 상품, 구매 당시 단가와 총액을 보존합니다.
- 저장소 구현은 동시 접근을 보호합니다. 판매 재고 확인과 차감은 하나의 잠금 안에서 수행하여 초과 판매를 방지합니다. 자기 상품 구매와 0 이하 수량, 총액 오버플로를 거부합니다.
- 저장 실패 시 재고를 보상 복구합니다. 구매 저장소는 오류 반환 시 주문이 저장되지 않았다는 계약을 지켜야 합니다.

## API

| 메서드 | 경로 | 용도 |
| --- | --- | --- |
| GET | `/health` | 상태 확인 |
| POST | `/api/v1/members` | 회원 등록 |
| GET | `/api/v1/members/{id}` | 회원 조회 |
| POST | `/api/v1/products` | 상품 등록 |
| GET | `/api/v1/products/{id}` | 상품 조회 |
| POST | `/api/v1/sales` | 판매 등록 |
| GET | `/api/v1/sales/{id}` | 판매 및 남은 재고 조회 |
| POST | `/api/v1/purchases` | 구매 생성 |
| GET | `/api/v1/purchases/{id}` | 구매 조회 |

생성은 201, 조회는 200을 반환합니다. 입력 오류는 400, 없는 리소스는 404, 이메일 중복이나 재고 부족은 409입니다. 오류 본문은 `{"error":"..."}`입니다. JSON 요청은 최대 1 MiB이며, 알 수 없는 필드와 여러 JSON 값은 거부합니다.

가격은 원화 정수 금액(`int64`)입니다. 예를 들어 `12000`은 12,000원입니다. 구매는 단일 판매 건에 대한 주문이며 결제 처리는 포함하지 않습니다.

### 요청 예시

아래 순서로 호출하고 응답의 `id`를 다음 요청에 넣습니다.

```sh
# 판매자 등록 → SELLER_ID
curl -sS localhost:8080/api/v1/members -H 'Content-Type: application/json' \
  -d '{"name":"판매자","email":"seller@example.com"}'

# 구매자 등록 → BUYER_ID
curl -sS localhost:8080/api/v1/members -H 'Content-Type: application/json' \
  -d '{"name":"구매자","email":"buyer@example.com"}'

# 상품 등록 → PRODUCT_ID
curl -sS localhost:8080/api/v1/products -H 'Content-Type: application/json' \
  -d '{"name":"키보드","description":"기계식 키보드"}'

# 판매 등록 → SALE_ID
curl -sS localhost:8080/api/v1/sales -H 'Content-Type: application/json' \
  -d '{"seller_id":"SELLER_ID","product_id":"PRODUCT_ID","unit_price":12000,"stock":10}'

# 구매 → PURCHASE_ID, total_price: 24000
curl -sS localhost:8080/api/v1/purchases -H 'Content-Type: application/json' \
  -d '{"buyer_id":"BUYER_ID","sale_id":"SALE_ID","quantity":2}'

# 남은 재고: 8
curl -sS localhost:8080/api/v1/sales/SALE_ID
curl -sS localhost:8080/api/v1/purchases/PURCHASE_ID
```

## 저장소 교체 및 현재 범위

메모리 저장소이므로 서버 재시작 시 모든 데이터가 사라집니다. DB 도입 시 각 모듈의 `Repository`를 구현하고 `app.New`의 조립을 변경합니다. 모듈별 테이블 소유권을 유지하고 다른 모듈의 테이블을 직접 조회하지 않습니다.

현재 보상 복구는 프로세스 내부 실행에 한정됩니다. DB를 도입할 때는 구매 저장과 재고 차감의 원자성을 위한 트랜잭션 경계 또는 영속적인 예약·보상 처리가 필요합니다. 다중 서버에서는 DB 수준의 조건부 재고 갱신도 필요합니다.

인증·인가, 결제, 취소, 배송, 목록·검색, 요청 멱등성은 아직 구현하지 않았습니다. 현재 ID는 요청 본문으로 받으므로 운영 서비스에서는 인증한 사용자 ID로 대체해야 합니다. 같은 구매 요청을 재전송하면 별도 주문이 만들어집니다.
