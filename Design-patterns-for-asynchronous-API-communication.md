# Design-patterns-for-asynchronous-API-communication
- Event-Driven Architecture có thể giảm sự phụ thuộc lẫn nhau, tăng sự an toàn, làm ứng dụng của bạn dễ dàng scale. Nhưng thiết kế này là một chủ đề và nhiệm vụ không hề dễ dàng.

- EDA là mô hình cho phép giảm thiểu các service phụ thuộc lẫn nhau trong kiến trúc microservice, nó loại bỏ kiểu mô hình các service phải biết cũng như query lẫn nhau bằng việc sử dụng asynchronours APIs. Điều này cho phép các người thiết kế hệ thống giảm latency và làm việc thay đổi các service một cách độc lập.

- Tóm lại, cách tiếp cận là sử dụng `message broker`, quản lý số lượng topic. Một service sẽ producers to topic, while other service consume from the topic. Mục tiêu khi giới thiệu  `message broker` là giảm sự phụ thuộc giữa các service này. Có nghĩa là khi service A có một vài thông tin, và service B sẽ không cần lấy thông tin trực tiếp từ service A. Thay vào đó A sẽ public thông tin vào topic. Service B có thể consume thông tin và lưu trữ data trong chính nó, hoặc làm một số công việc riêng của A.

- Bỏ qua nhiều thử thách về công nghệ với cách tiếp cận này thì một message từ nơi sản xuất (producing) có thể đến nơi tiêu thụ(consuming). Cũng như việc kết hợp trong mô hình này với kiểu database truyền thống. Những data nào là hữu ích để có thể gửi qua topic và làm thế nào để gửi data từ nơi này đến nơi khác một cách dễ dàng và an toàn. Trong những tình huống các các service là thực sự khác nhau.

## Giải Quyết Vấn Đề
- Để có thể giải quyết được những vấn đề trên, chúng ta sẽ sử dụng `Apache Kafka` là một công nghệ cơ bản và cũng hữu ích cho `message broker` ở phần trên. Bên cạnh một số thuận lợi mà tất cả ai cũng biết về kafka là khả năng `scale, recover và process 1 cách nhanh chống`, thì kafka có một vài thuốc tính làm topíc hữu ích hơn những queue truyền thống hoặc message bus:
  - Mỗi message trong kafka sẽ có 1 key liên kết với nó. Những key này có thể được sử dụng hữu ích trong những usecase khác nhau.
  - Mỗi message chỉ có thể delete bởi key-compaction hoặc retention setting.
  - Bất kì mội consumer(or consumer group) điều được quản lý bởi 1 offset của chính nó, có nghĩa là sẽ đọc lập với nhau.
    - ví dụ: topic A, producer 1 message M1, thì lúc này cónusmer G1 và G2 sẽ consumer đọc lập không liên quan gì đến offset của nhau.

-> Kết quả cuối cùng của những thuộc tính này là những message trong topic sẽ ở lại vĩnh viễn cho đến khi được dọn rác bởi retention hoặc compaction.

Với những thuộc tính trên bạn sẽ design topic theo tình huống sau:
- Entity topics - "Hiện tại là API X"
- Event topics - "X xảy ra"
- Request and response topics" - "làm API X" → "API xong"

## Entity topics - Sự Thật của vấn đề
- Là một cách cực kì hữu ích để sử dụng kafka cho việc di chuyển data giữa service, và nó sẽ báo cáo trạng thái hiện tại của object. 
- Một ví dụ thế này, bạn có một Service Customers, và biết tất cả thông tin về khách hàng của công ty. Nếu thiết kế api theo synchoronous, bạn sẽ 1 service và cấu hình api để những service khác gọi api nếu cần thông tin. Kafka, thay vào đó sẽ publish to Customer Topic, và tất cả những service khác sẽ consumer từ nó.
- Như vậy customer service có thể dump data trực tiếp trên hệ thống của của chính nó. Cách suy nghĩ này sẽ phù hợp với mô hình mỗi lưu data như họ muốn. Có thể bỏ vài field hoặc mapping thông tin lại theo cách của riêng mình. Lúc này mỗi service sẽ có data của chính nó khi cần có thể sử dụng luôn.
- 