function getInitials(name) {
    const parts = name.split(' ');
    if (parts.length > 1) {
        // 英文姓名取首字母
        return parts[0].charAt(0).toUpperCase() + parts[1].charAt(0).toUpperCase();
    } else {
        // 中文姓名处理
        if (name.length > 2) {
            return name.charAt(0) + '.' + name.charAt(name.length - 1);
        } else {
            return name.charAt(0) + name.charAt(name.length - 1);
        }
    }
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}