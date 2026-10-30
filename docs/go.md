# Go 学習メモ

Go で API を作るための基礎をまとめたもの。

---

## 1. Go のバックエンド事情

Go には「全部入り」の公式フレームワークはなく、**標準ライブラリ（net/http）＋必要な部品を組み合わせる**のが Go の文化。

### 標準ライブラリ

- Go 1.22 で `net/http` のルーティングが強化され、`GET /users/{id}` のようなメソッド指定・パスパラメータが標準で書ける
- 小〜中規模の API ならフレームワークなしでも十分

### 主な Web フレームワーク

| 名前 | 特徴 |
|---|---|
| Gin | 一番利用者が多い定番。情報量が多い |
| Echo | Gin と同じ立ち位置。好みで選ぶ |
| chi | net/http 完全互換の軽量ルーター。標準寄り |
| Fiber | Express 風。net/http 非互換なので既存ミドルウェアが使えない点に注意 |

### その他よく使う部品

| 用途 | ライブラリ |
|---|---|
| DI | 手動で依存を渡すのが普通。大きくなったら Wire / Fx |
| DB | GORM（ORM）、sqlc（SQL から型安全なコード生成）、ent、pgx（PostgreSQL ドライバ） |
| 設定 | Viper、または環境変数を直接読む |

SQL を練習したいなら、ORM より **pgx（生 SQL）や sqlc** が向いている。

---

## 2. パッケージとモジュール

```go
package main

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)
```

- `package main` ＋ `func main()` がプログラムの起動地点（C# の `static void Main()`）。ファイル名 `main.go` は慣習で必須ではない
- `import` は C# の `using`
- **使っていない import や変数はコンパイルエラー**
- `context` や `net/http` などの標準ライブラリは Go 本体に同梱されていてダウンロード不要

### ライブラリはどこにあるのか

`github.com/...` という import パスは URL っぽいが、**ライブラリを一意に識別する名前**で、初回ダウンロード時の取得先を示すだけ。ビルドのたびにネットワークへ取りに行くわけではない。

| | npm | Go |
|---|---|---|
| 置き場所 | プロジェクトごとの `node_modules` | PC 全体で共有するモジュールキャッシュ（`~/go/pkg/mod`） |
| 依存の記録 | `package.json` | `go.mod` |
| ロックファイル | `package-lock.json` | `go.sum`（中身のハッシュで改ざん検証） |

- 一度 `go get` すればオフラインでもビルドできる
- `go mod vendor` でプロジェクト内の `vendor` フォルダにコピーすることも可能（閉じたネットワークでのビルド用）

---

## 3. 変数宣言

```go
var count int       // 宣言のみ → 0 が入る
var name string     // 宣言のみ → "" が入る
var p Product       // 宣言のみ → 全フィールドがゼロ値の構造体

var count2 int = 10 // 型を明示して宣言＋代入
var count3 = 10     // 型推論で宣言＋代入
count4 := 10        // 関数内の宣言＋代入の定番

count = 5           // 代入のみ

var count5 := 10    // NG：var と := の二重指定でコンパイルエラー
```

- `var` を書いたら `=`、書かないなら `:=`
- `:=` は **関数の中でしか使えない**。パッケージレベルでは `var`
- 型は**変数名の後ろ**に書く（`int x` ではなく `x int`）

### ゼロ値

宣言だけした変数には型ごとのゼロ値が自動で入る。C# のように「未初期化で使うとエラー」にはならない。

| 型 | ゼロ値 |
|---|---|
| 数値 | `0` |
| 文字列 | `""` |
| bool | `false` |
| ポインタ・スライス・map・error など | `nil` |
| 構造体 | 各フィールドがそれぞれのゼロ値 |

### `:=` で複数の変数を扱うとき

左辺に**少なくとも 1 つ新しい変数があれば** `:=` が使える。既存の変数には代入だけが行われる。全部既存なら `=` を使う。

```go
rows, err := pool.Query(...)  // rows も err も新規
count, err := 0, nil          // count が新規なので OK。err は再代入
```

