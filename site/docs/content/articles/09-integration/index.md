---
weight: 9
bookFlatSection: false
title: "Интеграция рекламы"
---

# 09 Интеграция рекламы
<p>&nbsp;</p>

Рекламный сервер отдаёт рекламу по нескольким маршрутам. Полный список пайплайнов и их параметров описан в `dist/server/configs/config.yaml`.

| Маршрут | Назначение |
|---|---|
| `/banner/{placement}/{action}` | Статический баннер: `img` — картинка, `click` — переход |
| `/native/{placement}` | Нативный формат (JSON) |
| `/dsp/{placement}` | RTB/DSP-запрос |
| `/inapp/{placement}` | Формат для мобильных приложений |
| `/tracker/{data}.gif` | Трекинг-пиксель |
| `/postback` | Приём постбеков от сетей |

Во всех примерах ниже `https://ads.example.com` — публичный адрес вашего рекламного сервера (порт `8071`), а `{placement}` — **числовой ID плейсмента** из колонки `ID` в админке.

> В `dist/server/configs/config.yaml` для output-плагинов `web` и `inapp` задаётся параметр `base` — публичный адрес сервера, который подставляется в ссылки трекинга. Укажите там свой домен.

> Если сайт плейсмента выключен, сервер отвечает `200 OK` с пустым телом — это значит «рекламы нет». Клиент должен трактовать такой ответ как отсутствие рекламы (noad) и ничего не показывать.

## Баннер

Баннер добавляется как обычный HTML-элемент: картинка формируется на лету и оборачивается в ссылку перехода.

```html
<a href="https://ads.example.com/banner/4/click">
  <img src="https://ads.example.com/banner/4/img" width="350">
</a>
```

Размер тега `img` зависит от размеров картинок в рекламных кампаниях. Стилизовать элемент можно как угодно — главное, чтобы картинка была внутри ссылки. На баннер автоматически наносятся описание и маркировка рекламодателя.

## Нативный формат

JSON доступен по адресу `https://ads.example.com/native/4`. Ответ — массив объектов:

```json
[
  {
    "description": "Веб-разработка на Go",
    "information": "Реклама. Рекламодатель ООО «Пример», ИНН 0000000000. erid 2SDnjeP817J",
    "image": "https://ads.example.com/files/course.png",
    "target": "https://advertiser.example/landing",
    "impressions": ["https://ads.example.com/tracker/....gif?revenue={revenue}"],
    "clicks": ["https://ads.example.com/tracker/....gif?revenue={revenue}"],
    "data": "",
    "network": "coffee",
    "format": "banner",
    "revenue": 1.5
  }
]
```

Поля:

- `description` — текст объявления;
- `information` — маркировка рекламы (для display-креативов);
- `image` — ссылка на изображение;
- `target` — целевая ссылка для перехода;
- `impressions` — трекинг-пиксели показа;
- `clicks` — трекинг-пиксели клика;
- `data` — произвольные данные сети (например, JSON с `block` для Yandex);
- `network` — источник рекламы (`coffee`, `yandex` и т. д.);
- `format` — `banner` или `interstitial`;
- `revenue` — доход за показ (для сетей, подключённых через медиацию), опционально.

Для учёта статистики: вызовите все ссылки из `impressions` при показе и из `clicks` при клике. В трекер-ссылках есть макрос `{revenue}` — подставьте в него доход показа (для собственных креативов он не обязателен). Само объявление ведёт на `target`.

## JS-виджет

Для сайтов есть готовый скрипт `server/static/js/coffee.js`. Он сам загружает объявление, рендерит его в выбранный контейнер и отправляет трекинг:

```html
<div id="ads-container"></div>
<script src="https://ads.example.com/static/js/coffee.js" async></script>
<script>
  window.adsCoffeeCb = window.adsCoffeeCb || [];
  window.adsCoffeeCb.push(function () {
    window.adsCoffee.render({
      renderTo: "ads-container",
      placementId: "4",
      template: "horizontal" // "horizontal" | "horizontal-marker" | "default"
    });
  });
</script>
```

Виджет умеет показывать собственные (`coffee`) креативы и блоки Yandex: если сервер вернул `network: "yandex"`, виджет подключает Yandex Ads SDK и использует `data.block` как идентификатор блока. При отключённом сайте (пустой ответ) контейнер остаётся пустым.

> Скрипт запрашивает объявления с адреса `https://platform.ads.coffee`. Для самостоятельного хостинга измените этот URL в `server/static/js/coffee.js` на адрес вашего сервера.

## Мобильные приложения

Для Android есть SDK в каталоге `sdk/android/` (модуль `coffeesdk`). Он умеет загружать объявление, показывать баннер в `InAppAdView`, интерстишел в полноэкранном диалоге и отправлять события показа и клика. Формат для приложений отдаётся сервером по маршруту `/inapp/{placement}`.

```kotlin
CoffeeAds.setBaseUrl("https://ads.example.com")
CoffeeAds.initialize(context) {
    val loader = InAppAdLoader(context)
    val request = InAppAdRequest.Builder(PLACEMENT_ID).build()

    loader.loadAd(request, object : InAppAdLoadListener {
        override fun onAdLoaded(ad: InAppAd) {
            ad.setInAppAdEventListener(object : InAppAdEventListener {
                override fun onAdClicked() { ad.target?.let { openTarget(it) } }
                override fun onImpression(data: ImpressionData?) {}
            })

            if (ad.format?.equals("interstitial", ignoreCase = true) == true) {
                InAppInterstitial.show(context, ad) { /* закрыт */ }
            } else {
                val binder = InAppAdViewBinder.Builder(adView).build()
                ad.bindInAppAd(binder)
                ad.fireImpressionTrackers()
            }
        }

        override fun onAdFailedToLoad(error: AdRequestError) {
            // Рекламы нет (noad) — ничего не показываем
        }
    })
}
```

Особенности:

- **noad.** Пустой ответ (`200 OK` без объявлений, в том числе для выключенного сайта) приводит к `onAdFailedToLoad` — ничего не рисуйте и не показывайте интерстишел.
- **Интерстишел.** Клик в любом месте, кроме крестика, открывает `target`; трекинг клика отправляется автоматически.
- **Yandex.** SDK поддерживает креативы сетей, включая Yandex Ads через медиацию (идентификатор блока передаётся в поле `data`).

## Постбеки

Маршрут `/postback` принимает уведомления о конверсиях от партнёрских сетей. Если сеть поддерживает постбеки, укажите в её кабинете URL вашего сервера — так конверсии будут попадать в общую аналитику.
