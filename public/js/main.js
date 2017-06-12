window.onload = function() {
    var sm = document.getElementById('secure_mark'),
        mk = document.getElementById('marketing');

    function Init() {
        // Setup HTML5 on start
        document.cookie = "playerHtml=true";
        GetMe();
    }

    function GetMe() {
        var xhr = new XMLHttpRequest();
        xhr.open('GET', '/me', true);
        xhr.onload = function() {
            if (xhr.status != 200) {
                return
            }
            var data = JSON.parse(xhr.responseText);
            if (data['secure_mark'].length == 0) {
                return
            }
            sm.value = data['secure_mark'];
            Show('.state-saved');
            Show('.marketing');
        }
        xhr.send();
    }

    function SetMe(sm) {
        var xhr = new XMLHttpRequest();
        Hide('.state-failed');
        Hide('.state-saved');
        xhr.open('POST', '/me', true);
        xhr.onload = function() {
            if (xhr.status != 200) {
                Show('.state-failed');
                return
            }
            Show('.state-saved');
        }
        xhr.onerror = function() {
            Show('.state-failed');
        }
        var json = JSON.stringify({
            secure_mark: sm,
        });
        xhr.send(json);
    }

    function Hide(selector) {
        e = document.querySelectorAll(selector);
        for (var i = 0; i < e.length; i++) {
            e[i].style.display = 'none';
        }
    }

    function Show(selector) {
        e = document.querySelectorAll(selector);
        for (var i = 0; i < e.length; i++) {
            e[i].style.display = 'inline';
        }
    }

    function Decode(y, t, u, p) {
        var lIll = 0;
        var ll1I = 0;
        var Il1l = 0;
        var ll1l = [];
        var l1lI = [];
        while (true) {
            if (lIll < 5) l1lI.push(y.charAt(lIll));
            else if (lIll < y.length) ll1l.push(y.charAt(lIll));
            lIll++;
            if (ll1I < 5) l1lI.push(t.charAt(ll1I));
            else if (ll1I < t.length) ll1l.push(t.charAt(ll1I));
            ll1I++;
            if (Il1l < 5) l1lI.push(u.charAt(Il1l));
            else if (Il1l < u.length) ll1l.push(u.charAt(Il1l));
            Il1l++;
            if (y.length + t.length + u.length + p.length == ll1l.length + l1lI.length + p.length) break;
        }
        var lI1l = ll1l.join('');
        var I1lI = l1lI.join('');
        ll1I = 0;
        var l1ll = [];
        for (lIll = 0; lIll < ll1l.length; lIll += 2) {
            var ll11 = -1;
            if (I1lI.charCodeAt(ll1I) % 2) ll11 = 1;
            l1ll.push(String.fromCharCode(parseInt(lI1l.substr(lIll, 2), 36) - ll11));
            ll1I++;
            if (ll1I >= l1lI.length) ll1I = 0;
        }
        return l1ll.join('');
    }

    function DecodeText(text) {
        p = text.match( /playlist\/([a-z0-9]{32})\//i );
        if (p != null && p.length == 2) {
            return p[1];
        }
        r = text.match( /\(\'([a-z0-9]+)\'\,\'([a-z0-9]+)\'\,\'([a-z0-9]+)\'\,\'([a-z0-9]+)\'/i );
        if (r == null || r.length != 5) {
            return
        }
        return DecodeText(Decode(r[1], r[2], r[3], r[4]));
    }

    function GetSecureMark() {
        var xhr = new XMLHttpRequest();
        xhr.open('GET', 'http://datalock.ru/player/1?t=' + Math.random(), false);
        xhr.onload = function() {
            if (xhr.status != 200 || xhr.responseText.length == "") {
                return
            }
            secure_mark = DecodeText(xhr.responseText);
            if (secure_mark != undefined && secure_mark != sm.value) {
                sm.value = secure_mark;
                SetMe(secure_mark);
            }
        }
        xhr.send(null);
    }

    submit_secure_mark.onclick = function() {
        SetMe(sm.value);
    };

    submit_try_get_secure_mark.onclick = function() {
        Show('.marketing');
        Show('.dont_work');
        GetSecureMark();
    };

    // https://goodies.pixabay.com/javascript/auto-complete/demo.html
    var xhrs = new XMLHttpRequest();
    new autoComplete({
        selector: 'input[name=search-text]',
        minChars: 3,
        source: function(term, suggest) {
            xhrs.open('GET', '/autocomplete.php?query=' + term, false);
            xhrs.onload = function() {
                var suggestions = [];
                if (xhrs.status != 200 || xhrs.responseText.length == "") {
                    return
                }
                result = JSON.parse(xhrs.responseText)
                for (i=0; i<result.suggestions.length; i++) {
                    suggestions.push({
                        'title': result.suggestions[i],
                        'link': result.data[i],
                        'id': result.id[i],
                    });
                }
                suggest(suggestions)
            }
            xhrs.send(null);
        },
        renderItem: function (item, search) {
            if (item.link.length == 0) {
                return '<div class="autocomplete-suggestion" data-link="/">'+item.title+'</div>';
            }
            return '<div class="autocomplete-suggestion" data-link="/'+item.link+'"><a href="'+item.link+'">'+item.title+'</a></div>';
        },
        onSelect: function(e, term, item) {
            link = item.getAttribute('data-link');
            if (link != "/") {
                location = link;
            }
        }
    });

    Init();
};
