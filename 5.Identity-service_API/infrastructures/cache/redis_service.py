import json
from typing import Optional, Any

from redis.asyncio import Redis, from_url

from core.interfaces.cache_interface import ICacheService

class RedisService(ICacheService):
    def __init__(self, redis_url: str):
        self.redis: Redis = from_url(redis_url, decode_responses=True)


    async def get(self, key: str) -> Optional[Any]:
        value = await self.redis.get(key)
        if value is None:
            return None

        # Redis lưu dữ liệu dưới dạng bytes, cần decode ra string
        return value.decode("utf-8")

    async def set(self, key: str, value: Any, ttl: Optional[int] = 1500) -> bool:
        """
        Ghi đè hàm set.
        Nếu có truyền tham số ttl, dữ liệu sẽ tự động bốc hơi sau 'ttl' giây.
        """
        try:
            # Xử lý an toàn dữ liệu: Redis chỉ nhận string, bytes hoặc số.
            # Nếu truyền vào Dict/List thì tự động chuyển thành chuỗi JSON.
            if not isinstance(value, (str, bytes, int, float)):
                value = json.dumps(value)

            # Tham số 'ex' trong redis-py đại diện cho expire (tính bằng giây)
            result = await self.redis.set(name=key, value=value, ex=ttl)

            return bool(result)
        except Exception as e:
            # Log lỗi tại đây (Redis sập, rớt mạng...)
            print(f"Lỗi khi lưu Redis Key '{key}': {e}")
            return False

    async def delete(self, key: str) -> bool:
        try:
            # redis.delete trả về số lượng key đã xóa thành công
            result = await self.redis.delete(key)
            return result > 0
        except Exception as e:
            print(f"Lỗi khi xóa Redis Key '{key}': {e}")
            return False

    async def close(self):
        """Dọn dẹp và đóng kết nối Redis pool khi tắt ứng dụng"""
        # Lưu ý: Tùy phiên bản thư viện redis, có thể dùng .close() hoặc .aclose()
        if hasattr(self.redis, "aclose"):
            await self.redis.aclose()
        else:
            await self.redis.close()