const apiBase = location.origin + '/api/v1'

function logEl() {
  return document.getElementById('results')
}

function setStatus(text, cls = '') {
  const el = document.getElementById('status')
  el.textContent = text
  el.className = cls
}

function appendCase(title, ok, detail) {
  const wrap = document.createElement('div')
  wrap.className = 'case'
  const h = document.createElement('div')
  h.innerHTML = `<strong>${title}</strong> — <span class="${ok ? 'ok' : 'fail'}">${ok ? 'OK' : 'FAIL'}</span>`
  wrap.appendChild(h)
  if (detail) {
    const pre = document.createElement('pre')
    pre.textContent = typeof detail === 'string' ? detail : JSON.stringify(detail, null, 2)
    wrap.appendChild(pre)
  }
  logEl().appendChild(wrap)
}

async function http(method, path, body, headers = {}) {
  const res = await fetch(apiBase + path, {
    method,
    headers: { 'Content-Type': 'application/json', ...headers },
    body: body ? JSON.stringify(body) : undefined,
  })
  let data = null
  try { data = await res.json() } catch (_) { /* ignore */ }
  return { ok: res.ok, status: res.status, data }
}

function randStr(n = 6) {
  const chars = 'abcdefghijklmnopqrstuvwxyz0123456789'
  let s = ''
  for (let i = 0; i < n; i++) s += chars[Math.floor(Math.random() * chars.length)]
  return s
}

async function runTests() {
  logEl().innerHTML = ''
  setStatus('Running...', '')

  const suffix = randStr(8)
  const email = `u_${suffix}@example.com`
  const username = `user_${suffix}`
  const password = 'pass1234!'

  // health
  try {
    const res = await fetch(location.origin + '/health')
    appendCase('GET /health', res.ok, { status: res.status })
    if (!res.ok) throw new Error('health failed')
  } catch (e) {
    appendCase('GET /health exception', false, String(e))
    setStatus('Failed (server not running?)', 'fail')
    return
  }

  // create user
  const create1 = await http('POST', '/users', { email, username, password })
  appendCase('POST /users (create)', create1.ok && create1.status === 201, create1)

  // duplicate email
  const createDupEmail = await http('POST', '/users', { email, username: username + '_x', password })
  appendCase('POST /users (duplicate email)', !createDupEmail.ok, createDupEmail)

  // duplicate username
  const createDupUsername = await http('POST', '/users', { email: 'x_' + email, username, password })
  appendCase('POST /users (duplicate username)', !createDupUsername.ok, createDupUsername)

  // login newly created user
  const login = await http('POST', '/auth/login', { identifier: email, password })
  appendCase('POST /auth/login (new user)', login.ok, login)

  // fetch me using access token (if login ok)
  if (login.ok) {
    const me = await http('GET', '/users/me', null, { Authorization: 'Bearer ' + login.data.access_token })
    appendCase('GET /users/me (new user)', me.ok, me)
  }

  const anyFail = Array.from(logEl().querySelectorAll('.fail')).length > 0
  setStatus(anyFail ? 'Completed with failures' : 'All tests passed', anyFail ? 'fail' : 'ok')
}

document.getElementById('run').addEventListener('click', runTests)
