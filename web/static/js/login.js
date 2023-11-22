/**
 * Submits a login form via AJAX fetch request.
 * @param {Event} event - The form submission event.
 * @param {string} method - The method to use for the form
 */
async function submitLogin(event, method) {
    event.preventDefault();

    // Get form and error message element
    const form = event.target;
    const errorMessageElement = document.getElementById("error");

    try {
        // Send a fetch request to the server with the specified HTTP method and form data
        const response = await fetch(form.action, {
            method: method, // HTTP method (POST)
            body: new FormData(form), // Serialize form data
        });

        // Handle different response statuses
        if (response.status === 200) {
            // Redirect to the "/owner" page if the server responds with a 200 status code
            window.location.href = "/owner";
        } else if (response.status === 401) {
            // Display the error message if the response status is 401 (unauthorized)
            errorMessageElement.textContent = await response.text();
        } else if (response.status === 400) {
            // Parse response JSON for validation errors if the status is 400
            const errors = await response.json();
            const errorsArray = errors.errors;

            // Display error messages in the error message element
            errorMessageElement.textContent = errorsArray.join(", ");
        } else {
            // Log the response for debugging if it's not a 200, 401, or 400 status
            console.log(response);
        }
    } catch (error) {
        // Handle network errors during the fetch request
        console.error('Network error:', error);
    }
}
