/**
 * Created by chenning on 2015/12/2.
 */

function upload() {
    new AjaxUpload({
        fromId: "form",
        url: "/upload?timestamp=123456&nonce=123456&token=90f8d5RyvmxPeHJ4RL6h",
        method: "POST",
        timeout: 5000,
        sync: true,
        beforeSend: function () {
            document.getElementById("progress").style.display = "block";
        },
        onProgress: function (loaded, total) {
            var complete = (loaded / total * 100 | 0);
            var progress = document.getElementById("progress");
            progress.value = complete;
            progress.innerHTML = complete;
        },
        onComplete: function (result) {
            alert(result);
            window.location.reload();
        },
        onTimeout: function (event) {
            alert("timeout.");
        },
        onError: function (e) {
            alert(e);
        }
    });
}

var AjaxUpload = function (cfg) {
    if (!window.FormData) {
        alert("Can't support HTML5!");
        return;
    }

    this.isEmpty = function (obj) {
        if (!obj) {
            return true;
        }
        return false;
    };

    var cfg = cfg || {};
    if (this.isEmpty(cfg.fromId)) {
        return;
    }
    if (this.isEmpty(cfg.url)) {
        return;
    }

    this.id = cfg.fromId;
    this.method = cfg.method || "POST";
    this.url = cfg.url;
    this.async = !cfg.sync;
    this.resultType = cfg.resultType || "text";
    this.formData = new FormData(document.getElementById(this.id));
    this.xhr = new XMLHttpRequest();

    /** Progress **/
    this.xhr.upload.onprogress = function (event) {
        if (event.lengthComputable) {
            if (cfg.onProgress) {
                cfg.onProgress(event.loaded, event.total);
            }
        }
    };

    /** Synchronous requests must not set a timeout. **/
    if (!this.async) {
        console.log("setTimeout: " + cfg.timeout);
        if (cfg.timeout) {
            this.xhr.timeout = cfg.timeout;
        }
    }

    /** Timeout **/
    this.xhr.upload.ontimeout = function () {
        console.log("onTimeout: timeout.");
        if (cfg.onTimeout) {
            cfg.onTimeout();
        }
    };

    this.xhr.upload.onloadstart = function () {
        console.log("onLoadStart.");
    };

    /**  **/
    this.xhr.upload.onload = function () {
        setTimeout(function () {
            console.log("setTimeout: timeout.");
        }, 1000);
    };

    /** Error **/
    this.xhr.upload.onerror = function (e) {
        console.log("onError: " + e);
        if (cfg.onError) {
            cfg.onError(e);
        }
    };

    /** Ready State Change **/
    this.xhr.onreadystatechange = function () {
        console.log(xhr.readyState);
        if (xhr.readyState == 4) {
            console.log(xhr.status);
        }
        if (xhr.readyState === 4 && xhr.status === 200) {
            console.log(xhr.responseText);
            var res = xhr.responseText;
            if (xhr.resultType === 'json') {
                if ((typeof JSON) === 'undefine') {
                    res = eval("(" + res + ")");
                } else {
                    res = JSON.parse(res);
                }
            }
            console.log(res);
        }
    };

    /** Complete **/
    this.xhr.onload = function (event) {
        var res = event.target.responseText;
        if (this.resultType === 'json') {
            if ((typeof JSON) === 'undefine') {
                res = eval("(" + res + ")");
            } else {
                res = JSON.parse(res);
            }
        }
        if (cfg.onComplete) {
            cfg.onComplete(res);
        }
    };

    if (cfg.beforeSend) {
        cfg.beforeSend();
    }

    this.xhr.open(this.method, this.url, this.async);
    this.xhr.setRequestHeader("Cache-Control", "no-cache");
    this.xhr.setRequestHeader("X-Requested-With", "XMLHttpRequest");
    this.xhr.send(this.formData);

};
