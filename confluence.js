const axios = require('axios');
const { confluence } = require("./credentials.json");

const apiKey = confluence.api_key;

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
        const apiEndpoint = `${confluence.host}/api/content/${pageId}?expand=body.storage,version`;
        const { data } = await axios.get(apiEndpoint, this.axiosConfig);
        // extract necessary details
        let { type, title, body, version } = data;
        return { type, title, body, version };
    }
    async updateByPageId(pageId) {
        const apiEndpoint = `${confluence.host}/api/content/${pageId}?expand=body.storage`;
        // fetch Content
        const page = await this.getByPageId(pageId);
        // Update content with sonar stats
        let updatedContent = page.body.storage.value; // update content
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
}

module.exports = new Confluence();