import json
import datetime
from openai import OpenAI

# OpenAIクライアントを初期化
client = OpenAI(
    base_url="https://llm-proxy.trap.jp/"
)


def main():
    # batch_info.jsonを読み込み
    try:
        with open("batch_info.json", "r") as f:
            batch_info = json.load(f)
        batch_id = batch_info["batch_id"]
        print(f"バッチID: {batch_id}")
    except FileNotFoundError:
        print("batch_info.jsonが見つかりません")
        return
    except KeyError:
        print("batch_info.jsonにbatch_idが見つかりません")
        return

    # バッチのステータスを確認
    print("バッチステータスを確認中...")
    batch = client.batches.retrieve(batch_id)

    print(f"ステータス: {batch.status}")
    print(f"リクエスト数: {batch.request_counts}")

    if batch.completed_at:
        print(f"完了時刻: {batch.completed_at}")

    # 完了している場合、結果をダウンロード
    if batch.status == "completed":
        print("\nバッチが完了しました。結果をダウンロード中...")

        if batch.output_file_id:
            # 結果ファイルをダウンロード
            file_response = client.files.content(batch.output_file_id)

            # ローカルに保存
            output_filename = f"batch_results_{batch_id}.jsonl"
            with open(output_filename, "wb") as f:
                f.write(file_response.content)
            print(f"結果を保存しました: {output_filename}")

        if batch.error_file_id:
            # エラーファイルをダウンロード
            error_response = client.files.content(batch.error_file_id)

            error_filename = f"batch_errors_{batch_id}.jsonl"
            with open(error_filename, "wb") as f:
                f.write(error_response.content)
            print(f"エラー情報を保存しました: {error_filename}")

    elif batch.status == "failed":
        print(f"バッチが失敗しました。失敗時刻: {batch.failed_at}")

        if batch.error_file_id:
            error_response = client.files.content(batch.error_file_id)
            error_filename = f"batch_errors_{batch_id}.jsonl"
            with open(error_filename, "wb") as f:
                f.write(error_response.content)
            print(f"エラー情報を保存しました: {error_filename}")

    else:
        print("バッチはまだ処理中です")

    # ステータス情報を更新して保存
    batch_info.update({
        "last_checked": datetime.datetime.now().isoformat(),
        "status": batch.status,
        "request_counts": batch.request_counts,
        "output_file_id": batch.output_file_id,
        "error_file_id": batch.error_file_id,
        "completed_at": str(batch.completed_at) if batch.completed_at else None,
        "failed_at": str(batch.failed_at) if batch.failed_at else None
    })

    with open("batch_info.json", "w") as f:
        json.dump(batch_info, f, indent=2)

    print("batch_info.jsonを更新しました")

if __name__ == "__main__":
    main()
