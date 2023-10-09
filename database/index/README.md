# advanced index strategies in postgres
- indexing trong postgres là một tiến trình tạo một data-structure hỗ trợ cho việc optimized để search cũng như retrieve data từ table.
- index thực tế là copy 1 phần của table, cách này chúng ta active việc locate và retrieve những hàng thoả điều kiện query.
- khi một câu query được thực hiện, postgres sẽ kiểm tra xem indexes có không, và xác định nếu câu query thoả mãn điều kiện và xác định những hàng liên quang trong table 1 cách nhanh chống. Như vậy kết quả là sẽ nhanh hơn 1 cách đáng kể, đặc biệt là trong tình huống table lớn hoặc câu query phức tạp.
- postgres cung cấp 1 vài loại index như `B-tree`, `hash`, `GiST`, `SP-GiST`, and `BRIN`, Mỗi loại sẽ có những cách process khác nhau, nên cần xác định rõ usecase để tạo index 1 cách hợp lý.
- Một điểm cực kì quan trọng cần chú ý đó là khi index được tạo, sẽ ảnh hưởng đến hiệu năng trong việc write operation, như `insert` `delete` `update`.

## B-Tree Index
- Btree là loại index phổ biến nhất trong việc store và query trong postgres, nó cũng là index mặc định. Bất cứ khi nào bạn tạo index mà không chỉ định loại index thì postgres sẽ tạo Btree index cho table hoặc column.
- Btree được tổ chức như tree structure, index bắt đầu với root-node, với các node con là pointer. Mỗi node trong tree chứa nhiều key-value, nơi key là được dùng cho việc indexing còn value là con trỏ để liên kết đến data trong table.
- Để tạo B-tree index, sử dụng câu lệnh `create index`.
    ```sql
    CREATE INDEX index_name ON table_name;
    ```
### Column Indexing
- Để có thể tạo index trên 1 column của table thay vì toàn bộ table, dùng câu lệnh
    ```sql
    CREATE INDEX index_name ON table_name (column_name);
    ```
    - `index_name` là tên của index bạn tạo
    - `table_name` là index trên table nào.
    - `column_name` là tên của column trên table mà bạn tạo
- Ví dụ
    - Bây giờ sẽ tạo 1 table và dummy data cho table.
        ```sql
            CREATE TABLE info (
                id integer NOT NULL, 
                email VARCHAR, 
                location VARCHAR,
                name VARCHAR,
                uuid VARCHAR,
                age VARCHAR,
                ref_id integer
            );
        ```
    - insert data cho table
        ```sql
            INSERT INTO info (id, email, location, name,uuid, age,ref_id) 
            VALUES 
            (
                1, 'halie416@gmail.com', 'london', 'Headphone1', '4960d495-5c0b-43e2-8b79-4aed7f50be0d', '50',2
            ), 
            (
                2, 'romaine21@gmail.com', 'Australia', 'Webcam','bfa44785-adbd-4972-be3f-0988bbaa4a13', '50',2
            ), 
            (
                3, 'frederique19@gmail.com', 'Canada', 'iPhone 14 pro','e70e05de-312f-4497-bdd9-d612fd3ba0fc', '1259',1
            ), 
            (
                4, 'kenton_macejkovic80@hotmail.com', 'London', 'Wireless Mouse','b1c052f5-c274-4c24-84be-3775a4b08e22', '20',2
            ), 
            (
                5, 'alexis62@hotmail.com', 'Switzerland', 'Dell Charger','4c338f2a-71a0-4078-ba6c-37914b1badd2', '15',3
            ), 
            (
                6, 'concepcion_kiehn@hotmail.com', 'Canada', 'Longitech Keyboard','9b0888ac-e707-475c-b0fe-e882a2d5cac6', '499',3
            );
        ```
    - bạn cũng có thể random insert bằng câu lệnh dưới
        ```sql
            DO $FN$
            BEGIN
            FOR counter IN 1..1000 LOOP
                insert into info(id, email, location, name,uuid, age,ref_id)  values(
                    counter, 
                    md5(random()::text), 
                    md5(random()::text), 
                    md5(random()::text), 
                    uuid_in(overlay(overlay(md5(random()::text || ':' || random()::text) placing '4' from 13) placing to_hex(floor(random()*(11-8+1) + 8)::int)::text from 17)::cstring), 
                    counter,
                    counter
                );
            END LOOP;
            END;
            $FN$
        ```

- Bài toán đặt ra là trong query theo điều kiện uuid + ref_id

### Testing
- sẽ search bằng 1 uuid và ref_id =1
    ```sql
    EXPLAIN select * from info where uuid = '05aee2f5-aa10-4a15-9e2c-44166380e3a4' and ref_id =3;
    ```
  - kết quả explain sẽ là:
    ```sql
    "Seq Scan on info  (cost=0.00..38.00 rows=1 width=147)"
    " Filter: (((uuid)::text = '05aee2f5-aa10-4a15-9e2c-44166380e3a4'::text) AND (ref_id = 3))"
    ```

- Bây giờ bạn sẽ tạo index trên column `ref_id` bằng câu lệnh
    ```sql
        CREATE INDEX idx_info_ref_id ON info (ref_id);
    ```
