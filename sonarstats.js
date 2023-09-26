
const axios = require('axios');
const { sonarConfig } = require("./credentials.json");

const apiKey = sonarConfig.api_key;

class Sonar {
    #serialize(data) {
        const stats = {};
        data.measures.forEach(m => {
            stats[m.metric] = m
        });
        return stats;
    }
    async sonarStats(projectKey) {
        const apiEndpoint = `${sonarConfig.host}/api/measures/component?component=${projectKey}&metricKeys=${sonarConfig.metricKeys.join(",")}`;
        const axiosConfig = {
            method: "GET",
            headers: {
                'Authorization': `Basic ${Buffer.from(apiKey + ':').toString('base64')}`
            }
        };
        const { data: { component } } = await axios.get(apiEndpoint, axiosConfig);
        return component
    }
}

module.exports = new Sonar();