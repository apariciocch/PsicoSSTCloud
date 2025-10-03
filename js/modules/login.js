export function initLogin({ onLogout } = {}) {
  const loginForm = document.getElementById('login-form');
  if (!loginForm) return;

  const loginScreen = document.getElementById('login-screen');
  const loginError = document.getElementById('login-error');
  const app = document.getElementById('app');
  const userField = document.getElementById('login-username');
  const passField = document.getElementById('login-password');
  const logoutButton = document.getElementById('btnLogout');
  const credentials = { username: 'Armando', password: 'admin' };

  const showApp = () => {
    if (loginScreen) loginScreen.style.display = 'none';
    if (app) app.style.display = 'grid';
    if (loginError) loginError.textContent = '';
    if (passField) passField.value = '';
  };

  const showLogin = () => {
    if (app) app.style.display = 'none';
    if (loginScreen) loginScreen.style.display = 'grid';
    loginForm.reset();
    if (loginError) loginError.textContent = '';
    setTimeout(() => userField?.focus(), 100);
  };

  const logout = () => {
    sessionStorage.removeItem('psicosst-session');
    showLogin();
    if (typeof onLogout === 'function') {
      onLogout();
    }
  };

  logoutButton?.addEventListener('click', logout);

  const tryAutoLogin = () => {
    const savedSession = sessionStorage.getItem('psicosst-session');
    if (savedSession === 'active') {
      userField.value = credentials.username;
      showApp();
    } else {
      setTimeout(() => userField.focus(), 100);
    }
  };

  tryAutoLogin();

  loginForm.addEventListener('submit', (event) => {
    event.preventDefault();
    const username = userField.value.trim();
    const password = passField.value;

    if (username.toLowerCase() === credentials.username.toLowerCase() && password === credentials.password) {
      sessionStorage.setItem('psicosst-session', 'active');
      showApp();
    } else {
      if (loginError) {
        loginError.textContent = 'Usuario o contraseña incorrectos. (Armando / admin)';
      }
      passField.value = '';
      passField.focus();
    }
  });
}
