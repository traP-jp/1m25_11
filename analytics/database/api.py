import json
import csv
import datetime
from typing import List, Dict, Any

def convert_jsonl_to_csv(input_file: str, output_file: str) -> None:
    """
    JSONLファイルをCSV形式に変換する

    Args:
        input_file: 入力JSONLファイルのパス
        output_file: 出力CSVファイルのパス
    """
    fixed_creator_id = "3b261ff3-f940-4e2c-a626-27387b6dd71b"
    current_time = datetime.datetime.now().isoformat()

    with open(output_file, 'w', newline='', encoding='utf-8') as csvfile:
        writer = csv.writer(csvfile)

        # CSVヘッダーを書き込み
        writer.writerow(['stamp_id', 'description', 'creator_id', 'created_at', 'updated_at'])

        # JSONLファイルを読み込み、各行を処理
        with open(input_file, 'r', encoding='utf-8') as jsonlfile:
            for line_number, line in enumerate(jsonlfile, 1):
                line = line.strip()
                if not line:
                    continue

                try:
                    # JSONを安全に解析
                    data = json.loads(line)

                    # 必要なフィールドを取得
                    stamp_id = data.get('id', '')
                    description = data.get('description', '')

                    # CSV行を書き込み
                    writer.writerow([
                        stamp_id,
                        description,
                        fixed_creator_id,
                        current_time,
                        current_time
                    ])

                except json.JSONDecodeError as e:
                    print(f"行 {line_number} でJSONデコードエラー: {e}")
                    continue
                except Exception as e:
                    print(f"行 {line_number} で予期しないエラー: {e}")
                    continue

def main():
    """メイン処理"""
    input_file = "input.jsonl"  # 入力ファイル名を適宜変更
    output_file = "output.csv"  # 出力ファイル名を適宜変更

    try:
        convert_jsonl_to_csv(input_file, output_file)
        print(f"変換完了: {input_file} -> {output_file}")
    except FileNotFoundError:
        print(f"ファイルが見つかりません: {input_file}")
    except Exception as e:
        print(f"エラーが発生しました: {e}")

if __name__ == "__main__":
    main()
