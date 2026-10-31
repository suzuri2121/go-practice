package product

import (
	"context"
	"os"
	"slices"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// テスト用DBに接続する。TEST_DATABASE_URL が無ければテストをスキップ
func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatalf("接続失敗: %v", err)
	}
	t.Cleanup(pool.Close) // テスト終了時に閉じる

	return pool
}

// products を空にし、IDの採番も1からに戻す
func resetProducts(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	if _, err := pool.Exec(context.Background(), "TRUNCATE products RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("TRUNCATE失敗: %v", err)
	}
}

func TestStore_List(t *testing.T) {
	pool := newTestPool(t)
	ctx := context.Background()

	t.Run("登録済みの商品がID順で返る", func(t *testing.T) {
		resetProducts(t, pool)

		_, err := pool.Exec(ctx, `
			INSERT INTO products (name, category, price) VALUES
			('りんご', 'fruit', 150),
			('にんじん', 'vegetable', 80)
		`)
		if err != nil {
			t.Fatalf("テストデータ投入失敗: %v", err)
		}

		store := NewStore(pool)
		got, err := store.List(ctx)
		if err != nil {
			t.Fatalf("List がエラー: %v", err)
		}

		want := []Product{
			{ID: 1, Name: "りんご", Category: "fruit", Price: 150},
			{ID: 2, Name: "にんじん", Category: "vegetable", Price: 80},
		}
		if !slices.Equal(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("0件なら空スライス(nilではない)が返る", func(t *testing.T) {
		resetProducts(t, pool)

		store := NewStore(pool)
		got, err := store.List(ctx)
		if err != nil {
			t.Fatalf("List がエラー: %v", err)
		}

		if got == nil {
			t.Error("nil が返った。JSONで null になってしまう")
		}
		if len(got) != 0 {
			t.Errorf("0件のはずが %d 件", len(got))
		}
	})
}