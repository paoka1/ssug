import time
import requests

from urllib.parse import urljoin


class ShortURL:
    """
    短链实例
    """

    def __init__(self, ori_url: str, short_url: str, expire_time: int) -> None:
        self.ori_url = ori_url
        self.short_url = short_url
        self.expire_time = expire_time

    def __str__(self) -> str:
        expire_t = time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(self.expire_time))
        return f"{self.short_url} -> {self.ori_url}, 将于{expire_t}过期"


class SSUG:
    """
    Shauio's short URL generator 调用接口示例
    """

    def __init__(self, acc_key: str, service_url: str, ssug_url: str) -> None:
        """
        acc_key: API key
        service_url: 服务地址
        ssug_url: 短链接服务地址
        """
        self.__service_url = service_url
        self.__acc_key = acc_key
        self.__ssug_url = ssug_url

    def __get_sl(self, url: str) -> dict:
        """
        获取短链接
        """
        data = {
            "key": self.__acc_key,
            "url": url
        }
        r = requests.post(urljoin(self.__service_url, "add"), data=data)
        if r.status_code != 200 and r.json()["data"] == "":
            raise Exception("获取短链接失败：" + r.json()["msg"])
        return r.json()["data"]

    def get_short(self, url: str) -> ShortURL:
        """
        获取短链接
        """
        try:
            data = self.__get_sl(url)
        except Exception as e:
            raise e
        return ShortURL(url, urljoin(self.__ssug_url, data["short_url"]), data["expiration_time"])

    def is_running(self) -> bool:
        """
        检查服务是否正常运行
        """
        try:
            r = requests.get(self.__service_url)
        except requests.ConnectionError:
            return False
        if r.status_code != 200:
            return False
        return True


if __name__ == "__main__":
    # 使用默认配置
    ssug = SSUG("key123456", "http://127.0.0.1:8000/", "http://127.0.0.1:8000/")
    # 检测服务是否运行
    assert ssug.is_running(), "SSUG未在指定地址端口运行"
    # 获取短链接
    print(ssug.get_short("https://www.bilibili.com/video/BV1hq4y1s7VH"))
