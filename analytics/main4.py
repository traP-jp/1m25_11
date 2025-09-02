import requests
import os
import sys

# --- 設定項目 ---
# 環境変数 'BEARER_TOKEN' を読み込む
bearer_token = os.environ.get("BEARER_TOKEN")

# 環境変数が設定されていない場合は、エラーメッセージを出して終了する
if not bearer_token:
    print("エラー: 環境変数 'BEARER_TOKEN' が設定されていません。")
    sys.exit(1)

# APIのエンドポイントURL
url = "https://q.trap.jp/api/v3/messages"

# クエリパラメータ
params = {
    "word": '"るるん"',
    "sort": "createdAt"
}

# リクエストヘッダー
headers = {
    "accept": "application/json",
    # 読み込んだ環境変数の値をヘッダーに設定
    "Authorization": f"Bearer {bearer_token}"
}


# --- リクエストの実行 ---
try:
    # GETリクエストを送信
    response = requests.get(url, params=params, headers=headers)

    # ステータスコードが200番台でない場合に例外を発生させる
    response.raise_for_status()

    # レスポンスをJSON形式で取得して表示
    data = response.json()
    print(data)

# --- エラー処理 ---
except requests.exceptions.HTTPError as e:
    print(f"HTTPエラーが発生しました: {e}")
    print(f"レスポンスボディ: {e.response.text}")
except requests.exceptions.RequestException as e:
    print(f"リクエスト中にエラーが発生しました: {e}")
