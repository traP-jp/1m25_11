# create jsonl

import json
import os
import sys
import requests

bearer_token = os.environ.get("BEARER_TOKEN")
auth_token = os.getenv('TRAQ_AUTH_TOKEN')

try:
    # with文でファイルを開く
    with open('targeted_stamps.json', 'r', encoding='utf-8') as f_stamps:
        stamps = json.load(f_stamps)

    all_requests = []

    for stamp in stamps:
        # traQからデータ取得
        traQ_url = "https://q.trap.jp/api/v3/messages"
        traQ_headers = {
            "accept": "application/json",
            "Authorization": f"Bearer {bearer_token}"
        }
        traQ_params = {
            "word": f"\":{stamp['name']}:\"",
            "after": "2006-01-02T15:04:05Z",
            "before": "2006-01-02T15:04:05Z",
            "bot": "false",
            "limit": "100",
            "offset": "0",
            "sort": "createdAt",
        }

        traQ_response = requests.get(traQ_url, params=traQ_params, headers=traQ_headers)
        traQ_response.raise_for_status()
        traQ_data = traQ_response.json()['hits']

        # traQingからデータ取得
        traQing_url = "https://traqing.cp20.dev/api/stamps"
        traQing_params = {
            "stampId": stamp['id'],
            "isBot": "false",
            "groupBy": "message",
            "orderBy": "date",
            "order": "asc",
            "limit": "101",
            "offset": "0",
            "after": "2006-01-02T15:00:00.000Z",
            "before": "2026-01-02T14:59:59.999Z"
        }
        traQing_cookies = {
            'traq-auth-token': auth_token
        }

        traQing_response = requests.get(traQing_url, params=traQing_params, cookies=traQing_cookies)
        traQing_response.raise_for_status()
        traQing_data = traQing_response.json()

        # traQing_dataからmessageIdを取得してtraQ APIで詳細情報を取得
        traQing_messages = []
        for item in traQing_data:
            message_id = item['message']

            # traQ APIでメッセージの詳細を取得
            message_url = f"https://q.trap.jp/api/v3/messages/{message_id}"
            message_response = requests.get(message_url, headers=traQ_headers)

            if message_response.status_code == 200:
                traQing_messages.append(message_response.json())
            else:
                print(f"メッセージ {message_id} の取得に失敗: {message_response.status_code}")

        # 整形

        # TODO

        request = {
            "model": "gpt-5-nano",
            "messages": [
                {
                    "role": "system",
                    "content": "あなたは日本語で簡潔かつ客観的に記述するコンテンツ生成エンジンです。与えられた絵文字の画像と使用例（本文・リアクション）から、その絵文字の**概要・見た目・主な用法**を抽出し、「説明文」（約200字）と「キーワード」（検索用語の集合）を生成します。事実不明な細部は断定せず、汎用的で再利用可能な表現を優先します。"
                },
                {
                    "role": "developer",
                    "content": """出力は**厳密なJSONのみ**（余計な文言・改行・コードフェンス禁止）。
フォーマット:
```json
{"description":"…約200字の日本語…","keywords":["…","…"]}
```
## 生成ルール
* 言語: 日本語。敬体（です・ます）。
* 説明文（description）:
  * 文字数目安: **170〜230字**（200字±30字）。
  * 構成: ①見た目（色/表情/図柄/文字）→②主な使いどころ（気分・状況・ニュアンス）→③類似表現との差（強さ/丁寧さ/カジュアルさ等）。
  * 画像内の**文字は正確に転記**（例: 「OK」「了解」など）。アニメGIFなら「動く/点滅する」等も明記。
  * 投稿本文の長文引用は禁止。
* キーワード（keywords）:
  * 配列長: **25〜35個**。**重複禁止**。文字種は自由（ひらがな/カタカナ/漢字/英語小文字）。
  * 1要素は**1〜12文字程度**。ハッシュタグ・絵文字の記号は含めない。
  * 含める観点: ①画像の要素（色/形/表情/記号/文字そのもの）②用途・感情（例: うれしい/謝罪/了解/急ぎ）③口語ゆらぎ（例: りょ/オケ/OK/おけ）④検索で役立つ英語・和製英語（ok, agree, thumbs up等）⑤同義語・反対語の代表（控えめ/強め）。
  * 画像内テキストがある場合は**原文**と**読みのバリエーション**（例: OK, ok, オーケー）を優先追加。
* 推論手順（内部で実施）:
  1. 画像から: 物体/表情/配色/スタイル/文字の検出。
  2. 投稿・リアクションから: 使用場面・感情極性・礼儀度を推測（ただし一般化して記述）。
  3. 説明文を上記情報で組み立て、断定が難しい箇所は汎用表現で包む。
  4. キーワードを収集→正規化（小文字化/全半角統一）→重複除去→関連度順に並べ替え。
* 品質チェック（出力直前）:
  * 文字数範囲OK / JSON妥当 / キーワード重複なし / 断定のしすぎがないか。"""
                },
                {
                    "role": "user",
                    "content": [
                        {
                            "type": "text",
                            "text": f"""次の入力をもとに**JSONのみ**を出力してください。

## 入力
* 絵文字の名前: `{stamp['name']}`
* 本文で使われた投稿（配列）:
  `{traQ_data}`
* リアクションとして使われた投稿（配列）:
  `{traQing_data}`
* 絵文字画像: image_urlとして添付

### データ仕様
* `message`型:
```json
{
  "content": "本文",
  "stamps": ["絵文字の名前1","絵文字の名前2"]
}
```
* `messages`は本文中で`:name:`の形で使われた例、`reactions`はリアクションとして付いた例。

### 注意

* 投稿は**参考**です。用途はこれらに**限定しない**でください。
* 画像に描かれた**文字・図柄は最重要手がかり**として反映してください。
* 出力は次のキーのみ: `description`, `keywords`。**追加キー禁止**。"""
                        },
                        {
                            "type": "image_url",
                            "image_url": { "url": f"https://q.trap.jp/api/1.0/public/emoji/{stamp['id']}" },
                        }
                    ],
                }
            ],
            "stream": False,
            "reasoning": { "effort": "low" },
            "response_format": { "type": "json_object" },
            "max_tokens": 2000
        }
        request_formatted = {
            "custom_id": stamp['id'],
            "method": "POST",
            "url": "/v1/chat/completions",
            "body": request
        }
        all_requests.append(request_formatted)


    # JSONLファイルに書き出し
    with open('requests.jsonl', 'w', encoding='utf-8') as f_output:
        for request in all_requests:
            json.dump(request, f_output, ensure_ascii=False)
            f_output.write('\n')

    print(f"処理完了: {len(all_requests)}個のスタンプのデータをrequests.jsonlに保存しました。")

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
