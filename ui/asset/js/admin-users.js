function openCreditModal(userId, userName) {
    document.getElementById('modal-user-id').value = userId;
    document.getElementById('modal-user-name').value = userName;
    document.getElementById('credit-modal').classList.remove('hidden');
}
function closeCreditModal() {
    document.getElementById('credit-modal').classList.add('hidden');
}

function openEditUserModal(id, name, email, role, isActivated) {
    document.getElementById('edit-user-id').value = id;
    document.getElementById('edit-user-name').value = name;
    document.getElementById('edit-user-email').value = email;
    document.getElementById('edit-user-role').value = role;
    document.getElementById('edit-user-activated').value = isActivated;
    document.getElementById('editUserModal').classList.remove('hidden');
}

window.openCreditModal = openCreditModal;
window.closeCreditModal = closeCreditModal;
window.openEditUserModal = openEditUserModal;
