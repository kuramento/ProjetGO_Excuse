// server/static/script.js

document.getElementById('generateBtn').addEventListener('click', generateExcuse);
document.getElementById('addExcuseForm').addEventListener('submit', addExcuse);

function generateExcuse() {
    const button = document.getElementById('generateBtn');
    const excuseDiv = document.getElementById('excuse');
    const categorySelect = document.getElementById('categorySelect').value;

    // Désactiver le bouton et afficher un spinner pendant la requête
    button.disabled = true;
    button.innerHTML = `
        <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
        Chargement...
    `;
    excuseDiv.innerHTML = '';

    // Construire l'URL avec le paramètre de catégorie si sélectionné
    let url = '/api/excuse';
    if (categorySelect !== "") {
        url += `?category=${encodeURIComponent(categorySelect)}`;
    }

    fetch(url)
        .then(response => {
            if (!response.ok) {
                if (response.status === 404) {
                    throw new Error('Aucune excuse disponible pour cette catégorie.');
                }
                throw new Error('Erreur du serveur');
            }
            return response.json();
        })
        .then(data => {
            excuseDiv.innerHTML = `
                <div class="alert alert-success fade show" role="alert">
                    <strong>${data.category} :</strong> ${data.excuse}
                </div>
            `;
        })
        .catch(error => {
            excuseDiv.innerHTML = `
                <div class="alert alert-danger fade show" role="alert">
                    ${error.message}
                </div>
            `;
            console.error('Erreur:', error);
        })
        .finally(() => {
            // Réactiver le bouton et restaurer son texte
            button.disabled = false;
            button.innerHTML = '<i class="fa-solid fa-refresh"></i> Générer une excuse';
        });
}

function addExcuse(event) {
    event.preventDefault();

    const categoryInput = document.getElementById('excuseCategory');
    const excuseInput = document.getElementById('excuseText');
    const feedbackDiv = document.getElementById('addExcuseFeedback');

    const newExcuse = {
        category: categoryInput.value.trim(),
        excuse: excuseInput.value.trim()
    };

    // Validation rapide
    if (newExcuse.category === "" || newExcuse.excuse === "") {
        feedbackDiv.innerHTML = `
            <div class="alert alert-warning" role="alert">
                Veuillez remplir tous les champs.
            </div>
        `;
        return;
    }

    // Envoyer la requête POST pour ajouter l'excuse
    fetch('/api/excuse/add', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(newExcuse)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Erreur lors de l\'ajout de l\'excuse');
        }
        return response.json();
    })
    .then(data => {
        feedbackDiv.innerHTML = `
            <div class="alert alert-success" role="alert">
                Excuse ajoutée avec succès !
            </div>
        `;
        // Réinitialiser le formulaire
        document.getElementById('addExcuseForm').reset();
    })
    .catch(error => {
        feedbackDiv.innerHTML = `
            <div class="alert alert-danger" role="alert">
                Impossible d'ajouter l'excuse.
            </div>
        `;
        console.error('Erreur:', error);
    });
}
