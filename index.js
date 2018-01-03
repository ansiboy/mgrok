

window.onhashchange = function () {
    loadContent();
}

function loadContent() {
    let file_name = !location.hash ? 'index' : location.hash.substr(1);
    $.get(file_name + '.md', function (text) {
        debugger;
        var html_content = marked(text);
        $('.container').html(html_content);
    })
}

loadContent();