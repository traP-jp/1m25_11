import json

min_content_length = 10
max_content_length = 200

with open('stamps.json', 'r', encoding='utf-8') as f:
    stamps = json.load(f)
with open('traQ_data.json', 'r', encoding='utf-8') as traQ_f:
    all_traQ_messages = json.load(traQ_f)
with open('traQing_data.json', 'r', encoding='utf-8') as traQing_f:
    all_traQing_messages = json.load(traQing_f)

all_requests = []

for stamp in stamps:
    
    traQ_message_contents = []
    for item in all_traQ_messages:
        if item['stamp_name'] == stamp['name']:
            traQ_messages = item['messages']
            traQ_message_contents = [msg['content'] for msg in traQ_messages]
            break
    traQ_message_contents = [
        msg for msg in traQ_message_contents 
        if min_content_length <= len(msg) <= max_content_length
    ][:5]

    traQing_message_contents = []
    for item in all_traQing_messages:
        if item['stamp_name'] == stamp['name']:
            traQing_messages = item['messages']
            traQing_message_contents = [msg['content'] for msg in traQing_messages]
            break
    traQing_message_contents = [
        msg for msg in traQing_message_contents 
        if min_content_length <= len(msg) <= max_content_length
    ][:5]

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
`{traQ_message_contents}`
* リアクションとして使われた投稿（配列）:
`{traQing_message_contents}`
* 絵文字画像: image_urlとして添付

### データ仕様
* `message`型:
```json
{{
"content": "本文",
"stamps": ["絵文字の名前1","絵文字の名前2"]
}}
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

with open('requests.jsonl', 'w', encoding='utf-8') as f:
    for req in all_requests:
        json.dump(req, f, ensure_ascii=False)
        f.write('\n')