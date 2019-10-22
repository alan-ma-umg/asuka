function loadScript(src, callback) {
    let script = document.createElement('script');
    script.src = src;
    script.onload = function () {
        script = null;
        callback && callback()
    };
    document.body.appendChild(script);
}

function ajax(option) {
    let url = option.url || '',
        method = option.method || 'POST',
        data = option.data,
        headers = option.headers || {},
        timeout = option.timeout || 10000,
        success = option.success,
        error = option.error;

    let xhr = new XMLHttpRequest();
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
    };
    xhr.open(method, url, true);
    for (let k in headers) {
        xhr.setRequestHeader(k, headers[k]);
    }
    data ? xhr.send(typeof data == "string" ? data : JSON.stringify(data)) : xhr.send();
}

function createConfig(labels) {
    return {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{label: '', data: [], borderColor: 'rgb(97, 100, 102)', borderWidth: 0.5, fill: false}]
        },
        options: {
            layout: {padding: {left: 10, right: 10, top: 5, bottom: 0}},
            elements: {point: {pointStyle: "crossRot"}, line: {tension: 0.2}},
            responsive: false,
            legend: {display: false},
            scales: {
                xAxes: [{
                    display: true,
                    gridLines: {display: false, drawBorder: false},
                    ticks: {fontColor: "#dadada", fontSize: 10,}
                }],
                yAxes: [{
                    display: true,
                    gridLines: {display: false, drawBorder: false},
                    ticks: {fontColor: "#dadada", fontSize: 10,}
                }]
            },
            title: {display: false,}
        }
    };
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
    timer_d90g8df987g9dfg7df9gdfj = setTimeout(function () {
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
        ws = new WebSocket(wsUrl);
    } catch (e) {
        console.log(e);
        reconnectSocket();
        return
    }
    ws.onmessage = function (evt) {
        let data = JSON.parse(evt.data);
        if (data.hasOwnProperty("projects")) {
            vueContent.$data.projects = data.projects;
        } else {
            vueContent.$data[data.type] = data;
        }
        vueContent.$data.basic = data.basic;
        if (data.basic.loads) {
            document.title = "Asuka " + data.basic.loads[5].toFixed(2) + " / " + data.basic.loads[60].toFixed(2) + " / " + data.basic.time;
        }
        ws.send(vueContent.$data.action);

        //chart
        if (vueContent.chart1) {
            vueContent.chart1.data.datasets[0].data = [];
            for (let s in data.basic.loads) {
                vueContent.chart1.data.datasets[0].data.push(data.basic.loads[s]);
            }
            vueContent.chart1.update();
        } else if (data.basic.loads && typeof Chart !== 'undefined') {
            let labels = [];
            for (let k in data.basic.loads) {
                labels.push(timestampHumanReadable(k));
            }
            vueContent.chart1 = new Chart(document.getElementById('chart-legend-top').getContext('2d'), createConfig(labels))
        }
    };
    ws.onopen = function () {
        sendMessage(vueContent.$data.action);
        document.body.className = "";
        document.title = "Asuka connected";
    };
    ws.onerror = function () {
        document.body.className = "ws-closed";
        document.title = "Asuka Error !";
        ws && ws.close();
        reconnectSocket()
    };
    ws.onclose = function () {
        document.body.className = "ws-closed";
        document.title = "Asuka Closed !";
        ws && ws.close();
        reconnectSocket()
    };
    window.onbeforeunload = function () {
        ws && ws.close()
    };
}
