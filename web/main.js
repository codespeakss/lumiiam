import { createApp, ref } from 'https://unpkg.com/vue@3/dist/vue.esm-browser.prod.js'

const App = {
  setup() {
    const api_base = location.origin + '/api/v1'
    const identifier = ref('admin@example.com')
    const password = ref('admin123')
    const access_token = ref('')
    const refresh_token = ref('')
    const me = ref(null)
    const users = ref([])
    const error = ref('')

    const login = async () => {
      error.value = ''
      const res = await fetch(api_base + '/auth/login', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ identifier: identifier.value, password: password.value }) })
      const data = await res.json()
      if (!res.ok) { error.value = data.error || 'login failed'; return }
      access_token.value = data.access_token
      refresh_token.value = data.refresh_token
    }

    const fetch_me = async () => {
      error.value = ''
      const res = await fetch(api_base + '/users/me', { headers: { 'Authorization': 'Bearer ' + access_token.value }})
      const data = await res.json()
      if (!res.ok) { error.value = data.error || 'fetch me failed'; return }
      me.value = data
    }

    const fetch_users = async () => {
      const res = await fetch(api_base + '/users', { headers: { 'Authorization': 'Bearer ' + access_token.value }})
      const data = await res.json()
      if (!res.ok) { error.value = data.error || 'fetch users failed'; return }
      users.value = data.items
    }

    const do_refresh = async () => {
      const res = await fetch(api_base + '/auth/refresh', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ refresh_token: refresh_token.value }) })
      const data = await res.json()
      if (!res.ok) { error.value = data.error || 'refresh failed'; return }
      access_token.value = data.access_token
      refresh_token.value = data.refresh_token
    }

    const logout = async () => {
      const res = await fetch(api_base + '/auth/logout', { method: 'POST', headers: { 'Authorization': 'Bearer ' + access_token.value }})
      if (!res.ok) { const data = await res.json(); error.value = data.error || 'logout failed'; return }
      access_token.value = ''
      refresh_token.value = ''
      me.value = null
    }

    return { identifier, password, access_token, refresh_token, me, users, error, login, fetch_me, fetch_users, do_refresh, logout }
  },
  template: `
  <div style="max-width:680px;margin:40px auto;font-family:system-ui">
    <h2>lumiiam test</h2>
    <div style="display:flex;gap:8px;align-items:end">
      <label style="display:flex;flex-direction:column;">
        <span>identifier (email or username)</span>
        <input v-model="identifier" />
      </label>
      <label style="display:flex;flex-direction:column;">
        <span>password</span>
        <input type="password" v-model="password" />
      </label>
      <button @click="login">login</button>
      <button @click="do_refresh" :disabled="!refresh_token">refresh</button>
      <button @click="logout" :disabled="!access_token">logout</button>
    </div>

    <p v-if="error" style="color:#c00">{{ error }}</p>

    <div>
      <p><strong>access token:</strong> {{ access_token }}</p>
      <p><strong>refresh token:</strong> {{ refresh_token }}</p>
    </div>

    <div style="display:flex; gap:8px">
      <button @click="fetch_me" :disabled="!access_token">fetch me</button>
      <button @click="fetch_users" :disabled="!access_token">fetch users</button>
    </div>

    <pre v-if="me">me: {{ JSON.stringify(me, null, 2) }}</pre>
    <pre v-if="users && users.length">users: {{ JSON.stringify(users, null, 2) }}</pre>
  </div>
  `
}

createApp(App).mount('#app')
