# create jsonl

import json
import os
import sys
import requests
import tempfile

bearer_token = os.environ.get("BEARER_TOKEN")
auth_token = os.getenv('TRAQ_AUTH_TOKEN')

stamps_file = 'targeted_stamps.json'

since = "2006-01-02T15:04:05Z"
until = "2026-01-02T15:04:05Z"

# 出力ディレクトリ（スタンプごとにファイルを分ける）
OUT_DIR_TRAQ = 'out/traQ'
OUT_DIR_TRAQING = 'out/traQing'

os.makedirs(OUT_DIR_TRAQ, exist_ok=True)
os.makedirs(OUT_DIR_TRAQING, exist_ok=True)


def atomic_write(path, data):
    # 一時ファイルに書いてから置換することで中途半端なファイルを書き残さない
    dirpath = os.path.dirname(path)
    fd, tmp = tempfile.mkstemp(dir=dirpath)
    try:
        with os.fdopen(fd, 'w', encoding='utf-8') as f:
            json.dump(data, f, ensure_ascii=False, indent=2)
        os.replace(tmp, path)
    except Exception:
        try:
            os.remove(tmp)
        except Exception:
            pass
        raise


try:
    # stamps ファイルを読む
    with open(stamps_file, 'r', encoding='utf-8') as f_stamps:
        stamps = json.load(f_stamps)

    for stamp in stamps:
        stamp_id = stamp['id']
        stamp_name = stamp.get('name')

        traq_path = os.path.join(OUT_DIR_TRAQ, f"{stamp_id}.json")
        traqing_path = os.path.join(OUT_DIR_TRAQING, f"{stamp_id}.json")

        # 両方のファイルが存在すれば既に取得済みとみなしてスキップ
        if os.path.exists(traq_path) and os.path.exists(traqing_path):
            print(f"スタンプ {stamp_name} ({stamp_id}) は既に取得済みなのでスキップ")
            continue

        # traQからデータ取得
        traQ_url = "https://q.trap.jp/api/v3/messages"
        traQ_headers = {
            "accept": "application/json",
            "Authorization": f"Bearer {bearer_token}"
        }

        traQ_params = {
            "word": f"\":{stamp_name}:\"",
            "after": since,
            "before": until,
            "bot": "false",
            "limit": "10",
            "offset": "0",
            "sort": "createdAt"
        }

        traQ_messages = []
        traQ_response = requests.get(
            traQ_url, params=traQ_params, headers=traQ_headers)
        traQ_response.raise_for_status()
        traQ_messages.extend(traQ_response.json().get('hits', []))

        # ファイルに書き出す（成功したレスポンスのみ）
        traq_out = {
            "stamp_name": stamp_name,
            "stamp_id": stamp_id,
            "messages": traQ_messages
        }
        atomic_write(traq_path, traq_out)

        # traQingからデータ取得
        traQing_url = "https://traqing.cp20.dev/api/stamps"
        traQing_cookies = {
            'traq-auth-token': auth_token
        }

        traQing_params = {
            "stampId": stamp_id,
            "isBot": "false",
            "groupBy": "message",
            "orderBy": "date",
            "order": "asc",
            "limit": "10",
            "offset": "0",
            "after": "2006-01-02T15:00:00.000Z",
            "before": "2026-01-02T14:59:59.999Z"
        }

        traQing_messageIds = []
        traQing_response = requests.get(
            traQing_url, params=traQing_params, cookies=traQing_cookies)
        traQing_response.raise_for_status()
        traQing_messageIds.extend(traQing_response.json())

        # traQing_messageIdsからmessageIdを取得してtraQ APIで詳細情報を取得
        traQing_messages = []
        for item in traQing_messageIds:
            message_id = item['message']
            message_url = f"https://q.trap.jp/api/v3/messages/{message_id}"
            message_response = requests.get(message_url, headers=traQ_headers)
            if message_response.status_code == 200:
                traQing_messages.append(message_response.json())
            else:
                print(
                    f"メッセージ {message_id} の取得に失敗: {message_response.status_code}")

        traqing_out = {
            "stamp_name": stamp_name,
            "stamp_id": stamp_id,
            "messages": traQing_messages
        }
        atomic_write(traqing_path, traqing_out)

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
