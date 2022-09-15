## できること

**ログの保存**
* それぞれのログが1時間で何回発生したかを保存...(1)
* 保守の観点より前の1時間前の(1)の情報を保存しておく
* main.goの`StoreLog`関数参照(楽観ロック対応)
https://github.com/sirogamichandayo/study-redis/blob/main/log_management/main.go#L39

**ログの取得**
* todo

---

### ログの内容
* 名前(name)
* メッセージ(message)
* レベル(level)

---

### 補足

保存はクリーンアーキテクチャに組み込むならドメインサービスにすると思う。
取得の方はユースケースで呼び出す。
