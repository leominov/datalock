package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/leominov/datalock/seasonvar"
)

const (
	PlayerPageTemplate = `
    <!DOCTYPE html>
    <html>
        <head>
            <title>{{.Title}}</title>
            <meta name="keywords" content="{{.Keywords}}">
            <meta name="description" content="{{.Description}}">
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
            <iframe src="http://datalock.ru/player/{{.ID}}" allowfullscreen></iframe>
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

func IndexHandle(s *seasonvar.Seasonvar) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI := r.URL.RequestURI()
		seriesLink := s.AbsoluteLink(requestURI)
		if strings.Index(requestURI, ".html") == -1 {
			http.Redirect(w, r, seriesLink, http.StatusFound)
			return
		}
		seasonMeta, err := s.GetSeasonMeta(seriesLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		masterTmpl, err := template.New("master").Parse(PlayerPageTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("X-Cache", fmt.Sprintf("%d.%s", seasonMeta.CacheHitCounter, s.NodeName))
		err = masterTmpl.Execute(w, seasonMeta)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
