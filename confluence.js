const axios = require('axios');
const { confluenceConfig } = require("./credentials.json");

const apiKey = confluenceConfig.api_key;

class Confluence {
    constructor() {
        this.axiosConfig = {
            headers: {
                'Authorization': `Basic ${Buffer.from(apiKey).toString('base64')}`,
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            }
        }
    }
    async getByPageId(pageId) {
        const apiEndpoint = `${confluenceConfig.host}/api/content/${pageId}?expand=body.storage,version`;
        const { data } = await axios.get(apiEndpoint, this.axiosConfig);
        // extract necessary details
        let { type, title, body, version } = data;
        return { type, title, body, version };
    }
    async updateByPageId(pageId, sonarStats = []) {
        const apiEndpoint = `${confluenceConfig.host}/api/content/${pageId}?expand=body.storage`;
        // fetch Content
        const page = await this.getByPageId(pageId);
        // Update content with sonar stats
        let updatedContent = this.#generateHTML(page.body, sonarStats); // update content
        // Generate Request Body
        const bodyData = {
            title: page.title,
            type: page.type,
            version: {
                number: page.version.number + 1, message: "Updated by CronJob"
            },
            body: {
                storage: {
                    value: updatedContent,
                    representation: 'storage'
                }
            }
        }
        await axios.put(apiEndpoint, bodyData, this.axiosConfig);
        return true;
    }
    #generateHTML({ storage }, sonarStats) {
        const columns = [
            { name: 'Product', key: "name" },
            { name: 'Quality Gate', key: "alert_status" },
            { name: 'Code Smells', key: "code_smells" },
            { name: 'Bugs', key: "bugs" },
            { name: 'Vulnerabilities', key: "critical_severity_vulns" },
        ];

        let th = columns.map(c => {
            return `<th><p><strong>${c.name}</strong></p></th>`
        }).join("");

        let trs = sonarStats.map(stat => {
            let tds = columns.map(c => {
                let value = stat.measures[c.key];
                if (c.key == "name") value = stat.name;
                return `<td>${value == 'ERROR' ? 'Failed' : value}</td>`
            }).join("");
            return `<tr>${tds}</tr>`;
        }).join("");

        const html = `<table data-table-width="760" data-layout="default" ac:local-id="091ca39e-2b3b-4a0c-8720-7ee499fc6d65">
                        <tbody>
                            <tr>${th}</tr>
                            ${trs}
                        </tbody>
                    </table>`;
        return html;
    }
}

module.exports = new Confluence();