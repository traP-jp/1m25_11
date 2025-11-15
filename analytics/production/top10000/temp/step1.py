# create targeted_stamps.json

import json

try:
    # with文にカンマで区切って2つのファイルを開く
    with open('stamps.json', 'r', encoding='utf-8') as f_stamps, \
         open('traQ_data.json', 'r', encoding='utf-8') as f_allowed:

        # それぞれのファイルからデータを読み込む
        all_stamps = json.load(f_stamps)
        allowed_stamps = json.load(f_allowed)

    # データ処理
    allowed_id_set = {stamp['stamp_id'] for stamp in allowed_stamps}
    targeted_stamps = []

    for stamp in all_stamps:
        if stamp['id'] in allowed_id_set:
            targeted_stamps.append(stamp)

    # targeted_stampsをJSONファイルに書き出し
    with open('targeted_stamps.json', 'w', encoding='utf-8') as f_output:
        json.dump(targeted_stamps, f_output, ensure_ascii=False, indent=2)

    print(f"処理完了: {len(targeted_stamps)}個のスタンプをtargeted_stamps.jsonに保存しました。")

except FileNotFoundError as e:
    print(f"エラー: ファイルが見つかりません。 ({e.filename})")
    exit(1)
except json.JSONDecodeError as e:
    print(f"エラー: JSONの形式が正しくありません。 {e}")
    exit(1)
except KeyError as e:
    print(f"エラー: 必要なキー '{e}' が見つかりません。JSONデータの構造を確認してください。")
    exit(1)
except PermissionError as e:
    print(f"エラー: ファイルの書き込み権限がありません。 ({e.filename})")
    exit(1)
except Exception as e:
    print(f"予期しないエラーが発生しました: {e}")
    exit(1)
