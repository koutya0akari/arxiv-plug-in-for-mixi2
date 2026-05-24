# arXiv plug-in for mixi2

arXiv の新着論文を取得し、mixi2 のコミュニティに自動投稿する Plugin です。デフォルトでは次のカテゴリを監視します。

- `math.CT`
- `math.AG`
- `math.AT`
- `math.RT`
- `math.NT`
- `math.AC`
- `math.KT`
- `math.OA`
- `math.FA`
- `math.RA`

カテゴリごとに別々の mixi2 Plugin アプリケーションアカウントと投稿先コミュニティを設定できます。

## 動作

- `https://arxiv.org/list/{category}/new` から新着一覧を取得します。
- `New submissions` セクションだけを投稿対象にします。
- Cross-lists と Replacements は投稿しません。
- 1 論文につき 1 投稿を作成します。
- 投稿済みの arXiv ID は `data/posted.json` に保存します。
- 初回実行時は投稿せず、その時点の論文 ID を記録するだけにします。
- あるカテゴリで取得や投稿に失敗しても、他のカテゴリの処理は継続します。
- 投稿に失敗した論文 ID は投稿済みにせず、次回以降に再試行できるようにします。

mixi2 の投稿本文は 149 文字以内に制限されているため、長いタイトルは省略します。arXiv の URL は必ず残します。

## ローカルでの実行

Go 1.24.6 以上が必要です。

投稿せずに取得結果と投稿予定を確認します。

```sh
go run ./cmd/arxiv-mixi2-plugin --dry-run
```

投稿せずに現在の新着論文を記録します。

```sh
go run ./cmd/arxiv-mixi2-plugin --initialize-only
```

カテゴリを指定して実行します。

```sh
go run ./cmd/arxiv-mixi2-plugin --categories math.CT,math.AG
```

インストール済みコミュニティの ID を確認します。このコマンドは `MIXI2_<CATEGORY>_COMMUNITY_ID` を必要とせず、Client ID / Client Secret / Token URL / API Address だけで実行できます。

```sh
go run ./cmd/list-communities --categories math.CT
```

## GitHub Actions

`.github/workflows/post.yml` は毎日 12:00 / 13:30 / 15:00 JST に実行されます。新しい投稿済み ID が記録された場合は、`data/posted.json` を自動で commit します。

## mixi2 Plugin 側の設定

mixi2 Developer Platform でカテゴリごとに Plugin アプリケーションを作成し、対応する投稿先コミュニティにインストールします。この Plugin は arXiv から取得した内容をコミュニティに投稿するだけなので、イベント受信用の Webhook URL や Stream Address は使いません。

カテゴリごとに次の設定を行ってください。

1. mixi2 Developer Platform にログインし、「新規アプリケーション」から Plugin アプリケーションを作成します。
2. ID と表示名を入力します。ID は mixi2 上で表示され、後から変更できません。
3. Requirement で `Community.Post.Create` パーミッションを追加します。
4. 作成後、アプリケーション詳細画面の「認証情報」を開きます。
5. Client Secret を生成します。
6. 投稿先コミュニティへ Plugin をインストールします。
7. Client ID、Client Secret、投稿先 Community ID を控えます。Token URL と API Address は共通の値として控えます。

```text
Client ID
Client Secret
Token URL
API Address
Community ID
```

カテゴリごとに以下の GitHub Secrets を設定してください。`MATH_CT` の部分はカテゴリに応じて `MATH_AG`, `MATH_AT`, `MATH_RT`, `MATH_NT`, `MATH_AC`, `MATH_KT`, `MATH_OA`, `MATH_FA`, `MATH_RA` に置き換えます。

```text
MIXI2_MATH_CT_CLIENT_ID
MIXI2_MATH_CT_CLIENT_SECRET
MIXI2_MATH_CT_COMMUNITY_ID
```

共通の GitHub Secrets として以下も設定してください。

```text
MIXI2_TOKEN_URL
MIXI2_API_ADDRESS
```

カテゴリごとの環境変数の prefix は、arXiv カテゴリ名の `.` を `_` に置き換えて大文字化したものです。たとえば `math.CT` は `MIXI2_MATH_CT` になります。

`STREAM_ADDRESS` と `SIGNATURE_PUBLIC_KEY` は設定不要です。この Plugin は DM、メンション、リプライなどのイベントを受信せず、`CreatePost` API を `community_id` 付きで使用します。

Client Secret は秘密情報です。README、ソースコード、ログ、`data/posted.json` には書かないでください。漏えいした場合は mixi2 Developer Platform で再発行し、GitHub Secrets を更新してください。

## CLI オプション

- `--categories`: カンマ区切りの arXiv カテゴリ一覧。デフォルトは `math.CT,math.AG,math.AT,math.RT,math.NT,math.AC,math.KT,math.OA,math.FA,math.RA` です。
- `--state`: 投稿済み ID を保存する JSON ファイルのパス。デフォルトは `data/posted.json` です。
- `--dry-run`: 投稿も state 保存も行わず、処理内容だけをログ出力します。
- `--initialize-only`: 投稿せず、取得した論文 ID を投稿済みとして記録します。
- `--initialize-on-empty`: `data/posted.json` が空のときに投稿せず初期化だけ行うかを指定します。デフォルトは `true` です。
- `--post-interval`: 同一カテゴリ内の投稿間隔。デフォルトは `4s` です。
- `--request-timeout`: arXiv 取得と mixi2 投稿のタイムアウト。デフォルトは `30s` です。

## 謝辞

このプロジェクトは arXiv の公開ページから取得した情報を利用します。arXiv 公式に承認・提供されているプロジェクトではありません。

参考:

- https://github.com/Krypf/arxiv_bot
- https://developer.mixi.social/docs/getting-started/quickstart
- https://developer.mixi.social/docs/getting-started/concepts
- https://developer.mixi.social/docs/guides/plugin
- https://developer.mixi.social/docs/guides/api-usage
- https://developer.mixi.social/docs/guides/sdk
- https://info.arxiv.org/help/api/index.html
