document.addEventListener('DOMContentLoaded', function() {
    const flashcardForm = document.getElementById('flashcardForm');
    const flashcardList = document.getElementById('flashcardList');

    // Add event listener to the form submission
    flashcardForm.addEventListener('submit', function(event) {
        event.preventDefault();

        const deckIdInput = document.getElementById('deckId');
        const questionInput = document.getElementById('question');
        const answerInput = document.getElementById('answer');

        const deckId = deckIdInput.value.trim();
        const question = questionInput.value.trim();
        const answer = answerInput.value.trim();

        if (deckId === '' || question === '' || answer === '') {
            return;
        }

        // Create a new flashcard item and add it to the list
        const flashcardItem = document.createElement('li');
        flashcardItem.classList.add('flashcardItem');
        flashcardItem.innerHTML = `
            <div class="question">${question}</div>
            <div class="answer">${answer}</div>
        `;
        flashcardList.appendChild(flashcardItem);

        // Reset the form inputs
        deckIdInput.value = '';
        questionInput.value = '';
        answerInput.value = '';

        // Send the flashcard data to the server
        const flashcardData = {
            deckId: deckId,
            question: question,
            answer: answer
        };

        fetch('/api/flashcards', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(flashcardData)
        })
            .then(response => response.json())
            .then(data => console.log(data))
            .catch(error => console.error(error));
    });

    // Fetch and display existing flashcards on page load
    fetch('/api/flashcards')
        .then(response => response.json())
        .then(flashcards => {
            flashcards.forEach(flashcard => {
                const flashcardItem = document.createElement('li');
                flashcardItem.classList.add('flashcardItem');
                flashcardItem.innerHTML = `
                    <div class="question">${flashcard.Question}</div>
                    <div class="answer">${flashcard.Answer}</div>
                `;
                flashcardList.appendChild(flashcardItem);
            });
        })
        .catch(error => console.error(error));
});
