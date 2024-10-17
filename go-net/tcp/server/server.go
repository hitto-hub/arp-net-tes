// メインパッケージを宣言
package main

// 必要な標準ライブラリをインポート
import (
    "fmt"  // 標準出力やエラーメッセージを表示するため
    "net"  // ネットワーク機能を提供するパッケージ（TCP接続に必要）
)

func main() {
    // サーバーソケットを作成して、TCPプロトコルで特定のポート（8080）を待ち受けるように設定
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        // サーバーの起動に失敗した場合、エラーメッセージを表示してプログラムを終了
        fmt.Println("サーバー起動エラー:", err)
        return
    }
    // main関数が終了する前にリスナー（サーバーソケット）を閉じるように設定
    defer listener.Close()

    // サーバーが正常に起動したことを確認するためのメッセージを表示
    fmt.Println("サーバーがポート8080で待ち受けています...")

    // クライアントからの接続を待ち受ける無限ループ
    for {
        // クライアントからの新しい接続を受け入れる
        conn, err := listener.Accept()
        if err != nil {
            // 接続の受け入れに失敗した場合はエラーメッセージを表示して、次の接続を待ち受ける
            fmt.Println("接続エラー:", err)
            continue
        }

        // クライアントが接続したことを表示
        fmt.Println("クライアントが接続しました")

        // 新しい接続に対して、handleConnection関数で処理を行う
        // goキーワードを使って、別のゴルーチン（軽量スレッド）で非同期に処理を実行
        go handleConnection(conn)
    }
}

// 接続ごとのクライアント処理を行う関数
func handleConnection(conn net.Conn) {
    // 関数が終了する前に接続を確実に閉じるように設定
    defer conn.Close()

    // クライアントからのデータを格納するためのバッファ（1024バイトの配列）を作成
    buffer := make([]byte, 1024)

    // クライアントからデータを受信
    n, err := conn.Read(buffer)
    if err != nil {
        // データ受信に失敗した場合はエラーメッセージを表示して関数を終了
        fmt.Println("データ受信エラー:", err)
        return
    }

    // 受信したデータを文字列として表示（nは受信したバイト数）
    fmt.Println("受信:", string(buffer[:n]))

    // サーバーからクライアントへの応答メッセージを作成
    response := "サーバーからの応答です"

    // クライアントに応答メッセージを送信
    _, err = conn.Write([]byte(response))
    if err != nil {
        // 応答の送信に失敗した場合はエラーメッセージを表示
        fmt.Println("応答送信エラー:", err)
    }
}
