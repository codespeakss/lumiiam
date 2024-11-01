document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault();
    
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    
    // 简单的登录验证示例（请替换为实际的验证逻辑）
    if (username === 'admin' && password === 'password') {
        alert('登录成功！');
        // 可以在这里重定向到另一个页面
    } else {
        document.getElementById('errorMessage').innerText = '用户名或密码错误';
    }
});
