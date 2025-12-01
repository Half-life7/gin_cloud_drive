// 网盘功能
let currentPath = '';
let isAdmin = false;
// 排序相关变量
let currentSortBy = 'name';
let currentSortOrder = 'asc';
// 待上传文件列表
let pendingFiles = [];

// 页面加载完成后初始化
window.addEventListener('load', function() {
	// 检查是否已登录
	checkLogin();
	// 初始化排序图标
	updateSortIcons();
	// 加载文件列表
	loadFileList();
	// 初始化上传表单
	initUploadForm();
});

// 检查登录状态
function checkLogin() {
	// 简单检查Cookie是否存在
	const cookies = document.cookie.split(';');
	isAdmin = false;
	for (const cookie of cookies) {
		const [name, value] = cookie.trim().split('=');
		if (name === 'auth_token' && value === 'admin_auth_token') {
			isAdmin = true;
			break;
		}
	}
	console.log('Login status checked:', isAdmin);
}

// 加载文件列表
function loadFileList() {
	fetch(`/api/file/list?path=${currentPath}&sort_by=${currentSortBy}&sort_order=${currentSortOrder}`)
		.then(response => response.json())
		.then(data => {
			if (data.code === 200) {
				renderFileList(data.data);
			} else {
				showMessage('获取文件列表失败', 'error');
			}
		})
		.catch(error => {
			console.error('获取文件列表失败:', error);
			showMessage('获取文件列表失败', 'error');
		});
}

// 排序文件
function sortFiles(sortBy) {
	// 如果点击的是当前排序字段，则切换排序方向
	if (currentSortBy === sortBy) {
		currentSortOrder = currentSortOrder === 'asc' ? 'desc' : 'asc';
	} else {
		// 否则，设置新的排序字段和默认升序
		currentSortBy = sortBy;
		currentSortOrder = 'asc';
	}

	// 更新排序图标
	updateSortIcons();

	// 重新加载文件列表
	loadFileList();
}

// 更新排序图标
function updateSortIcons() {
	// 重置所有图标
	const icons = ['sortNameIcon', 'sortSizeIcon', 'sortTimeIcon', 'sortTypeIcon'];
	icons.forEach(iconId => {
		document.getElementById(iconId).textContent = '↑';
	});

	// 设置当前排序字段的图标
	const currentIconId = `sort${currentSortBy.charAt(0).toUpperCase() + currentSortBy.slice(1)}Icon`;
	document.getElementById(currentIconId).textContent = currentSortOrder === 'asc' ? '↑' : '↓';

	// 更新按钮样式
	const buttons = ['sortName', 'sortSize', 'sortTime', 'sortType'];
	buttons.forEach(btnId => {
		document.getElementById(btnId).classList.remove('active');
	});
	document.getElementById(`sort${currentSortBy.charAt(0).toUpperCase() + currentSortBy.slice(1)}`).classList.add('active');
}

