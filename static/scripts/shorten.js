// shorten.js

document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("shorten-form");
    const resultDiv = document.getElementById("result");

    form.addEventListener("submit", async function (e) {
        e.preventDefault(); // stop form from reloading page

        const originalURL = document.getElementById("original-url").value;
        const customCode = document.getElementById("custom-code").value;
        const expiryDate = document.getElementById("expiry-date").value;
        const isPublic = document.getElementById("is-public").checked;
        const createdBy = getToken(); // from token.js

        const payload = {
            original_url: originalURL,
            custom_code: customCode || undefined,
            is_public: isPublic,
            created_by: createdBy,
        };

        if (expiryDate) {
            payload.expires_at = expiryDate;
        }

        try {
            const res = await fetch("/shorten", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(payload),
            });

            const data = await res.json();

            if (res.ok) {
                const fullLink = `${window.location.origin}/${data.code}`;
                resultDiv.innerHTML = `
                    <p>✅ Short link created:</p>
                    <input type="text" value="${fullLink}" readonly style="width: 300px;" />
                    <button onclick="navigator.clipboard.writeText('${fullLink}')">Copy</button>
                `;
            } else {
                resultDiv.innerHTML = `<p style="color: red;">❌ ${data.message || "Something went wrong"}</p>`;
            }
        } catch (err) {
            console.error(err);
            resultDiv.innerHTML = `<p style="color: red;">❌ Error: ${err.message}</p>`;
        }
    });
});
