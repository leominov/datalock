package handlers

import "html/template"

const (
	SecuredPlayerPageCode = `
    <!DOCTYPE html>
    <html>
        <head>
            <title>{{.Meta.Title}}</title>
            <meta name="keywords" content="{{.Meta.Keywords}}">
            <meta name="description" content="{{.Meta.Description}}">
            <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
            <link rel="stylesheet" type="text/css" href="/tpl/asset/css/main.css?17.4.4">
            <link rel="stylesheet" type="text/css" href="/tpl/asset/css/pg.season.css?17.4.4">
            <style type="text/css">
                body {
                    margin: 0;
                    background-color: #000;
                }
                .pgs-player {
                    display: block;
                    background: #000;
                    border: none;
                    height: 100%;
                    width: 98%;
                }
                .pgs-player-block {
                    background-image: none;
                }
            </style>
        </head>
        <body>
            <div>
                <div class="pgs-player" data-id-season="{{.Meta.ID}}" data-id-serial="{{.Meta.Serial}}">
                    <script type="text/javascript">
                        var mark = {
                            'href': '',
                            'trans': ''
                        }
                    </script>
                    <ul class="pgs-mark_line" data-player="marks" data-info="count">
                        <li class="pgs-mark_line-cur"></li>
                        <li class="pgs-mark_line-switch act" data-click="playerSwich"><span class="html svico-html5">HTML5</span><span class="flash svico-flash">Flash</span></li>
                    </ul>
                    <script type="text/javascript">
                        var data4play = {
                            'secureMark': '{{.User.SecureMark}}',
                            'time': 1494234314
                        }
                    </script>
                    <div id="player_wrap" data-player="wrap">
                        <div data-player="inside" class="pgs-player-inside">
                            <div class="svload"></div>
                        </div>
                        <div id="errorPlayer" class="pgs-player-block error">Ошибка при загрузке плеера, попробуйте перезагрузить страницу</div>
                    </div>
                </div>
            </div>
            <!-- Yandex.Metrika counter -->
            <script type="text/javascript">
                (function (d, w, c) {
                    (w[c] = w[c] || []).push(function() {
                        try {
                            w.yaCounter41725489 = new Ya.Metrika({
                                id:41725489,
                                clickmap:true,
                                trackLinks:true,
                                accurateTrackBounce:true
                            });
                        } catch(e) { }
                    });

                    var n = d.getElementsByTagName("script")[0],
                        s = d.createElement("script"),
                        f = function () { n.parentNode.insertBefore(s, n); };
                    s.type = "text/javascript";
                    s.async = true;
                    s.src = "https://mc.yandex.ru/metrika/watch.js";

                    if (w.opera == "[object Opera]") {
                        d.addEventListener("DOMContentLoaded", f, false);
                    } else { f(); }
                })(document, window, "yandex_metrika_callbacks");
            </script>
            <noscript>
                <div>
                    <img src="https://mc.yandex.ru/watch/41725489" style="position:absolute; left:-9999px;" alt="" />
                </div>
            </noscript>
            <!-- /Yandex.Metrika counter -->
            <script src="/tpl/asset/vendor/jquery.min.js"></script>
            <script src="/tpl/asset/vendor/js.cookie.min.js"></script>
            <script src="/tpl/asset/vendor/jquery.tooltipster.min.js"></script>
            <script src="/tpl/asset/js/main.min.js?17.4.4"></script>
            <script src="/tpl/asset/vendor/swfobject.min.js"></script>
            <script src="/tpl/asset/js/pg.marks.min.js?17.4.4"></script>
            <script src="/tpl/asset/js/pg.player.min.js?17.4.4"></script>
        </body>
    </html>
    `
	PlayerPageCode = `
    <!DOCTYPE html>
    <html>
        <head>
            <title>{{.Meta.Title}}</title>
            <meta name="keywords" content="{{.Meta.Keywords}}">
            <meta name="description" content="{{.Meta.Description}}">
            <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
            <style type="text/css">
                body {
                    margin: 0;
                    background-color: #000;
                }
                iframe {
                    display: block;
                    background: #000;
                    border: none;
                    height: 100vh;
                    width: 100vw;
                }
            </style>
        </head>
        <body>
            <iframe src="http://datalock.ru/player/{{.Meta.ID}}" allowfullscreen></iframe>
            <!-- Yandex.Metrika counter -->
            <script type="text/javascript">
                (function (d, w, c) {
                    (w[c] = w[c] || []).push(function() {
                        try {
                            w.yaCounter41725489 = new Ya.Metrika({
                                id:41725489,
                                clickmap:true,
                                trackLinks:true,
                                accurateTrackBounce:true
                            });
                        } catch(e) { }
                    });

                    var n = d.getElementsByTagName("script")[0],
                        s = d.createElement("script"),
                        f = function () { n.parentNode.insertBefore(s, n); };
                    s.type = "text/javascript";
                    s.async = true;
                    s.src = "https://mc.yandex.ru/metrika/watch.js";

                    if (w.opera == "[object Opera]") {
                        d.addEventListener("DOMContentLoaded", f, false);
                    } else { f(); }
                })(document, window, "yandex_metrika_callbacks");
            </script>
            <noscript>
                <div>
                    <img src="https://mc.yandex.ru/watch/41725489" style="position:absolute; left:-9999px;" alt="" />
                </div>
            </noscript>
            <!-- /Yandex.Metrika counter -->
        </body>
    </html>
    `
)

var (
	SecuredPlayerTemplate *template.Template
	PlayerPageTemplate    *template.Template
)

func ParseTemplates() error {
	var err error
	if PlayerPageTemplate, err = template.New("master").Parse(PlayerPageCode); err != nil {
		return err
	}
	if SecuredPlayerTemplate, err = template.New("master").Parse(SecuredPlayerPageCode); err != nil {
		return err
	}
	return nil
}
