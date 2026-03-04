// Auth stub — Keycloak optionnel
// En dev sans Keycloak, on bypasse l'auth

export async function login() {
  // Redirect vers dashboard directement si pas de Keycloak
  const keycloakUrl = (import.meta as any).env?.VITE_KEYCLOAK_URL
  if (!keycloakUrl) {
    window.location.href = '/dashboard'
    return
  }
  try {
    const { UserManager, WebStorageStateStore } = await import('oidc-client-ts')
    const um = new UserManager({
      authority:    keycloakUrl + '/realms/' + (import.meta as any).env.VITE_KEYCLOAK_REALM,
      client_id:    (import.meta as any).env.VITE_KEYCLOAK_CLIENT_ID || 'green-frontend',
      redirect_uri: window.location.origin + '/callback',
      response_type: 'code',
      scope: 'openid profile email',
      userStore: new WebStorageStateStore({ store: window.sessionStorage }),
    })
    await um.signinRedirect()
  } catch (e) {
    window.location.href = '/dashboard'
  }
}

export async function logout() {
  window.location.href = '/login'
}

export async function getUser() { return null }
export async function getToken() { return null }
