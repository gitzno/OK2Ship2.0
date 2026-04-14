FROM chrislusf/seaweedfs:latest

# Tạo thư mục cấu hình
RUN mkdir -p /etc/seaweedfs

# Copy file security vào (file này phải nằm cùng thư mục với Dockerfile)
COPY security.toml /etc/seaweedfs/security.toml

EXPOSE 9333 8080 8888 8333

ENTRYPOINT ["/usr/bin/weed"]