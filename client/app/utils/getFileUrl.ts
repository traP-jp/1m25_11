// スタンプのidから画像のURLを取得する。ない場合は:loading:を返す
export default function (fileId: string | undefined) {
  return fileId ? `https://q.trap.jp/api/1.0/public/emoji/${fileId}` : 'https://q.trap.jp/api/1.0/public/emoji/bc9a3814-f185-4b3d-ac1f-3c8f12ad7b52';
};
