from abc import ABC, abstractmethod
from typing import Optional, Any


class ICacheService(ABC):
    @abstractmethod
    async def get(self, key: str) -> Optional[Any]:
        """Lấy giá trị từ Cache theo Key."""
        pass

    @abstractmethod
    async def set(self, key: str, value: Any, ttl: Optional[int] = None) -> bool:
        """
        Lưu giá trị vào Cache.
        :param key: Khóa định danh (vd: 'user:123:stamp')
        :param value: Giá trị cần lưu
        :param ttl: Thời gian hết hạn (Time To Live) tính bằng GIÂY.
        """
        pass

    @abstractmethod
    async def delete(self, key: str) -> bool:
        """Xóa một Key khỏi Cache."""
        pass

    @abstractmethod
    async def close(self):
        pass