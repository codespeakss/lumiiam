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
function setCookie(name, value, options = {}) {
    let cookieString = `${encodeURIComponent(name)}=${encodeURIComponent(value)}`;
    if (options.expires) {
        cookieString += `; expires=${options.expires}`;
    }
    if (options.path) {
        cookieString += `; path=${options.path}`;
    }
    if (options.secure) {
        cookieString += `; Secure`;
    }
    if (options.sameSite) {
        cookieString += `; SameSite=${options.sameSite}`;
    }
    document.cookie = cookieString;
}

function getCookie(name) {
    const match = document.cookie.match(`(?:^|; )${encodeURIComponent(name)}=([^;]*)`);
    return match ? decodeURIComponent(match[1]) : null;
}