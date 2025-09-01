from openai import OpenAI

client = OpenAI(
    base_url="https://llm-proxy.trap.jp/"
)

try:
    MESSAGE = '''# 概要
「絵文字」の画像データと情報を受け取り、「説明文」を作成してください。

# 背景
私たちはサークル内向けSNSを運用しています。このSNSにはSlackやDiscord同様の「絵文字」機能があり、作成された絵文字はメッセージ内で自由に使ったり、リアクションとしてメッセージにつけたりできます。現在、このSNSで使われている絵文字をまとめたポータルサイトを作成しています。このポータルサイトに掲載する絵文字の説明文を、自動で生成したいと考えています。

# 手順
- はじめに、下記チェックリストを確認しながら全体の流れを意識してください。

1. 絵文字の画像データから読み取れる内容や文字や写真、図柄を把握する。
2. 把握した情報と概要情報をもとに、絵文字の使われ方を分析する。
3. 絵文字の説明文を作成する。

# チェックリスト
- 絵文字の画像データを参照し、見た目や特徴、読み取れる文字・内容を把握する
- 利用例や画像データから、主な使い方・用途を推測する
- 含めるべき情報：概要、由来、意味、どんなときに使うか
- 説明文は自然な日本語で200字程度にまとめる
- あまり利用例に特化しすぎず、汎用的な説明にする

# 絵文字の概要
- 名前：`ranpuro_4`

## この絵文字が使われている投稿例
絵文字は、名前を `:` で囲って記述されています。以下は、この絵文字が使われている直近5件の投稿例です。

```md
{$MESSAGE_1}
```
```md
{$MESSAGE_2}
```
```md
{$MESSAGE_3}
```
```md
{$MESSAGE_4}
```
```md
{$MESSAGE_5}
```

# 出力形式
説明文のみを自然な日本語で1つ出力してください。JSONやCSVなど構造化出力や補足説明は不要です。
'''

    response = client.responses.create(
        model="gpt-5-nano",
        input=[
            {
                "role": "user",
                "content": [
                    {
                        "type": "input_text",
                        "text": MESSAGE
                    },
                    {
                        "type": "input_image",
                        "image_url": "https://q.trap.jp/api/1.0/public/emoji/e77c3b8a-9ac2-45b1-b16b-11b1f3dcbc31"
                    }
                ]
            }
        ],
        stream=False,
        reasoning={ "effort": "low" },
    )

    print(response.output_text)
    print(response)
except Exception as e:
    print(f"エラーが発生しました: {e}")
