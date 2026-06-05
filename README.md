# go_CRUD_accounts
Account CRUD. 

# Bảng `accounts` với:
id (uuid),  
name,  
currency (ISO 4217 — "USD", "VND", "JPY"),  
balance (int64, minor units, không phải float64),  
created_at,  
updated_at,  
deleted_at (soft delete)  

# 5 endpoints: 
/POST   
/GET one  
/GET list  
/PATCH name  
/DELETE (soft).  

# Requirements
- SQLite
- Có Test

# Does it stick?
Cái thực sự đo "đã fluent chưa":
Khi build, để ý 6 thứ này. Có hesitate khi type không?

- func (h *Handler) Create(w http.ResponseWriter, r *http.Request) — handler signature
- ctx := r.Context() — context propagation từ request
- json.NewDecoder(r.Body).Decode(&req) với proper error handling
- db.QueryContext(ctx, ...) + defer rows.Close() — DB pattern with context
- func TestCreateAccount(t *testing.T) với httptest.NewServer
- Graceful shutdown: srv.Shutdown(ctx) on SIGTERM

4/6 flow tự nhiên = fluent, move on. 
< 3/6 = Tour of Go chưa stick, re-do warmup. 
Đây là honest self-check, không gian lận với chính mình ở đây — vì M5 (Ledger concurrency) sẽ punish brutally nếu Go syntax vẫn còn cản trở.

# Can it be attacked?
Mở rộng service của bạn:

Thêm POST /accounts dùng đầy đủ pattern trên: MaxBytesReader, DisallowUnknownFields, DTO riêng (createAccountRequest), validate Name + Currency, server set Balance = 0, trả 201.
Thêm GET /accounts?currency=USD&limit=10 với query parsing: reject duplicate currency key, parse + bound limit trong [1,100], filter accounts theo currency.
Test 5 attack scenarios bằng curl — mỗi cái phải bị reject đúng cách:
- Body thừa field: curl -X POST -d '{"name":"x","currency":"USD","balance":99999}' → expect 400 (unknown field)
- Body không phải JSON: curl -X POST -d 'not json' → expect 400
- Currency lạ: curl -X POST -d '{"name":"x","currency":"XYZ"}' → expect 400
- Duplicate query: curl 'localhost:8080/accounts?currency=USD&currency=EUR' → expect 400
- Limit quá lớn: curl 'localhost:8080/accounts?limit=999999' → expect 400
- Và 1 happy path: tạo account hợp lệ → expect 201, balance = 0 dù bạn cố gửi balance khác


Trả lời 1 câu sau khi xong:
> Bạn vừa dùng DTO riêng (createAccountRequest) thay vì decode thẳng vào account. Nhưng điều này tốn công: bạn phải maintain hai struct, và copy field từ DTO sang domain struct bằng tay. Một junior engineer sẽ nói "thừa, decode thẳng vào account cho nhanh." Bạn phản biện thế nào? Đưa một lý do mà junior đó không thể cãi lại.

Câu này test xem bạn đã internalize trust boundary, hay mới chỉ copy code tôi. Câu trả lời tốt sẽ chạm tới một điều sâu hơn cả security — nó chạm tới tại sao input type và domain type là hai concern khác nhau về bản chất.