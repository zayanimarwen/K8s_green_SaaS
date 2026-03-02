import { UserManager, WebStorageStateStore, User } from 'oidc-client-ts'

const userManager = new UserManager({
  authority:    import.meta.env.VITE_KEYCLOAK_URL + '/realms/' + import.meta.env.VITE_KEYCLOAK_REALM,
  client_id:    import.meta.env.VITE_KEYCLOAK_CLIENT_ID || 'green-frontend',
  redirect_uri: window.location.origin + '/callback',
  post_logout_redirect_uri: window.location.origin + '/login',
  response_type: 'code',
  scope: 'openid profile email',
  userStore: new WebStorageStateStore({ store: window.sessionStorage }),
  automaticSilentRenew: true,
  silent_redirect_uri: window.location.origin + '/silent-renew.html',
})

export async function login() {
  await userManager.signinRedirect()
}

export async function handleCallback(): Promise<User> {
  return userManager.signinRedirectCallback()
}

export async function logout() {
  await userManager.signoutRedirect()
}

export async function getUser(): Promise<User | null> {
  return userManager.getUser()
}

export async function getToken(): Promise<string | null> {
  const user = await userManager.getUser()
  return user?.access_token ?? null
}

export { userManager }
