// Создаем глобальный объект и очередь колбэков
window.adsCoffeeCb = window.adsCoffeeCb || [];

// Основная функция рендеринга
window.adsCoffee = {
  render: function(params) {
    const { 
      renderTo, 
      placementId, 
      template = 'default',
      contentWidth = template === 'default' ? 220 : 380
    } = params;
    
    if (!renderTo || !placementId) {
      console.error('Не указаны обязательные параметры renderTo или placementId');
      return;
    }

    const container = document.getElementById(renderTo);
    if (!container) {
      console.error(`Элемент с id="${renderTo}" не найден`);
      return;
    }

    // get data from params

    // Загружаем данные
    fetch(`https://platform.ads.coffee/native/${encodeURIComponent(placementId)}`)
    // fetch(`http://localhost:8081/native/${encodeURIComponent(placementId)}`)
      .then(response => {
        if (!response.ok) throw new Error('Ошибка загрузки данных');
        return response.json();
      })
      .then(data => {
        if (data.length == 0) {
          return;
        }

        data = data[0];
        // Создаем стили
        const styleId = 'ads-coffee-styles';
        if (!document.getElementById(styleId)) {
          const style = document.createElement('style');
          style.id = styleId;
          style.textContent = `
            /* Базовые стили */
            .ads-coffee {
              font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Roboto, 
                         Helvetica Neue, Arial, sans-serif;
              box-sizing: border-box;
            }
            
            /* Стили для тултипа */
            .ads-coffee-tooltip {
              position: relative;
              display: inline-block;
              cursor: pointer;
              text-decoration: none;
              color: #666;
            }
            .ads-coffee-tooltip:hover::after {
              content: attr(data-tooltip);
              position: absolute;
              left: 50%;
              transform: translateX(-50%);
              bottom: 100%;
              background: #333;
              color: #fff;
              padding: 5px 10px;
              border-radius: 4px;
              font-size: 14px;
              white-space: normal;
              width: 200px;
              text-align: center;
              margin-bottom: 10px;
              z-index: 100;
            }
            .ads-coffee-tooltip:hover::before {
              content: "";
              position: absolute;
              left: 50%;
              transform: translateX(-50%);
              bottom: calc(100% - 5px);
              border: 5px solid transparent;
              border-top-color: #333;
            }
            
            /* Общие стили контента */
            .ads-coffee-content {
              background: rgba(0, 0, 0, .03);
              color: #505050;
              padding: 1em;
              overflow: hidden;
              text-align: center;
              display: flex;
              flex-direction: column;
              justify-content: center;
              align-items: center;
              border-radius: 4px;
              box-sizing: border-box;
            }
            .ads-coffee-content a {
              text-decoration: none;
              color: #505050;
            }
            .ads-coffee-content a:hover {
              text-decoration: underline;
            }
            .ads-coffee-content img {
              display: block;
              height: 140px;
              width: 180px;
              object-fit: contain;
              margin-bottom: 10px;
            }
            .ads-coffee-text {
              font-size: 14px;
              line-height: 1.4;
              width: 100%;
            }
            
            /* Горизонтальные варианты */
            .ads-coffee-horizontal .ads-coffee-content {
              flex-direction: row;
              flex-wrap: wrap;
              align-items: center;
              justify-content: flex-start;
              text-align: left;
              padding: 15px;
            }
            .ads-coffee-horizontal .ads-coffee-img {
              margin-right: 15px;
              flex-shrink: 0;
            }
            .ads-coffee-horizontal .ads-coffee-text {
              align-self: center;
              margin: 0;
              flex: 1;
              text-align: left;
              font-size: 15px;
            }
            .ads-coffee-horizontal .ads-coffee-content img {
              margin-bottom: 0;
              width: 140px;
              height: auto;
            }
            
            /* Специальные стили для horizontal-marker */
            .ads-coffee-horizontal-marker .ads-coffee-content {
              padding-bottom: 25px;
            }
            .ads-coffee-marker {
              font-size: 12px;
              color: #999;
              width: 100%;
              margin-top: 15px;
              padding-top: 10px;
              border-top: 1px solid #eee;
              line-height: 1.3;
              text-align: left;
            }
            
            /* Подпись */
            .ads-coffee-callout {
              font-style: italic;
              margin: 0 1em 1em;
              padding: 0 1em;
              text-align: right;
              font-size: 12px;
              color: #999;
            }
            
            /* Ошибка */
            .ads-coffee-error {
              color: #999;
              font-size: 12px;
              padding: 10px;
              text-align: center;
            }
          `;
          document.head.appendChild(style);
        }

        // Генерируем HTML в зависимости от шаблона
        let html;
        switch(template) {
          case 'horizontal':
            html = `
              <div class="ads-coffee ads-coffee-horizontal" style="max-width: ${contentWidth}px;">
                <div class="ads-coffee-content" style="max-width: ${contentWidth}px;">
                  <div class="ads-coffee-img">
                    <a href="${data.target}" rel="nofollow noopener" target="_blank" class="ads-coffee-click">
                      <img src="${data.image}" loading="lazy" alt="${data.description}">
                    </a>
                  </div>
                  <div class="ads-coffee-text">
                    <a href="${data.target}" rel="nofollow noopener" target="_blank" class="ads-coffee-click">
                      ${data.description}
                    </a>
                  </div>
                </div>
                <div class="ads-coffee-callout" style="max-width: ${contentWidth}px;">
                  <a rel="nofollow noopener" href="${data.target}" 
                     class="ads-coffee-tooltip ads-coffee-click" 
                     data-tooltip="${data.description}">Реклама</a>
                </div>
              </div>
            `;
            break;
            
          case 'horizontal-marker':
            html = `
              <div class="ads-coffee ads-coffee-horizontal ads-coffee-horizontal-marker" style="max-width: ${contentWidth}px;">
                <div class="ads-coffee-content" style="max-width: ${contentWidth}px">
                  <div class="ads-coffee-img">
                    <a href="${data.target}" rel="nofollow noopener" target="_blank" class="ads-coffee-click">
                      <img src="${data.image}" loading="lazy" alt="${data.description}">
                    </a>
                  </div>
                  <div class="ads-coffee-text">
                    <a href="${data.target}" rel="nofollow noopener" target="_blank" class="ads-coffee-click">
                      ${data.description}
                    </a>
                  </div>
                  <div class="ads-coffee-marker">
                    ${data.description}
                  </div>
                </div>
              </div>
            `;
            break;
            
          default:
            html = `
              <div class="ads-coffee" style="max-width: ${contentWidth + 60}px;">
                <div class="ads-coffee-content" style="max-width: ${contentWidth}px">
                  <div class="ads-coffee-img">
                    <a href="${data.target}" rel="nofollow noopener" target="_blank" class="ads-coffee-click">
                      <img src="${data.image}" loading="lazy" alt="${data.description}">
                    </a>
                  </div>
                  <div class="ads-coffee-text">
                    <a href="${data.target}" rel="nofollow noopener" target="_blank" class="ads-coffee-click">
                      ${data.description}
                    </a>
                  </div>
                </div>
                <div class="ads-coffee-callout" style="max-width: ${contentWidth}px">
                  <a rel="nofollow noopener" href="${data.target}" 
                     class="ads-coffee-tooltip ads-coffee-click" 
                     data-tooltip="${data.information}">Реклама</a>
                </div>
              </div>
            `;
        }

        container.innerHTML = html;

        // Трекинг показа
        if (data.impression) {
          const img = new Image();
          img.src = data.impression;
        }

        // Трекинг кликов
        container.querySelectorAll('.ads-coffee-click').forEach(el => {
          el.addEventListener('click', (e) => {
            if (data.click) {
              const img = new Image();
              img.src = data.click;
            }
          });
        });
      })
      .catch(error => {
        console.error('AdsCoffee error:', error);
        container.innerHTML = '<div class="ads-coffee-error">Рекламный блок</div>';
      });
  }
};

// Обрабатываем очередь
if (window.adsCoffeeCb) {
  window.adsCoffeeCb.forEach(cb => cb());
  window.adsCoffeeCb.push = function(c) {
    c();
    return this.length;
  };
}