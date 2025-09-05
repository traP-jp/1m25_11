# create jsonl

import json
import os
import sys
import requests

bearer_token = os.environ.get("BEARER_TOKEN")
auth_token = os.getenv('TRAQ_AUTH_TOKEN')

max_messages = 10000 # スタンプごとに取得する最大メッセージ数。最終的にはwhileの条件式が常にTrueになる程度に大きくしたいが、offsetは9900なでしか指定できないので10000にした
messages_per_request = 100 # 1回のリクエストで取得するメッセージ数。テスト用に小さくしているが、本番では100にする
stamps_file = 'stamps.json'

since = "2006-01-02T15:04:05Z"
until = "2026-01-02T15:04:05Z"

try:
    # with文でファイルを開く
    with open(stamps_file, 'r', encoding='utf-8') as f_stamps:
        stamps = json.load(f_stamps)

    all_traQ_messages = []
    all_traQing_messages = []

    for stamp in stamps:
        # traQからデータ取得
        traQ_url = "https://q.trap.jp/api/v3/messages"
        traQ_headers = {
            "accept": "application/json",
            "Authorization": f"Bearer {bearer_token}"
        }
        traQ_messages_count = 0
        traQ_messages = []
        while traQ_messages_count < max_messages:
            traQ_params = {
                "word": f"\":{stamp['name']}:\"",
                "after": "2006-01-02T15:04:05Z",
                "before": "2026-01-02T15:04:05Z",
                "bot": "false",
                "limit": str(messages_per_request),
                "offset": str(traQ_messages_count),
                "sort": "createdAt"
            }

            traQ_response = requests.get(traQ_url, params=traQ_params, headers=traQ_headers)
            traQ_response.raise_for_status()
            traQ_messages.extend(traQ_response.json()['hits'])
            traQ_messages_count += messages_per_request # 404の場合もカウントするのでlen(traQ_response.json()['hits'])ではない
            if len(traQ_response.json()['hits']) < messages_per_request:
                break
        all_traQ_messages.append({
            "stamp_name": stamp['name'],
            "stamp_id": stamp['id'],
            "messages": traQ_messages
        })

        # traQingからデータ取得
        traQing_url = "https://traqing.cp20.dev/api/stamps"
        traQing_cookies = {
            'traq-auth-token': auth_token
        }

        traQing_messageIds_count = 0
        traQing_messageIds = []
        while traQing_messageIds_count < max_messages:
            traQing_params = {
                "stampId": stamp['id'],
                "isBot": "false",
                "groupBy": "message",
                "orderBy": "date",
                "order": "asc",
                "limit": str(messages_per_request),
                "offset": str(traQing_messageIds_count),
                "after": "2006-01-02T15:00:00.000Z",
                "before": "2026-01-02T14:59:59.999Z"
            }

            traQing_response = requests.get(traQing_url, params=traQing_params, cookies=traQing_cookies)
            traQing_response.raise_for_status()
            traQing_messageIds.extend(traQing_response.json())
            traQing_messageIds_count += messages_per_request
            if len(traQing_response.json()) < messages_per_request:
                break

        # traQing_messageIdsからmessageIdを取得してtraQ APIで詳細情報を取得
        traQing_messages = []
        for item in traQing_messageIds:
            message_id = item['message']

            # traQ APIでメッセージの詳細を取得
            message_url = f"https://q.trap.jp/api/v3/messages/{message_id}"
            message_response = requests.get(message_url, headers=traQ_headers)

            if message_response.status_code == 200:
                traQing_messages.append(message_response.json())
            else:
                print(f"メッセージ {message_id} の取得に失敗: {message_response.status_code}")
        all_traQing_messages.append({
            "stamp_name": stamp['name'],
            "stamp_id": stamp['id'],
            "messages": traQing_messages
        })

        with open('traQ_data.json', 'w', encoding='utf-8') as f:
            json.dump(all_traQ_messages, f, ensure_ascii=False, indent=2)
        with open('traQing_data.json', 'w', encoding='utf-8') as f:
            json.dump(all_traQing_messages, f, ensure_ascii=False, indent=2)


except FileNotFoundError as e:
    print(f"エラー: ファイルが見つかりません。 ({e.filename})")
    sys.exit(1)
except requests.exceptions.RequestException as e:
    print(f"エラー: API リクエストに失敗しました。 {e}")
    sys.exit(1)
except json.JSONDecodeError as e:
    print(f"エラー: JSONの形式が正しくありません。 {e}")
    sys.exit(1)
except Exception as e:
    print(f"予期しないエラーが発生しました: {e}")
    sys.exit(1)
