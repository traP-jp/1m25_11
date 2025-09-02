import requests
import os
import sys

# --- 設定項目 ---
# 環境変数 'BEARER_TOKEN' を読み込む
auth_token = os.environ.get("TRAQ_AUTH_TOKEN")

# 環境変数が設定されていない場合は、エラーメッセージを出して終了する
if not auth_token:
    print("エラー: 環境変数 'TRAQ_AUTH_TOKEN' が設定されていません。")
    sys.exit(1)

# APIのエンドポイントURL
url = "https://traqing.cp20.dev/api/stamps"

params = {
    "stampId": "099ef742-c9de-4b8f-9201-82acfad96521",
    "isBot": "false",
    "groupBy": "message",
    "orderBy": "date",
    "order": "asc",
    "limit": "11",
    "offset": "0",
    "after": "2024-09-01T15:00:00.000Z",
    "before": "2025-09-01T14:59:59.999Z"
}

cookies = {
    'traq-auth-token': auth_token
}

try:
    # 4. Cookie付きでGETリクエストを送信
    response = requests.get(url, params=params, cookies=cookies)

    # 5. エラーがあれば例外を発生させる (ステータスコードが2xxでない場合)
    response.raise_for_status()

    # 6. レスポンスの内容をJSON形式で表示
    print("リクエストに成功しました。")
    print("ステータスコード:", response.status_code)
    print("レスポンス:")
    print(response.json())

except requests.exceptions.RequestException as e:
    print(f"リクエスト中にエラーが発生しました: {e}")
except requests.exceptions.JSONDecodeError:
    print("レスポンスをJSONとしてデコードできませんでした。")
    print("レスポンス (テキスト):", response.text)
