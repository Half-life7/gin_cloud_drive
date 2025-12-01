// 用户认证和个人中心功能

// 页面加载完成后初始化
window.addEventListener('load', function() {
	// 检查登录状态
	checkLoginStatus();
	
	// 初始化主题
	const savedTheme = localStorage.getItem('theme');
	if (savedTheme) {
		// 如果有保存的主题，应用它
		setTheme(savedTheme);
	} else {
		// 否则检查系统偏好
		const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
		setTheme(prefersDark ? 'dark' : 'light');
	}
});

// 检查登录状态
function checkLoginStatus() {
	const cookies = document.cookie.split(';');
	let isLoggedIn = false;
	
	for (const cookie of cookies) {
		const [name, value] = cookie.trim().split('=');
		if (name === 'auth_token' && value === 'admin_auth_token') {
			isLoggedIn = true;
			break;
		}
	}

	if (isLoggedIn) {
		// 已登录状态
		document.getElementById('loginSection').style.display = 'none';
		document.getElementById('userSection').style.display = 'block';
		document.getElementById('userInitial').textContent = '管';
		document.getElementById('avatar').style.backgroundColor = '#28a745';
	} else {
		// 未登录状态
		document.getElementById('loginSection').style.display = 'block';
		document.getElementById('userSection').style.display = 'none';
		document.getElementById('userInitial').textContent = '游';
		document.getElementById('avatar').style.backgroundColor = '#007bff';
	}
}

// 显示登录表单
function showLoginForm() {
	document.getElementById('loginModal').style.display = 'flex';
}

// 关闭登录表单
function closeLoginModal() {
	document.getElementById('loginModal').style.display = 'none';
	document.getElementById('loginStatus').innerHTML = '';
	document.getElementById('username').value = '';
	document.getElementById('password').value = '';
}

// 登录
function login() {
	const username = document.getElementById('username').value;
	const password = document.getElementById('password').value;
	const loginStatus = document.getElementById('loginStatus');

	if (!username || !password) {
		loginStatus.innerHTML = '<p style="color: red;">请输入用户名和密码</p>';
		return;
	}

	fetch('/api/auth/login', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({ username, password }),
	})
	.then(response => response.json())
	.then(data => {
		if (data.code === 200) {
			loginStatus.innerHTML = '<p style="color: green;">登录成功，页面将自动刷新...</p>';
			// 关闭登录表单并刷新页面
			setTimeout(() => {
				closeLoginModal();
				// 刷新页面，确保所有组件都能获取到最新的登录状态
				window.location.reload();
			}, 1000);
		} else {
			loginStatus.innerHTML = `<p style="color: red;">登录失败: ${data.message}</p>`;
		}
	})
	.catch(error => {
		console.error('登录失败:', error);
		loginStatus.innerHTML = '<p style="color: red;">登录失败</p>';
	});
}

// 登出
function logout() {
	fetch('/api/auth/logout', {
		method: 'POST',
	})
	.then(response => response.json())
	.then(data => {
		if (data.code === 200) {
			// 更新登录状态
			checkLoginStatus();
			// 刷新页面
			window.location.reload();
		}
	})
	.catch(error => {
		console.error('登出失败:', error);
	});
}

// 显示个人中心
function showUserCenter() {
	alert('个人中心功能开发中...');
	// 这里可以添加个人中心的实现，比如显示用户信息、系统设置等
}

// 点击模态框外部关闭登录表单
document.addEventListener('click', function(e) {
	const modal = document.getElementById('loginModal');
	if (e.target === modal) {
		closeLoginModal();
	}
});

// 暗色模式切换功能

// 设置主题
function setTheme(theme) {
	console.log('setTheme called with theme:', theme);
	const body = document.body;
	console.log('Body element:', body);
	
	// 获取所有太阳和月亮图标
	const sunIcons = document.querySelectorAll('.theme-icon.sun');
	const moonIcons = document.querySelectorAll('.theme-icon.moon');
	console.log('Sun icons:', sunIcons);
	console.log('Moon icons:', moonIcons);
	
	if (theme === 'dark') {
		console.log('Adding dark class to body');
		body.classList.add('dark');
		// 隐藏所有太阳图标，显示所有月亮图标
		sunIcons.forEach(icon => icon.style.display = 'none');
		moonIcons.forEach(icon => icon.style.display = 'inline');
	} else {
		console.log('Removing dark class from body');
		body.classList.remove('dark');
		// 显示所有太阳图标，隐藏所有月亮图标
		sunIcons.forEach(icon => icon.style.display = 'inline');
		moonIcons.forEach(icon => icon.style.display = 'none');
	}
	
	// 保存主题偏好到本地存储
	localStorage.setItem('theme', theme);
	console.log('Theme saved to localStorage:', theme);
}

// 切换主题
function toggleTheme() {
	console.log('toggleTheme called');
	const currentTheme = localStorage.getItem('theme') || 'light';
	console.log('Current theme:', currentTheme);
	const newTheme = currentTheme === 'light' ? 'dark' : 'light';
	console.log('New theme:', newTheme);
	setTheme(newTheme);
}

// 监听系统主题变化
window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function(e) {
	const savedTheme = localStorage.getItem('theme');
	// 只有在用户没有明确设置主题时，才跟随系统变化
	if (!savedTheme) {
		setTheme(e.matches ? 'dark' : 'light');
	}
});