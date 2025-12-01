// 轮播图功能
let currentSlide = 0;
let slides = [];

// 初始化轮播图
function initCarousel() {
	// 从服务器获取轮播图图片
		fetch('/api/file/carousel')
			.then(response => response.json())
			.then(data => {
				if (data.code === 200) {
					// 过滤出图片文件，允许所有文件类型，确保所有图片都能显示
					slides = data.data.filter(item => {
						// 允许所有文件，或者根据扩展名判断
						const ext = item.name.toLowerCase().split('.').pop();
						const imageExts = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg'];
						return imageExts.includes(ext);
					});
					// 渲染轮播图
					renderCarousel();
					// 启动自动轮播
					startAutoCarousel();
				}
			})
			.catch(error => {
				console.error('获取轮播图失败:', error);
			});
}

// 渲染轮播图
function renderCarousel() {
	const carouselContainer = document.getElementById('carouselContainer');
	const carouselControls = document.getElementById('carouselControls');

	// 清空容器
	carouselContainer.innerHTML = '';
	carouselControls.innerHTML = '';

	console.log('Rendering carousel with slides:', slides);
	
	// 添加轮播图图片
	slides.forEach((slide, index) => {
		// 创建图片元素
		const img = document.createElement('img');
		// 直接使用slide.path作为图片路径，确保路径正确
		img.src = `/upload/${slide.path}`;
		img.className = 'carousel-item';
		img.onclick = () => openImageModal(`/upload/${slide.path}`);
		// 添加图片加载错误处理
		img.onerror = function() {
			console.error('Failed to load image:', this.src);
			this.style.display = 'none';
		};
		carouselContainer.appendChild(img);

		// 创建控制按钮
		const btn = document.createElement('button');
		btn.className = `carousel-control ${index === 0 ? 'active' : ''}`;
		btn.onclick = () => goToSlide(index);
		carouselControls.appendChild(btn);
	});

	// 设置初始位置
	updateCarouselPosition();
}

// 打开大图预览
function openImageModal(imageUrl) {
	const modal = document.getElementById('imageModal');
	const modalImg = document.getElementById('modalImage');
	modal.style.display = 'flex';
	modalImg.src = imageUrl;
}

// 关闭大图预览
function closeImageModal() {
	const modal = document.getElementById('imageModal');
	modal.style.display = 'none';
}

// 点击模态框外部关闭
window.addEventListener('click', function(e) {
	const modal = document.getElementById('imageModal');
	if (e.target === modal) {
		closeImageModal();
	}
});

// 更新轮播图位置
function updateCarouselPosition() {
	const carouselContainer = document.getElementById('carouselContainer');
	carouselContainer.style.transform = `translateX(-${currentSlide * 100}%)`;

	// 更新控制按钮状态
	const controls = document.querySelectorAll('.carousel-control');
	controls.forEach((control, index) => {
		control.className = `carousel-control ${index === currentSlide ? 'active' : ''}`;
	});
}

// 跳转到指定幻灯片
function goToSlide(index) {
	currentSlide = index;
	if (currentSlide < 0) {
		currentSlide = slides.length - 1;
	} else if (currentSlide >= slides.length) {
		currentSlide = 0;
	}
	updateCarouselPosition();
}

// 下一张幻灯片
function nextSlide() {
	goToSlide(currentSlide + 1);
}

// 上一张幻灯片
function prevSlide() {
	goToSlide(currentSlide - 1);
}

// 启动自动轮播
function startAutoCarousel() {
	setInterval(nextSlide, 5000); // 每5秒切换一次
}

// 页面加载完成后初始化轮播图
window.addEventListener('load', initCarousel);