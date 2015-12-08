$("textarea[data-provide='markdown']").markdown({
    language: 'zh',
    fullscreen: {
        enable: true
    },
    resize: 'vertical',
    localStorage: 'md',
    imgurl: '/upload/ajaxUpload',
    base64url: '/upload/base64Upload'
});
