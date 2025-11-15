import json # JSONパース
from openai import OpenAI # APIリクエスト
from pydantic import BaseModel # JSON構造化
import re # 正規表現
import time # 待機時間
import sys # システム終了
import os # ファイル操作
import requests # HTTP リクエスト
import base64 # base64エンコード
import cairosvg # SVG to PNG変換
from io import BytesIO # バイナリデータ処理

# LLM Output Format
class StampInfo(BaseModel):
    description: str
    keywords: list[str]

# APIキー残高不足時に投げる専用例外
class InsufficientQuotaError(Exception):
    pass

def get_image_as_base64(emoji_id):
    """
    絵文字画像を取得し、SVGの場合はPNGに変換してbase64エンコードしたdata URLを返す
    """
    try:
        url = f"https://q.trap.jp/api/1.0/public/emoji/{emoji_id}"
        response = requests.get(url, timeout=30)
        response.raise_for_status()

        content_type = response.headers.get('content-type', '').lower()
        image_data = response.content

        # SVG形式の検出
        is_svg = (
            'svg' in content_type or
            image_data.startswith(b'<svg') or
            image_data.startswith(b'<?xml') and b'<svg' in image_data[:1000]
        )

        if is_svg:
            # SVGをPNGに変換
            png_data = cairosvg.svg2png(bytestring=image_data, output_width=128, output_height=128)
            mime_type = 'image/png'
            final_data = png_data
        else:
            # PNG/JPEG/GIFなどはそのまま使用
            if 'png' in content_type:
                mime_type = 'image/png'
            elif 'jpeg' in content_type or 'jpg' in content_type:
                mime_type = 'image/jpeg'
            elif 'gif' in content_type:
                mime_type = 'image/gif'
            elif 'webp' in content_type:
                mime_type = 'image/webp'
            else:
                mime_type = 'image/png'  # デフォルト
            final_data = image_data

        # base64エンコード
        encoded = base64.b64encode(final_data).decode('utf-8')
        return f"data:{mime_type};base64,{encoded}"

    except Exception as e:
        print(f"警告: 画像の取得/変換に失敗 ID: {emoji_id}, エラー: {e}")
        # フォールバック: 元のURL形式
        return f"https://q.trap.jp/api/1.0/public/emoji/{emoji_id}"

def createRequest(id, prompt):
    image_url = get_image_as_base64(id)

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
                    { "type": "image_url", "image_url": { "url": image_url } }
                ],
            }
        ],
        "stream": False,
        #"temperature": 0,
        "max_tokens": 2000
    }

def load_processed_ids():
    """既に処理済みのIDセットを取得（IDを文字列として正規化）"""
    processed_ids = set()
    if os.path.exists('llm_output.jsonl'):
        try:
            with open('llm_output.jsonl', 'r', encoding='utf-8') as f:
                for line in f:
                    line = line.strip()
                    if not line:
                        continue
                    try:
                        data = json.loads(line)
                    except json.JSONDecodeError:
                        # 不正行は無視
                        continue
                    if 'id' in data:
                        # 型の不一致防止のため文字列化して保持
                        processed_ids.add(str(data['id']))
        except (IOError, Exception) as e:
            print(f"警告: llm_output.jsonl の読み込み中にエラー: {e}")
    return processed_ids

def append_output(output_data):
    """出力ファイルに1行追加（id を文字列化、fsyncで確実に同期）"""
    try:
        # id を文字列に正規化してから書き込む
        if 'id' in output_data:
            output_data['id'] = str(output_data['id'])
        with open('llm_output.jsonl', 'a', encoding='utf-8') as f:
            json.dump(output_data, f, ensure_ascii=False)
            f.write('\n')
            f.flush()
            try:
                os.fsync(f.fileno())
            except Exception:
                # fsync に失敗しても致命的ではない（例: Windows 互換等）
                pass
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

            # 残高不足/課金上限系エラーは非再試行で即時中断
            quota_signals = (
                'insufficient_quota',
                'insufficient quota',
                'quota exceeded',
                'exceeded your current quota',
                'hard limit reached',
                'out of credits',
                'insufficient credits',
                'insufficient balance',
                'payment required',  # 402系
                'billing',
                '402',
            )
            if any(s in error_str for s in quota_signals):
                print("エラー: APIキーの残高不足（または課金上限）を検知。再試行しません。")
                raise InsufficientQuotaError(error_str)

            # Rate Limitは指数バックオフで再試行
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
                # それ以外は再試行せず即座に上位へ投げて中断判断させる
                print(f"APIエラー: {e}")
                raise

    raise Exception("予期しないエラー: 再試行ループを抜けました")

