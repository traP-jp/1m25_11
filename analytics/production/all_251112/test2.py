import json

def extract_missing_ids():
    # llm_output.jsonl から id を抽出
    output_ids = set()
    with open('llm_output.jsonl', 'r', encoding='utf-8') as f:
        for line in f:
            data = json.loads(line.strip())
            if 'id' in data:
                output_ids.add(data['id'])

    # llm_input.jsonl から id を抽出
    input_ids = set()
    with open('llm_input.jsonl', 'r', encoding='utf-8') as f:
        for line in f:
            data = json.loads(line.strip())
            if 'id' in data:
                input_ids.add(data['id'])

    # output にあって input にない id を抽出
    missing_ids = output_ids - input_ids

    # 結果を表示
    print(f"llm_output.jsonl にのみ存在する ID の数: {len(missing_ids)}")
    print(f"該当する ID: {sorted(missing_ids)}")

    return missing_ids

if __name__ == '__main__':
    missing_ids = extract_missing_ids()
