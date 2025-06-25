// token.js

// Generates and stores a unique token in localStorage
function getToken() {
    let token = localStorage.getItem("urlify-token");

    if (!token) {
        token = crypto.randomUUID(); // creates a unique ID
        localStorage.setItem("urlify-token", token);
    }

    return token;
}
