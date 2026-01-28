package main

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/presets/gorm2op"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/login"
	. "github.com/qor5/x/v3/ui/vuetify"
	. "github.com/theplant/htmlgo"
	"go.uber.org/fx"
	"golang.org/x/text/language"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/config"
	"go.ads.coffee/platform/admin/internal/database"
	"go.ads.coffee/platform/admin/internal/internat"
	"go.ads.coffee/platform/admin/internal/logger"
	"go.ads.coffee/platform/admin/internal/modules/ads"
	"go.ads.coffee/platform/admin/internal/modules/media"
	"go.ads.coffee/platform/admin/internal/modules/users"
	"go.ads.coffee/platform/admin/internal/s3storage"
	"go.ads.coffee/platform/admin/internal/server"
)

func main() {
	fx.New(
		fx.Provide(
			func() prometheus.Registerer {
				// default prometheus
				return prometheus.DefaultRegisterer
			},
		),
		config.Module,
		database.Module,
		logger.Module,
		s3storage.Module,
		server.Module,

		ads.Module,
		users.Module,
		media.Module,

		fx.Provide(
			configure,
			auth,
		),
		fx.Invoke(
			serve,
		),
	).Run()
}

func serve(lc fx.Lifecycle, srv *server.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return srv.Serve()
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}

func auth(users *users.Users, pb *presets.Builder) *login.Builder {
	return users.Auth(pb)
}

func configure(
	db *gorm.DB,
	media *media.Media,
	users *users.Users,
	ads *ads.Ads,
) *presets.Builder {
	b := presets.New()

	// Set up the project name, ORM and Homepage
	b.URIPrefix("/admin").
		// BrandTitle("AdEx").
		DataOperator(gorm2op.DataOperator(db)).
		HomePageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
			r.Body = VContainer(
				H1("Реклама"),
				P().Text("Лучшая DSP"))
			return
		})

	media.Configure(b)
	ads.Configure(b)
	users.Configure(b)

	b.MenuOrder(
		"advertisers",
		"campaigns",
		"bgroups",
		"banners",
		"separator",
		"media-library",
		"users",
		"separator",
		"networks",
	)

	i18nB := b.GetI18n()

	i18nB.SupportLanguages(language.Russian, language.English)
	i18nB.
		RegisterForModule(language.English, presets.ModelsI18nModuleKey, internat.Messages_en_EN_ModelsI18nModuleKey).
		RegisterForModule(language.Russian, presets.ModelsI18nModuleKey, internat.Messages_ru_RU_ModelsI18nModuleKey).
		RegisterForModule(language.English, presets.CoreI18nModuleKey, internat.Messages_en_EN).
		RegisterForModule(language.Russian, presets.CoreI18nModuleKey, internat.Messages_ru_RU).
		GetSupportLanguagesFromRequestFunc(func(r *http.Request) []language.Tag {
			return b.GetI18n().GetSupportLanguages()
		})

	return b
}
