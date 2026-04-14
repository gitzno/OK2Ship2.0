#!/bin/bash

# 1. Khởi động SQL Server ở chế độ background
/opt/mssql/bin/sqlservr &
pid=$!

echo "Đang chờ SQL Server khởi động..."

# 2. Vòng lặp kiểm tra trạng thái SQL Server
for i in {1..60}; do
    /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P "$MSSQL_SA_PASSWORD" -C -Q "SELECT 1" > /dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        echo "SQL Server đã sẵn sàng sau $i lần thử!"
        
        # 3. Thực thi tự động toàn bộ file .sql
        echo "Bắt đầu hợp nhất và thực thi toàn bộ script trong thư mục scripts..."
        
        # Kỹ thuật: 
        # - cat /path/*.sql: Nối tất cả nội dung file theo thứ tự alphabet (0... 1... 2...)
        # - | : Đẩy toàn bộ luồng dữ liệu vào duy nhất một phiên sqlcmd
        cat /usr/src/app/scripts/*.sql | /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P "$MSSQL_SA_PASSWORD" -C
        
        if [ $? -eq 0 ]; then
            echo "Khởi tạo toàn bộ Database, Partition và Tables hoàn tất thành công!"
        else
            echo "LỖI: Có lỗi xảy ra trong quá trình thực thi các file SQL."
        fi
        
        break
    fi
    
    echo "SQL Server chưa sẵn sàng, tiếp tục chờ (lần thử $i)..."
    sleep 2
done

# 4. Giữ container luôn chạy
echo "Hệ thống đang vận hành..."
wait $pid