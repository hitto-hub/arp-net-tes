// メインパッケージを宣言
package main

// 必要な標準ライブラリをインポート
import (
    "fmt"  // 標準出力やエラーメッセージを表示するため
    "net"  // ネットワーク機能を提供するパッケージ（TCP接続に必要）
    "os"   // オペレーティングシステムに関連する機能（プログラム終了など）を提供するパッケージ
)

func main() {
    // サーバーにTCP接続を確立
    conn, err := net.Dial("tcp", "localhost:8080")  // "localhost:8080"は接続先のアドレスとポート
    if err != nil {
        // 接続に失敗した場合、エラーメッセージを表示し、プログラムを終了
        fmt.Println("接続エラー:", err)
        os.Exit(1)  // 異常終了コード1でプログラムを終了
    }
    // main関数が終了する前に接続を閉じるように設定
    defer conn.Close()

    // サーバーに送信するメッセージを作成
    message := "こんにちは、サーバー！"

    // サーバーにメッセージを送信
    _, err = conn.Write([]byte(message))  // 文字列をバイト配列に変換して送信
    if err != nil {
        // メッセージの送信に失敗した場合はエラーメッセージを表示
        fmt.Println("メッセージ送信エラー:", err)
        return
    }

    // サーバーからの応答を受信するためのバッファを作成
    buffer := make([]byte, 1024)  // 1024バイトのバッファを確保

    // サーバーからデータを受信
    n, err := conn.Read(buffer)
    if err != nil {
        // 応答の受信に失敗した場合はエラーメッセージを表示
        fmt.Println("応答受信エラー:", err)
        return
    }

    // サーバーから受信したデータを表示
    fmt.Println("サーバーからの応答:", string(buffer[:n]))
}
