import json # JSONパース
from openai import OpenAI # APIリクエスト
from pydantic import BaseModel # JSON構造化
import re # 正規表現
import time # 待機時間
import sys # システム終了
import os # ファイル操作

# LLM Output Format
class StampInfo(BaseModel):
    description: str
    keywords: list[str]

def createRequest(id, prompt):
    return {
        "model": "gpt-4.1-nano",
        "messages": [
            {
                "role": "system",
                "content": "あなたは日本語で簡潔かつ客観的に記述するコンテンツ生成エンジンです。与えられた絵文字の画像と使用例（本文・リアクション）から、その絵文字の概要・見た目・主な用法を抽出し、説明文とキーワードのみをJSONで返します。思考過程・推論の説明や前置きは一切出力しないでください。"
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
                    { "type": "text", "text": prompt },
                    { "type": "image_url", "image_url": { "url": f"https://q.trap.jp/api/1.0/public/emoji/{id}" } }
                ],
            }
        ],
        "stream": False,
        #"temperature": 0,
        "max_tokens": 2000
    }

def load_processed_ids():
    """既に処理済みのIDセットを取得"""
    processed_ids = set()
    if os.path.exists('llm_output.jsonl'):
        try:
            with open('llm_output.jsonl', 'r', encoding='utf-8') as f:
                for line in f:
                    line = line.strip()
                    if line:
                        data = json.loads(line)
                        if 'id' in data:
                            processed_ids.add(data['id'])
        except (json.JSONDecodeError, IOError) as e:
            print(f"警告: llm_output.jsonl の読み込み中にエラー: {e}")
    return processed_ids

def append_output(output_data):
    """出力ファイルに1行追加"""
    try:
        with open('llm_output.jsonl', 'a', encoding='utf-8') as f:
            json.dump(output_data, f, ensure_ascii=False)
            f.write('\n')
    except IOError as e:
        print(f"エラー: 出力ファイルへの書き込みに失敗: {e}")
        raise

def make_api_request_with_retry(client, request_data, max_retries=5):
    """Rate Limit対応の再試行機能付きAPIリクエスト"""
    for attempt in range(max_retries):
        try:
            response = client.chat.completions.create(**request_data)
            return response
        except Exception as e:
            error_str = str(e).lower()
            if 'rate limit' in error_str or 'too many requests' in error_str:
                if attempt < max_retries - 1:
                    wait_time = (2 ** attempt) * 60
                    print(f"Rate Limitに達しました。{wait_time}秒後に再試行します。")
                    time.sleep(wait_time)
                    continue
                else:
                    print(f"Rate Limitエラー: 最大再試行回数に達しました")
                    raise
            else:
                print(f"APIエラー: {e}")
                raise

    raise Exception("予期しないエラー: 再試行ループを抜けました")

client = OpenAI(
    base_url="https://llm-proxy.trap.jp/"
)

##### main

### ファイルを開く

try:
    with open('llm_input.jsonl', 'r', encoding='utf-8') as f_in:
        inputs = []
        for line in f_in:
            line = line.strip()
            if line:
                try:
                    inputs.append(json.loads(line))
                except json.JSONDecodeError as e:
                    print(f"警告: 不正なJSON行をスキップ: {e}")
                    continue
except FileNotFoundError:
    print("エラー: llm_input.jsonl が見つかりません")
    sys.exit(1)
except IOError as e:
    print(f"エラー: llm_input.jsonl の読み込みに失敗: {e}")
    sys.exit(1)

# 既に処理済みのIDを取得
processed_ids = load_processed_ids()
print(f"既に処理済み: {len(processed_ids)} 件")

### ループ

processed_count = 0
skipped_count = 0

for llm_input in inputs:
    input_id = llm_input.get('id')
    if not input_id:
        print("警告: IDが見つからない入力をスキップ")
        continue

    if input_id in processed_ids:
        print(f"skip {input_id}")
        skipped_count += 1
        continue

    print(f"start {input_id}")

    try:
        request_data = createRequest(input_id, llm_input['prompt'])
        response = make_api_request_with_retry(client, request_data)

        # APIレスポンスにIDを追加
        try:
            content = response.choices[0].message.content
            response_json = json.loads(content)
            response_json['id'] = input_id
        except Exception:
            print(f"警告: APIレスポンスのJSONパースに失敗 ID: {input_id}")
            response_json = {
                "id": input_id,
                "description": f"パースエラーのため生成できませんでした\n{response}",
                "keywords": ["エラー", "パース失敗"]
            }

        # 出力ファイルに追加
        append_output(response_json)
        processed_ids.add(input_id)
        processed_count += 1
        print(f"完了 {input_id}")

    except KeyError as e:
        print(f"エラー: 必要なキーが見つかりません ID: {input_id}, キー: {e}")
        continue
    except Exception as e:
        print(f"エラー: ID {input_id} の処理中に予期しないエラー: {e}")
        continue

print(f"\n処理完了:")
print(f"  新規処理: {processed_count} 件")
print(f"  スキップ: {skipped_count} 件")
print(f"  合計入力: {len(inputs)} 件")