---

## 4. 構造体（struct）

```go
type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}
```

- Go にクラスはなく、struct が C# の class / record の役割を担う
- バッククォートの `json:"id"` は**タグ**。C# の `[JsonPropertyName("id")]` と同じ役割
- **大文字始まり＝public、小文字始まり＝private**。小文字にすると json パッケージから見えず JSON に出力されない

---

## 5. エラー処理

```go
pool, err := pgxpool.New(ctx, dsn)
if err != nil {
	log.Fatal(err)
}

if err := pool.Ping(ctx); err != nil {
	log.Fatal(err)
}
```

- **Go には例外（try/catch）がない**。失敗はエラーを戻り値で返し、呼び出し側が `if err != nil` で毎回チェックする
- 関数は**複数の戻り値**を返せる
- `nil` は C# の `null`
- `if 代入; 条件` の形で、変数をその if 文の中だけで使うのが定番
- `if` の条件に丸括弧は付けない
- `log.Fatal` はエラーを表示してプログラムを終了する

---

## 6. defer

**defer を書いた関数が終了するタイミング**で実行される。C# の `using` や `finally` に近い。

```go
func example() {
	defer fmt.Println("1")
	defer fmt.Println("2")
	defer fmt.Println("3")
	fmt.Println("本体")
}
// 出力: 本体 → 3 → 2 → 1
```

- 常に main の終了時ではない。main に書けば main の終了時、ハンドラに書けば**そのリクエストの処理が終わるたび**
- 複数ある場合は**後に書いたものから**実行される（スタック順）
- `return` やエラーで途中で抜けても実行される。だから「開いた直後に defer で閉じる」のが定番
- **`log.Fatal` や `os.Exit` では実行されない**。main.go の `defer pool.Close()` は実は実行されずに終了している。きちんと作るならグレースフルシャットダウンを書く

---

## 7. 接続プール（pool）

DB への接続を作るのは重い処理（TCP 接続・認証・セッション準備）。**接続を開いたまま何本か保持しておき、借りて返して使い回す**のが接続プール。

- 同時に来たリクエストに空いている接続を割り当て、足りなければ作り、上限を超えたら待たせる
- C# の ADO.NET でも裏で動いている（`SqlConnection.Close()` は実際にはプールへの返却）。Go の pgx ではそれを `pool` という変数として明示的に持っている

| 処理 | 意味 |
|---|---|
| `pgxpool.New` | プールを作る（起動時に 1 回） |
| `pool.Query` | プールから接続を借りて SQL を実行 |
| `rows.Close()` | **接続をプールに返す** |

`defer rows.Close()` を忘れると接続を借りっぱなしになり、プールが枯渇してサーバーが固まる。

### Ping

`pool.Ping(ctx)` は DB に軽い問い合わせを送って、実際に接続できるかを確認する。

- `pgxpool.New` はその時点で接続しないことがある（遅延接続）ので、パスワード違いや DB 停止でもサーバーが起動してしまう
- 起動時に Ping して、問題をその場で表面化させる（フェイルファスト）
- `/health` などのヘルスチェック用エンドポイントでも使う

---

## 8. context

**「この処理はもう続けなくていい」という合図を、呼び出しの奥まで伝える仕組み**。C# の `CancellationToken` に相当。

例：重い SQL の実行中にユーザーがブラウザを閉じた場合、`r.Context()` がキャンセル状態になり、pgx が DB にクエリの中止を伝える。

```go
// 何の制限もない大元の context（起動処理などで使う）
ctx := context.Background()

// リクエストに紐付いた context（クライアントが切断したらキャンセル）
ctx := r.Context()

// 制限時間付きの context（3 秒経ったら自動でキャンセル）
ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
defer cancel()
rows, err := pool.Query(ctx, "SELECT ...")
```

- キャンセルや時間制限が関係しそうな関数は**第 1 引数に ctx を取る**のが慣習
- context は親子関係になっていて、親がキャンセルされると子も全部キャンセルされる
- 覚え方：**ハンドラの中では `r.Context()`、それ以外は `context.Background()`**

