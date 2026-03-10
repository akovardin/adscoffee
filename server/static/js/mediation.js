// Загружаем дополнительные скрипты
function loadScript(src) {
    const script = document.createElement('script');
    script.src = src;
    script.async = true;
    document.head.appendChild(script);
}

loadScript('https://yandex.ru/ads/system/context.js');
loadScript('https://platform.coffee/static/js/coffee.js');

// Инициализируем глобальные объекты и очереди колбэков
window.yaContextCb = window.yaContextCb || [];
window.adsCoffeeCb = window.adsCoffeeCb || [];

// Основная функция для работы с медиацией
window.adsCoffeeMediation = {
    render: function (options) {

        console.log("error");

        if (!options || !options.placement || !options.renderTo) {
            console.error('adsCoffeeMediation.render: required options (block, renderTo) are missing');
            return;
        }

        // fetch(`https://ads.coffee/mediation/${options.placement}`)
        fetch(`http://127.0.0.1:8090/native/${options.placement}`)
            .then(response => {

                console.log(response);

                if (!response.ok) throw new Error('Network response was not ok');
                return response.json();
            })
            .then(data => {

                console.log(data);

                // {
                //     "unit": "3qj9qddp1n4e08q",
                //     "network": "yandex",
                //     "data": {
                //         "block": "R-A-14476736-1"
                //     }
                // }

                if (data.network === 'yandex') {
                    window.yaContextCb.push(() => {
                        Ya.Context.AdvManager.render({
                            blockId: data.data.placement,
                            renderTo: options.renderTo
                        });
                    });
                } else if (data.network === 'coffee') {
                    window.adsCoffeeCb.push(() => {
                        window.adsCoffee.render({
                            template: "horizontal",
                            renderTo: options.renderTo,
                            blockId: data.data.placement
                        });
                    });
                }
            })
            .catch(error => {
                console.error('adsCoffeeMediation.render error:', error);
            });
    }
};

// Обрабатываем очередь колбэков, если они были добавлены до загрузки скрипта
if (window.adsCoffeeMediationCb) {
    window.adsCoffeeMediationCb.forEach(callback => callback());
    window.adsCoffeeMediationCb.push = function (c) {
        c();
        return this.length;
    };
}