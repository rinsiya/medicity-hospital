/* ==========================================================================
   MEDICITY DOCTOR DASHBOARD - INTERACTIVE SCRIPT
   ========================================================================== */

document.addEventListener('DOMContentLoaded', () => {

    // --- DOM Elements ---
    const notifBtn = document.getElementById('notifBtn');
    const notifDropdown = document.getElementById('notifDropdown');
    const editProfileBtn = document.getElementById('editProfileBtn');
    const editProfileModal = document.getElementById('editProfileModal');
    const closeProfileModal = document.getElementById('closeProfileModal');
    const cancelProfileModal = document.getElementById('cancelProfileModal');
    const profileEditForm = document.getElementById('profileEditForm');

    const editBioBtn = document.getElementById('editBioBtn');
    const dispBio = document.getElementById('dispBio');

    const dropzone = document.getElementById('dropzone');
    const fileInput = document.getElementById('fileInput');
    const selectedFileTag = document.getElementById('selectedFileTag');
    const selectedFileName = document.getElementById('selectedFileName');
    const cancelFileBtn = document.getElementById('cancelFileBtn');
    const uploadForm = document.getElementById('uploadForm');
    const remarksInput = document.getElementById('remarksInput');
    const documentList = document.getElementById('documentList');

    const addExpBtn = document.getElementById('addExpBtn');
    const addExpModal = document.getElementById('addExpModal');
    const closeExpModal = document.getElementById('closeExpModal');
    const cancelExpModal = document.getElementById('cancelExpModal');
    const expAddForm = document.getElementById('expAddForm');
    const experienceList = document.getElementById('experienceList');

    const completionFill = document.getElementById('completionFill');
    const completionText = document.getElementById('completionText');
    const completionBadge = document.getElementById('completionBadge');
    const signOutBtn = document.getElementById('signOutBtn');

    // State
    let profileCompletionPercentage = 85;

    // --- 1. Notification Dropdown Toggle ---
    notifBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        notifDropdown.classList.toggle('active');
    });

    document.addEventListener('click', () => {
        notifDropdown.classList.remove('active');
    });

    notifDropdown.addEventListener('click', (e) => e.stopPropagation());

    // --- 2. Edit Profile Modal Handler ---
    editProfileBtn.addEventListener('click', () => {
        editProfileModal.classList.add('active');
    });

    const closeProfileModalFunc = () => {
        editProfileModal.classList.remove('active');
    };

    closeProfileModal.addEventListener('click', closeProfileModalFunc);
    cancelProfileModal.addEventListener('click', closeProfileModalFunc);

    profileEditForm.addEventListener('submit', (e) => {
        e.preventDefault();

        // Get values
        const name = document.getElementById('editFullName').value;
        const role = document.getElementById('editRole').value;
        const email = document.getElementById('editEmail').value;
        const phone = document.getElementById('editPhone').value;
        const dept = document.getElementById('editDept').value;
        const exp = document.getElementById('editExp').value;
        const hospital = document.getElementById('editHospital').value;

        // Update DOM
        document.getElementById('dispFullName').textContent = name;
        document.getElementById('dispRole').textContent = role;
        document.getElementById('dispEmail').textContent = email;
        document.getElementById('dispPhone').textContent = phone;
        document.getElementById('dispDept').textContent = dept;
        document.getElementById('dispExp').textContent = exp;
        document.getElementById('dispHospital').textContent = hospital;

        // Update Preview Sidebar
        document.getElementById('prevName').textContent = name;
        document.getElementById('prevRole').textContent = role;

        closeProfileModalFunc();
        showToast('Registration details updated successfully!', 'success');
    });

    // --- 3. Editable Biography ---
    editBioBtn.addEventListener('click', () => {
        const currentBio = dispBio.textContent.trim();
        const newBio = prompt("Edit Professional Biography:", currentBio);
        if (newBio && newBio.trim() !== '') {
            dispBio.textContent = newBio.trim();
            showToast('Biography updated!', 'success');
        }
    });

    // --- 4. Drag & Drop File Upload ---
    dropzone.addEventListener('dragover', (e) => {
        e.preventDefault();
        dropzone.classList.add('dragover');
    });

    dropzone.addEventListener('dragleave', () => {
        dropzone.classList.remove('dragover');
    });

    dropzone.addEventListener('drop', (e) => {
        e.preventDefault();
        dropzone.classList.remove('dragover');
        if (e.dataTransfer.files.length > 0) {
            handleFileSelection(e.dataTransfer.files[0]);
        }
    });

    fileInput.addEventListener('change', (e) => {
        if (e.target.files.length > 0) {
            handleFileSelection(e.target.files[0]);
        }
    });

    function handleFileSelection(file) {
        selectedFileName.textContent = file.name;
        selectedFileTag.style.display = 'inline-flex';
    }

    cancelFileBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        fileInput.value = '';
        selectedFileTag.style.display = 'none';
    });

    uploadForm.addEventListener('submit', (e) => {
        e.preventDefault();
        const file = fileInput.files[0];
        if (!file && selectedFileTag.style.display === 'none') {
            showToast('Please select a file to upload.', 'info');
            return;
        }

        const fileName = file ? file.name : selectedFileName.textContent;
        const submitBtn = document.getElementById('submitUploadBtn');
        const originalBtnHTML = submitBtn.innerHTML;

        // Show uploading state
        submitBtn.innerHTML = '<i class="fa-solid fa-spinner fa-spin"></i> Uploading...';
        submitBtn.disabled = true;

        setTimeout(() => {
            // Append file to Uploaded Documents list
            const li = document.createElement('li');
            li.className = 'document-item';
            li.setAttribute('data-doc', fileName);
            li.innerHTML = `
                <div class="doc-info">
                    <i class="fa-solid fa-file-circle-check doc-icon"></i>
                    <span class="doc-name" title="${fileName}">${fileName}</span>
                </div>
                <div class="doc-actions">
                    <button class="btn-action download" title="Download Document" onclick="downloadSimulatedDoc('${fileName}')">
                        <i class="fa-solid fa-download"></i>
                    </button>
                    <button class="btn-action remove" title="Remove Document" onclick="removeDocument(this)">
                        <i class="fa-solid fa-xmark"></i>
                    </button>
                </div>
            `;
            documentList.appendChild(li);

            // Reset Form
            fileInput.value = '';
            remarksInput.value = '';
            selectedFileTag.style.display = 'none';
            submitBtn.innerHTML = originalBtnHTML;
            submitBtn.disabled = false;

            // Recalculate Completion
            updateCompletionProgress(100);

            showToast(`Document "${fileName}" uploaded successfully!`, 'success');
        }, 1200);
    });

    // --- 5. Add Professional Experience ---
    addExpBtn.addEventListener('click', () => {
        addExpModal.classList.add('active');
    });

    const closeExpModalFunc = () => {
        addExpModal.classList.remove('active');
    };

    closeExpModal.addEventListener('click', closeExpModalFunc);
    cancelExpModal.addEventListener('click', closeExpModalFunc);

    expAddForm.addEventListener('submit', (e) => {
        e.preventDefault();
        const role = document.getElementById('newExpRole').value.trim();
        const period = document.getElementById('newExpPeriod').value.trim();

        if (role && period) {
            const li = document.createElement('li');
            li.className = 'exp-item';
            li.innerHTML = `
                <div class="exp-info">
                    <strong class="exp-role">${role.toLowerCase()}</strong>
                    <span class="exp-period">${period}</span>
                </div>
                <button class="btn-remove-exp" onclick="removeExpItem(this)" title="Remove">
                    <i class="fa-solid fa-xmark"></i>
                </button>
            `;
            experienceList.appendChild(li);

            document.getElementById('newExpRole').value = '';
            document.getElementById('newExpPeriod').value = '';
            closeExpModalFunc();
            showToast('Professional experience added!', 'success');
        }
    });

    // --- Sign Out Action ---
    signOutBtn.addEventListener('click', () => {
        if (confirm('Are you sure you want to sign out of MediCity?')) {
            showToast('Signing out...', 'info');
            setTimeout(() => {
                alert('Signed out successfully.');
            }, 800);
        }
    });

    // Change Avatar Trigger
    document.getElementById('changeAvatarBtn').addEventListener('click', () => {
        const url = prompt("Enter Image URL for Profile Photo:", "doctor_avatar.png");
        if (url) {
            document.getElementById('profileAvatarImg').src = url;
            showToast('Profile photo updated!', 'success');
        }
    });

    // Progress bar updater
    function updateCompletionProgress(newVal) {
        profileCompletionPercentage = Math.min(newVal, 100);
        completionFill.style.width = `${profileCompletionPercentage}%`;
        completionText.textContent = `${profileCompletionPercentage}% Profile Completion`;

        if (profileCompletionPercentage === 100) {
            completionFill.style.backgroundColor = '#10b981';
            completionBadge.textContent = 'COMPLETE';
            completionBadge.className = 'status-badge complete';
            completionBadge.style.backgroundColor = '#d1fae5';
            completionBadge.style.color = '#059669';
        }
    }

});

// Global functions for inline onclick bindings
function removeDocument(btnElement) {
    const item = btnElement.closest('.document-item');
    const docName = item.getAttribute('data-doc') || 'Document';
    if (confirm(`Remove ${docName}?`)) {
        item.remove();
        showToast(`Removed ${docName}`, 'info');
    }
}

function removeExpItem(btnElement) {
    const item = btnElement.closest('.exp-item');
    item.remove();
    showToast('Experience entry removed', 'info');
}

function downloadSimulatedDoc(fileName) {
    showToast(`Downloading ${fileName}...`, 'info');
}

// Toast notification helper
function showToast(message, type = 'info') {
    const container = document.getElementById('toastContainer');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    
    const iconClass = type === 'success' ? 'fa-solid fa-circle-check' : 'fa-solid fa-circle-info';
    toast.innerHTML = `<i class="${iconClass}"></i> <span>${message}</span>`;
    
    container.appendChild(toast);
    
    setTimeout(() => {
        toast.remove();
    }, 4000);
}

