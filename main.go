package main

import (
	"flag"

	"encoding/json"

	astilectron "github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	astilog "github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

const htmlAbout = `Welcome !!`

var (
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", true, "Enables the debug mode")
	w       *astilectron.Window
)

func main() {
	flag.Parse()
	astilog.FlagInit()

	astilog.Debugf("Running app built at %s", BuiltAt)

	var bootstrapOptions = bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug: *debug,
		MenuOptions: []*astilectron.MenuItemOptions{
			{
				Label: astilectron.PtrStr("File"),
				SubMenu: []*astilectron.MenuItemOptions{
					{
						Label: astilectron.PtrStr("About"),
						OnClick: func(e astilectron.Event) (deleteListener bool) {
							if err := bootstrap.SendMessage(w, "about", htmlAbout, func(m *bootstrap.MessageIn) {
								var s string
								if err := json.Unmarshal(m.Payload, &s); err != nil {
									astilog.Error(errors.Wrap(err, "unmarshaling payload failed"))
									return
								}
								astilog.Infof("About modal has been displayed and payload is %s!", s)
							}); err != nil {
								astilog.Error(errors.Wrap(err, "Sending about Event failed"))
							}
							return
						},
					},
					{Role: astilectron.MenuItemRoleClose},
				},
			},
		},
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]

			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{
			{
				Homepage:       "index.html",
				MessageHandler: handleMessages,
				Options: &astilectron.WindowOptions{
					BackgroundColor: astilectron.PtrStr("#FFF"),
					Center:          astilectron.PtrBool(true),
					Height:          astilectron.PtrInt(700),
					Width:           astilectron.PtrInt(1050),
				},
			},
		},
	}

	if err := bootstrap.Run(bootstrapOptions); err != nil {
		astilog.Fatal(errors.Wrap(err, "Running bootstrap failed"))
	}
}
