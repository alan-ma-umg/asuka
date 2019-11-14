function loadScript(src, callback) {
    let script = document.createElement('script');
    script.src = src;
    script.onload = () => {
        script = null;
        callback && callback()
    };
    document.body.appendChild(script);
}

function ajax({url, method = "POST", data, headers = {}, timeout = 10000, success, error, complete}) {
    const xhr = new XMLHttpRequest();
    xhr.timeout = timeout;
    xhr.onreadystatechange = function () {
        if (this.readyState !== 4) {
            return
        }
        if (this.status === 200 || this.status === 304) {
            success && success(this);
        } else {
            error && error(this);
        }
        complete && complete(this);
    };
    xhr.open(method, url, true);
    xhr.setRequestHeader("X-Requested-With", "XMLHttpRequest");
    for (let k in headers) {
        xhr.setRequestHeader(k, headers[k]);
    }
    data ? xhr.send((typeof data == "string" || data instanceof FormData) ? data : JSON.stringify(data)) : xhr.send();
}

function timestampHumanReadable(timestamp) {
    if (timestamp < 60) {
        return timestamp + 's';
    }
    if (timestamp < 3600) {
        return timestamp / 60 + 'm';
    }
    if (timestamp < 3600) {
        return timestamp / 60 + 'm';
    }
    if (timestamp < 86400) {
        return timestamp / 3600 + 'h';
    }
    return timestamp / 86400 + 'd';
}

function showLoading() {
    document.getElementById('loading').style.display = 'block';
    document.getElementById('layer').style.display = 'block';
}

function hideLoading() {
    setTimeout(() => {
        document.getElementById('layer').style.display = 'none';
        document.getElementById('loading').style.display = 'none';
    }, 300);
}

function sendCommand(cmd, projectName) {
    showLoading();
    ajax({
        method: "POST", url: "/cmd", data: {"projectName": projectName, "cmd": cmd}, complete: () => {
            hideLoading()
        }, success: () => {
            if (typeof vueContent === "undefined") {
                return
            }

            let i = 0;
            let t = setInterval(() => {
                if (++i > 100) {
                    clearInterval(t);
                }
                if (vueContent.$data.payload.basic.log_mod >= vueContent.$data.payload.basic.log_check) {
                    showLoading();
                    ajax({
                        url: "/log", success: (res) => {
                            hideLoading();
                            popupWindow('<h1 class="text-green">Logging <small class="text-gray" ></small></h1>', res.response)
                        }
                    });
                    clearInterval(t);
                } else if (vueContent.$data.payload.basic.tcp_filter.LogMod >= vueContent.$data.payload.basic.tcp_filter.LogCheck) {
                    showLoading();
                    ajax({
                        url: "/log/tcp", success: (res) => {
                            hideLoading();
                            popupWindow('<h1 class="text-green">TCP Logging <small class="text-gray" ></small></h1>', res.response)
                        }
                    });
                    clearInterval(t);
                }
            }, 200)
        }
    })
}

function sendMessage(msg) {
    if (!ws) {
        return;
    }
    if (ws.readyState !== 1) {
        return;
    }
    ws.send(msg);
}

function reconnectSocket() {
    clearTimeout(timer_d90g8df987g9dfg7df9gdfj);
    if (manualFlag) {
        return;
    }
    timer_d90g8df987g9dfg7df9gdfj = setTimeout(() => {
        handlerSocket()
    }, 2000)
}

function manualOpen() {
    if (!ws) {
        return;
    }
    if (ws.readyState !== 3) {
        return;
    }
    manualFlag = false;
    handlerSocket()
}

function manualClose() {
    if (!ws) {
        return
    }
    if (ws.readyState !== 1) {
        return;
    }
    manualFlag = true;
    ws.close();
}

function handlerSocket() {
    try {
        document.title = "Asuka connecting...";
        ws = new WebSocket(wsUrl);
    } catch (e) {
        console.log(e);
        document.title = "Asuka exception";
        reconnectSocket();
        return
    }
    ws.onmessage = (evt) => {
        let data = JSON.parse(evt.data);
        vueContent.$data.payload = data;
        document.title = "Asuka " + data.basic.loads[5].toFixed(2) + " / " + data.basic.loads[60].toFixed(2) + " / " + data.basic.time;
        ws.send("");

        //chart
        if (data.basic.hasOwnProperty("loads")) {
            let cacheKey = "";
            for (let k in data.basic.loads) {
                cacheKey += data.basic.loads[k].toFixed(3) + k
            }
            if (cacheKey !== window.chartUpdatecacheCheck) {
                window.chartUpdatecacheCheck = cacheKey;
                lineChart(vueContent.$data.canvas, data.basic.loads, timestampHumanReadable, chartCeil);
            }
        }
    };
    ws.onopen = () => {
        sendMessage("");
        document.title = "Asuka connected";
    };
    ws.onerror = () => {
        document.title = "Asuka Error !";
        ws && ws.close();
        reconnectSocket()
    };
    ws.onclose = () => {
        document.title = "Asuka Closed !";
        ws && ws.close();
        reconnectSocket()
    };
}

// device detection
function isMobile() {
    if (window.isMobileCache === undefined) {
        if (/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|ipad|iris|kindle|Android|Silk|lge |maemo|midp|mmp|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows (ce|phone)|xda|xiino/i.test(navigator.userAgent)
            || /1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i.test(navigator.userAgent.substr(0, 4))) {
            window.isMobileCache = true
        } else {
            window.isMobileCache = false
        }
    }

    return window.isMobileCache;
}

Number.prototype.fileSizeH = function () {
    return fileSizeH(this.valueOf())
};

