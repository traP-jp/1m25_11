#!/usr/bin/env python3
import argparse
import json
import sys
from typing import Set

def count_unique_ids(path: str):
    unique_ids: Set[str] = set()
    total_lines = 0
    malformed = 0
    missing_id = 0

    try:
        with open(path, "r", encoding="utf-8") as f:
            for lineno, line in enumerate(f, start=1):
                total_lines += 1
                line = line.strip()
                if not line:
                    continue
                try:
                    obj = json.loads(line)
                except json.JSONDecodeError:
                    malformed += 1
                    print(line)
                    continue
                if 'id' in obj and obj['id'] is not None:
                    unique_ids.add(str(obj['id']))
                else:
                    missing_id += 1
    except FileNotFoundError:
        print(f"エラー: ファイルが見つかりません: {path}", file=sys.stderr)
        return 2
    except Exception as e:
        print(f"エラー: ファイル読み込み中に例外: {e}", file=sys.stderr)
        return 3

    print(f"ファイル: {path}")
    print(f"総行数: {total_lines}")
    print(f"JSONパース失敗行: {malformed}")
    print(f"id 欠如行: {missing_id}")
    print(f"ユニーク id 数: {len(unique_ids)}")
    return 0

def main():
    parser = argparse.ArgumentParser(description="Count unique ids in a JSONL file (llm_output.jsonl).")
    parser.add_argument("path", nargs="?", default="llm_output.jsonl", help="path to the jsonl file")
    args = parser.parse_args()
    rc = count_unique_ids(args.path)
    sys.exit(rc)

if __name__ == "__main__":
    main()