def validate_response_format(content, input_id):
    """
    APIレスポンスの形式をチェックし、不正な場合はFalseを返す
    """
    try:
        response_json = json.loads(content)

        # 必須フィールドの存在チェック
        if 'description' not in response_json or 'keywords' not in response_json:
            print(f"警告: 必須フィールドが不足 ID: {input_id}")
            return False, None

        # descriptionに改行が含まれていないかチェック
        if '\n' in response_json['description']:
            print(f"警告: descriptionに改行が含まれています ID: {input_id}")
            return False, None

        # keywordsが配列かチェック
        if not isinstance(response_json['keywords'], list):
            print(f"警告: keywordsが配列ではありません ID: {input_id}")
            return False, None

        return True, response_json

    except json.JSONDecodeError as e:
        print(f"警告: JSONパースエラー ID: {input_id}, エラー: {e}")
        return False, None

def process_single_request(client, input_id, prompt, max_format_retries=3):
    """
    単一のリクエストを処理し、形式チェックして必要に応じて再試行
    """
    for format_attempt in range(max_format_retries):
        try:
            request_data = createRequest(input_id, prompt)
            response = make_api_request_with_retry(client, request_data)
            print(response)

            content = response.choices[0].message.content
            is_valid, response_json = validate_response_format(content, input_id)

            if is_valid:
                response_json['id'] = input_id
                return response_json
            else:
                if format_attempt < max_format_retries - 1:
                    print(f"形式エラーのため再試行します ({format_attempt + 1}/{max_format_retries}) ID: {input_id}")
                    time.sleep(2)  # 短い待機時間
                    continue
                else:
                    print(f"最大再試行回数に達しました。エラー応答を生成します ID: {input_id}")
                    return {
                        "id": input_id,
                        "description": "形式エラーのため正常な応答を生成できませんでした",
                        "keywords": ["エラー", "形式不正"]
                    }

        except Exception as e:
            # API関連のエラーは上位に投げる
            raise e

    # ここには到達しないはず
    raise Exception("予期しないエラー: 形式再試行ループを抜けました")

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
aborted_due_to_quota = False  # 残高不足で中断したかのフラグ
aborted_due_to_error = False  # 一般エラーで中断したかのフラグ

for llm_input in inputs:
    input_id = llm_input.get('id')
    if not input_id:
        print("警告: IDが見つからない入力をスキップ")
        continue

    # 比較用キーは常に文字列に正規化する
    id_key = str(input_id)

    if id_key in processed_ids:
        print(f"skip {id_key}")
        skipped_count += 1
        continue

    print(f"start {id_key}")

    try:
        response_json = process_single_request(client, input_id, llm_input['prompt'])

        # 出力IDを文字列で正規化してから書き込む
        response_json['id'] = id_key

        # 出力ファイルに追加
        append_output(response_json)
        processed_ids.add(id_key)
        processed_count += 1
        print(f"完了 {id_key}")

    except InsufficientQuotaError:
        print("APIキーの残高が不足しているため処理を中断します。残高を追加して再実行してください。")
        aborted_due_to_quota = True
        break
    except KeyError as e:
        print(f"エラー: 必要なキーが見つかりません ID: {id_key}, キー: {e}。処理を中断します。")
        aborted_due_to_error = True
        break
    except Exception as e:
        print(f"エラー: ID {id_key} の処理中に予期しないエラーが発生しました。処理を中断します。詳細: {e}")
        aborted_due_to_error = True
        break

if aborted_due_to_quota:
    print("\n注意: APIキーの残高不足により処理を途中で中断しました。")
if aborted_due_to_error:
    print("\n注意: Rate Limit以外のエラーにより処理を途中で中断しました。")

print(f"\n処理完了:")
print(f"  新規処理: {processed_count} 件")
print(f"  スキップ: {skipped_count} 件")
print(f"  合計入力: {len(inputs)} 件")
