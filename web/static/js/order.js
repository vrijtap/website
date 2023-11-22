/**
 * Function to submit a payment request.
 * @param {string} ID - The unique identifier associated with the payment.
 */
function submitPayment(ID) {
    // Gather relevant data from the user input field.
    const userInput = document.getElementById('userInput').value;
    
    // Create a data object with payment information.
    const paymentData = {
        quantity: userInput,
        id: ID,
    };

    // Configure the HTTP request options.
    const requestOptions = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(paymentData),
    };
    
    // Send a POST request to the backend for payment processing.
    fetch('/order', requestOptions)
        .then(response => response.json()) // Parse response JSON data
        .then(data => {
            // Check if the response contains a 'url' field
            if (data.url) {
                // If the response contains a URL, redirect to it.
                window.location.href = data.url;
            } else {
                // Handle cases where the response does not contain a valid URL.
                console.error('Invalid URL in response:', data);
            }
        })
        .catch(error => {
            // Handle network errors or exceptions.
            console.error('Error:', error);
        });
}