- Với câu lệnh query như sau:
    ```sql
    EXPLAIN ANALYZE select * from info where uuid = '05aee2f5-aa10-4a15-9e2c-44166380e3a4' and ref_id =3;
    ```

- kết quả:
  - total 1k record
    ```sql
    "Seq Scan on info  (cost=0.00..38.00 rows=1 width=147) (actual time=0.079..0.145 rows=1 loops=1)"
    "  Filter: (((uuid)::text = '05aee2f5-aa10-4a15-9e2c-44166380e3a4'::text) AND (ref_id = 3))"
    "  Rows Removed by Filter: 999"
    "Planning Time: 0.217 ms"
    "Execution Time: 0.157 ms"
    -- after create
    "Bitmap Heap Scan on info  (cost=5.65..31.65 rows=1 width=147) (actual time=0.071..0.072 rows=1 loops=1)"
    "  Recheck Cond: (ref_id = 3)"
    "  Filter: ((uuid)::text = '05aee2f5-aa10-4a15-9e2c-44166380e3a4'::text)"
    "  Rows Removed by Filter: 199"
    "  Heap Blocks: exact=6"
    "  ->  Bitmap Index Scan on idx_info_ref_id  (cost=0.00..5.65 rows=200 width=0) (actual time=0.018..0.019 rows=200 loops=1)"
    "        Index Cond: (ref_id = 3)"
    "Planning Time: 0.165 ms"
    "Execution Time: 0.095 ms"
    ```
  - total 10k record
    ```sql
    "Seq Scan on info  (cost=0.00..373.00 rows=1 width=148) (actual time=0.087..1.326 rows=1 loops=1)"
    "  Filter: (((uuid)::text = '05aee2f5-aa10-4a15-9e2c-44166380e3a4'::text) AND (ref_id = 3))"
    "  Rows Removed by Filter: 9999"
    "Planning Time: 0.051 ms"
    "Execution Time: 1.343 ms"
    -- after create
    "Bitmap Heap Scan on info  (cost=40.29..311.29 rows=1 width=148) (actual time=0.107..0.593 rows=1 loops=1)"
    "  Recheck Cond: (ref_id = 3)"
    "  Filter: ((uuid)::text = '05aee2f5-aa10-4a15-9e2c-44166380e3a4'::text)"
    "  Rows Removed by Filter: 3199"
    "  Heap Blocks: exact=74"
    "  ->  Bitmap Index Scan on idx_info_ref_id  (cost=0.00..40.28 rows=3200 width=0) (actual time=0.064..0.064 rows=3200 loops=1)"
    "        Index Cond: (ref_id = 3)"
    "Planning Time: 0.251 ms"
    "Execution Time: 0.614 ms"
    ```

## Hash Index
- Hash index là được thiết kế cho query data một cách cực nhanh, khi điểu kiện query là bằng trên index column nào đó, hash index có thể làm việc truy vấn data rất nhanh, và đồng thời hash index cũng xác định trực tiếp vùng lưu trữ dữ liệu mong muốn, chỉ phù hợp cho những tình huống query trên so sánh bằng, như `=` or `in`.
- không như những loại index khác, hash index là khi có một hoạt động thay đổi data ` (inserts, updates, and deletes)` thì sẽ cần rehash lại và quá trình này là tốn kém hơn `BTree` index.
- Để tạo hash index trong postgres, chúng ta sẽ xử dụng `using hash`, ví dụ
```sql
CREATE INDEX hash_name ON table_name USING HASH (column_name);
```
  - câu lệnh trên là thực hiện tại index name `hash_name` trên table `table_name`
  - `column_name` là cột muốn index

- Một số điểm cần chú ý khi tạo hash index là không phù hợp cho việc truy vấn theo `range` hoặc `sorting`. Nhưng trong trường hợp này thì `BTree` sẽ phù hợp hơn.
- Rõ ràng là hash index sẽ có trường hợp cụ thể và giới hạn của nó.

### Ví dụ: 
- Tạo index bằng statement:
```sql
CREATE INDEX idx_info_ref_id_hash ON info USING HASH(uuid);
```

### Kết quả
- với lượng data là 10k record
- trước khi tạo hash index:
```sql
"Seq Scan on info  (cost=0.00..348.00 rows=1 width=148) (actual time=0.076..1.547 rows=1 loops=1)"
"Filter: ((uuid)::text = '05aee2f5-aa10-4a15-9e2c-44166380e3a4'::text)"
"Rows Removed by Filter: 9999"
"Planning Time: 0.302 ms"
"Execution Time: 1.565 ms"
```
- sau khi tạo index:
```sql
"Index Scan using idx_info_ref_id_hash on info  (cost=0.00..8.02 rows=1 width=148) (actual time=0.015..0.016 rows=1 loops=1)"
"Index Cond: ((uuid)::text = '05aee2f5-aa10-4a15-9e2c-44166380e3a4'::text)"
"Planning Time: 2.123 ms"
"Execution Time: 0.039 ms"
```
- với kết quả trên, thì có thể thấy khi đánh hash index kết quả đã nhanh hơn gấp nhiều lần. 
