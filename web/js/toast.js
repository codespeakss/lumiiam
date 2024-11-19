// toast.js
export function showToast(message, duration = 5000) {
    const toast = document.getElementById("toast");
    if (!toast) {
        console.error("Toast element not found!");
        return;
    }

    toast.textContent = message; // Update toast message
    toast.classList.add("show");

    setTimeout(() => {
        toast.classList.remove("show");
    }, duration); // Toast will disappear after the specified duration
}