---

## 9. ServeMux（ルーター）

**URL を見て、どの関数に処理させるかを振り分けるルーター**。Mux はマルチプレクサ（multiplexer）の略。

```go
mux := http.NewServeMux()

mux.HandleFunc("GET /products", listProducts)
mux.HandleFunc("GET /products/{id}", getProduct)
mux.HandleFunc("POST /products", createProduct)

http.ListenAndServe(":8080", mux)

func getProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") // /products/3 なら "3"（文字列）
	fmt.Fprintf(w, "product id: %s", id)
}
```

- どれにも一致しない → 自動で 404、パスは一致するがメソッドが違う → 自動で 405
- 複数に一致する場合は、より具体的なパターンが優先
- パスパラメータは常に**文字列**。数値にするなら `strconv.Atoi`
- `ListenAndServe` の第 2 引数に mux を渡すことで、全リクエストが mux 経由で各関数に振り分けられる
- Gin や chi なども、根本的には ServeMux の高機能版

### 登録と実行は別

- `mux.HandleFunc(...)` 自体は**起動時に 1 回だけ**実行され、関数を登録するだけ
- 渡した無名関数の中身は**リクエストが来るたびに**実行される

---

## 10. ハンドラの引数

```go
func(w http.ResponseWriter, r *http.Request) { ... }
```

- **「変数名 型」の順**。`w` と `r` はただの変数名だが、Go では全員この名前を使う慣習
- `http.ResponseWriter` は「net/http パッケージの ResponseWriter 型」。別パッケージの型は毎回パッケージ名を付ける
- `*http.Request` は Request 構造体への**ポインタ**。大きな構造体をコピーせず参照で渡している
- `ResponseWriter` に `*` がないのは**インターフェース**だから
- `HandleFunc` はこの形の関数を要求しているので、型や順番を変えるとコンパイルエラー

### w（ResponseWriter）＝ 返事を書く紙

```go
w.Header().Set("Content-Type", "application/json") // ヘッダを設定
w.WriteHeader(http.StatusCreated)                   // ステータスコード（省略すると 200）
w.Write([]byte("hello"))                            // 本文を書き込む
```

### r（Request）＝ 届いた手紙

```go
r.Method               // "GET" や "POST"
r.URL.Path             // "/products/3"
r.URL.Query().Get("q") // ?q=xxx の値
r.PathValue("id")      // ルートの {id} の値
r.Header.Get("Authorization")
r.Body                 // POST された本文
r.Context()            // このリクエストの context
```

---

## 11. スライスとループ

```go
products := []Product{}
for rows.Next() {
	var p Product
	if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	products = append(products, p)
}
if err := rows.Err(); err != nil {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}
```

- `[]Product{}` は**スライス**（C# の `List<Product>`）。空で初期化すると 0 件のとき JSON が `null` ではなく `[]` になる
- Go のループは `for` だけ。`while` も `for 条件 {}` で書く
- `&` は変数のアドレスを渡す記号。`rows.Scan` に「ここに書き込んで」と伝える（C# の `out` に近い）
- SELECT の列順と Scan の引数順を合わせる
- **`append` は新しいスライスを返すので再代入が必要**（C# の `list.Add()` と違う）
- ループ中のエラーは `rows.Next()` が false になるだけなので、ループ後に `rows.Err()` で確認する

---

## 12. 完成した main.go

```go
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://app:app@localhost:5432/shop?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		rows, err := pool.Query(r.Context(), "SELECT id, name, category, price FROM products ORDER BY id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		products := []Product{}
		for rows.Next() {
			var p Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

全体の流れ：**DB に接続 → `/products` にハンドラを登録 → サーバーを起動して待つ → リクエストが来たら SQL を実行して JSON で返す**

---

## C# から来て最初に慣れるべきこと

1. **`if err != nil` によるエラー処理**（例外がない）
2. **大文字・小文字による公開範囲**
3. **`append` の再代入**
4. **型は名前の後ろ**
5. **defer は書いた関数の終了時に実行**