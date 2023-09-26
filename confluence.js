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
        let updatedContent = this.#generateHTML(sonarStats); // update content
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
    #generateHTML(sonarStats) {
        //TODO: generate HTML

        return JSON.stringify(sonarStats);
    }
}

module.exports = new Confluence();