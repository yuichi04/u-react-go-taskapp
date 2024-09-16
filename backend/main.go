/*
main.goの責務：アプリケーションのエントリーポイント
依存関係の設定、ルーターの初期化、サーバーの起動など、アプリケーション全体の設定と起動を担当
*/

package main

import (
	"backend/db"
)

func main() {
	db := db.NewDB()
}
