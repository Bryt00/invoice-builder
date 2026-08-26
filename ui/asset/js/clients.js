function toggleAddClientForm(e) {
    if (e && e.preventDefault) e.preventDefault();
    var card = document.getElementById('add-client-card');
    if (card) {
        if (card.classList.contains('hidden')) {
            card.classList.remove('hidden');
        } else {
            card.classList.add('hidden');
        }
    }
}
window.toggleAddClientForm = toggleAddClientForm;

function toggleEditClient(clientId) {
    var editRow = document.getElementById('client-edit-' + clientId);
    if (editRow) {
        if (editRow.classList.contains('hidden')) {
            // Close all other edit rows first
            var allEdits = document.querySelectorAll('[id^="client-edit-"]');
            allEdits.forEach(function(row) { row.classList.add('hidden'); });
            editRow.classList.remove('hidden');
        } else {
            editRow.classList.add('hidden');
        }
    }
}
window.toggleEditClient = toggleEditClient;
