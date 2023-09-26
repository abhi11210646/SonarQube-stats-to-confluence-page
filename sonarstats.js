
const axios = require('axios');
const { sonar } = require("./credentials.json");

const apiKey = sonar.api_key;

class Sonar {
    async sonarStats(projectKey) {
        const apiEndpoint = `${sonar.host}/api/measures/component?component=${projectKey}&metricKeys=quality_gate_details`;
        const axiosConfig = {
            method: "GET",
            headers: {
                'Authorization': `Basic ${Buffer.from(apiKey + ':').toString('base64')}`,
            }
        };
        const { data } = await axios.get(apiEndpoint, axiosConfig);
        return data.component;
    }
}

module.exports = new Sonar();