import json # JSONパース
from openai import OpenAI # APIリクエスト
from pydantic import BaseModel # JSON構造化
import re # 正規表現
import sys # システム終了
import os  # ディレクトリ走査

min_length = 10
max_length = 250

# 日本語の文字が含まれているか
def contains_japanese(text: str) -> bool:
    # ひらがな、カタカナ、漢字のUnicode範囲
    pattern = re.compile(r'[\u3040-\u30FF\u4E00-\u9FFF]')
    return bool(pattern.search(text))

# LLM Output Format
class StampInfo(BaseModel):
    description: str
    keywords: list[str]

def createRequest(stamp, traQ_content, traQing_content):
    return {
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
                        "text": createPrompt(stamp, traQ_content, traQing_content)
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
        "text_format": StampInfo,
        "max_tokens": 2000
    }

def createPrompt(stamp, traQ_content, traQing_content):
    return f"""次の入力をもとに**JSONのみ**を出力してください。

## 入力
* 絵文字の名前: `{stamp['name']}`
* 本文で使われた投稿（配列）:
{traQ_content}
* リアクションとして使われた投稿（配列）:
{traQing_content}
* 絵文字画像: image_urlとして添付

### 注意

* 投稿は**参考**です。用途はこれらに**限定しない**でください。
* 画像に描かれた**文字・図柄は最重要手がかり**として反映してください。
* 出力は次のキーのみ: `description`, `keywords`。**追加キー禁止**。"""

# メッセージの判定
def is_usable_message(message):
    return contains_japanese(message['content']) and min_length <= len(message['content']) and len(message['content']) <= max_length

# メッセージの整形
def format_message(message):
    # message: { "id", "userId", "channelId", "content", "createdAt", "updatedAt", "pinned", "stamps", "threadId" }
    # message: { "name": str, "stamps": [str], }
    res = {}
    res["content"] = message['content']
    res["stamps"] = list(dict.fromkeys(
        all_stamps_dict[item["stampId"]]
        for item in message['stamps']
        if item["stampId"] in all_stamps_dict
    ))
    return res

# メッセージの配列の処理
def format_messages(messages):
    res = []
    for message in messages:
        if is_usable_message(message):
            res.append(format_message(message))
    return res[:5]

##### main

### ファイルを開く

try:
    # すべてのスタンプの一覧
    with open('stamps.json', 'r', encoding='utf-8') as f_all:
        all_stamps = json.load(f_all)
        all_stamps_dict = {stamp['id']: stamp['name'] for stamp in all_stamps}
except FileNotFoundError:
    print("エラー: stamps.json が見つかりません")
    sys.exit(1)
except json.JSONDecodeError:
    print("エラー: stamps.json のJSONフォーマットが不正です")
    sys.exit(1)

# --- ここから変更: targeted_stamps.json を廃止し、traQ/traQing ディレクトリを走査してスタンプIDを収集する ---

def load_messages_from_dir(dirpath):
    """
    dirpath 内の *.json を読み込んで { id: messages_list } を返す。
    ファイルは配列 (messages) か、{"messages": [...]} のどちらにも対応する。
    """
    res = {}
    if not os.path.isdir(dirpath):
        # 空ディレクトリ扱いにするが警告を出す
        print(f"警告: ディレクトリ {dirpath} が見つかりません。スキップします")
        return res

    for filename in os.listdir(dirpath):
        if not filename.lower().endswith('.json'):
            continue
        id_ = filename[:-5]  # strip .json
        filepath = os.path.join(dirpath, filename)
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                data = json.load(f)
                if isinstance(data, list):
                    messages = data
                elif isinstance(data, dict) and 'messages' in data and isinstance(data['messages'], list):
                    messages = data['messages']
                else:
                    # 想定外のフォーマットは空リストにして警告
                    print(f"警告: ファイル {filepath} のフォーマットが想定外です。空のメッセージリストとして扱います")
                    messages = []
                res[id_] = messages
        except json.JSONDecodeError:
            print(f"警告: JSONパースエラー: {filepath} をスキップします")
            continue
        except Exception as e:
            print(f"警告: ファイル読み込み中にエラー: {filepath} - {e}")
            continue
    return res

# traQ / traQing ディレクトリから読み込む
all_traQ_messages_dict = load_messages_from_dir('traQ')
all_traQing_messages_dict = load_messages_from_dir('traQing')

# targeted_stamps.json を使わず、両ディレクトリに現れるファイル名（ID）を処理対象とする
all_ids = set(all_traQ_messages_dict.keys()) | set(all_traQing_messages_dict.keys())
if not all_ids:
    print("エラー: 処理対象となるスタンプが traQ / traQing ディレクトリに見つかりません")
    sys.exit(1)

# stamps は既存コードが期待する形に合わせて { "id": ..., "name": ... } のリストにする
stamps = []
for id_ in sorted(all_ids):
    stamps.append({
        "id": id_,
        "name": all_stamps_dict.get(id_, id_)  # 名前がなければ id を代用
    })

# --- 変更終わり ---

### 各項目を処理

inputs = []

for stamp in stamps:
    try:
        traQ_content = format_messages(all_traQ_messages_dict.get(stamp['id'], []))
        traQing_content = format_messages(all_traQing_messages_dict.get(stamp['id'], []))
        inputs.append({
            "id": stamp['id'],
            "prompt": createPrompt(stamp, str(traQ_content), str(traQing_content))
        })
    except KeyError as e:
        print(f"警告: スタンプID {stamp.get('id', 'unknown')} の処理中にキーエラー: {e}")
        continue
    except Exception as e:
        print(f"警告: スタンプID {stamp.get('id', 'unknown')} の処理中にエラー: {e}")
        continue

### 出力
try:
    with open('llm_input.jsonl', 'w', encoding='utf-8') as f_in:
        for llm_input in inputs:
            json.dump(llm_input, f_in, ensure_ascii=False)
            f_in.write('\n')
    print(f"成功: {len(inputs)} 件のデータを llm_input.jsonl に出力しました")
except IOError as e:
    print(f"エラー: ファイル書き込みに失敗しました: {e}")
    sys.exit(1)
