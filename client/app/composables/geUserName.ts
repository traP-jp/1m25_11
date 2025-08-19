export const useForwardedUser = () => {
  const headers = useRequestHeaders(['x-forwarded-user']);

  // ヘッダーの値を取得して返す
  return headers['x-forwarded-user'];
};
