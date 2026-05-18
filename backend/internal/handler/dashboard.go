package handler

import (
	"net/http"
)

// DashboardHandler sirve el dashboard HTML
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(dashboardHTML))
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>PsicoSST Cloud - Dashboard</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        
        .container {
            width: 100%;
            max-width: 800px;
            background: white;
            border-radius: 10px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            overflow: hidden;
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        
        .header h1 {
            font-size: 28px;
            margin-bottom: 10px;
        }
        
        .header p {
            opacity: 0.9;
            font-size: 14px;
        }
        
        .content {
            padding: 30px;
        }
        
        .section {
            margin-bottom: 30px;
        }
        
        .section h2 {
            font-size: 18px;
            margin-bottom: 15px;
            color: #333;
            border-bottom: 2px solid #667eea;
            padding-bottom: 10px;
        }
        
        .form-group {
            margin-bottom: 15px;
        }
        
        label {
            display: block;
            margin-bottom: 5px;
            color: #555;
            font-weight: 500;
        }
        
        input, button {
            font-family: inherit;
            font-size: 14px;
        }
        
        input[type="text"],
        input[type="password"] {
            width: 100%;
            padding: 10px 12px;
            border: 1px solid #ddd;
            border-radius: 5px;
            transition: border-color 0.3s;
        }
        
        input[type="text"]:focus,
        input[type="password"]:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }
        
        button {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-weight: 600;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        
        button:hover:not(:disabled) {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.3);
        }
        
        button:disabled {
            opacity: 0.6;
            cursor: not-allowed;
        }
        
        .status {
            padding: 12px 15px;
            border-radius: 5px;
            margin-bottom: 15px;
            font-size: 14px;
        }
        
        .status.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        
        .status.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        
        .status.info {
            background: #d1ecf1;
            color: #0c5460;
            border: 1px solid #bee5eb;
        }
        
        .token-display {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 5px;
            border: 1px solid #dee2e6;
            word-break: break-all;
            font-family: 'Courier New', monospace;
            font-size: 12px;
            max-height: 150px;
            overflow-y: auto;
            margin: 10px 0;
        }
        
        .endpoint-list {
            list-style: none;
        }
        
        .endpoint-list li {
            padding: 10px;
            margin-bottom: 8px;
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            border-radius: 3px;
        }
        
        .endpoint-list code {
            background: #e9ecef;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
            font-size: 12px;
        }
        
        .auth-section {
            display: none;
        }
        
        .auth-section.active {
            display: block;
        }
        
        .server-status {
            display: flex;
            align-items: center;
            gap: 10px;
            margin-bottom: 20px;
        }
        
        .status-indicator {
            width: 12px;
            height: 12px;
            border-radius: 50%;
            background: #28a745;
            animation: pulse 2s infinite;
        }
        
        @keyframes pulse {
            0% {
                box-shadow: 0 0 0 0 rgba(40, 167, 69, 0.4);
            }
            70% {
                box-shadow: 0 0 0 8px rgba(40, 167, 69, 0);
            }
            100% {
                box-shadow: 0 0 0 0 rgba(40, 167, 69, 0);
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🛡️ PsicoSST Cloud</h1>
            <p>Sistema de Seguridad Basada en el Comportamiento</p>
        </div>
        
        <div class="content">
            <!-- Server Status -->
            <div class="section">
                <div class="server-status">
                    <div class="status-indicator"></div>
                    <span id="status-text">Verificando servidor...</span>
                </div>
            </div>
            
            <!-- Login Section -->
            <div class="section" id="login-section">
                <h2>🔐 Iniciar Sesión</h2>
                <div id="login-message"></div>
                <form id="login-form" onsubmit="handleLogin(event)">
                    <div class="form-group">
                        <label for="email">Email:</label>
                        <input type="text" id="email" value="admin@example.com" required>
                    </div>
                    <div class="form-group">
                        <label for="password">Contraseña:</label>
                        <input type="password" id="password" value="password123" required>
                    </div>
                    <button type="submit" id="login-btn">Iniciar Sesión</button>
                </form>
            </div>
            
            <!-- Auth Section -->
            <div class="section auth-section" id="auth-section">
                <h2>✅ Autenticación Activa</h2>
                <button onclick="logout()" style="background: #dc3545;">Cerrar Sesión</button>
                <div id="user-info" style="margin-top: 20px; padding: 15px; background: #f8f9fa; border-radius: 5px;"></div>
            </div>
            
            <!-- API Endpoints -->
            <div class="section">
                <h2>📡 Endpoints Disponibles</h2>
                <ul class="endpoint-list">
                    <li><code>GET</code> <code>/health</code> - Estado del servidor</li>
                    <li><code>POST</code> <code>/api/v1/auth/login</code> - Autenticación</li>
                    <li><code>POST</code> <code>/api/v1/auth/logout</code> - Cerrar sesión</li>
                    <li><code>POST</code> <code>/api/v1/auth/refresh</code> - Renovar token</li>
                    <li><code>POST</code> <code>/api/v1/auth/change-password</code> - Cambiar contraseña</li>
                    <li><code>GET</code> <code>/api/v1/users</code> - Listar usuarios</li>
                    <li><code>GET</code> <code>/api/v1/employees</code> - Listar empleados</li>
                    <li><code>GET</code> <code>/api/v1/observations</code> - Listar observaciones</li>
                    <li><code>GET</code> <code>/api/v1/incidents</code> - Listar incidentes</li>
                    <li><code>GET</code> <code>/api/v1/analytics/dashboard</code> - Dashboard analítico</li>
                </ul>
            </div>
        </div>
    </div>
    
    <script>
        let authToken = null;
        let currentUser = null;
        
        // Verificar estado del servidor
        async function checkServerStatus() {
            try {
                const response = await fetch('/health');
                if (response.ok) {
                    document.getElementById('status-text').textContent = '✅ Servidor en línea';
                    document.querySelector('.status-indicator').style.background = '#28a745';
                } else {
                    throw new Error('Server error');
                }
            } catch (error) {
                document.getElementById('status-text').textContent = '❌ Servidor fuera de línea';
                document.querySelector('.status-indicator').style.background = '#dc3545';
            }
        }
        
        // Manejar login
        async function handleLogin(event) {
            event.preventDefault();
            
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            const messageDiv = document.getElementById('login-message');
            const loginBtn = document.getElementById('login-btn');
            
            messageDiv.innerHTML = '';
            loginBtn.disabled = true;
            
            try {
                const response = await fetch('/api/v1/auth/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ email, password })
                });
                
                const data = await response.json();
                
                if (response.ok && data.access_token) {
                    authToken = data.access_token;
                    currentUser = data.user;
                    localStorage.setItem('authToken', authToken);
                    localStorage.setItem('currentUser', JSON.stringify(currentUser));
                    showAuthSection();
                    messageDiv.innerHTML = '<div class="status success">✅ Autenticación exitosa</div>';
                } else {
                    messageDiv.innerHTML = '<div class="status error">❌ ' + (data.message || 'Error de autenticación') + '</div>';
                }
            } catch (error) {
                messageDiv.innerHTML = '<div class="status error">❌ Error: ' + error.message + '</div>';
            } finally {
                loginBtn.disabled = false;
            }
        }
        
        // Mostrar sección de autenticación
        function showAuthSection() {
            document.getElementById('login-section').style.display = 'none';
            document.getElementById('auth-section').classList.add('active');
            
            if (currentUser) {
                const userHTML = '<strong>Usuario Autenticado:</strong><br>' +
                    '📧 Email: ' + currentUser.email + '<br>' +
                    '👤 Nombre: ' + currentUser.first_name + ' ' + currentUser.last_name + '<br>' +
                    '🔑 ID: ' + currentUser.id + '<br>' +
                    '👑 Super Admin: ' + (currentUser.is_super_admin ? 'Sí ✅' : 'No') + '<br>' +
                    '<br>' +
                    '<strong>Token de Acceso:</strong><br>' +
                    '<div class="token-display">' + authToken + '</div>';
                document.getElementById('user-info').innerHTML = userHTML;
            }
        }
        
        // Logout
        function logout() {
            authToken = null;
            currentUser = null;
            localStorage.removeItem('authToken');
            localStorage.removeItem('currentUser');
            document.getElementById('login-section').style.display = 'block';
            document.getElementById('auth-section').classList.remove('active');
            document.getElementById('login-form').reset();
            document.getElementById('login-message').innerHTML = '';
        }
        
        // Inicializar
        window.addEventListener('load', function() {
            checkServerStatus();
            
            // Restaurar sesión si existe
            const savedToken = localStorage.getItem('authToken');
            const savedUser = localStorage.getItem('currentUser');
            if (savedToken && savedUser) {
                authToken = savedToken;
                currentUser = JSON.parse(savedUser);
                showAuthSection();
            }
        });
    </script>
</body>
</html>
`
