import json
import datetime
from openai import OpenAI

# OpenAIクライアントを初期化
client = OpenAI(
    base_url="https://llm-proxy.trap.jp/"
)

def main():
    # requests.jsonlファイルをアップロード
    print("requests.jsonlをアップロード中...")
    with open("requests.jsonl", "rb") as file:
        batch_input_file = client.files.create(
            file=file,
            purpose="batch"
        )

    print(f"ファイルID: {batch_input_file.id}")

    # バッチジョブを作成
    print("バッチジョブを作成中...")
    batch = client.batches.create(
        input_file_id=batch_input_file.id,
        endpoint="/v1/chat/completions",
        completion_window="24h"
    )

    print(f"バッチID: {batch.id}")
    print(f"ステータス: {batch.status}")

    # バッチ情報をローカルに保存
    batch_info = {
        "batch_id": batch.id,
        "input_file_id": batch_input_file.id,
        "status": batch.status,
        "created_at": datetime.datetime.now().isoformat(),
        "openai_created_at": str(batch.created_at)
    }

    with open("batch_info.json", "w") as f:
        json.dump(batch_info, f, indent=2)

    print("バッチ情報をbatch_info.jsonに保存しました")
    print(f"バッチIDをメモしてください: {batch.id}")

if __name__ == "__main__":
    main()
