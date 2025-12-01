// 日志管理页面JavaScript
let currentPage = 1;
let currentParams = {};

// 页面加载完成后初始化
window.addEventListener('load', function() {
    // 初始化日志列表
    loadLogs();
    // 初始化日志统计
    loadLogStats();
    // 绑定查询表单提交事件
    document.getElementById('logQueryForm').addEventListener('submit', function(e) {
        e.preventDefault();
        currentPage = 1;
        loadLogs();
    });
});

// 加载日志列表
function loadLogs() {
    // 显示加载状态
    const tbody = document.getElementById('logTableBody');
    tbody.innerHTML = '<tr><td colspan="8" style="text-align: center; padding: 20px;">加载中...</td></tr>';
    
    // 收集查询参数
    const form = document.getElementById('logQueryForm');
    const formData = new FormData(form);
    const params = {
        page: currentPage,
        page_size: formData.get('page_size') || 20
    };
    
    // 添加其他查询参数
    const startDate = formData.get('start_date');
    const endDate = formData.get('end_date');
    const level = formData.get('level');
    const type = formData.get('type');
    const ip = formData.get('ip');
    const action = formData.get('action');
    const file = formData.get('file');
    
    if (startDate) params.start_date = startDate;
    if (endDate) params.end_date = endDate;
    if (level) params.level = level;
    if (type) params.type = type;
    if (ip) params.ip = ip;
    if (action) params.action = action;
    if (file) params.file = file;
    
    currentParams = params;
    
    // 发送请求获取日志列表
    fetch(`/api/log/list?${new URLSearchParams(params)}`)
        .then(response => response.json())
        .then(data => {
            if (data.code === 200) {
                renderLogs(data.data);
                renderPagination(data.data);
            } else {
                tbody.innerHTML = `<tr><td colspan="8" style="text-align: center; padding: 20px; color: red;">${data.message || '获取日志失败'}</td></tr>`;
            }
        })
        .catch(error => {
            console.error('获取日志失败:', error);
            tbody.innerHTML = '<tr><td colspan="8" style="text-align: center; padding: 20px; color: red;">获取日志失败，请检查网络连接</td></tr>';
        });
}

// 渲染日志列表
function renderLogs(result) {
    const tbody = document.getElementById('logTableBody');
    
    if (result.logs.length === 0) {
        tbody.innerHTML = '<tr><td colspan="8" style="text-align: center; padding: 20px;">没有找到匹配的日志记录</td></tr>';
        return;
    }
    
    let html = '';
    result.logs.forEach(log => {
        // 格式化时间
        const time = new Date(log.timestamp).toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        });
        
        // 格式化文件大小
        const size = formatFileSize(log.size || 0);
        
        // 生成日志行
        html += `
            <tr>
                <td>${time}</td>
                <td class="log-level ${log.level.toLowerCase()}">${log.level}</td>
                <td>${log.type}</td>
                <td>${log.ip}</td>
                <td>${log.action}</td>
                <td>${log.file || '-'}</td>
                <td>${size}</td>
                <td>${log.details || '-'}</td>
            </tr>
        `;
    });
    
    tbody.innerHTML = html;
}

// 渲染分页控件
function renderPagination(result) {
    const pagination = document.getElementById('pagination');
    const totalPages = Math.ceil(result.total / result.page_size);
    
    if (totalPages <= 1) {
        pagination.innerHTML = '<div style="text-align: center; padding: 10px;">共 ' + result.total + ' 条记录</div>';
        return;
    }
    
    let html = '<div class="pagination-controls">';
    
    // 上一页按钮
    html += `<button class="btn btn-secondary ${currentPage === 1 ? 'disabled' : ''}" onclick="changePage(${currentPage - 1})">上一页</button>`;
    
    // 页码按钮
    const startPage = Math.max(1, currentPage - 2);
    const endPage = Math.min(totalPages, currentPage + 2);
    
    if (startPage > 1) {
        html += `<button class="btn btn-secondary" onclick="changePage(1)">1</button>`;
        if (startPage > 2) {
            html += '<span class="pagination-ellipsis">...</span>';
        }
    }
    
    for (let i = startPage; i <= endPage; i++) {
        html += `<button class="btn ${i === currentPage ? 'btn-primary' : 'btn-secondary'}" onclick="changePage(${i})">${i}</button>`;
    }
    
    if (endPage < totalPages) {
        if (endPage < totalPages - 1) {
            html += '<span class="pagination-ellipsis">...</span>';
        }
        html += `<button class="btn btn-secondary" onclick="changePage(${totalPages})"></button>`;
    }
    
    // 下一页按钮
    html += `<button class="btn btn-secondary ${currentPage === totalPages ? 'disabled' : ''}" onclick="changePage(${currentPage + 1})">下一页</button>`;
    
    // 页码信息
    html += `<span class="pagination-info">第 ${currentPage} / ${totalPages} 页，共 ${result.total} 条记录</span>`;
    
    html += '</div>';
    
    pagination.innerHTML = html;
}

// 切换页码
function changePage(page) {
    if (page < 1) return;
    currentPage = page;
    loadLogs();
}

// 加载日志统计信息
function loadLogStats() {
    fetch('/api/log/stats')
        .then(response => response.json())
        .then(data => {
            if (data.code === 200) {
                updateLogStats(data.data);
            }
        })
        .catch(error => {
            console.error('获取日志统计失败:', error);
        });
}

// 更新日志统计信息
function updateLogStats(stats) {
    document.getElementById('totalLogs').textContent = stats.total_logs;
    document.getElementById('fileLogs').textContent = stats.type_stats.FILE || 0;
    document.getElementById('userLogs').textContent = stats.type_stats.USER || 0;
    document.getElementById('accessLogs').textContent = stats.type_stats.ACCESS || 0;
}

// 重置查询表单
function resetQueryForm() {
    document.getElementById('logQueryForm').reset();
    currentPage = 1;
    loadLogs();
}

// 清理旧日志
function clearOldLogs() {
    if (confirm('确定要清理旧日志吗？此操作不可恢复。')) {
        const days = prompt('请输入要保留的日志天数（默认7天）:', '7');
        if (days === null) return;
        
        const daysInt = parseInt(days) || 7;
        
        fetch(`/api/log/clear?days=${daysInt}`, {
            method: 'DELETE'
        })
        .then(response => response.json())
        .then(data => {
            if (data.code === 200) {
                alert(data.message);
                // 重新加载日志列表和统计信息
                loadLogs();
                loadLogStats();
            } else {
                alert(data.message || '清理旧日志失败');
            }
        })
        .catch(error => {
            console.error('清理旧日志失败:', error);
            alert('清理旧日志失败，请检查网络连接');
        });
    }
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
