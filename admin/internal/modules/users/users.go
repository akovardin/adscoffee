package users

import (
	plogin "github.com/qor5/admin/v3/login"
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/login"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/modules/users/models"
)

var loginSecret = "edfkdehfkdfhu"

type Users struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Users {
	return &Users{
		db: db,
	}
}

func (u *Users) Auth(pb *presets.Builder) *login.Builder {
	lb := plogin.New(pb).
		DB(u.db).
		UserModel(&models.User{}).
		Secret(loginSecret).
		TOTP(false)

	pb.ProfileFunc(func(ctx *web.EventContext) h.HTMLComponent {
		return h.A(h.Text("Выход")).Href(lb.LogoutURL).Style("margin-left:100px")
	})

	return lb
}

func (u *Users) Configure(b *presets.Builder) {
	m := b.Model(&models.User{}).
		MenuIcon("mdi-account-multiple")
	// Label("Пользователи")

	m.Editing("Name", "Account", "Password").Field("Password").
		SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
			u := obj.(*models.User)
			if v := ctx.R.FormValue(field.Name); v != "" {
				u.Password = v
				u.EncryptPassword()
			}
			return nil
		})
}

func (u *Users) Migrate() {
	err := u.db.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		panic(err)
	}
}