// 渲染文件列表
function renderFileList(files) {
	const fileList = document.getElementById('fileList');
	fileList.innerHTML = '';

	// 添加返回上一级目录按钮
	if (currentPath !== '') {
		const backItem = document.createElement('div');
		backItem.className = 'file-item';
		backItem.innerHTML = `
			<div class="file-info">
				<span class="file-icon">📁</span>
				<span>..</span>
			</div>
			<div class="file-actions">
				<button class="btn btn-secondary" onclick="navigateTo('${getParentPath(currentPath)}')">进入</button>
			</div>
		`;
		fileList.appendChild(backItem);
	}

	// 添加文件和目录
	files.forEach(file => {
		const fileItem = document.createElement('div');
		fileItem.className = 'file-item';

		// 文件图标
		let icon = '📄';
		if (file.is_directory) {
			icon = '📁';
		} else if (file.type === 'image') {
			icon = '🖼️';
		} else if (file.type === 'video') {
			icon = '🎬';
		} else if (file.type === 'audio') {
			icon = '🎵';
		} else if (file.type === 'document') {
			icon = '📋';
		}

		// 基本信息
		let itemHTML = `
			<div class="file-info">
				<span class="file-icon">${icon}</span>
				<span>${file.name}</span>
				<span style="color: #666; font-size: 0.8rem;">
					${file.is_directory ? '目录' : formatFileSize(file.size)} • ${file.modified_time}
				</span>
			</div>
			<div class="file-actions">
		`;

		// 添加操作按钮
		if (file.is_directory) {
			// 目录操作
			itemHTML += `<button class="btn btn-secondary" onclick="navigateTo('${file.path}')">进入</button>`;
			if (isAdmin) {
				itemHTML += `<button class="btn btn-primary" onclick="showFolderSelector('move', '${file.path}')">移动</button>`;
				itemHTML += `<button class="btn btn-danger" onclick="deleteFile('${file.path}')">删除</button>`;
			}
		} else {
			// 文件操作
			itemHTML += `<button class="btn btn-secondary" onclick="downloadFile('${file.path}')">下载</button>`;
			itemHTML += `<button class="btn btn-primary" onclick="previewFile('${file.path}')">预览</button>`;
			if (isAdmin) {
				itemHTML += `<button class="btn btn-primary" onclick="showFolderSelector('move', '${file.path}')">移动</button>`;
				itemHTML += `<button class="btn btn-danger" onclick="deleteFile('${file.path}')">删除</button>`;
			}
		}

		itemHTML += `</div>`;
		fileItem.innerHTML = itemHTML;
		fileList.appendChild(fileItem);
	});
}

// 获取父路径
function getParentPath(path) {
	const parts = path.split('/');
	parts.pop();
	return parts.join('/');
}

// 导航到目录
function navigateTo(path) {
	currentPath = path;
	loadFileList();
}

// 下载文件
function downloadFile(path) {
	window.location.href = `/api/file/download/${path}`;
}

// 预览文件
function previewFile(path) {
	window.open(`/api/file/preview/${path}`, '_blank');
}

// 删除文件
function deleteFile(path) {
	if (!confirm('确定要删除这个文件/目录吗？')) {
		return;
	}

	console.log('Deleting file:', path);
	fetch(`/api/file/delete/${encodeURIComponent(path)}`, {
		method: 'DELETE',
	})
	.then(response => {
		console.log('Delete response status:', response.status);
		return response.json();
	})
	.then(data => {
		console.log('Delete response data:', data);
		if (data.code === 200) {
			showMessage('删除成功', 'success');
			loadFileList();
		} else {
			showMessage(`删除失败: ${data.message}`, 'error');
		}
	})
	.catch(error => {
		console.error('删除文件失败:', error);
		showMessage(`删除失败: ${error.message}`, 'error');
	});
}

// 初始化上传表单
function initUploadForm() {
	const uploadForm = document.getElementById('uploadForm');
	uploadForm.addEventListener('submit', function(e) {
		e.preventDefault();
		uploadFile();
	});
	
	// 初始化拖拽上传
	initDragAndDrop();
}

// 初始化拖拽上传
function initDragAndDrop() {
	const dropZone = document.getElementById('dropZone');
	const fileInput = document.getElementById('file');
	
	// 拖拽事件处理
	dropZone.addEventListener('dragenter', handleDragEnter);
	dropZone.addEventListener('dragover', handleDragOver);
	dropZone.addEventListener('dragleave', handleDragLeave);
	dropZone.addEventListener('drop', handleDrop);
	
	// 点击选择文件
	dropZone.addEventListener('click', function(e) {
		// 避免重复触发：如果点击的是label或其子元素，就不再次触发fileInput.click()
		let target = e.target;
		while (target) {
			if (target.tagName === 'LABEL') {
				return; // label会自动触发fileInput.click()，不需要再次触发
			}
		target = target.parentElement;
		}
		fileInput.click();
	});
	
	// 文件选择变化时添加到缓冲区域
	fileInput.addEventListener('change', function() {
		// 只有选择了文件才添加到缓冲区域，取消选择时不触发
		if (this.files.length > 0) {
			addFilesToBuffer(this.files);
			// 清空文件输入，以便可以再次选择相同的文件
			this.value = '';
		}
	});
}

