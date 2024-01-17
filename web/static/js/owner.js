/**
 * Function to fetch the weight.
 * @param {string} endpoint - The API endpoint to fetch the weight.
 * @returns {Promise<number>} A promise that resolves to the capacity value as a number.
 */
async function getWeight(endpoint) {
    // Configure the HTTP request options for a GET request.
    const requestOptions = {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        },
    };

    try {
        // Send a GET request to the API endpoint for fetching the weight.
        const response = await fetch(endpoint, requestOptions);

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        // Use response.text() to get the response as a string
        const data = await response.text();

        // Use the 'data' value as needed (e.g., convert it to a number)
        const capacityValue = parseInt(data, 10);

        return capacityValue;
    } catch (error) {
        console.error('Error:', error);
        throw error; // Re-throw the error to be caught by the caller
    }
}

// Function to scale percentage values
function scalePercentage(percentage, maxScale) {
    if (percentage >= 0 && percentage <= 100) {
        return (percentage / 100) * maxScale;
    } else {
        console.error('Invalid percentage');
        return 0;
    }
}

// Function to update the height and top values
function updateBarStyle(element_id, percentage) {
    const bar = document.getElementById(element_id);

    // Validate and set the height
    if (percentage >= 0 && percentage <= 100) {
        // Scale the percentages to the desired range
        heightPercentage = scalePercentage(percentage, 96);
        topPercentage = 96 - heightPercentage;

        bar.style.height = heightPercentage + '%';
        bar.style.top = topPercentage + '%';
    } else {
        console.error('Invalid percentage');
    }
}
