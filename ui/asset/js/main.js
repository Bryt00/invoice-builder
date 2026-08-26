document.addEventListener('DOMContentLoaded', () => {
    // 1. Password visibility toggle
    const toggleButtons = document.querySelectorAll('.password-toggle-btn');
    toggleButtons.forEach(btn => {
        btn.addEventListener('click', () => {
            const inputId = btn.getAttribute('data-target');
            if (!inputId) return;

            const inputEl = document.getElementById(inputId);
            const iconEl = btn.querySelector('span');

            if (inputEl && iconEl) {
                if (inputEl.type === 'password') {
                    inputEl.type = 'text';
                    iconEl.textContent = 'visibility_off';
                } else {
                    inputEl.type = 'password';
                    iconEl.textContent = 'visibility';
                }
            }
        });
    });

    // 2. Auto-disappearing Session Flash Alerts
    const flashAlerts = document.querySelectorAll('.flash-alert');
    flashAlerts.forEach(alert => {
        // Add smooth transition styles
        alert.style.transition = 'opacity 0.4s ease, transform 0.4s ease, max-height 0.4s ease, margin 0.4s ease, padding 0.4s ease';
        alert.style.opacity = '1';
        alert.style.maxHeight = '200px';

        // Add close button if not already present
        if (!alert.querySelector('.flash-close-btn')) {
            const closeBtn = document.createElement('button');
            closeBtn.type = 'button';
            closeBtn.className = 'flash-close-btn ml-auto opacity-70 hover:opacity-100 transition-opacity p-1 rounded-lg shrink-0 cursor-pointer';
            closeBtn.innerHTML = '<span class="material-symbols-outlined text-[18px]">close</span>';
            closeBtn.onclick = () => dismissFlash(alert);
            alert.appendChild(closeBtn);
        }

        // Auto-dismiss after 4 seconds
        setTimeout(() => {
            dismissFlash(alert);
        }, 4000);
    });

    function dismissFlash(element) {
        if (!element || element.dataset.dismissed) return;
        element.dataset.dismissed = 'true';
        element.style.opacity = '0';
        element.style.transform = 'translateY(-10px)';
        element.style.maxHeight = '0px';
        element.style.marginTop = '0px';
        element.style.marginBottom = '0px';
        element.style.paddingTop = '0px';
        element.style.paddingBottom = '0px';
        element.style.overflow = 'hidden';
        setTimeout(() => {
            element.remove();
        }, 450);
    }
});

// 3. Custom Confirmation Modal Handler (Alpine.js integration)
function showCustomConfirm(e, formElement, message, title) {
    if (e && e.preventDefault) e.preventDefault();
    
    window.dispatchEvent(new CustomEvent('open-confirm', {
        detail: {
            form: formElement,
            message: message || 'Are you sure you want to delete this item?',
            title: title || 'Confirm Action'
        }
    }));
    
    return false;
}
window.showCustomConfirm = showCustomConfirm;

// Navigation UI Logic
function toggleMobileNav() {
    const menu = document.getElementById('mobile-nav-menu');
    if (menu) menu.classList.toggle('hidden');
}
window.toggleMobileNav = toggleMobileNav;

function toggleAdminMobileNav() {
    const menu = document.getElementById('admin-mobile-nav-menu');
    if (menu) menu.classList.toggle('hidden');
}
window.toggleAdminMobileNav = toggleAdminMobileNav;

document.addEventListener('DOMContentLoaded', () => {
    const path = window.location.pathname;

    // Header nav (user_nav.html)
    document.querySelectorAll('header nav a').forEach(link => {
        const href = link.getAttribute('href');
        if (!href) return;
        
        let isActive = false;
        if (href === '/user/dashboard' && (path === '/user/dashboard' || path === '/user/dashboard/')) {
            isActive = true;
        } else if (href === '/user/invoices' && path.indexOf('/user/invoices') === 0) {
            isActive = true;
        } else if (href === '/user/clients' && path.indexOf('/user/clients') === 0) {
            isActive = true;
        } else if (href === '/user/profile/setup' && path.indexOf('/user/profile') === 0) {
            isActive = true;
        }

        if (isActive) {
            link.classList.remove('text-on-surface-variant', 'border-transparent');
            link.classList.add('text-primary', 'border-primary', 'font-bold');
        } else {
            link.classList.remove('text-primary', 'border-primary', 'font-bold');
            if (!link.classList.contains('text-on-surface-variant')) {
                link.classList.add('text-on-surface-variant');
            }
        }
    });

    // Sidebar nav (nav.html)
    document.querySelectorAll('nav a.sidebar-nav-item').forEach(link => {
        const href = link.getAttribute('href');
        if (!href) return;

        let isActive = false;
        if (href === '/user/dashboard' && (path === '/user/dashboard' || path === '/user/dashboard/')) {
            isActive = true;
        } else if (href === '/user/invoices' && path.indexOf('/user/invoices') === 0) {
            isActive = true;
        } else if (href === '/user/clients' && path.indexOf('/user/clients') === 0) {
            isActive = true;
        } else if (href === '/user/profile/setup' && path.indexOf('/user/profile') === 0) {
            isActive = true;
        }

        if (isActive) {
            link.classList.remove('text-on-surface-variant', 'hover:bg-surface-container-low');
            link.classList.add('bg-primary-container/30', 'text-primary', 'font-bold');
        } else {
            link.classList.remove('bg-primary-container/30', 'text-primary', 'font-bold');
            if (!link.classList.contains('text-on-surface-variant')) {
                link.classList.add('text-on-surface-variant');
            }
        }
    });

    // Admin sidebar nav active state
    document.querySelectorAll('.admin-sidebar-item, .admin-nav-link').forEach(link => {
        const href = link.getAttribute('href');
        if (!href) return;
        const isActive = path === href || (href !== '/admin/dashboard' && path.startsWith(href));
        if (isActive) {
            link.classList.add('bg-amber-500/15', 'text-amber-600', 'dark:text-amber-400', 'font-bold');
            link.classList.remove('text-on-surface-variant');
        } else {
            link.classList.remove('bg-amber-500/15', 'text-amber-600', 'dark:text-amber-400', 'font-bold');
            link.classList.add('text-on-surface-variant');
        }
    });

    // Auto-set progress bar widths from data-usage attributes
    document.querySelectorAll('[data-usage]').forEach(el => {
        const usage = el.getAttribute('data-usage');
        if (usage !== null) {
            el.style.width = usage + '%';
        }
    });
});
