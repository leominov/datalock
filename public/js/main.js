window.onload = function() {
    // https://goodies.pixabay.com/javascript/auto-complete/demo.html
    var xhrs = new XMLHttpRequest();
    new autoComplete({
        selector: 'input[name=search-text]',
        minChars: 3,
        source: function(term, suggest) {
            xhrs.open('GET', '/autocomplete.php?query=' + term, true);
            xhrs.onload = function() {
                var suggestions = [];
                if (xhrs.status != 200 || xhrs.responseText.length == "") {
                    return
                }
                result = JSON.parse(xhrs.responseText)
                for (i=0; i<result.suggestions.valu.length; i++) {
                    if (/serial\-/i.test(result.data[i]) == false) {
                        continue;
                    }
                    suggestions.push({
                        'title': result.suggestions.valu[i],
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
};