// 拖拽进入事件
function handleDragEnter(e) {
	e.preventDefault();
	e.stopPropagation();
	const dropZone = document.getElementById('dropZone');
	dropZone.classList.add('drag-over');
}

// 拖拽悬停事件
function handleDragOver(e) {
	e.preventDefault();
	e.stopPropagation();
	const dropZone = document.getElementById('dropZone');
	dropZone.classList.add('drag-over');
}

// 拖拽离开事件
function handleDragLeave(e) {
	e.preventDefault();
	e.stopPropagation();
	const dropZone = document.getElementById('dropZone');
	dropZone.classList.remove('drag-over');
}

// 拖拽放下事件
function handleDrop(e) {
	e.preventDefault();
	e.stopPropagation();
	const dropZone = document.getElementById('dropZone');
	dropZone.classList.remove('drag-over');
	
	// 获取拖拽的文件
	const files = e.dataTransfer.files;
	if (files.length > 0) {
		addFilesToBuffer(files);
	}
}

// 添加文件到缓冲区域
function addFilesToBuffer(files) {
	for (let i = 0; i < files.length; i++) {
		pendingFiles.push(files[i]);
	}
	updateFileBufferDisplay();
}

// 更新文件缓冲区域显示
function updateFileBufferDisplay() {
	const bufferArea = document.getElementById('fileBuffer');
	const bufferList = document.getElementById('pendingFilesList');
	const startUploadBtn = document.getElementById('startUploadBtn');
	
	// 清空当前列表
	bufferList.innerHTML = '';
	
	if (pendingFiles.length === 0) {
		bufferArea.style.display = 'none';
		startUploadBtn.disabled = true;
		return;
	}
	
	// 显示缓冲区域
	bufferArea.style.display = 'block';
	startUploadBtn.disabled = false;
	
	// 添加文件到列表
	pendingFiles.forEach((file, index) => {
		const fileItem = document.createElement('div');
		fileItem.className = 'pending-file-item';
		fileItem.innerHTML = `
			<div class="pending-file-info">
				<span class="file-icon">${getFileIcon(file.name)}</span>
				<span class="pending-file-name">${file.name}</span>
				<span class="pending-file-size">${formatFileSize(file.size)}</span>
			</div>
			<button class="btn btn-danger btn-sm" onclick="removeFileFromBuffer(${index})">删除</button>
		`;
		bufferList.appendChild(fileItem);
	});
}

// 从缓冲区域移除文件
function removeFileFromBuffer(index) {
	pendingFiles.splice(index, 1);
	updateFileBufferDisplay();
}

// 获取文件图标
function getFileIcon(filename) {
	const ext = filename.split('.').pop().toLowerCase();
	const iconMap = {
		// 图片文件
		'jpg': '🖼️', 'jpeg': '🖼️', 'png': '🖼️', 'gif': '🖼️', 'bmp': '🖼️', 'svg': '🖼️',
		// 文档文件
		'doc': '📋', 'docx': '📋', 'pdf': '📋', 'txt': '📋', 'md': '📋', 'rtf': '📋',
		// 视频文件
		'mp4': '🎬', 'avi': '🎬', 'mov': '🎬', 'wmv': '🎬', 'flv': '🎬',
		// 音频文件
		'mp3': '🎵', 'wav': '🎵', 'flac': '🎵', 'aac': '🎵',
		// 压缩文件
		'zip': '📦', 'rar': '📦', '7z': '📦', 'tar': '📦', 'gz': '📦',
		// 代码文件
		'js': '💻', 'html': '💻', 'css': '💻', 'go': '💻', 'py': '💻', 'java': '💻', 'c': '💻', 'cpp': '💻',
		// 其他文件
		'default': '📄'
	};
	return iconMap[ext] || iconMap.default;
}

