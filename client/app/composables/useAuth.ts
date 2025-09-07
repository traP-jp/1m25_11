export const useAuth = () => {
  const user = useState<Schemas['UserStatus'] | null>('user');

  // ログインしているかどうか判断するためのプロパティ
  const isLoggedIn = computed(() => !!user.value);

  return {
    user: readonly(user),
    isLoggedIn: readonly(isLoggedIn),
  };
};
