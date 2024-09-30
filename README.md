# 概要

店舗経営スマホゲーム開発プロジェクト『Ringo』バックエンド実装となります。

gRPCを用いてクライアントと通信を行い、MySQLにデータを保存します。

インフラストラクチャについては抽象化していますが，開発環境ではDockerを用いてMySQLを起動しています。

# 使い方

## ローカル環境での起動方法

- リポジトリをクローンします
- プロジェクトルートに移動します
- ./dockerに移動して `$ docker-compose up -d` でコンテナを起動します
- プロジェクトルートに戻り `$ go run ./cmd/.` でサーバーを起動します
  - 4444ポートでサーバーが起動します

### gRPCの実行

現在サーバリフレクションを有効にしているためgRPCurlが利用できます．

#### ユーザ登録

`$ grpcurl -plaintext localhost:4444 ringosu.Ringo.SignUp`

サーバからuserIdとrowPasswordを受け取り，次のリクエストでログインできます．

#### ユーザログイン

`$ grpcurl -plaintext -d '{"userId": "UUUU", "rowPassword": "XXXX"
}' localhost:4444 ringosu.Ringo.Login`

ログインに成功するとサーバからアクセストークンを受け取ることができます．

ゲーム内での通信は全てアクセストークンを用いて行います．

## テスト

プロジェクトルートで `$ go test ./...` でテストを実行します

現在以下のテストが実装されています

- ユニットテスト
- Dockerを用いたMySQLへのリクエストテスト
- Dockerを用いたE2Eテスト
