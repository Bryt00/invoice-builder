document.addEventListener('DOMContentLoaded', function() {
    var btn = document.getElementById('copy-link-btn');
    var toast = document.getElementById('copy-toast');
    if (btn && toast) {
        btn.addEventListener('click', function() {
            var link = btn.getAttribute('data-link');
            if (!link) return;
            var fullUrl = window.location.origin + link;
            navigator.clipboard.writeText(fullUrl).then(function() {
                toast.style.opacity = '1';
                toast.style.pointerEvents = 'auto';
                setTimeout(function() {
                    toast.style.opacity = '0';
                    toast.style.pointerEvents = 'none';
                }, 2500);
            });
        });
    }
});
