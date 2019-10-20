function ajaxDo_(option) {
    let url = option.url || '',
        data = option.data,
        method = option.method || 'get',
        headers = option.headers || {},
        success = option.success,
        timeout = option.timeout || 10000,
        error = option.error;

    let xhr = new XMLHttpRequest();
    xhr.timeout = timeout;
    xhr.onreadystatechange = function () {
        if (this.readyState === 4) {
            if (this.status === 200 || this.status === 304) {
                success && success(this);
            } else {
                error && error(this);
            }
        }
    };
    xhr.open(method, url, true);
    for (let k in headers) {
        xhr.setRequestHeader(k, headers[k]);
    }
    if (data) {
        xhr.send(data);
    } else {
        xhr.send();
    }
}

function pixivHenTaiStart_(domain) {
    let bottomTimes = 0;
    let fsdfsdfgdfg = setInterval(function () {
        if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
            if (bottomTimes++ > 300) {
                clearInterval(fsdfsdfgdfg);
                hentaiStart();
                window.scrollTo(0, 0); //to top
            }
        } else {
            window.scrollTo(0, document.body.scrollHeight); //to bottom
            bottomTimes = 0;
        }
    }, 10);

    function hentaiStart() {
        let postJson = [];
        let urlParams = new URLSearchParams(window.location.search);
        let IllustId = urlParams.get('id');
        Array.from(document.body.querySelectorAll('img[src*="1200.jpg"]')).forEach(function (element) {
            postJson.push({
                url: element.src,
                illustId: IllustId,
            });
        });

        postJson.length && ajaxDo_({
            method: "POST",
            url: domain + "/project/pixiv/crawl/upload",
            data: JSON.stringify(postJson),
        });
    }
}
