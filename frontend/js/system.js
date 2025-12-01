// 系统状态监控
let systemChart = null;

// 页面加载完成后初始化
window.addEventListener('load', function() {
	// 初始化图表
	initChart();
	// 加载系统信息
	loadSystemInfo();
	// 加载历史数据
	loadHistoryData(1);
	// 定时更新系统信息
	setInterval(loadSystemInfo, 5000);
});

// 初始化图表
function initChart() {
	// 检查Chart是否可用
	if (typeof Chart === 'undefined' || Chart === null) {
		console.warn('Chart.js library is not available, chart functionality will be disabled');
		systemChart = null;
		return;
	}
	
	const ctx = document.getElementById('systemChart').getContext('2d');
	systemChart = new Chart(ctx, {
		type: 'line',
		data: {
			labels: [],
			datasets: [
				{
					label: 'CPU使用率 (%)',
					data: [],
					borderColor: 'rgba(255, 99, 132, 1)',
					backgroundColor: 'rgba(255, 99, 132, 0.2)',
					fill: false,
					borderWidth: 2
				},
				{
					label: '内存使用率 (%)',
					data: [],
					borderColor: 'rgba(54, 162, 235, 1)',
					backgroundColor: 'rgba(54, 162, 235, 0.2)',
					fill: false,
					borderWidth: 2
				},
				{
					label: '磁盘使用率 (%)',
					data: [],
					borderColor: 'rgba(75, 192, 192, 1)',
					backgroundColor: 'rgba(75, 192, 192, 0.2)',
					fill: false,
					borderWidth: 2
				}
			]
		},
		options: {
			responsive: true,
			maintainAspectRatio: false,
			height: 400,
			scales: {
				y: {
					beginAtZero: true,
					max: 100
				}
			}
		}
	});
}

// 加载系统信息
function loadSystemInfo() {
	fetch('/api/system/info')
		.then(response => response.json())
		.then(data => {
			updateSystemInfo(data);
		})
		.catch(error => {
			console.error('获取系统信息失败:', error);
		});
}

// 更新系统信息
function updateSystemInfo(info) {
	// 更新CPU信息
	document.getElementById('cpuUsage').textContent = info.cpu.usage_percent.toFixed(1) + '%';
	document.getElementById('cpuCores').textContent = info.cpu.cores;

	// 更新内存信息
	document.getElementById('memoryUsage').textContent = info.memory.usage_percent.toFixed(1) + '%';
	document.getElementById('memoryAvailable').textContent = formatFileSize(info.memory.available);

	// 更新磁盘信息
	document.getElementById('diskUsage').textContent = info.disk.usage_percent.toFixed(1) + '%';
	document.getElementById('diskAvailable').textContent = formatFileSize(info.disk.free);

	// 更新系统信息
	document.getElementById('os').textContent = info.os;
	document.getElementById('hostname').textContent = info.hostname;
	document.getElementById('currentTime').textContent = info.time;
}

// 加载历史数据
function loadHistoryData(hours) {
	fetch(`/api/system/history?hours=${hours}`)
		.then(response => response.json())
		.then(data => {
			updateChart(data.data);
		})
		.catch(error => {
			console.error('获取历史数据失败:', error);
		});
}

// 更新图表
function updateChart(data) {
	// 检查图表是否初始化
	if (!systemChart) {
		console.warn('Chart is not initialized, skipping chart update');
		return;
	}
	
	// 准备数据
	const labels = [];
	const cpuData = [];
	const memoryData = [];
	const diskData = [];

	// 处理数据点
	data.forEach(point => {
		// 格式化时间
		const date = new Date(point.timestamp * 1000);
		labels.push(date.toLocaleTimeString());
		// 添加数据
		cpuData.push(point.cpu);
		memoryData.push(point.memory);
		diskData.push(point.disk);
	});

	// 更新图表
	systemChart.data.labels = labels;
	systemChart.data.datasets[0].data = cpuData;
	systemChart.data.datasets[1].data = memoryData;
	systemChart.data.datasets[2].data = diskData;
	systemChart.update();
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