// 上传单个文件
function uploadSingleFile(file, path, uploadStatus, totalFiles, currentFile) {
	return new Promise((resolve, reject) => {
		// 检查文件大小（100MB限制）
		const maxSize = 100 * 1024 * 1024; // 100MB
		if (file.size > maxSize) {
			reject(new Error(`文件 ${file.name} 大小超过限制（最大100MB）`));
			return;
		}

		const formData = new FormData();
		formData.append('file', file);
		formData.append('path', path);

		// 使用XMLHttpRequest实现上传进度
		const xhr = new XMLHttpRequest();

		// 上传进度
		xhr.upload.addEventListener('progress', function(e) {
			if (e.lengthComputable) {
				const percent = Math.round((e.loaded / e.total) * 100);
				uploadStatus.className = 'info';
				uploadStatus.innerHTML = `正在上传文件 ${currentFile}/${totalFiles}：${file.name} (${percent}%)`;
			}
		});

		// 上传完成
		xhr.addEventListener('load', function() {
			if (xhr.status === 200) {
				try {
					const data = JSON.parse(xhr.responseText);
					if (data.code === 200) {
						resolve(file.name);
					} else {
						reject(new Error(`文件 ${file.name} 上传失败: ${data.message}`));
					}
				} catch (error) {
					reject(new Error(`文件 ${file.name} 上传失败：服务器返回无效响应`));
				}
			} else {
				let errorMsg = `文件 ${file.name} 上传失败：HTTP ${xhr.status}`;
				if (xhr.status === 401) {
					errorMsg += "（未授权，请先登录）";
				} else if (xhr.status === 500) {
					errorMsg += "（服务器内部错误）";
				} else if (xhr.status === 404) {
					errorMsg += "（上传接口不存在）";
				}
				reject(new Error(errorMsg));
			}
		});

		// 上传错误
		xhr.addEventListener('error', function() {
			reject(new Error(`文件 ${file.name} 上传失败：网络错误`));
		});

		// 上传超时
		xhr.addEventListener('timeout', function() {
			reject(new Error(`文件 ${file.name} 上传失败：超时`));
		});

		// 发送请求
		xhr.open('POST', '/api/file/upload');
		xhr.send(formData);
	});
}

// 上传文件（支持批量）
function uploadFile() {
	const pathInput = document.getElementById('path');
	const uploadStatus = document.getElementById('uploadStatus');

	// 检查缓冲区域是否有文件
	if (pendingFiles.length === 0) {
		uploadStatus.className = 'error';
		uploadStatus.innerHTML = '请先选择要上传的文件';
		return;
	}

	const path = pathInput.value;
	const totalFiles = pendingFiles.length;
	let uploadedFiles = 0;
	let failedFiles = 0;
	let errorMessages = [];

	// 显示上传状态
	uploadStatus.className = 'info';
	uploadStatus.innerHTML = `准备上传 ${totalFiles} 个文件...`;

	// 上传所有文件
	const uploadPromises = [];
	for (let i = 0; i < pendingFiles.length; i++) {
		const file = pendingFiles[i];
		uploadPromises.push(uploadSingleFile(file, path, uploadStatus, totalFiles, i + 1));
	}

	// 等待所有上传完成
	Promise.allSettled(uploadPromises).then(results => {
		results.forEach(result => {
			if (result.status === 'fulfilled') {
				uploadedFiles++;
			} else {
				failedFiles++;
				errorMessages.push(result.reason.message);
			}
		});

		// 显示上传结果
		if (failedFiles === 0) {
			uploadStatus.className = 'success';
			uploadStatus.innerHTML = `全部 ${totalFiles} 个文件上传成功`;
		} else if (uploadedFiles === 0) {
			uploadStatus.className = 'error';
			uploadStatus.innerHTML = `全部 ${totalFiles} 个文件上传失败：<br>${errorMessages.join('<br>')}`;
		} else {
			uploadStatus.className = 'info';
			uploadStatus.innerHTML = `上传完成：${uploadedFiles} 个成功，${failedFiles} 个失败<br>失败原因：<br>${errorMessages.join('<br>')}`;
		}

		// 刷新文件列表
		loadFileList();

		// 清空缓冲区域
		pendingFiles = [];
		updateFileBufferDisplay();

		// 3秒后清除状态
		setTimeout(() => {
			uploadStatus.innerHTML = '';
			uploadStatus.className = '';
		}, 5000);
	});
}



