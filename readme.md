# Hệ thống phần mềm OK2SHIP

## Requirement

    - Docker

## Background

Dự án này để chuyển đổi dự án cũ OK2Ship sang một hệ thống mới.

Trước đây chỉ tập chung vào giải quyết vấn đề trước mắt là:

- UI cho công nhân sử dụng => lưu database => xuất file
- Gói gọn trong một project winform. Không có kiến trúc.
- Database mssql lưu trữ nặng nề.
- Logic ngày càng phúc tạp.
-  => Vì vậy tôi quyết định chuyển đổi nó thành hệ thống hướng sự kiện thế hệ mới.


## Hướng dẫn cài đặt (Installation):

Khởi tạo docker file
```
docker-compose down && docker-compose build --no-cache && docker-compose up
```

## Tính năng (Features):

Tác giả (Authors): gitzno

### Triển khai nhanh

