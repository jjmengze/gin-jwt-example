# gin-jwt-example

 示範一個透過gin實現JWT的一個範例，對相關欄位進行簡易的解說。


 jwt全名為 [Json web token](https://jwt.io/)，一般會用來做Client的身份驗證。

JWT 由三部分組成，分別分由base64(<Header>).base64(<Payload>).base64(<Signature>) ，資料透過base64 decode 並且由`.`做資料分隔。

## Header

通常有兩個部分組成，分別是 Token 類型以及最後JWT要用什麼格式做簽章（例如：HMAC、SHA256或者RSA）
```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```
## Payload
這個部分會放一些使用者資料，例如User-ID、User-Email。
注意千萬不要放機密資料base64 decode資料就被看光光囉!

### Reserved
標準的JWT有預留一些欄位例如
```bash
iss (Issuer) - JWT簽發人
sub (Subject) - JWT的用戶
aud (Audience) - 接收JWT的者
exp (Expiration Time) - JWT過期時間，注意過期時間必須要大於簽發時間
nbf (Not Before) - 定義在什麼時間之前，該JWT都是不可用的（可以理解為開始營運時間）
iat (Issued At) - JWT的簽發時間
jti (JWT ID) - JWT的唯一ID，作為一次性token使用
```

### Public

JWT預設幫你做了一些命名，例如姓名，生日等等可以看[JWT 相關規範](https://www.iana.org/assignments/jwt/jwt.xhtml#claims)

## Private

如果上面說的都不符合你要的需求的話，我們可以自訂欄位。 例如使用者ID、使用者姓名等等。

```json
{
     "user_id": "1234567890",
     "username": "jason"
}
```

 ## Signature

這一個欄位是對上述兩個部分（Header與Payload）進行一個簽章，既然要簽章必定會有一個鑰匙（Secret）這把鑰匙通常放在JWT發行者的手上。

另外簽章過後如果header或是payload被串改那JWT發行者只要透過他的鑰匙簡單驗證一下就知道有沒有被串改了。

 基本上簽章的形式會如下所示。
 ```bash
 HMACSHA256(
   base64Encode(header) + "." +
   base64Encode(payload),
   secret)
 ```
 依照剛剛的範例我們來簽章一下吧！
 1. Header
 ```json
{
  "alg": "HS256",
  "typ": "JWT"
}
 ```
 2. Playload
 ```json
 {
      "user_id": "1234567890",
      "username": "jason"
 }
 ```
 3. key
 請自己產生xD,或是到這個[jwt網站](https://jwt.io/)他可以幫你生一個

 ```bash
 eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
 eyJ1c2VyX2lkIjoiMTIzNDU2Nzg5MCIsInVzZXJuYW1lIjoiamFzb24ifQ.
 _1NpLsmQuuXYm6LdE3Tpt7s3oahenR6-CyCXPFxCa-Y
 ```

 這三個資料會用`.`來區分開來。

 ## 使用時機
這邊只做很簡單的介紹，JWT還可以結合session token 或是結合oath2 做出比較複雜的應用，這裡不展開說明複雜的情境。


使用者透過帳號密碼進行登入

server 驗證帳號密碼,回應使用者一個JWT token

使用者可以拿這組 token 做其他事情，例如 post 文章等。

> 問題是使用者怎麼拿 token 去做事？？
> 透過header ?
> 透過body ?
> 還是透過什麼 ？

只要使用者發請求在Header帶入Authorization: Bearer <JWT-token> ，當Server收到並成功解析Token就可以讓使用者做後續的操作囉！


簡易的JWT說明就到這裡，剩下的就是結合golang web framework gin了，就自行看code吧。