// 格式化文件大小
function formatFileSize(size) {
	if (size < 1024) {
		return size + ' B';
	} else if (size < 1024 * 1024) {
		return (size / 1024).toFixed(2) + ' KB';
	} else if (size < 1024 * 1024 * 1024) {
		return (size / (1024 * 1024)).toFixed(2) + ' MB';
	} else {
		return (size / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
	}
}

// 全局变量，用于存储当前选中的文件夹路径
let currentFolderPath = '';
// 全局变量，用于存储待移动的文件路径
let fileToMove = '';

// 显示文件夹选择器
function showFolderSelector(mode = 'upload', filePath = '') {
	const modal = document.getElementById('folderModal');
	const title = document.getElementById('folderModalTitle');
	modal.style.display = 'flex';
	
	// 根据模式设置标题
	if (mode === 'move') {
		title.textContent = '选择移动目标路径';
		fileToMove = filePath;
	} else {
		title.textContent = '选择上传路径';
		fileToMove = '';
	}
	
	// 加载根目录下的文件夹
	loadFolderList('');
}

// 关闭文件夹选择器
function closeFolderModal() {
	document.getElementById('folderModal').style.display = 'none';
	fileToMove = '';
}

// 加载文件夹列表
function loadFolderList(path) {
	currentFolderPath = path;
	
	// 更新当前路径显示
	document.getElementById('currentPath').textContent = '/' + path;
	
	// 加载文件夹列表
	fetch(`/api/file/list?path=${path}`)
		.then(response => response.json())
		.then(data => {
			if (data.code === 200) {
				// 过滤出文件夹
				const folders = data.data.filter(item => item.is_directory);
				// 渲染文件夹列表
				renderFolderList(folders);
				// 更新面包屑导航
				updateBreadcrumb(path);
			} else {
				console.error('获取文件夹列表失败:', data.message);
			}
		})
		.catch(error => {
			console.error('获取文件夹列表失败:', error);
		});
}

// 渲染文件夹列表
function renderFolderList(folders) {
	const folderList = document.getElementById('folderList');
	folderList.innerHTML = '';

	// 添加返回上一级目录按钮（如果不是根目录）
	if (currentFolderPath !== '') {
		const backItem = document.createElement('div');
		backItem.style.cursor = 'pointer';
		backItem.style.padding = '1rem';
		backItem.style.border = '1px solid #eee';
		backItem.style.borderRadius = '5px';
		backItem.style.textAlign = 'center';
		backItem.style.transition = 'all 0.3s ease';
		backItem.innerHTML = `
			<div style="font-size: 2rem; margin-bottom: 0.5rem;">📁</div>
			<div>..</div>
		`;
		backItem.onclick = () => {
			const parentPath = getParentPath(currentFolderPath);
			loadFolderList(parentPath);
		};
		backItem.onmouseover = () => {
			backItem.style.backgroundColor = '#f0f0f0';
			backItem.style.transform = 'translateY(-2px)';
		};
		backItem.onmouseout = () => {
			backItem.style.backgroundColor = '#fff';
			backItem.style.transform = 'translateY(0)';
		};
		folderList.appendChild(backItem);
	}

	// 添加文件夹
	folders.forEach(folder => {
		const folderItem = document.createElement('div');
		folderItem.style.cursor = 'pointer';
		folderItem.style.padding = '1rem';
		folderItem.style.border = '1px solid #eee';
		folderItem.style.borderRadius = '5px';
		folderItem.style.textAlign = 'center';
		folderItem.style.transition = 'all 0.3s ease';
		folderItem.innerHTML = `
			<div style="font-size: 2rem; margin-bottom: 0.5rem;">📁</div>
			<div style="word-break: break-all;">${folder.name}</div>
		`;
		folderItem.onclick = () => {
			loadFolderList(folder.path);
		};
		folderItem.onmouseover = () => {
			folderItem.style.backgroundColor = '#f0f0f0';
			folderItem.style.transform = 'translateY(-2px)';
		};
		folderItem.onmouseout = () => {
			folderItem.style.backgroundColor = '#fff';
			folderItem.style.transform = 'translateY(0)';
		};
		folderList.appendChild(folderItem);
	});
}

// 更新面包屑导航
function updateBreadcrumb(path) {
	const breadcrumb = document.getElementById('folderBreadcrumb');
	breadcrumb.innerHTML = '';

	// 添加根目录
	const rootItem = document.createElement('span');
	rootItem.className = 'breadcrumb-item';
	rootItem.textContent = '/';
	rootItem.style.cursor = 'pointer';
	rootItem.style.color = '#007bff';
	rootItem.onclick = () => loadFolderList('');
	breadcrumb.appendChild(rootItem);

	// 添加路径中的各个目录
	if (path !== '') {
		const parts = path.split('/');
		let currentPath = '';
		for (let i = 0; i < parts.length; i++) {
			const part = parts[i];
			if (part === '') continue;
			currentPath += '/' + part;
			
			const separator = document.createElement('span');
			separator.textContent = ' > ';
			separator.style.color = '#666';
			breadcrumb.appendChild(separator);

			const item = document.createElement('span');
			item.className = 'breadcrumb-item';
			item.textContent = part;
			item.style.cursor = 'pointer';
			item.style.color = '#007bff';
			const fullPath = currentPath.substring(1); // 去掉开头的斜杠
			item.onclick = () => loadFolderList(fullPath);
			breadcrumb.appendChild(item);
		}
	}
}

// 选择当前文件夹
function selectCurrentFolder() {
	if (fileToMove) {
		// 移动文件
		moveFile(fileToMove, currentFolderPath);
	} else {
		// 设置上传路径
		document.getElementById('path').value = currentFolderPath;
	}
	closeFolderModal();
}

// 移动文件
function moveFile(oldPath, newPath) {
	console.log('Moving file:', oldPath, 'to', newPath);
	
	// 确保旧路径和新路径不相同
	if (oldPath === newPath) {
		showMessage('源路径和目标路径相同，无需移动', 'info');
		return;
	}
	
	// 确保新路径是目录时，构建正确的目标路径
	fetch('/api/file/move', {
		method: 'PUT',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({
			old_path: oldPath,
			new_path: newPath
		}),
	})
	.then(response => {
		console.log('Move response status:', response.status);
		return response.json();
	})
	.then(data => {
		console.log('Move response data:', data);
		if (data.code === 200) {
			showMessage('文件移动成功', 'success');
			loadFileList();
		} else {
			showMessage(`文件移动失败: ${data.message}`, 'error');
		}
	})
	.catch(error => {
		console.error('移动文件失败:', error);
		showMessage(`文件移动失败: ${error.message}`, 'error');
	});
}

// 导航到文件夹
function navigateToFolder(path) {
	loadFolderList(path);
}

// 显示消息
function showMessage(message, type) {
	const uploadStatus = document.getElementById('uploadStatus');
	const color = type === 'success' ? 'green' : 'red';
	uploadStatus.innerHTML = `<p style="color: ${color};">${message}</p>`;
	// 3秒后自动清除
	setTimeout(() => {
		uploadStatus.innerHTML = '';
	}, 3000);
}

// 新建文件夹
function createNewFolder() {
	const newFolderName = document.getElementById('newFolderName').value.trim();
	const path = document.getElementById('path').value;
	const statusDiv = document.getElementById('createFolderStatus');

	if (!newFolderName) {
		statusDiv.className = 'error';
		statusDiv.innerHTML = '请输入文件夹名称';
		return;
	}

	// 发送请求创建文件夹
	fetch('/api/file/mkdir', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify({
			Path: path + '/' + newFolderName
		}),
	})
	.then(response => response.json())
	.then(data => {
		if (data.code === 200) {
			statusDiv.className = 'success';
			statusDiv.innerHTML = '文件夹创建成功';
			// 清空输入框
			document.getElementById('newFolderName').value = '';
			// 刷新文件列表
			loadFileList();
			// 3秒后清除状态
			setTimeout(() => {
				statusDiv.innerHTML = '';
				statusDiv.className = '';
			}, 3000);
		} else {
			statusDiv.className = 'error';
			statusDiv.innerHTML = `创建失败: ${data.message}`;
		}
	})
	.catch(error => {
		console.error('创建文件夹失败:', error);
		statusDiv.className = 'error';
		statusDiv.innerHTML = '创建失败: 网络错误';
	});
}