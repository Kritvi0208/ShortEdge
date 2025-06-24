// dashboard.js

document.addEventListener("DOMContentLoaded", async function () {
    const dashboard = document.getElementById("dashboard");
    const token = getToken(); // from token.js

    try {
        const res = await fetch("/all");

        if (!res.ok) {
            dashboard.innerHTML = `<p style="color:red;">❌ Failed to load links</p>`;
            return;
        }

        const data = await res.json();

        // Filter links by current browser's token
        const myLinks = data.filter(link => link.created_by === token);

        if (myLinks.length === 0) {
            dashboard.innerHTML = `<p>No links found yet. Create one above 👆</p>`;
            return;
        }

        // Create table to display links
        let html = `
            <table border="1" cellpadding="8" cellspacing="0">
              <tr>
                <th>Short Link</th>
                <th>Original URL</th>
                <th>Public</th>
                <th>Expires</th>
              </tr>
        `;

        myLinks.forEach(link => {
            const shortURL = `${window.location.origin}/${link.code}`;
            html += `
              <tr>
                <td><a href="${shortURL}" target="_blank">${shortURL}</a></td>
                <td>${link.original_url}</td>
                <td>${link.is_public ? '✅' : '❌'}</td>
                <td>${link.expires_at ? link.expires_at.split("T")[0] : 'Never'}</td>
              </tr>
            `;
        });

        html += `</table>`;
        dashboard.innerHTML = html;

    } catch (err) {
        console.error(err);
        dashboard.innerHTML = `<p style="color:red;">❌ Error loading dashboard: ${err.message}</p>`;
    }
});
