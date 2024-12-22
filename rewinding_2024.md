# rewinding 2024

- Một năm 2024 sắp kết thúc, là lúc nhìn lại trong 1 năm qua Golang đã có những điểm nào mới.

## Go 1.23:
- Go 1.23 cập nhật một vài điểm mang điến trải nghiệp tốt hơn cho lâp trình viên và hiệu năng tối hợp cho ứng dụng.

### Các Điểm Chính.
1. Build nhanh hơn.
- Nhanh hơn khoảng 15%, tăng công xuất đáng kể trong quá trình CI/CD

2. Sync Package
- `sync.Map` được cải tiến với khả năng xử lý tốt hơn, giúp dễ dàng đáp ứng trong việc xử lý các ứng dụng cần độ trễ thấp.

3. Độc Lập Môi Trường.
- Cho phép lâp trình viên có thể độc lâp từng môi trường, dễ dàng thay đổi qua để chuyển môi trường production.

4. Observability
- Đã tích hợp các tool monitor vào như `Telemetry` `Prometheus` `Grafana`
- Giúp các lâp trình viên dễ debug cũng như tối ưu hiệu năng.

5. Xây dựng Opensource Lớn
- Với hiệu năng cao và khả năng biên dịch nhanh, Golang đã trờ thành ngôn ngữ phổ biến để xây dựng các open source lớn: `Docker`, `Kubernetes`, `Terraform`, ...