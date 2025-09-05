import { useState } from '#app';
import { computed } from 'vue';

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

/**
 * ユーザー一覧を管理し、検索機能を提供する composable。
 * アプリケーション全体でユーザーデータを共有し、ユーザー名やユーザーIDでの検索をサポートする。
 */
export const useUsers = () => {
  // ユーザー一覧のstate（配列とMap形式で保持）
  const userList = useState<Schemas['UserProfile'][]>('user-list', () => []);
  const userMap = useState<Map<string, Schemas['UserProfile']>>('user-map', () => new Map());

  /**
   * ユーザー名（traq_id）でユーザーを検索する
   * @param userName 検索するユーザー名（部分一致）
   * @returns マッチしたユーザーの配列
   */
  const searchByUserName = (userName: string): Schemas['UserProfile'][] => {
    if (!userName.trim()) return userList.value;

    const searchTerm = userName.toLowerCase();
    return userList.value.filter(user =>
      user.traq_id.toLowerCase().includes(searchTerm),
    );
  };

  /**
   * ユーザーID（user_id）でユーザーを検索する
   * @param userId 検索するユーザーID（部分一致）
   * @returns マッチしたユーザーの配列
   */
  const searchByUserId = (userId: string): Schemas['UserProfile'][] => {
    if (!userId.trim()) return userList.value;

    const searchTerm = userId.toLowerCase();
    return userList.value.filter(user =>
      user.user_id.toLowerCase().includes(searchTerm),
    );
  };

  /**
   * 表示名でユーザーを検索する
   * @param displayName 検索する表示名（部分一致）
   * @returns マッチしたユーザーの配列
   */
  const searchByDisplayName = (displayName: string): Schemas['UserProfile'][] => {
    if (!displayName.trim()) return userList.value;

    const searchTerm = displayName.toLowerCase();
    return userList.value.filter(user =>
      user.user_display_name.toLowerCase().includes(searchTerm),
    );
  };

  /**
   * 複合検索：ユーザー名、ユーザーID、表示名のいずれかにマッチするユーザーを検索
   * @param query 検索クエリ
   * @returns マッチしたユーザーの配列
   */
  const searchUsers = (query: string): Schemas['UserProfile'][] => {
    if (!query.trim()) return userList.value;

    const searchTerm = query.toLowerCase();
    return userList.value.filter(user =>
      user.traq_id.toLowerCase().includes(searchTerm)
      || user.user_id.toLowerCase().includes(searchTerm)
      || user.user_display_name.toLowerCase().includes(searchTerm),
    );
  };

  /**
   * ユーザーIDでユーザーを取得する（完全一致）
   * @param userId ユーザーID
   * @returns 見つかったユーザー、または undefined
   */
  const getUserById = (userId: string): Schemas['UserProfile'] | undefined => {
    return userMap.value.get(userId);
  };

  /**
   * ユーザー名（traq_id）でユーザーを取得する（完全一致）
   * @param userName ユーザー名
   * @returns 見つかったユーザー、または undefined
   */
  const getUserByName = (userName: string): Schemas['UserProfile'] | undefined => {
    return userList.value.find(user => user.traq_id === userName);
  };

  // 読み取り専用の computed プロパティ
  const users = computed(() => userList.value);
  const usersMap = computed(() => userMap.value);

  return {
    users,
    usersMap,
    searchByUserName,
    searchByUserId,
    searchByDisplayName,
    searchUsers,
    getUserById,
    getUserByName,
  };
};
