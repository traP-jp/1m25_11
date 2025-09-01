import litellm
from litellm import batch_completion

# Proxy URLを設定
litellm.api_base = "https://llm-proxy.trap.jp"

# モデルを指定してChatCompletionを呼び出す
try:
    MESSAGE = '''# 概要

「絵文字」の画像と情報を提供しますので、「説明文」を作成してください。

# 背景

私たちのサークルでは、部内向けのSNSを運用しています。このSNSには「絵文字」機能があります。これはSlackやDiscordにある機能と同様のものです。作成された絵文字は、メッセージの中で自由に使ったり、リアクションとしてメッセージにつけたりすることができます。
私たちのプロジェクトでは現在、このSNSで使用されている絵文字をまとめたポータルサイトを作成しています。このポータルサイトに掲載する絵文字の説明文を、自動で作成したいです。

# 手順

1. 絵文字の画像を参照し、写っているものや文字、情報を読み取る
2. 1で読み取った情報と概要情報をもとに絵文字の使われ方を分析する
3. 絵文字の説明文を作成する

# 説明文の条件

- 画像から読み取れる情報を組み込む
- 絵文字の概要、由来、意味、どんなときに使うか、などの情報を含める
- 自然な文章にする
- 200字程度にする
- 説明文が下記の投稿例に特化しすぎたものにしない

# 絵文字の情報

- 名前：`ranpuro_4`

## この絵文字が使われている投稿例
絵文字は、名前を `:` で囲って記述されています。以下は、この絵文字が使われている直近5件の投稿例です。

```md
# 合宿:ranpuro_1.large::ranpuro_2.large::ranpuro_3.large::ranpuro_4.large::ranpuro_5.large:
登壇者募集します！
夏合宿でもらんぷろが行われます。
1日目の夜21:30-22:30、場所は宴会場です。
登壇を希望する人は[らんぷろMD](https://md.trap.jp/CBG7JH3oTa63e0DIO43jaA)に必要事項入力をお願いします。
今回のらんぷろは1時間となりますので登壇希望者が多かった場合は先着順での発表になります:pray-nya:
```
```md
:ranpuro_1::ranpuro_2::ranpuro_3::ranpuro_4::ranpuro_5: の :ranpuro_1::ranpuro_2::ranpuro_3: の部分使われてなさすぎで可哀想というお気持ち
```
```md
雑談:ranpuro_1::ranpuro_2::ranpuro_3::ranpuro_4::ranpuro_5:開催！！
https://q.trap.jp/messages/0198f681-d6a8-7115-b5df-23ee3a5693ba
```
```md
## :ranpuro_4:班でUnity講習会をやります:blobwobwork:
8/18(明日)の21:00から:ranpuro_4:班向けにUnity講習会をやります。
誰でも参加OKです。ただし実習とかは無いです。
ワンマンソンでUnityを使ってる人にはもしかしたら助けになるかも？
資料: https://md.trap.jp/XCIETGMOShi08gGevBkWgw
```
```md
これ、俺がこれまでに押したスタンプランキングなんですけど、なぜか:ranpuro_5:の数が:ranpuro_4:より1つだけ小さくて:hatena_zubora:ってなってます。どこで押し忘れたんだろう。
ちなみにもらった数は同数なんですよね。皆さん:ranpuro_4::ranpuro_5:押してくれてありがとうございます。
```
'''
    response = litellm.completion(
        model="gpt-5-nano",
        messages=[
            {
                "role": "user",
                "content": [
                    {
                        "type": "text",
                        "text": MESSAGE
                    },
                    {
                        "type": "image_url",
                        "image_url": {
                            "url": "https://q.trap.jp/api/1.0/public/emoji/e77c3b8a-9ac2-45b1-b16b-11b1f3dcbc31"
                        }
                    }
                ]
            }
        ],
        stream=False,
        reasoning={ "effort": "low" },
    )

    # 応答を出力
    if response.choices:
        print(response.choices[0].message.content)
        print(response)
    else:
        print(response)
except Exception as e:
    print(f"エラーが発生しました: {e}")
