function toggleEditProfile(e) {
    if (e) e.preventDefault();
    var viewCard = document.getElementById('profile-view-card');
    var editForm = document.getElementById('profile-edit-form');
    if (viewCard && editForm) {
        if (viewCard.classList.contains('hidden')) {
            viewCard.classList.remove('hidden');
            editForm.classList.add('hidden');
        } else {
            viewCard.classList.add('hidden');
            editForm.classList.remove('hidden');
        }
    }
}
window.toggleEditProfile = toggleEditProfile;

document.addEventListener('DOMContentLoaded', function() {
    var editBtn = document.getElementById('edit-profile-btn');
    if (editBtn) {
        editBtn.onclick = toggleEditProfile;
    }
    var cancelBtn = document.getElementById('cancel-edit-btn');
    if (cancelBtn) {
        cancelBtn.onclick = toggleEditProfile;
    }
});
