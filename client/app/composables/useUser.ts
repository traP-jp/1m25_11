import { useState } from '#app';

/**
 * アプリケーション全体で利用可能な、リアクティブなユーザー名を管理する composable。
 * サーバーで X-Forwarded-User ヘッダーから値を取得し、useState に格納する。
 * @returns {Ref<string | null>} ユーザー名を保持するリアクティブな Ref オブジェクト
 */
export const useUser = () => {
  // 'user' というキーで state を定義。初期値は null。
  const user = useState<string | null>('user', () => null);

  // user.value がまだセットされていない場合のみ、ヘッダーから値を取得する処理を行う。
  // これにより、サーバーで一度だけヘッダーが読み込まれ、クライアントではその値が再利用される。
  if (import.meta.server && user.value === null) {
    const headers = useRequestHeaders(['x-forwarded-user']);
    const username = headers['x-forwarded-user'];
    if (username) {
      user.value = username;
    }
  }

  return user;
};