function fileSizeH(bytes) {
    if (bytes === 0) {
        return "0B";
    }
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return (bytes / Math.pow(1024, i)).toFixed(1) * 1 + ['B', 'K', 'M', 'G', 'T'][i];
}

function pad2(n) {
    return (n < 10 ? '0' : '') + n;
}

function chartCeil(n) {
    return Math.ceil(n * 1000) / 1000;
}

function chartTimeSince(t) {
    let dateTime = new Date();
    dateTime = dateTime.setSeconds(dateTime.getSeconds() - t);
    const d = new Date(dateTime);
    return pad2(d.getMinutes()) + ":" + pad2(d.getSeconds());
}

Number.prototype.timestamp2date = function () {
    const d = new Date(this.valueOf() * 1000);
    return pad2(d.getMonth() + 1) + "/" + pad2(d.getDate()) + "," + pad2(d.getHours()) + ":" + pad2(d.getMinutes()) + ":" + pad2(d.getSeconds());
};

Number.prototype.numFormat = function () {
    return numFormat(this.valueOf())
};

String.prototype.numFormat = function () {
    return numFormat(this.valueOf())
};

String.prototype.urlTruncate = function () {
    let str = this.valueOf();
    if (str.length > 40) {
        return ".." + str.substring(str.length - 40) + "(" + str.length + ")"
    }
    return str
};

function numFormat(v) {
    if (v === undefined || isNaN(v)) {
        return "0"
    }
    return new Intl.NumberFormat().format(v)
}

//listen for a link
document.addEventListener('click', (evt) => {
    const path = evt.path || (evt.composedPath && evt.composedPath());
    for (let i = 0; i < path.length; i++) {
        if (path[i] && path[i].tagName === 'A' && path[i].href.trim() !== "" && (path[i].target === undefined || path[i].target === "")) {
            evt.preventDefault();
            goToUrl(path[i].href);
            return
        }
    }
});

function goToUrl(dstUrl) {
    if (document.referrer && document.referrer.startsWith(location.origin) && document.referrer.replace(location.origin, '') === dstUrl.replace(location.origin, '') && ((dstUrl.startsWith("/") && !dstUrl.startsWith("//")) || dstUrl.indexOf(location.host) > -1)) {
        history.replaceState(dstUrl, document.title, dstUrl);
        history.back();// Browser cache
    } else {
        location.href = dstUrl
    }
}

function lineChart(canvasElement, loads, xConvert, yConvert) {
    canvasElement.width = canvasElement.offsetWidth;
    canvasElement.height = canvasElement.offsetHeight;

    let lineOffset = 50; //偏移量类似padding/margin的作用
    let context = canvasElement.getContext("2d");
    let lineCanvasWidth = canvasElement.width - lineOffset;
    let lineCanvasHeight = canvasElement.height - lineOffset;
    let minValue = Math.min(...Object.values(loads));
    let maxValue = Math.max(...Object.values(loads));
    let heightUnitPX = lineCanvasHeight / (maxValue - minValue);
    let len = (Object.values(loads).length);
    let widthUnitPx = lineCanvasWidth / (len - 1);
    let fontSize = 10;
    context.font = fontSize + "px 'open sans'";
    context.fillStyle = "#dadada";
    context.lineWidth = 0.5;
    context.strokeStyle = "#666666";

    let mod = 0;
    if (widthUnitPx < 100) {
        mod = Math.ceil(100 / widthUnitPx)
    }

    let i = 0;
    for (let k in loads) {
        //line chart
        let x = i * widthUnitPx, y = (maxValue - loads[k]) * heightUnitPX;
        x += lineOffset / 2;
        y += lineOffset / 2;

        // when minValue === maxValue
        if (heightUnitPX === Infinity) {
            y = canvasElement.height / 2
        }

        //line
        context.lineTo(x, y);

        if (mod === 0 || i % mod === 0 || len - 1 === i) {
            //dot
            context.arc(x, y, 1, 0, 2 * Math.PI);
            //y text
            context.fillText(yConvert ? yConvert(loads[k]) : loads[k], x - fontSize / 2, y - fontSize);
            //x text
            context.fillText(xConvert ? xConvert(k) : k, x - fontSize / 2, canvasElement.height - 1);
        }

        i++;
    }
    context.stroke();
}

function popupWindow(title, content) {
    const popupWindowEl = document.getElementById("popup-window");
    const popupWindowLayerEl = document.getElementById("popup-window-layer");
    const popupWindowTitleEl = document.getElementById("popup-window-title");
    const popupWindowContentEl = document.getElementById("popup-window-content");
    if (window.getComputedStyle(popupWindowEl).display === 'none') {
        popupWindowLayerEl.style.display = 'block';
        popupWindowTitleEl.innerHTML = title;
        popupWindowContentEl.innerHTML = content;
        popupWindowEl.style.display = "block";
        popupWindowContentEl.style.maxHeight = Math.min((window.innerHeight * .5), popupWindowContentEl.scrollHeight) + 'px';

        //run javascript
        Array.from(popupWindowContentEl.querySelectorAll("script")).forEach(oldScript => {
            const newScript = document.createElement("script");
            Array.from(oldScript.attributes)
                .forEach(attr => newScript.setAttribute(attr.name, attr.value));
            newScript.appendChild(document.createTextNode(oldScript.innerHTML));
            oldScript.parentNode.replaceChild(newScript, oldScript);
        });
    }
}

function popupWindowClose() {
    const popupWindowContentEl = document.getElementById("popup-window-content");
    while (popupWindowContentEl.firstChild) {
        popupWindowContentEl.removeChild(popupWindowContentEl.firstChild);
    }
    popupWindowContentEl.scrollTo(0, 0);

    document.getElementById("popup-window").style.display = "none";
    document.getElementById("popup-window-layer").style.display = 'none';
}