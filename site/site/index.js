// Очередь колбэков для рекламного скрипта AdsCoffee
window.adsCoffeeCb = window.adsCoffeeCb || [];

// Инициализация рекламного блока
window.adsCoffeeCb.push(() => {
    window.adsCoffee.render({
        renderTo: "ads-cf-container",
        blockId: "35qleqc71q48d3h",
        template: "horizontal"
        // template: "horizontal-marker"
    });
});

// Yandex.Metrika counter
(function(m, e, t, r, i, k, a) {
    m[i] = m[i] || function() { (m[i].a = m[i].a || []).push(arguments) };
    m[i].l = 1 * new Date();
    for (var j = 0; j < document.scripts.length; j++) { if (document.scripts[j].src === r) { return; } }
    k = e.createElement(t), a = e.getElementsByTagName(t)[0], k.async = 1, k.src = r, a.parentNode.insertBefore(k, a);
})(window, document, "script", "https://mc.yandex.ru/metrika/tag.js", "ym");

ym(100191868, "init", {
    clickmap: true,
    trackLinks: true,
    accurateTrackBounce: true
});
