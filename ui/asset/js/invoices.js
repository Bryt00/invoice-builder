function filterInvoices(status, btn) {
    var rows = document.querySelectorAll('#invoices-tbody tr[data-status]');
    rows.forEach(function(row) {
        if (status === 'all' || row.dataset.status === status) {
            row.style.display = '';
        } else {
            row.style.display = 'none';
        }
    });
    var tabs = document.querySelectorAll('.invoice-tab-btn');
    tabs.forEach(function(t) {
        t.classList.remove('bg-surface', 'text-primary', 'shadow-xs');
        t.classList.add('text-on-surface-variant');
    });
    if (btn) {
        btn.classList.add('bg-surface', 'text-primary', 'shadow-xs');
        btn.classList.remove('text-on-surface-variant');
    }
}
window.filterInvoices = filterInvoices;
