# Cài mod
- Chạy lệnh
  ```bash
  go get github.com/PuerkitoBio/goquery
  ```
# Sử dụng
- Tải
```bash
./downloader -url=<URL post cần tải | url dạng https://voz.vn/t/tong-hop-nhung-addon-chat-cho-firefox-chromium.682181/> -ref=<referer header> (-output=<tên thư mục chứa html | bỏ trống thì tự lấy phần "tong-hop-nhung-addon-chat-cho-firefox-chromium.682181">)
```
- Fix link dẫn tới post (fix các kiểu link dạng /p/post-id /goto/post?id= hoặc /t/bai-viet/post-id) thành link gốc dạng /t/baiviet/#post-id để dễ xử lý sau này
```bash
./downloader -scan=<thư mục chứa bài viết đã tải về> -domain=<tên miền của forum sài xenforo đã tải về>
```
# Lưu ý
- Chỉ lấy post header (phần tiêu đề bài viết) và post body (phần nội dung)
- Proxy muốn sử dụng SOCKS5 thì sài dạng `-proxy=socks5://ip:port` còn nếu proxy thì sài `-proxy=http://ip:port`
# Ý tưởng
- Sửa lại link ảnh public chạy qua proxy.php của voz
- Tự fetch ảnh attachment và trả về dạng base64
