document.addEventListener('DOMContentLoaded', () => {
    const lineItemsContainer = document.getElementById('line-items-container');
    const addLineItemBtn = document.getElementById('add-line-item-btn');
    const subtotalEl = document.getElementById('calc-subtotal');
    const taxEl = document.getElementById('calc-tax');
    const discountEl = document.getElementById('calc-discount');
    const totalEl = document.getElementById('calc-total');

    const taxRateInput = document.getElementById('tax_rate');
    const discountInput = document.getElementById('discount_amount');
    const clientSelect = document.getElementById('client_id');
    const currencySelect = document.getElementById('currency');
    const currencyBadgeEl = document.getElementById('currency-badge');

    // Currency symbol map (mirrors internal/currency/currency.go)
    const currencySymbols = {
        USD: '$', EUR: '€', GBP: '£', GHS: 'GH₵', NGN: '₦',
        CAD: 'CA$', AUD: 'A$', JPY: '¥', INR: '₹',
        ZAR: 'R', BRL: 'R$', AED: 'AED'
    };

    function getCurrencySymbol() {
        const code = currencySelect ? currencySelect.value : 'USD';
        return currencySymbols[code] || code;
    }

    // Live preview elements
    const invNumInput = document.getElementById('invoice_number');
    const issueDateInput = document.getElementById('issue_date');
    const dueDateInput = document.getElementById('due_date');
    const clientEmailInput = document.getElementById('client_email');
    const clientAddressInput = document.getElementById('client_address');
    const notesInput = document.getElementById('notes');

    const prevInvNum = document.getElementById('prev-inv-num');
    const editorInvNum = document.getElementById('editor-inv-num');
    const prevIssueDate = document.getElementById('prev-issue-date');
    const prevDueDate = document.getElementById('prev-due-date');
    const prevClientEmail = document.getElementById('prev-client-email');
    const prevClientAddress = document.getElementById('prev-client-address');
    const prevNotes = document.getElementById('prev-notes');
    const prevLineItems = document.getElementById('prev-line-items');

    function syncLivePreview() {
        if (invNumInput && prevInvNum) {
            prevInvNum.textContent = invNumInput.value || 'INV-0000';
            if (editorInvNum) editorInvNum.textContent = invNumInput.value || 'INV-0000';
        }
        if (issueDateInput && prevIssueDate) prevIssueDate.textContent = issueDateInput.value || 'YYYY-MM-DD';
        if (dueDateInput && prevDueDate) prevDueDate.textContent = dueDateInput.value || 'YYYY-MM-DD';
        if (clientEmailInput && prevClientEmail) prevClientEmail.textContent = clientEmailInput.value || 'Select or enter client email';
        if (clientAddressInput && prevClientAddress) prevClientAddress.textContent = clientAddressInput.value || 'Billing address preview';
        if (notesInput && prevNotes) prevNotes.textContent = notesInput.value || 'Thank you for your business!';
    }

    // Autofill client details if saved client selected
    if (clientSelect) {
        clientSelect.addEventListener('change', () => {
            const selectedOption = clientSelect.options[clientSelect.selectedIndex];
            if (!selectedOption) return;
            const email = selectedOption.getAttribute('data-email') || '';
            const address = selectedOption.getAttribute('data-address') || '';

            if (clientEmailInput && email) clientEmailInput.value = email;
            if (clientAddressInput && address) clientAddressInput.value = address;
            syncLivePreview();
        });
    }

    [invNumInput, issueDateInput, dueDateInput, clientEmailInput, clientAddressInput, notesInput].forEach(el => {
        if (el) el.addEventListener('input', syncLivePreview);
    });

    function updateCurrencyBadge() {
        const code = currencySelect ? currencySelect.value : 'USD';
        const sym = getCurrencySymbol();
        if (currencyBadgeEl) currencyBadgeEl.textContent = code + ' (' + sym + ')';
    }

    function calculateTotals() {
        if (!lineItemsContainer) return;
        const sym = getCurrencySymbol();
        let subtotal = 0;
        let previewHtml = '';

        const rows = lineItemsContainer.querySelectorAll('.line-item-row');
        rows.forEach(row => {
            const descInput = row.querySelector('.item-desc') || row.querySelector('input[name="item_description[]"]');
            const qtyInput = row.querySelector('.item-qty');
            const priceInput = row.querySelector('.item-price');
            const rowTotalEl = row.querySelector('.item-total');

            const desc = descInput?.value || 'Service/Product Item';
            const qty = parseFloat(qtyInput?.value) || 0;
            const price = parseFloat(priceInput?.value) || 0;
            const rowTotal = qty * price;

            if (rowTotalEl) {
                rowTotalEl.textContent = sym + rowTotal.toFixed(2);
            }
            subtotal += rowTotal;

            previewHtml += `
                <tr>
                    <td class="py-2 font-medium">${desc}</td>
                    <td class="py-2 text-center">${qty}</td>
                    <td class="py-2 text-right">${sym}${price.toFixed(2)}</td>
                    <td class="py-2 text-right font-semibold">${sym}${rowTotal.toFixed(2)}</td>
                </tr>
            `;
        });

        if (prevLineItems) {
            prevLineItems.innerHTML = previewHtml || '<tr><td colspan="4" class="py-2 text-center text-outline">No items added yet.</td></tr>';
        }

        const taxRate = parseFloat(taxRateInput?.value) || 0;
        const discount = parseFloat(discountInput?.value) || 0;
        const taxAmount = subtotal * (taxRate / 100);
        const grandTotal = Math.max(0, subtotal + taxAmount - discount);

        if (subtotalEl) subtotalEl.textContent = sym + subtotal.toFixed(2);
        if (taxEl) taxEl.textContent = sym + taxAmount.toFixed(2);
        if (discountEl) discountEl.textContent = sym + discount.toFixed(2);
        if (totalEl) totalEl.textContent = sym + grandTotal.toFixed(2);

        updateCurrencyBadge();
    }

    if (lineItemsContainer) {
        lineItemsContainer.addEventListener('input', calculateTotals);
    }
    if (taxRateInput) taxRateInput.addEventListener('input', calculateTotals);
    if (discountInput) discountInput.addEventListener('input', calculateTotals);
    if (currencySelect) {
        currencySelect.addEventListener('change', () => {
            calculateTotals();
            updateCurrencyBadge();
        });
    }

    if (addLineItemBtn && lineItemsContainer) {
        addLineItemBtn.addEventListener('click', () => {
            const newRow = document.createElement('div');
            newRow.className = 'line-item-row grid grid-cols-12 gap-2 items-center p-2.5 bg-surface-container-low/50 rounded-xl border border-outline-variant/40';
            newRow.innerHTML = `
                <div class="col-span-6">
                    <input type="text" name="item_description[]" placeholder="Description of service or product" required
                           class="item-desc w-full px-2 py-1.5 bg-transparent border-0 border-b border-outline-variant/60 focus:border-primary outline-none font-body text-sm text-on-surface">
                </div>
                <div class="col-span-2">
                    <input type="number" name="item_quantity[]" value="1" min="0.01" step="0.01" required
                           class="item-qty w-full px-2 py-1.5 bg-transparent border-0 border-b border-outline-variant/60 focus:border-primary outline-none font-body text-sm text-on-surface text-center">
                </div>
                <div class="col-span-2">
                    <input type="number" name="item_price[]" value="0.00" min="0" step="0.01" required
                           class="item-price w-full px-2 py-1.5 bg-transparent border-0 border-b border-outline-variant/60 focus:border-primary outline-none font-body text-sm text-on-surface text-right">
                </div>
                <div class="col-span-2 flex items-center justify-between pl-1">
                    <span class="item-total font-headline text-xs font-semibold text-on-surface">$0.00</span>
                    <button type="button" class="remove-line-item p-1 text-outline hover:text-error transition-colors">
                        <span class="material-symbols-outlined text-[16px]">close</span>
                    </button>
                </div>
            `;
            lineItemsContainer.appendChild(newRow);
            calculateTotals();
        });

        lineItemsContainer.addEventListener('click', (e) => {
            const removeBtn = e.target.closest('.remove-line-item');
            if (removeBtn) {
                const rows = lineItemsContainer.querySelectorAll('.line-item-row');
                if (rows.length > 1) {
                    removeBtn.closest('.line-item-row').remove();
                    calculateTotals();
                }
            }
        });
    }

    const saveDraftBtn = document.getElementById('save-draft-btn');
    const invoiceForm = document.getElementById('invoice-form');
    const invoiceAction = document.getElementById('invoice-action');

    if (saveDraftBtn && invoiceForm && invoiceAction) {
        saveDraftBtn.addEventListener('click', () => {
            invoiceAction.value = 'draft';
            invoiceForm.submit();
        });
    }

    syncLivePreview();
    calculateTotals();
});
