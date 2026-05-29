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

# Tasks
You still owe me:

- Build v2 with /accounts/{id}, fake delayed fetchAccount, 1-second timeout
- Test 3 scenarios với curl
- Answer the sequencing question — 3 things happen when timeout fires, in order, what are they?