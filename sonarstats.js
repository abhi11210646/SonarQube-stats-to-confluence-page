
const axios = require('axios');
const { sonarConfig } = require("./config");

const apiKey = sonarConfig.api_key;

class Sonar {
    #serialize(data) {
        const obj = {};
        for (let measure of data.measures) {
            obj[measure.metric] = measure.value;
        }
        data.measures = obj;
        return data;
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
        return this.#serialize(component);
    }
}

module.exports = new Sonar();