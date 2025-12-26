package scripts

import (
	log "Sowhp/concert/logger"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ReportGenerator struct {
	resultName string
	resultDir  string
}

func NewReportGenerator(resultName string) *ReportGenerator {
	return &ReportGenerator{
		resultName: resultName,
		resultDir:  "./result",
	}
}

func CreateHtml(resultMap map[string]map[string][]string) error {
	if len(resultMap) == 0 {
		return fmt.Errorf("结果数据为空，无法生成报告")
	}

	var resultName string
	for name := range resultMap {
		resultName = name
		break
	}

	generator := NewReportGenerator(resultName)
	return generator.generateReports(resultMap[resultName])
}

func (rg *ReportGenerator) generateReports(data map[string][]string) error {
	if err := rg.generateTextReport(data); err != nil {
		log.Error(fmt.Sprintf("生成文本报告失败: %v", err))
		return err
	}

	if err := rg.generateHTMLReport(data); err != nil {
		log.Error(fmt.Sprintf("生成HTML报告失败: %v", err))
		return err
	}

	return nil
}

func (rg *ReportGenerator) generateTextReport(data map[string][]string) error {
	csvPath := filepath.Join(rg.resultDir, rg.resultName+".csv")

	file, err := os.OpenFile(csvPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("创建CSV报告文件失败: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Warning(fmt.Sprintf("关闭CSV报告文件失败: %v", closeErr))
		}
	}()

	header := "Website URL Address,Title Name,Status,Screenshot Path\n"
	if _, err := file.WriteString(header); err != nil {
		return fmt.Errorf("写入CSV报告表头失败: %w", err)
	}

	for url, info := range data {
		if len(info) < 3 {
			log.Warning(fmt.Sprintf("URL %s 的信息不完整，跳过", url))
			continue
		}

		titleName := info[0]
		status := info[1]
		screenshotPath := info[2]

		line := fmt.Sprintf("%s,%s,%s,%s\n",
			url,
			titleName,
			status,
			screenshotPath)

		if _, err := file.WriteString(line); err != nil {
			return fmt.Errorf("写入数据行失败: %w", err)
		}
	}

	return nil
}

func (rg *ReportGenerator) generateHTMLReport(data map[string][]string) error {
	htmlPath := filepath.Join(rg.resultDir, rg.resultName+".html")

	// 构建HTML内容
	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Sowhp 截图报告 - %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 10px; background-color: #f5f5f5; }
        .container { max-width: 1800px; margin: 0 auto; background-color: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #333; text-align: center; margin-bottom: 30px; }
        table { width: 100%%; border-collapse: collapse; margin-top: 20px; table-layout: fixed; }
        th, td { padding: 8px; text-align: left; border-bottom: 1px solid #ddd; vertical-align: top; word-wrap: break-word; }
        th { background-color: #4CAF50; color: white; font-weight: bold; font-size: 14px; }
        td { font-size: 12px; }
        th:nth-child(1), td:nth-child(1) { width: 20%%; }
        th:nth-child(2), td:nth-child(2) { width: 15%%; }
        th:nth-child(3), td:nth-child(3) { width: 10%%; }
        th:nth-child(4), td:nth-child(4) { width: 15%%; }
        th:nth-child(5), td:nth-child(5) { width: 40%%; }
        tr:hover { background-color: #f5f5f5; }
        .url-link { color: #1976D2; text-decoration: none; word-break: break-all; }
        .url-link:hover { text-decoration: underline; }
        .screenshot { width: 100%%; height: auto; border: 1px solid #ddd; border-radius: 4px; cursor: pointer; transition: transform 0.2s; display: block; }
        .screenshot:hover { transform: scale(1.05); }
        .status-success { color: #4CAF50; font-weight: bold; }
        .status-error { color: #f44336; font-weight: bold; }
        .status-timeout { color: #ff9800; font-weight: bold; }
        .status-dns { color: #9c27b0; font-weight: bold; }
        .status-ssl { color: #795548; font-weight: bold; }
        .summary { background-color: #e3f2fd; padding: 15px; border-radius: 4px; margin-bottom: 20px; }
        .pagination { text-align: center; margin: 20px 0; }
        .pagination button { margin: 0 5px; padding: 8px 12px; border: 1px solid #ddd; background: white; cursor: pointer; border-radius: 4px; }
        .pagination button:hover { background: #f5f5f5; }
        .pagination button.active { background: #4CAF50; color: white; border-color: #4CAF50; }
        .pagination button:disabled { background: #f5f5f5; color: #999; cursor: not-allowed; }
        .response-content { max-width: 100%%; max-height: 300px; overflow: auto; font-family: 'Courier New', monospace; font-size: 10px; line-height: 1.2; background: #f8f8f8; padding: 6px; border-radius: 4px; white-space: pre-wrap; border: 1px solid #ddd; word-break: break-all; }
        .modal { display: none; position: fixed; z-index: 1000; left: 0; top: 0; width: 100%%; height: 100%%; background-color: rgba(0,0,0,0.9); backdrop-filter: blur(2px); overflow-y: auto; }
        .modal-content { position: absolute; top: 50%%; left: 50%%; transform: translate(-50%%, -50%%); max-width: 98%%; display: flex; align-items: center; justify-content: center; }
        .modal img { max-width: 100%%; max-height: 100%%; border-radius: 8px; box-shadow: 0 4px 20px rgba(0,0,0,0.5); transition: transform 0.3s ease; cursor: zoom-in; }
        .modal img.zoomed { cursor: zoom-out; }
        .loading-spinner { position: absolute; top: 50%%; left: 50%%; transform: translate(-50%%, -50%%); color: #fff; font-size: 18px; }
        .zoom-controls { position: fixed; bottom: 30px; left: 50%%; transform: translateX(-50%%); z-index: 1001; display: flex; gap: 10px; }
        .zoom-btn { background: rgba(0,0,0,0.7); color: white; border: none; padding: 10px 15px; border-radius: 25px; cursor: pointer; font-size: 16px; transition: all 0.3s ease; }
        .zoom-btn:hover { background: rgba(255,255,255,0.2); transform: scale(1.05); }
    </style>
</head>
<body>
    <div class="container">
        <h1>网站截图报告 - %s</h1>`, rg.resultName, rg.resultName)

	// 计算统计信息
	totalCount := len(data)
	successCount := 0
	for _, info := range data {
		if len(info) >= 2 && info[1] != "连接失败" && info[1] != "ERROR" {
			successCount++
		}
	}

	htmlContent += fmt.Sprintf(`
        <div class="summary">
            <p>总计: %d 个地址，成功: %d 个，失败: %d 个</p>
        </div>
        <div class="pagination" id="pagination"></div>
        <table id="dataTable">
            <thead>
                <tr>
                    <th>URL地址</th>
                    <th>网站标题</th>
                    <th>状态码</th>
                    <th>截图</th>
                    <th>响应内容</th>
                </tr>
            </thead>
            <tbody id="tableBody"></tbody>
        </table>
        <div id="imageModal" class="modal">
            <div class="modal-content">
                <div class="loading-spinner" id="loadingSpinner">加载中...</div>
                <img id="modalImage" style="display: none;">
            </div>
            <div class="zoom-controls" id="zoomControls" style="display: none;">
                <button class="zoom-btn" onclick="zoomImage(-0.2)">-</button>
                <button class="zoom-btn" onclick="resetZoom()">重置</button>
                <button class="zoom-btn" onclick="zoomImage(0.2)">+</button>
            </div>
        </div>
    </div>
    <script>
        // 使用安全的数据传递方式
        window.reportData = {`, totalCount, successCount, totalCount-successCount)

	// 构建数据结构
	type ReportItem struct {
		URL        string `json:"url"`
		Title      string `json:"title"`
		Status     string `json:"status"`
		Screenshot string `json:"screenshot"`
		Response   string `json:"response"`
	}

	var items []ReportItem
	for url, info := range data {
		if len(info) < 4 {
			continue
		}

		titleName := info[0]
		status := info[1]
		screenshotPath := info[2]
		responseContent := ""
		if len(info) >= 4 {
			responseContent = info[3]
		}

		if screenshotPath != "data/" && screenshotPath != "" {
			screenshotPath = fmt.Sprintf("%s/%s", rg.resultName, screenshotPath)
		}

		items = append(items, ReportItem{
			URL:        url,
			Title:      titleName,
			Status:     status,
			Screenshot: screenshotPath,
			Response:   responseContent,
		})
	}

	// 使用JSON编码确保数据安全
	jsonData, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("JSON编码失败: %w", err)
	}

	htmlContent += fmt.Sprintf("\n            items: %s\n        };", string(jsonData))

	htmlContent += `

         const itemsPerPage = 20;
         let currentPage = 1;
         let totalPages = Math.ceil(window.reportData.items.length / itemsPerPage);

        function renderTable(page) {
            try {
                const tbody = document.getElementById('tableBody');
                if (!tbody) {
                    console.error('tableBody element not found');
                    return;
                }
                
                tbody.innerHTML = '';
                
                const start = (page - 1) * itemsPerPage;
                const end = start + itemsPerPage;
                const pageData = window.reportData.items.slice(start, end);
                
                pageData.forEach(function(item) {
                    const row = tbody.insertRow();
                    
                    const urlCell = row.insertCell();
                    const link = document.createElement('a');
                    link.href = item.url;
                    link.target = '_blank';
                    link.className = 'url-link';
                    link.textContent = truncateString(item.url, 50);
                    urlCell.appendChild(link);
                    
                    const titleCell = row.insertCell();
                    titleCell.textContent = item.title;
                    
                    const statusCell = row.insertCell();
                    let statusClass = 'status-success';
                    if (item.status === 'TIMEOUT') {
                        statusClass = 'status-timeout';
                    } else if (item.status === 'DNS_ERROR') {
                        statusClass = 'status-dns';
                    } else if (item.status === 'SSL_ERROR') {
                        statusClass = 'status-ssl';
                    } else if (item.status === 'CONNECTION_REFUSED' || item.status === 'ERROR' || item.status === '连接失败') {
                        statusClass = 'status-error';
                    } else if (item.status.indexOf('4') === 0 || item.status.indexOf('5') === 0) {
                        statusClass = 'status-error';
                    }
                    const statusSpan = document.createElement('span');
                    statusSpan.className = statusClass;
                    statusSpan.textContent = item.status;
                    statusCell.appendChild(statusSpan);
                    
                    const screenshotCell = row.insertCell();
                    if (item.screenshot && item.screenshot !== 'data/') {
                        const img = document.createElement('img');
                        img.src = item.screenshot;
                        img.className = 'screenshot';
                        img.alt = '网站截图';
                        img.onclick = function() { openModal(this.src); };
                        screenshotCell.appendChild(img);
                    } else {
                        screenshotCell.textContent = '无截图';
                    }
                    
                    const responseCell = row.insertCell();
                    const responseDiv = document.createElement('div');
                    responseDiv.className = 'response-content';
                    responseDiv.textContent = item.response;
                    responseCell.appendChild(responseDiv);
                });
            } catch (e) {
                console.error('Error rendering table:', e);
            }
        }

        function renderPagination() {
            try {
                const pagination = document.getElementById('pagination');
                if (!pagination) return;
                
                pagination.innerHTML = '';
                
                const prevBtn = document.createElement('button');
                prevBtn.textContent = '上一页';
                prevBtn.disabled = currentPage === 1;
                prevBtn.onclick = function() {
                    if (currentPage > 1) {
                        currentPage--;
                        renderTable(currentPage);
                        renderPagination();
                    }
                };
                pagination.appendChild(prevBtn);
                
                const startPage = Math.max(1, currentPage - 2);
                const endPage = Math.min(totalPages, currentPage + 2);
                
                for (let i = startPage; i <= endPage; i++) {
                    const pageBtn = document.createElement('button');
                    pageBtn.textContent = i;
                    pageBtn.className = i === currentPage ? 'active' : '';
                    pageBtn.onclick = function() {
                        currentPage = i;
                        renderTable(currentPage);
                        renderPagination();
                    };
                    pagination.appendChild(pageBtn);
                }
                
                const nextBtn = document.createElement('button');
                nextBtn.textContent = '下一页';
                nextBtn.disabled = currentPage === totalPages;
                nextBtn.onclick = function() {
                    if (currentPage < totalPages) {
                        currentPage++;
                        renderTable(currentPage);
                        renderPagination();
                    }
                };
                pagination.appendChild(nextBtn);
            } catch (e) {
                console.error('Error rendering pagination:', e);
            }
        }

        function truncateString(str, maxLen) {
            if (str.length <= maxLen) {
                return str;
            }
            return str.substring(0, maxLen - 3) + '...';
        }

        let currentZoom = 1;
        let isDragging = false;
        let startX, startY, translateX = 0, translateY = 0;

        function openModal(imageSrc) {
            try {
                const modal = document.getElementById('imageModal');
                const modalImg = document.getElementById('modalImage');
                const loadingSpinner = document.getElementById('loadingSpinner');
                const zoomControls = document.getElementById('zoomControls');
                
                if (modal && modalImg && loadingSpinner && zoomControls) {
                    modal.style.display = 'block';
                    loadingSpinner.style.display = 'block';
                    modalImg.style.display = 'none';
                    zoomControls.style.display = 'none';
                    
                    // 重置缩放和位置
                    currentZoom = 1.5;
                    translateX = 0;
                    translateY = 0;
                    modalImg.style.transform = 'scale(1.5) translate(0px, 0px)';
                    modalImg.classList.add('zoomed');
                    
                    modalImg.onload = function() {
                        loadingSpinner.style.display = 'none';
                        modalImg.style.display = 'block';
                        zoomControls.style.display = 'flex';
                    };
                    
                    modalImg.onerror = function() {
                        loadingSpinner.innerHTML = '图片加载失败';
                        setTimeout(() => {
                            modal.style.display = 'none';
                        }, 2000);
                    };
                    
                    modalImg.src = imageSrc;
                }
            } catch (e) {
                console.error('Error opening modal:', e);
            }
        }

        function zoomImage(delta) {
            const modalImg = document.getElementById('modalImage');
            if (!modalImg) return;
            
            currentZoom = Math.max(0.5, Math.min(5, currentZoom + delta));
            updateImageTransform();
            
            if (currentZoom > 1) {
                modalImg.classList.add('zoomed');
            } else {
                modalImg.classList.remove('zoomed');
            }
        }

        function resetZoom() {
            currentZoom = 1.5;
            translateX = 0;
            translateY = 0;
            updateImageTransform();
            const modalImg = document.getElementById('modalImage');
            if (modalImg) {
                modalImg.classList.add('zoomed');
            }
        }

        function updateImageTransform() {
            const modalImg = document.getElementById('modalImage');
            if (modalImg) {
                modalImg.style.transform = 'scale(' + currentZoom + ') translate(' + translateX + 'px, ' + translateY + 'px)';
            }
        }
        function initializeReport() {
            try {
                console.log('Initializing report, data items:', window.reportData.items.length);
                renderTable(1);
                renderPagination();
                
                const modal = document.getElementById('imageModal');
                const modalImg = document.getElementById('modalImage');
                
                window.onclick = function(event) {
                    if (event.target === modal) {
                        modal.style.display = 'none';
                        resetZoom();
                    }
                };
                
                // 添加滚轮缩放功能
                if (modalImg) {
                    modalImg.addEventListener('wheel', function(e) {
                        e.preventDefault();
                        const delta = e.deltaY > 0 ? -0.1 : 0.1;
                        zoomImage(delta);
                    });
                    
                    // 添加拖拽功能
                    modalImg.addEventListener('mousedown', function(e) {
                        if (currentZoom > 1) {
                            isDragging = true;
                            startX = e.clientX - translateX;
                            startY = e.clientY - translateY;
                            modalImg.style.cursor = 'grabbing';
                        }
                    });
                    
                    document.addEventListener('mousemove', function(e) {
                        if (isDragging && currentZoom > 1) {
                            translateX = e.clientX - startX;
                            translateY = e.clientY - startY;
                            updateImageTransform();
                        }
                    });
                    
                    document.addEventListener('mouseup', function() {
                        isDragging = false;
                        if (modalImg && currentZoom > 1) {
                            modalImg.style.cursor = 'grab';
                        }
                    });
                    
                    // 双击重置缩放
                    modalImg.addEventListener('dblclick', function() {
                        resetZoom();
                    });
                }
            } catch (e) {
                console.error('Error initializing report:', e);
            }
        }

        // 多重初始化确保页面正常加载
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', initializeReport);
        } else {
            initializeReport();
        }
        
        // 备用初始化
        setTimeout(function() {
            if (document.getElementById('tableBody').children.length === 0) {
                console.log('Fallback initialization');
                initializeReport();
            }
        }, 100);
    </script>
</body>
</html>`

	// 写入文件
	file, err := os.OpenFile(htmlPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("创建报告文件失败: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Warning(fmt.Sprintf("关闭报告文件失败: %v", closeErr))
		}
	}()

	if _, err := file.WriteString(htmlContent); err != nil {
		return fmt.Errorf("写入报告失败: %w", err)
	}

	log.Info(fmt.Sprintf("生成报告成功: %s", htmlPath))
	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
